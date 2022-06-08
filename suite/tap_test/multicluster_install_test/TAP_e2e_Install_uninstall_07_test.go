//go:build all || multicluster_install

package multicluster_install_tests

import (
	"testing"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
)

func TestTapInstallRunProfile(t *testing.T) {
	t.Log("************** TestCase START: TestTapInstallRunProfile **************")
	validatePackagesList := []string{"appliveview-connector", "appsso", "cartographer", "cnrs", "cert-manager", "contour", "fluxcd-source-controller", "image-policy-webhook", "ootb-delivery-basic", "ootb-templates", "service-bindings", "services-toolkit", "source-controller", "tap", "tap-auth", "tap-telemetry"}

	testenv.Test(t,
		common_features.InstallPackage(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, suiteConfig.Tap.Namespace, suiteConfig.Multicluster.RunTapValuesFile, suiteConfig.Tap.PollTimeout),
		common_features.ValidateListofInstalledPackage(t, suiteConfig.Tap.Namespace, validatePackagesList),
	)
	t.Log("************** TestCase END: TestTapInstallRunProfile **************")
}
