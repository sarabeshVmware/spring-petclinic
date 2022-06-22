package common_tests

import (
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/models"
	"sigs.k8s.io/e2e-framework/pkg/env"
)

// func Innerloop_multicluster_test(t *testing.T, testenv env.Environment, suiteConfig models.SuiteConfig) {
// 	t.Log("************** TestMulticlusterInnerloopBasicSupplychainGitSourceLiveUpdate **************")
// 	testenv.Test(t,
// 		// switch to Iterate cluster
// 		common_features.ChangeContext(t, suiteConfig.Multicluster.IterateClusterContext),
// 		common_features.TanzuDeployWorkload(t, suiteConfig.Innerloop.Workload.YamlFile, suiteConfig.Innerloop.Workload.Namespace),
// 		common_features.VerifyGitRepoStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
// 		common_features.VerifyTanzuJavaWebAppBuildStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.BuildNameSuffix, suiteConfig.Innerloop.Workload.Namespace),
// 		common_features.VerifyTanzuJavaWebAppImagesKpacStatus(t, suiteConfig.Innerloop.Workload.Namespace),
// 		common_features.VerifyTanzuJavaWebAppPodIntentStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
// 		common_features.VerifyTanzuJavaWebAppImageRepositoryDelivery(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.ImageDeliverySuffix, suiteConfig.Innerloop.Workload.Namespace),
// 		common_features.VerifyTanzuJavaWebAppDeliverable(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
// 		common_features.VerifyTanzuJavaWebAppRevisionStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
// 		common_features.VerifyTanzuJavaWebAppKsvcStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
// 		common_features.VerifyTanzuWorkloadStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
// 		common_features.VerifyWorkloadResponse(t, suiteConfig.Innerloop.Workload.URL, suiteConfig.Innerloop.Workload.OriginalString, ""),
// 		common_features.GitClone(t, suiteConfig.GitCredentials.Username, suiteConfig.GitCredentials.Email, suiteConfig.Innerloop.Workload.Gitrepository),
// 		common_features.UpdateTiltFile(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace, "iterate"),
// 		common_features.TiltUp(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
// 		common_features.VerifyTanzuWorkloadStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
// 		common_features.ReplaceStringInFile(t, suiteConfig.Innerloop.Workload.OriginalString, suiteConfig.Innerloop.Workload.NewString, suiteConfig.Innerloop.Workload.ApplicationFilePath, suiteConfig.Innerloop.Workload.Name),
// 		common_features.VerifyWorkloadResponse(t, suiteConfig.Innerloop.Workload.URL, suiteConfig.Innerloop.Workload.NewString, ""),
// 		common_features.InnerloopCleanUp(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
// 	)
// 	t.Log("************** TestMulticlusterInnerloopBasicSupplychainGitSourceLiveUpdate **************")
// }

func Outerloop_basic_supplychain_test(t *testing.T, testenv env.Environment, suiteConfig models.SuiteConfig, outerloopConfig models.OuterloopConfig) {
	t.Log("************** TestOuterloopBasicSupplychainGitSource **************")
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

	t.Log("************** TestOuterloopBasicSupplychainGitSource **************")
}

func Outerloop_testing_supplychain_test(t *testing.T, testenv env.Environment, suiteConfig models.SuiteConfig, outerloopConfig models.OuterloopConfig) {
	t.Log("************** TestOuterloopBasicSupplychainGitSource **************")
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

	t.Log("************** TestOuterloopBasicSupplychainGitSource **************")
}

func Outerloop_scanning_supplychain_test(t *testing.T, testenv env.Environment, suiteConfig models.SuiteConfig, outerloopConfig models.OuterloopConfig) {
	t.Log("************** TestOuterloopBasicSupplychainGitSource **************")
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

	t.Log("************** TestOuterloopBasicSupplychainGitSource **************")
}
