//go:build all || multicluster_innerloop || multicluster_innerloop_basic_git_source
package multicluster

import (
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
	"testing"
)

func TestInnerloopBasicSupplychainGitSourceLiveUpdate(t *testing.T) {
	t.Log("************** TestCase START: TestMulticlusterInnerloopBasicSupplychainGitSourceLiveUpdate **************")
	testenv.Test(t,
		// switch to Iterate cluster
		common_features.ChangeContext(t, suiteConfig.Multicluster.IterateClusterContext),
		common_features.TanzuDeployWorkload(t, suiteConfig.Innerloop.Workload.YamlFile, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyTanzuJavaWebAppGitRepository(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyTanzuJavaWebAppBuildStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.BuildNameSuffix, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyTanzuJavaWebAppImagesKpacStatus(t, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyTanzuJavaWebAppPodIntentStatus(t, suiteConfig.Innerloop.Workload.Name,  suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyTanzuJavaWebAppImageRepositoryDelivery(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.ImageDeliverySuffix, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyTanzuJavaWebAppDeliverable(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyTanzuJavaWebAppRevisionStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyTanzuJavaWebAppKsvcStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyTanzuWorkloadStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyWorkloadResponse(t, suiteConfig.Innerloop.Workload.URL, suiteConfig.Innerloop.Workload.OriginalString, ""),
		common_features.GitClone(t, suiteConfig.GitCredentials.Username, suiteConfig.GitCredentials.Email, suiteConfig.Innerloop.Workload.Gitrepository),
		common_features.UpdateTiltFile(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace, "iterate"),
		common_features.TiltUp(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyTanzuWorkloadStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.ReplaceStringInFile(t, suiteConfig.Innerloop.Workload.OriginalString, suiteConfig.Innerloop.Workload.NewString, suiteConfig.Innerloop.Workload.ApplicationFilePath, suiteConfig.Innerloop.Workload.Name),
		common_features.VerifyWorkloadResponse(t, suiteConfig.Innerloop.Workload.URL, suiteConfig.Innerloop.Workload.NewString, ""),
		common_features.InnerloopCleanUp(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
	)
	t.Log("************** TestCase END: TestMulticlusterInnerloopBasicSupplychainGitSourceLiveUpdate **************")
}