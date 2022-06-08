//go:build all || multicluster_install

package multicluster_install_tests

import (
	"testing"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
)

func TestTapInstallViewProfile(t *testing.T) {
	t.Log("************** TestCase START: TestTapInstallViewProfile **************")
	validatePackagesList := []string{"accelerator", "api-portal", "appliveview", "cert-manager", "contour", "fluxcd-source-controller", "learningcenter", "learningcenter-workshops", "metadata-store", "source-controller", "tap", "tap-gui", "tap-telemetry"}
	
	testenv.Test(t,
		common_features.InstallPackage(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, suiteConfig.Tap.Namespace, suiteConfig.Multicluster.ViewTapValuesFile, suiteConfig.Tap.PollTimeout),
		common_features.ValidateListofInstalledPackage(t, suiteConfig.Tap.Namespace, validatePackagesList),
	)
	t.Log("************** TestCase END: TestTapInstallViewProfile **************")
}
