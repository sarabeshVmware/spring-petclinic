//go:build all || multicluster_install

package multicluster_install_tests

import (
	"testing"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
)

func TestTapInstallIterateProfile(t *testing.T) {
	t.Log("************** TestCase START: TestTapInstallIterateProfile **************")
	validatePackagesList := []string{"appliveview", "appliveview-connector", "appsso", "appliveview-conventions", "buildservice", "cartographer", "cert-manager", "cnrs", "contour", "conventions-controller", "developer-conventions", "fluxcd-source-controller", "image-policy-webhook", "ootb-delivery-basic", "ootb-supply-chain-basic", "ootb-templates",  "policy-controller", "service-bindings", "services-toolkit", "source-controller", "spring-boot-conventions", "tap", "tap-auth", "tap-telemetry", "tekton-pipelines"}
	testenv.Test(t,
		common_features.InstallPackage(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, suiteConfig.Tap.Namespace, suiteConfig.Multicluster.IterateTapValuesFile, suiteConfig.Tap.PollTimeout),
		common_features.ValidateListofInstalledPackage(t, suiteConfig.Tap.Namespace, validatePackagesList),
	)
	t.Log("************** TestCase END: TestTapInstallIterateProfile **************")
}
