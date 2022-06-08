//go:build all || multicluster_install

package multicluster_install_tests

import (
	"testing"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
)

func TestTapInstallBuildProfile(t *testing.T) {
	t.Log("************** TestCase START: TestTapInstallBuildProfile **************")
	validatePackagesList := []string{ "appliveview-conventions", "buildservice", "cartographer", "cert-manager", "contour", "conventions-controller", "fluxcd-source-controller","grype", "ootb-supply-chain-basic", "ootb-templates", "scanning", "source-controller", "spring-boot-conventions", "tap", "tap-auth", "tap-telemetry", "tekton-pipelines"}

	testenv.Test(t,
		common_features.InstallPackage(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, suiteConfig.Tap.Namespace, suiteConfig.Multicluster.BuildTapValuesFile, suiteConfig.Tap.PollTimeout),
		common_features.ValidateListofInstalledPackage(t, suiteConfig.Tap.Namespace, validatePackagesList),
	)
	t.Log("************** TestCase END: TestTapInstallBuildProfile **************")
}
