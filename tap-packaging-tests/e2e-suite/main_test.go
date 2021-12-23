package suite

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gitlab.eng.vmware.com/tap/tap-packages/tap-packaging-tests/e2e-suite/envfuncs"
	e2e "gitlab.eng.vmware.com/tap/tap-packages/tap-packaging-tests/e2e-suite/pkg"
	tap "gitlab.eng.vmware.com/tap/tap-packages/tap-packaging-tests/pkg"
	"gopkg.in/yaml.v3"
	"sigs.k8s.io/e2e-framework/pkg/env"
)

var testenv env.Environment

type config struct {
	Namespaces        []string `yaml:"namespaces"`
	PackageRepository struct {
		Image     string `yaml:"image"`
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"package_repository"`
	TapFull struct {
		Name        string `yaml:"name"`
		Namespace   string `yaml:"namespace"`
		PackageName string `yaml:"package_name"`
		ValuesFile  string `yaml:"values_file"`
		Version     string `yaml:"version"`
		PollTimeout string `yaml:"poll_timeout"`
	} `yaml:"tap-full"`
	TapLight struct {
		Name        string `yaml:"name"`
		Namespace   string `yaml:"namespace"`
		PackageName string `yaml:"package_name"`
		ValuesFile  string `yaml:"values_file"`
		Version     string `yaml:"version"`
		PollTimeout string `yaml:"poll_timeout"`
	} `yaml:"tap-light"`
	Secret1 struct {
		Export    bool   `yaml:"export"`
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
		Registry  string `yaml:"registry"`
		Username  string `yaml:"username"`
		Password  string `yaml:"password"`
	} `yaml:"secret-1"`
	Secret2 struct {
		Export    bool   `yaml:"export"`
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
		Registry  string `yaml:"registry"`
		Username  string `yaml:"username"`
		Password  string `yaml:"password"`
	} `yaml:"secret-2"`
}

func setLogger() {
	os.MkdirAll(filepath.Join(e2e.GetFileDir(), "logs"), 0755)
	logFilePath := filepath.Join("logs", fmt.Sprintf("log_%s.log", time.Now().Format(time.RFC3339Nano)))
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	tap.CheckError(err)
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func TestMain(m *testing.M) {
	setLogger()
	tap.CheckPrerequisites() // TODO: create function in this package

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(fmt.Errorf("error while getting user home directory: %w", err))
	}
	testenv = env.NewWithKubeConfig(filepath.Join(home, ".kube", "config"))

	configBytes, err := os.ReadFile(filepath.Join(e2e.GetFileDir(), "suite-config.yaml"))
	if err != nil {
		log.Fatal(fmt.Errorf("error while reading config file: %w", err))
	}

	config := config{}
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		log.Fatal(fmt.Errorf("error while unmarshalling config file: %w", err))
	}

	testenv.Setup(
		envfuncs.CreateNamespaces(config.Namespaces),
		envfuncs.CreateSecret(config.Secret1.Name, config.Secret1.Registry, config.Secret1.Username, config.Secret1.Password, config.Secret1.Namespace, config.Secret1.Export),
		envfuncs.CreateSecret(config.Secret2.Name, config.Secret2.Registry, config.Secret2.Username, config.Secret2.Password, config.Secret2.Namespace, config.Secret2.Export),
		envfuncs.AddPackageRepository(config.PackageRepository.Name, config.PackageRepository.Image, config.PackageRepository.Namespace),
		envfuncs.CheckIfPackageRepositoryReconciled(config.PackageRepository.Name, config.PackageRepository.Namespace, 10),
		envfuncs.InstallPackage(config.TapFull.Name, config.TapFull.PackageName, config.TapFull.Version, config.TapFull.Namespace, filepath.Join(e2e.GetFileDir(), "values", config.TapFull.ValuesFile), config.TapFull.PollTimeout),
		envfuncs.CheckIfPackageInstalled(config.TapFull.Name, config.TapFull.Namespace, 10),
	)

	testenv.Finish(
		envfuncs.UninstallPackage(config.TapFull.Name, config.TapFull.Namespace),
		envfuncs.DeletePackageRepository(config.PackageRepository.Name, config.PackageRepository.Image, config.PackageRepository.Namespace),
		envfuncs.DeleteSecret(config.Secret2.Name, config.Secret2.Registry, config.Secret2.Username, config.Secret2.Password, config.Secret2.Namespace, config.Secret2.Export),
		envfuncs.DeleteSecret(config.Secret1.Name, config.Secret1.Registry, config.Secret1.Username, config.Secret1.Password, config.Secret1.Namespace, config.Secret1.Export),
		envfuncs.DeleteNamespaces(config.Namespaces),
	)

	os.Exit(testenv.Run(m))
}
