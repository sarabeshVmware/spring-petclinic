package common_features

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/git"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/imgpkg"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_helpers"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_libs"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzu_libs"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func UpdateGitRepository(t *testing.T, gitUsername string, gitEmail string, gitRepository string, projectName string, accessToken string, fileName string, originalString string, newString string, commitMessage string) features.Feature {
	return features.New("git-update").
		Assess("git-config", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("setting git config")

			// set git config
			err := git.GitConfig(gitUsername, gitEmail)
			if err != nil {
				t.Error("error while setting git config")
				t.FailNow()
			} else {
				t.Log("set git config")
			}

			return ctx
		}).
		Assess("git-clone", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("cloning git repo")

			// clone
			err := git.GitClone(rootDir, gitRepository)
			if err != nil {
				t.Error("error while cloning git repo")
				t.FailNow()
			} else {
				t.Log("cloned git repo")
			}

			return ctx
		}).
		Assess("git-seturl", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("setting git remote url")

			// set remote url
			err := git.GitSetRemoteUrl(filepath.Join(rootDir, projectName), accessToken, gitRepository)
			if err != nil {
				t.Error("error while setting git remote url")
				t.FailNow()
			} else {
				t.Log("set git remote url")
			}

			return ctx
		}).
		Assess("replace-string-in-file", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("replacing string in file")

			// replace string
			err := utils.ReplaceStringInFile(filepath.Join(rootDir, projectName, fileName), originalString, newString)
			if err != nil {
				t.Error("error while replacing string in file")
				t.FailNow()
			} else {
				t.Log("replaced string in file")
			}

			return ctx
		}).
		Assess("git-add", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("adding files to git index")

			// add files
			err := git.GitAdd(filepath.Join(rootDir, projectName), []string{fileName})
			if err != nil {
				t.Error("error while adding files to git index")
				t.FailNow()
			} else {
				t.Log("added files to git index")
			}

			return ctx
		}).
		Assess("git-commit", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("git committing index")

			// commit
			err := git.GitCommit(filepath.Join(rootDir, projectName), commitMessage)
			if err != nil {
				t.Error("error while committing git index")
				t.FailNow()
			} else {
				t.Log("committed git index")
			}

			return ctx
		}).
		Assess("git-push", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("pushing changes to repo")

			// push
			err := git.GitPush(filepath.Join(rootDir, projectName), false)
			if err != nil {
				t.Error("error while pushing changes to repo")
				t.Fail() // DON'T DO t.FailNow() AS WE WANT TO CLEAN UP REGARDLESS OF THE STATE OF THE TEST
			} else {
				t.Log("pushed changes to repo")
			}

			return ctx
		}).
		Assess("remove-dir", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			dir := filepath.Join(rootDir, projectName)
			t.Logf("removing directory %s", dir)
			err := os.RemoveAll(dir)
			if err != nil {
				t.Error(fmt.Errorf("error while removing directory %s: %w", dir, err))
				t.FailNow()
			}
			t.Logf("directory %s removed", dir)
			return ctx
		}).
		Feature()
}

func OuterloopCleanUp(t *testing.T, workloadName string, projectName string, namespace string) features.Feature {
	return features.New("outerloop cleanup").
		Assess("delete-workload", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("Deleting workload")
			tanzu_libs.DeleteWorkload(workloadName, namespace)
			return ctx
		}).
		Feature()
}

func MulticlusterOuterloopCleanup(t *testing.T, workloadName string, projectName string, namespace string, buildContext string, runContext string) features.Feature {
	return features.New("outerloop cleanup").
		Assess("delete-deliverable", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			// changing to run cluster
			_, err := kubectl_libs.UseContext(runContext)
			if err != nil {
				t.Errorf("error while changing context to %s", runContext)
				t.FailNow()
			} else {
				t.Logf("context changed to %s", runContext)
			}

			t.Logf("Deleting deliverable")
			kubectl_libs.DeleteDeliverable(workloadName, namespace)
			return ctx
		}).
		Assess("delete-workload", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			// changing to build cluster
			_, err := kubectl_libs.UseContext(buildContext)
			if err != nil {
				t.Errorf("error while changing context to %s", buildContext)
				t.FailNow()
			} else {
				t.Logf("context changed to %s", buildContext)
			}

			t.Logf("Deleting workload")
			tanzu_libs.DeleteWorkload(workloadName, namespace)
			return ctx
		}).
		Feature()
}

func DeletePipeline(t *testing.T, pipeline string, namespace string) features.Feature {
	return features.New("delete-pipeline").
		Assess("delete-pipeline", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("deleting pipeline")

			// delete pipeline
			_, err := kubectl_libs.DeletePipeline(pipeline, namespace)
			if err != nil {
				t.Error("error while deleting pipeline")
				t.Fail() // DON'T DO t.FailNow() AS WE WANT TO CLEAN UP REGARDLESS OF THE STATE OF THE TEST
			} else {
				t.Log("deleted pipeline")
			}
			return ctx
		}).
		Feature()
}

