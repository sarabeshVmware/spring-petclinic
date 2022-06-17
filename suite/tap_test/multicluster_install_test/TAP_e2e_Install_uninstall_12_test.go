//go:build all || multicluster_install

package multicluster_install_tests

import (
	"testing"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
)

func TestTapUninstallFullProfilewithTbsSecret(t *testing.T) {
	t.Log("************** TestCase START: TestTapUninstallFullProfilewithTbsSecret **************")

	testenv.Test(t,
		common_features.DeletePackage(t, suiteConfig.Tap.Name, suiteConfig.Tap.Namespace),
		common_features.DeleteSecret(t, suiteConfig.TanzuNetCredentialsSecret.Name, suiteConfig.TanzuNetCredentialsSecret.Namespace),
	)
	t.Log("************** TestCase END: TestTapUninstallFullProfilewithTbsSecret **************")
}
