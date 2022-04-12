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

	// setup
	testenv.Setup(
		envfuncs.InstallClusterEssentials(suiteConfig.TanzuClusterEssentials.Bundle, suiteConfig.TanzuClusterEssentials.Registry, suiteConfig.TapRegistrySecret.Username, suiteConfig.TapRegistrySecret.Password, suiteConfig.TanzuClusterEssentials.Filename),
		envfuncs.CreateNamespaces(suiteConfig.CreateNamespaces),
	)

	// finish
	testenv.Finish(
		envfuncs.DeleteNamespaces(suiteConfig.CreateNamespaces),
	)

	os.Exit(testenv.Run(m))
}
