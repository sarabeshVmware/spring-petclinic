//go:build all || multicluster_upgrade

package multicluster_upgrade

import (
	// "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_tests"
	// "path/filepath"
	"testing"
)

func TestTapUpgrade(t *testing.T) {
	t.Log("************** TestCase START: TestTapUpgrade **************")

	//tap_1_0_2_values_file := filepath.Join(filepath.Join(utils.GetFileDir(), "../../resources/components"), "tap-values.yaml")
	//crud tests

	common_tests.Outerloop_scanning_supplychain_test(t, testenv, suiteConfig, outerloopConfig, suiteConfig.UpgradeVersions.TapVersion)
	//upgrading tap repo in all 3 clusters
	testenv.Test(t,
		common_features.ChangeContext(t, suiteConfig.Multicluster.ViewClusterContext),
		common_features.UpdatePackageRepository(t, suiteConfig.PackageRepository.Name, suiteConfig.UpgradeVersions.Image, suiteConfig.UpgradeVersions.UpgradeTapRepoVersion, suiteConfig.Tap.Namespace),
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.UpdatePackageRepository(t, suiteConfig.PackageRepository.Name, suiteConfig.UpgradeVersions.Image, suiteConfig.UpgradeVersions.UpgradeTapRepoVersion, suiteConfig.Tap.Namespace),
		common_features.ChangeContext(t, suiteConfig.Multicluster.RunClusterContext),
		common_features.UpdatePackageRepository(t, suiteConfig.PackageRepository.Name, suiteConfig.UpgradeVersions.Image, suiteConfig.UpgradeVersions.UpgradeTapRepoVersion, suiteConfig.Tap.Namespace),

		//upgrading tap in run cluster
		common_features.UpdateTapVersion(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Namespace, suiteConfig.Tap.ValuesSchemaFile, suiteConfig.UpgradeVersions.UpgradeTapVersion, suiteConfig.Tap.PollTimeout),
	)
	common_tests.Outerloop_scanning_supplychain_verify(t, testenv, suiteConfig, outerloopConfig, suiteConfig.UpgradeVersions.TapVersion)
	common_tests.Outerloop_scanning_supplychain_cleanup(t, testenv, suiteConfig, outerloopConfig, suiteConfig.UpgradeVersions.TapVersion)
	//crud tests
	common_tests.Outerloop_scanning_supplychain_test(t, testenv, suiteConfig, outerloopConfig, suiteConfig.UpgradeVersions.TapVersion)

	testenv.Test(t,
		//upgrading tap in build cluster
		common_features.ChangeContext(t, suiteConfig.Multicluster.BuildClusterContext),
		common_features.UpdateTapVersion(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Namespace, suiteConfig.Tap.ValuesSchemaFile, suiteConfig.UpgradeVersions.UpgradeTapVersion, suiteConfig.Tap.PollTimeout),
	)
	common_tests.Outerloop_scanning_supplychain_verify(t, testenv, suiteConfig, outerloopConfig, suiteConfig.UpgradeVersions.TapVersion)
	common_tests.Outerloop_scanning_supplychain_cleanup(t, testenv, suiteConfig, outerloopConfig, suiteConfig.UpgradeVersions.TapVersion)
	//crud tests
	common_tests.Outerloop_scanning_supplychain_test(t, testenv, suiteConfig, outerloopConfig, suiteConfig.UpgradeVersions.UpgradeTapVersion)

	testenv.Test(t,
		//upgrading tap in view cluster
		common_features.ChangeContext(t, suiteConfig.Multicluster.ViewClusterContext),
		common_features.UpdateTapVersion(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Namespace, suiteConfig.Tap.ValuesSchemaFile, suiteConfig.UpgradeVersions.UpgradeTapVersion, suiteConfig.Tap.PollTimeout),
	)
	common_tests.Outerloop_scanning_supplychain_verify(t, testenv, suiteConfig, outerloopConfig, suiteConfig.UpgradeVersions.UpgradeTapVersion)
	common_tests.Outerloop_scanning_supplychain_cleanup(t, testenv, suiteConfig, outerloopConfig, suiteConfig.UpgradeVersions.UpgradeTapVersion)
	//crud tests
	common_tests.Outerloop_scanning_supplychain_test(t, testenv, suiteConfig, outerloopConfig, suiteConfig.UpgradeVersions.UpgradeTapVersion)
	common_tests.Outerloop_scanning_supplychain_verify(t, testenv, suiteConfig, outerloopConfig, suiteConfig.UpgradeVersions.UpgradeTapVersion)
	common_tests.Outerloop_scanning_supplychain_cleanup(t, testenv, suiteConfig, outerloopConfig, suiteConfig.UpgradeVersions.UpgradeTapVersion)

	t.Log("************** TestCase END: TestTapUpgrade **************")
}