func VerifyRevisionStatus(t *testing.T, name string, namespace string) features.Feature {
	return features.New(fmt.Sprintf("verify-%s-revision-status", name)).
		Assess("verify-revision-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying revision ready status")

			revisionName = kubectl_helpers.GetLatestRevision(name, namespace, 1, 30)
			revisionReady := kubectl_helpers.ValidateRevisionStatus(revisionName, name, namespace, 5, 30)
			if !revisionReady {
				t.Error("revision not ready")
				t.FailNow()
			} else {
				t.Log("revision ready")
			}
			return ctx
		}).
		Feature()
}

func VerifyRevisionStatusAfterUpdate(t *testing.T, name string, namespace string) features.Feature {
	return features.New(fmt.Sprintf("verify-%s-revision-status", name)).
		Assess("verify-revision-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying revision ready status")

			revisionName = kubectl_helpers.GetNewerRevision(revisionName, name, namespace, 5, 30)
			revisionReady := kubectl_helpers.ValidateRevisionStatus(revisionName, name, namespace, 5, 30)
			if !revisionReady {
				t.Error("revision not ready")
				t.FailNow()
			} else {
				t.Log("revision ready")
			}
			return ctx
		}).
		Feature()
}

func VerifyKsvcStatus(t *testing.T, name string, namespace string) features.Feature {
	return features.New(fmt.Sprintf("verify-%s-ksvc-status", name)).
		Assess("verify-ksvc-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying ksvc ready status")

			ksvcReady := kubectl_helpers.VerifyKsvcStatus(name, namespace, revisionName, 5, 30)
			if !ksvcReady {
				t.Error("ksvc not ready")
				t.FailNow()
			} else {
				t.Log("ksvc ready")
			}

			return ctx
		}).
		Feature()
}

func VerifyKsvcStatusAfterUpdate(t *testing.T, name string, namespace string) features.Feature {
	return features.New(fmt.Sprintf("verify-%s-ksvc-status", name)).
		Assess("verify-ksvc-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying ksvc ready status")

			ksvcReady := kubectl_helpers.VerifyNewerKsvcStatus(name, namespace, revisionName, 5, 30)
			if !ksvcReady {
				t.Error("ksvc not ready")
				t.FailNow()
			} else {
				t.Log("ksvc ready")
			}

			return ctx
		}).
		Feature()
}

func VerifyServiceBindingsStatus(t *testing.T, name string, serviceBindingsSuffix string, namespace string) features.Feature {
	return features.New(fmt.Sprintf("verify-%s-ksvc-status", name)).
		Assess("verify-service-bindings-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying service bindings ready status")

			sbname := fmt.Sprintf("%[1]s-%[1]s%[2]s", name, serviceBindingsSuffix)
			if !kubectl_helpers.ValidateServiceBindings(sbname, namespace, 5, 30) {
				t.Error("service bindings not ready")
				t.FailNow()
			} else {
				t.Log("service bindings ready")
			}

			return ctx
		}).
		Feature()
}

func VerifyPipelineRunStatus(t *testing.T, name string, namespace string) features.Feature {
	return features.New(fmt.Sprintf("verify-%s-pipeline-status", name)).
		Assess("verify-pipeline-runs-succeded", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying pipeline runs status")

			// check
			pipelineRunSucceeded := kubectl_helpers.ValidatePipelineRuns(name, namespace, 5, 30)
			if !pipelineRunSucceeded {
				t.Error("pipeline runs not succeeded")
				t.FailNow()
			} else {
				t.Log("pipeline runs succeeded")
			}

			return ctx
		}).
		Feature()
}

func PatchServiceAccountSecrets(t *testing.T, sa string, namespace string, secret string) features.Feature {
	return features.New("patch-sa-secret").
		Assess("patch-sa-secret", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("Patching default sa secrets")
			res := kubectl_helpers.PatchServiceAccountWithNewSecret(sa, namespace, secret)
			if !res {
				t.Error("error while patching sa secret")
				t.Fail()
			} else {
				t.Log("patched sa secret")
			}
			return ctx
		}).
		Feature()
}

func VerifyImageskpac(t *testing.T, namespace string) features.Feature {
	return features.New("verify-images.kpac-status").
		Assess("verify-images.kpac-true", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verifying latest image status")

			// check
			if !kubectl_helpers.ValidateLatestImageStatus(namespace, 15, 60) {
				t.Error("image status is not true")
				t.FailNow()
			} else {
				t.Log("image status is true")
			}

			return ctx
		}).
		Feature()
}

func ImageCopyFromDeliverableToRepo(t *testing.T, name string, namespace string, target string) features.Feature {
	return features.New("image package copy to another repo").
		Assess("getting deliverable source image package", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			valid := kubectl_helpers.ValidateBuildClusterDeliverableStatus(name, namespace, 10, 30)
			if !valid {
				t.Errorf("error while getting deliverable %s", name)
				t.FailNow()
			} else {
				sourceRepo = kubectl_libs.GetDeliverables(name, namespace)[0].SOURCE
			}
			return ctx
		}).
		Assess("imgpkg copy to different registry", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			err := imgpkg.ImgpkgCopy(sourceRepo, target)
			if err != nil {
				t.Errorf("error while copying image bundles from %s to %s", sourceRepo, target)
				t.FailNow()
			} else {
				t.Logf("copied image bundles from %s to %s", sourceRepo, target)
			}
			return ctx
		}).Feature()
}
