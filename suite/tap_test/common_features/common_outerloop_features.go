package common_features

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/git"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzu_libs"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
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
		Feature()
}

func OuterloopCleanUp(t *testing.T, workloadName string, namespace string) features.Feature {
	return features.New("innerloop cleanup").
		Assess("delete-workload", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("Deleting workload")
			tanzu_libs.DeleteWorkload(workloadName, namespace)
			return ctx
		}).
		Assess("remove-dir", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			dir := filepath.Join(utils.GetFileDir(), workloadName)
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
