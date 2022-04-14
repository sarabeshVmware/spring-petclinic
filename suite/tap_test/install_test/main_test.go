package install_tests

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/envfuncs"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/models"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

var testenv env.Environment
var suiteConfig = models.SuiteConfig{}
var outerloopConfig = models.OuterloopConfig{}
var suiteResourcesDir = filepath.Join(utils.GetFileDir(), "../../resources/suite")

func TestMain(m *testing.M) {
	// set logger
	logFile, err := utils.SetLogger(filepath.Join(utils.GetFileDir(), "logs"))
	if err != nil {
		log.Fatal(fmt.Errorf("error while setting log file %s: %w", logFile, err))
	}

	home, _ := os.UserHomeDir()
	cfg, _ := envconf.NewFromFlags()
	cfg.WithKubeconfigFile(filepath.Join(home, ".kube", "config"))
	testenv = env.NewWithConfig(cfg)

	// read suite config
	suiteConfig = models.GetSuiteConfig()
	outerloopConfig, err = models.GetOuterloopConfig()
	developerNamespaceFile := filepath.Join(suiteResourcesDir, "developer-namespace.yaml")

	// setup
	testenv.Setup(
		envfuncs.InstallClusterEssentials(suiteConfig.TanzuClusterEssentials.TanzunetHost, suiteConfig.TanzuClusterEssentials.TanzunetApiToken, suiteConfig.TanzuClusterEssentials.ProductFileId, suiteConfig.TanzuClusterEssentials.ReleaseVersion, suiteConfig.TanzuClusterEssentials.ProductSlug, suiteConfig.TanzuClusterEssentials.DownloadBundle, suiteConfig.TanzuClusterEssentials.InstallBundle, suiteConfig.TanzuClusterEssentials.InstallRegistryHostname, suiteConfig.TanzuClusterEssentials.InstallRegistryUsername, suiteConfig.TanzuClusterEssentials.InstallRegistryPassword),
		envfuncs.AddFinalizersToKappControllerClusterRole(),
		envfuncs.CreateNamespaces(suiteConfig.CreateNamespaces),
		envfuncs.CreateSecret(suiteConfig.TapRegistrySecret.Name, suiteConfig.TapRegistrySecret.Registry, suiteConfig.TapRegistrySecret.Username, suiteConfig.TapRegistrySecret.Password, suiteConfig.TapRegistrySecret.Namespace, suiteConfig.TapRegistrySecret.Export),
		envfuncs.CreateSecret(suiteConfig.RegistryCredentialsSecret.Name, suiteConfig.RegistryCredentialsSecret.Registry, suiteConfig.RegistryCredentialsSecret.Username, suiteConfig.RegistryCredentialsSecret.Password, suiteConfig.RegistryCredentialsSecret.Namespace, suiteConfig.RegistryCredentialsSecret.Export),
		envfuncs.AddPackageRepository(suiteConfig.PackageRepository.Name, suiteConfig.PackageRepository.Image, suiteConfig.PackageRepository.Namespace),
		envfuncs.CheckIfPackageRepositoryReconciled(suiteConfig.PackageRepository.Name, suiteConfig.PackageRepository.Namespace, 10, 60),
		envfuncs.SetupDeveloperNamespace(developerNamespaceFile, suiteConfig.CreateNamespaces[0]),
		envfuncs.CreateClusterRoleBinding(),
	)

	// finish
	testenv.Finish(
		envfuncs.DeletePackageRepository(suiteConfig.PackageRepository.Name, suiteConfig.PackageRepository.Namespace),
		envfuncs.DeleteSecret(suiteConfig.RegistryCredentialsSecret.Name, suiteConfig.RegistryCredentialsSecret.Namespace),
		envfuncs.DeleteSecret(suiteConfig.TapRegistrySecret.Name, suiteConfig.TapRegistrySecret.Namespace),
		envfuncs.DeleteNamespaces(suiteConfig.CreateNamespaces),
	)

	os.Exit(testenv.Run(m))
}
