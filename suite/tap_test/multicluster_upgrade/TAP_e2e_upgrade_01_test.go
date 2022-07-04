//go:build all || multicluster_upgrade

package multicluster_upgrade

import (
	//"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_tests"
	// "path/filepath"
	"testing"
)

func TestTapUpgrade(t *testing.T) {
	t.Log("************** TestCase START: TestTapUpgrade **************")

	//tap_1_0_2_values_file := filepath.Join(filepath.Join(utils.GetFileDir(), "../../resources/components"), "tap-values.yaml")
	//crud tests
	common_tests.Outerloop_scanning_supplychain_test(t, testenv, suiteConfig, outerloopConfig)

	//upgrading tap repo in all 3 clusters
	common_features.ChangeContext(t, suiteConfig.Multicluster.ViewClusterContext),
	common_featues.Updso atePackageRepository(t, suiteConfig.PackageRepository.Name, suiteConfig.UpgradeVersions.Image, suiteConfig.UpgradeVersions.UpgradeTapRepoVersion, suiteConfig.Tap.Namespace),
	common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
	common_featues.UpdatePackageRepository(t, suiteConfig.PackageRepository.Name, suiteConfig.UpgradeVersions.Image, suiteConfig.UpgradeVersions.UpgradeTapRepoVersion, suiteConfig.Tap.Namespace),
	common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
	common_featues.UpdatePackageRepository(t, suiteConfig.PackageRepository.Name, suiteConfig.UpgradeVersions.Image, suiteConfig.UpgradeVersions.UpgradeTapRepoVersion, suiteConfig.Tap.Namespace),

	//upgrading tap in run cluster
	common_features.UpdateTapVersion(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Namespace, suiteConfig.Tap.ValuesSchemaFile, suiteConfig.UpgradeVersions.UpgradeTapVersion, suiteConfig.Tap.PollTimeout),

	//crud tests
	common_tests.Outerloop_scanning_supplychain_test(t, testenv, suiteConfig, OuterloopConfig)

	//upgrading tap in build cluster
	common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
	common_features.UpdateTapVersion(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Namespace, suiteConfig.Tap.ValuesSchemaFile, suiteConfig.UpgradeVersions.UpgradeTapVersion, suiteConfig.Tap.PollTimeout),

	//crud tests
	common_tests.Outerloop_scanning_supplychain_test(t, testenv, suiteConfig, OuterloopConfig)

	//upgrading tap in view cluster
	common_features.ChangeContext(t, suiteConfig.Multicluster.ViewClusterContext),
	common_features.UpdateTapVersion(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Namespace, suiteConfig.Tap.ValuesSchemaFile, suiteConfig.UpgradeVersions.UpgradeTapVersion, suiteConfig.Tap.PollTimeout),

	//crud tests
	common_tests.Outerloop_scanning_supplychain_test(t, testenv, suiteConfig, OuterloopConfig)

	// 	common_features.DeletePackage(t, suiteConfig.Tap.Name, suiteConfig.Tap.Namespace),
	)
	t.Log("************** TestCase END: TestTapUpgrade **************")
}
