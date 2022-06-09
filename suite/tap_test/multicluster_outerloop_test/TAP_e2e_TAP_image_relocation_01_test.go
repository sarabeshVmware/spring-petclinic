//go:build all || multicluster_outerloop || multicluster_image_relocation

package multicluster_outerloop_test

import (
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
)

func TestImageRelocation(t *testing.T) {
	t.Log("************** TestCase START: TestImageRelocation **************")
	testenv.Test(t,
		common_features.CreateGithubRepo(t, outerloopConfig.Project.Name, outerloopConfig.Project.RepoTemplate, outerloopConfig.Project.AccessToken),

		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.TanzuDeployWorkload(t, outerloopConfig.Workload.YamlFile, outerloopConfig.Namespace),
		common_features.VerifyBuildStatus(t, outerloopConfig.Workload.Name, suiteConfig.Innerloop.Workload.BuildNameSuffix, outerloopConfig.Namespace),
		common_features.VerifyPodIntentStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyTaskRunStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Workload.TaskRunInfix, outerloopConfig.Namespace),

		//imgpkg copy to target repo
		common_features.DockerLogin(t, suiteConfig.NonTanzuRepository[1].Server, suiteConfig.NonTanzuRepository[1].Username, suiteConfig.NonTanzuRepository[1].Password),
		common_features.ImageCopyFromDeliverableToRepo(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace, outerloopConfig.TestTargetRepo),

		//run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.ApplyKubectlConfigurationFile(t, outerloopConfig.Mysql.YamlFile, outerloopConfig.Namespace),
		common_features.DeleteSecret(t, suiteConfig.TapRegistrySecret.Name, outerloopConfig.Namespace),
		common_features.CreateSecret(t, suiteConfig.TapRegistrySecret.Name, suiteConfig.NonTanzuRepository[1].Server, suiteConfig.NonTanzuRepository[1].Username, suiteConfig.NonTanzuRepository[1].Password, suiteConfig.NonTanzuRepository[1].PasswordType, outerloopConfig.Namespace, false),

		//copying deliverable from build to run context
		common_features.ProcessDeliverable(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext, outerloopConfig.TestTargetRepo),

		//run context
		common_features.VerifyKsvcStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyWorkloadResponse(t, outerloopConfig.Project.Host, outerloopConfig.Project.OriginalString, outerloopConfig.Project.WebpageRelativePath),

		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.UpdateGitRepository(t, outerloopConfig.Project.Username, outerloopConfig.Project.Email, outerloopConfig.Project.Repository, outerloopConfig.Project.Name, outerloopConfig.Project.AccessToken, outerloopConfig.Project.File, outerloopConfig.Project.OriginalString, outerloopConfig.Project.NewString, outerloopConfig.Project.CommitMessage),
		common_features.VerifyBuildStatusAfterUpdate(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.ImageCopyFromDeliverableToRepo(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace, outerloopConfig.TestTargetRepo),

		//run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.VerifyKsvcStatusAfterUpdate(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyWorkloadResponse(t, outerloopConfig.Project.Host, outerloopConfig.Project.NewString, outerloopConfig.Project.WebpageRelativePath),

		common_features.DeleteGithubRepo(t, outerloopConfig.Project.Name, outerloopConfig.Project.AccessToken),
		common_features.MulticlusterOuterloopCleanup(t, outerloopConfig.Workload.Name, outerloopConfig.Project.Name, outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext),
	)

	t.Log("************** TestCase END: TestImageRelocation **************")
}
