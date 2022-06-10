//go:build all || multicluster_install

package multicluster_install_tests

import (
	"testing"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
)

func TestTapUninstallBuildProfile(t *testing.T) {
	t.Log("************** TestCase START: TestTapUninstallBuildProfile **************")

	testenv.Test(t,
		common_features.DeletePackage(t, suiteConfig.Tap.Name, suiteConfig.Tap.Namespace),
	)
	t.Log("************** TestCase END: TestTapUninstallBuildProfile **************")
}
