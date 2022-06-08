//go:build all || multicluster_outerloop || multicluster_outerloop_basic

package multicluster_outerloop_test

import (
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
	"testing"
)

func TestOuterloopBasicSupplychainGitSource(t *testing.T) {
	t.Log("************** TestCase START: TestOuterloopBasicSupplychainGitSource **************")
	testenv.Test(t,
		common_features.CreateGithubRepo(t, outerloopConfig.Project.Name, outerloopConfig.Project.RepoTemplate, outerloopConfig.Project.AccessToken),

		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.TanzuDeployWorkload(t, outerloopConfig.Workload.YamlFile, outerloopConfig.Namespace),
		common_features.VerifyGitRepoStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyBuildStatus(t, outerloopConfig.Workload.Name, suiteConfig.Innerloop.Workload.BuildNameSuffix, outerloopConfig.Namespace),
		common_features.VerifyPodIntentStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyTanzuWorkloadStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyTaskRunStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Workload.TaskRunInfix, outerloopConfig.Namespace),

		//run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.ApplyKubectlConfigurationFile(t, outerloopConfig.Mysql.YamlFile, outerloopConfig.Namespace),

		//copying deliverable from build to run context
		common_features.ProcessDeliverable(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext, ""),

		//run context
		common_features.VerifyRevisionStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyKsvcStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyServiceBindingsStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Workload.ServiceBindingSuffix, outerloopConfig.Namespace),
		common_features.VerifyWorkloadResponse(t, outerloopConfig.Project.Host, outerloopConfig.Project.OriginalString, outerloopConfig.Project.WebpageRelativePath),

		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.UpdateGitRepository(t, outerloopConfig.Project.Username, outerloopConfig.Project.Email, outerloopConfig.Project.Repository, outerloopConfig.Project.Name, outerloopConfig.Project.AccessToken, outerloopConfig.Project.File, outerloopConfig.Project.OriginalString, outerloopConfig.Project.NewString, outerloopConfig.Project.CommitMessage),
		common_features.VerifyBuildStatusAfterUpdate(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),

		//run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.VerifyRevisionStatusAfterUpdate(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyKsvcStatusAfterUpdate(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyWorkloadResponse(t, outerloopConfig.Project.Host, outerloopConfig.Project.NewString, outerloopConfig.Project.WebpageRelativePath),

		common_features.DeleteGithubRepo(t, outerloopConfig.Project.Name, outerloopConfig.Project.AccessToken),
		common_features.MulticlusterOuterloopCleanup(t, outerloopConfig.Workload.Name, outerloopConfig.Project.Name, outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext),
	)

	t.Log("************** TestCase END: TestOuterloopBasicSupplychainGitSource **************")
}
