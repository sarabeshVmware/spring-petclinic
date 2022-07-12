package common_tests

import (
	"strings"
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/models"
	"sigs.k8s.io/e2e-framework/pkg/env"
)

func Outerloop_scanning_supplychain_test(t *testing.T, testenv env.Environment, suiteConfig models.SuiteConfig, outerloopConfig models.OuterloopConfig, tapVersion string) {

	var scanPolicyFile = ""
	if strings.Split(tapVersion, ".")[1] == "1" {
		scanPolicyFile = outerloopConfig.ScanPolicy.OldYamlFile
	} else {
		scanPolicyFile = outerloopConfig.ScanPolicy.YamlFile
	}

	t.Log("************** TestOuterloopBasicSupplychainGitSource **************")
	testenv.Test(t,
		common_features.ChangeContext(t, suiteConfig.Multicluster.ViewClusterContext),
		common_features.UpdateDomainRecords(t),
		common_features.CreateGithubRepo(t, outerloopConfig.Project.Name, outerloopConfig.Project.RepoTemplate, outerloopConfig.Project.AccessToken),

		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.UpdateMetadataStoreScanning(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, tapVersion, "build", "testing_scanning", suiteConfig.Tap.Namespace, suiteConfig.Tap.PollTimeout, outerloopConfig.MetadataStore.Domain, suiteConfig.Multicluster.ViewClusterContext, suiteConfig.Multicluster.BuildClusterContext, outerloopConfig.MetadataStore.Namespace),
		common_features.ApplyKubectlConfigurationFile(t, outerloopConfig.SpringPetclinicPipeline.YamlFile, outerloopConfig.Namespace),
		common_features.ApplyKubectlConfigurationFile(t, scanPolicyFile, outerloopConfig.Namespace),
		common_features.TanzuDeployWorkload(t, outerloopConfig.Workload.TestYamlFile, outerloopConfig.Namespace),
		common_features.VerifyGitRepoStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifySourceScanStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyPipelineRunStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyBuildStatus(t, outerloopConfig.Workload.Name, suiteConfig.Innerloop.Workload.BuildNameSuffix, outerloopConfig.Namespace),
		common_features.VerifyPodIntentStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyTanzuWorkloadStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyImageScanStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
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
	)

	t.Log("************** TestOuterloopBasicSupplychainGitSource **************")
}

func Outerloop_scanning_supplychain_verify(t *testing.T, testenv env.Environment, suiteConfig models.SuiteConfig, outerloopConfig models.OuterloopConfig, tapVersion string) {

	t.Log("************** VerifyOuterloopBasicSupplychainGitSource **************")
	testenv.Test(t,

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
		common_features.DeleteNamespace(t, "metadata-store-secrets", suiteConfig.Multicluster.BuildClusterContext),
		common_features.DeletePipeline(t, outerloopConfig.SpringPetclinicPipeline.Name, outerloopConfig.Namespace),
		common_features.MulticlusterOuterloopCleanup(t, outerloopConfig.Workload.Name, outerloopConfig.Project.Name, outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext),
	)

	t.Log("************** VerifyOuterloopBasicSupplychainGitSource **************")
}

func Outerloop_scanning_supplychain_cleanup(t *testing.T, testenv env.Environment, suiteConfig models.SuiteConfig, outerloopConfig models.OuterloopConfig, tapVersion string) {

	t.Log("************** DeleteOuterloopBasicSupplychainGitSource **************")
	testenv.Test(t,
		common_features.DeleteGithubRepo(t, outerloopConfig.Project.Name, outerloopConfig.Project.AccessToken),
		common_features.DeleteNamespace(t, "metadata-store-secrets", suiteConfig.Multicluster.BuildClusterContext),
		common_features.DeletePipeline(t, outerloopConfig.SpringPetclinicPipeline.Name, outerloopConfig.Namespace),
		common_features.MulticlusterOuterloopCleanup(t, outerloopConfig.Workload.Name, outerloopConfig.Project.Name, outerloopConfig.Namespace, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext),
	)

	t.Log("************** DeleteOuterloopBasicSupplychainGitSource **************")
}
