//go:build all || multicluster_outerloop || multicluster_outerloop_scan_multiapps

package multicluster_outerloop_scan_test

import (
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
)

func TestOuterloopScanSupplychainMultipleApps(t *testing.T) {
	t.Log("************** TestCase START: TestOuterloopScanSupplychainMultipleApps **************")
	testenv.Test(t,
		common_features.ChangeContext(t, suiteConfig.Multicluster.ViewClusterContext),
		common_features.UpdateDomainRecords(t),

		//build context
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.UpdateMetadataStoreScanning(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, "build", "testing_scanning", suiteConfig.Tap.Namespace, suiteConfig.Tap.PollTimeout, outerloopConfig.MetadataStore.Domain, suiteConfig.Multicluster.ViewClusterContext, suiteConfig.Multicluster.BuildClusterContext, outerloopConfig.MetadataStore.Namespace),
		common_features.ApplyKubectlConfigurationFile(t, outerloopConfig.BuildPacks.PipelineYamlFile, outerloopConfig.Namespace),
		common_features.ApplyKubectlConfigurationFile(t, outerloopConfig.BuildPacks.ScanPolicy, outerloopConfig.Namespace),
		common_features.TanzuDeployWorkloads(t, outerloopConfig),
		common_features.VerifyBuildPackWorkloadsGitrepoStatus(t, outerloopConfig),
		common_features.VerifyBuildPackWorkloadsTestTaskrunStatus(t, outerloopConfig),
		common_features.VerifyBuildPackWorkloadsSourceScanStatus(t, outerloopConfig),
		common_features.VerifyBuildPackWorkloadsBuildStatus(t, outerloopConfig),
		common_features.VerifyBuildPackWorkloadsImageScanStatus(t, outerloopConfig),
		common_features.VerifyBuildPackWorkloadsPodintents(t, outerloopConfig),

		//run context
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.ApplyKubectlConfigurationFile(t, outerloopConfig.Mysql.YamlFile, outerloopConfig.Namespace),

		//copying deliverable from build to run context
		common_features.ProcessDeliverableForBuildPackWorkloads(t, outerloopConfig, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext, ""),

		// //run context
		common_features.VerifyBuildPackWorkloadsRevisionStatus(t, outerloopConfig),
		common_features.VerifyBuildPackWorkloadsKsvcStatus(t, outerloopConfig),
		common_features.VerifyBuildPackWorkloadsReachability(t, outerloopConfig),
		common_features.ListBuildPackWorkloadsVulnerabilities(t, outerloopConfig, true, outerloopConfig.MetadataStore.Domain, suiteConfig.Multicluster.ViewClusterContext, suiteConfig.Multicluster.BuildClusterContext),
		common_features.VerifyBuildPackWorkloadsDataExistInMetadata(t, outerloopConfig),
		common_features.DeleteNamespace(t, "metadata-store-secrets", suiteConfig.Multicluster.BuildClusterContext),
		common_features.MulticlusterOuterloopCleanupforBuildPackWorkloads(t, outerloopConfig, suiteConfig.Multicluster.BuildClusterContext, suiteConfig.Multicluster.RunClusterContext),
	)
	t.Log("************** TestCase END: TestOuterloopScanSupplychainMultipleApps **************")
}
