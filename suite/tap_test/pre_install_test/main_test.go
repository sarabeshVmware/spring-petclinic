package pre_install_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/envfuncs"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/models"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

var testenv env.Environment
var suiteConfig = models.SuiteConfig{}
var outerloopConfig = models.OuterloopConfig{}

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
	outerloopConfig, _ = models.GetOuterloopConfig()

	// setup
	testenv.Setup(
		envfuncs.InstallClusterEssentials(suiteConfig.TanzuClusterEssentials.TanzunetHost, suiteConfig.TanzuClusterEssentials.TanzunetApiToken, suiteConfig.TanzuClusterEssentials.ProductFileId, suiteConfig.TanzuClusterEssentials.ReleaseVersion, suiteConfig.TanzuClusterEssentials.ProductSlug, suiteConfig.TanzuClusterEssentials.DownloadBundle, suiteConfig.TanzuClusterEssentials.InstallBundle, suiteConfig.TanzuClusterEssentials.InstallRegistryHostname, suiteConfig.TanzuClusterEssentials.InstallRegistryUsername, suiteConfig.TanzuClusterEssentials.InstallRegistryPassword),
		envfuncs.CreateNamespaces(suiteConfig.CreateNamespaces),
	)

	// finish
	testenv.Finish(
		envfuncs.DeleteNamespaces(suiteConfig.CreateNamespaces),
	)

	os.Exit(testenv.Run(m))
}
