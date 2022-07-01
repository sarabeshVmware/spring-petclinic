//go:build all || multicluster_install

package multicluster_install_tests

import (
	"testing"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
)

func TestTapInstallFullProfilewithTbsSecret(t *testing.T) {
	t.Log("************** TestCase START: TestTapInstallFullProfilewithTbsSecret **************")
	validatePackagesList := []string{"accelerator", "api-portal", "appsso", "appliveview", "appliveview-connector", "appliveview-conventions", "buildservice", "cartographer", "cert-manager", "cnrs", "contour", "conventions-controller", "developer-conventions", "fluxcd-source-controller", "grype", "image-policy-webhook", "learningcenter", "learningcenter-workshops", "metadata-store", "ootb-delivery-basic", "ootb-supply-chain-basic", "ootb-templates", "policy-controller", "scanning", "service-bindings", "services-toolkit", "source-controller", "spring-boot-conventions", "tap", "tap-auth", "tap-gui", "tap-telemetry", "tekton-pipelines"}
	
	testenv.Test(t,
		common_features.CreateSecret(t, suiteConfig.TanzuNetCredentialsSecret.Name, suiteConfig.TanzuNetCredentialsSecret.Registry, suiteConfig.TanzuNetCredentialsSecret.Username, suiteConfig.TanzuNetCredentialsSecret.Password, "string", suiteConfig.TanzuNetCredentialsSecret.Namespace, suiteConfig.TanzuNetCredentialsSecret.Export),
		common_features.InstallPackage(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, suiteConfig.Tap.Namespace, suiteConfig.Multicluster.FullProfilewithTbsSecretFile, suiteConfig.Tap.PollTimeout),
		common_features.ValidateListofInstalledPackage(t, suiteConfig.Tap.Namespace, validatePackagesList),
	)
	t.Log("************** TestCase END: TestTapInstallFullProfilewithTbsSecret **************")
}
