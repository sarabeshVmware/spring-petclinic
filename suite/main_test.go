package suite

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/tap-packaging-tests/suite/envfuncs"
	"gopkg.in/yaml.v3"
	"sigs.k8s.io/e2e-framework/pkg/env"
)

var testenv env.Environment

var config = struct {
	Namespaces []string `yaml:"namespaces"`
	Outerloop  struct {
		CatalogInfoYaml string `yaml:"catalog_info_yaml"`
		Mysql           struct {
			Name      string `yaml:"name"`
			Namespace string `yaml:"namespace"`
			YamlFile  string `yaml:"yaml_file"`
		} `yaml:"mysql"`
		Namespace  string `yaml:"namespace"`
		ScanPolicy struct {
			Namespace string `yaml:"namespace"`
			YamlFile  string `yaml:"yaml_file"`
		} `yaml:"scan_policy"`
		SpringPetclinic struct {
			ImageRepositoryName string `yaml:"imagerepository_name"`
			Name                string `yaml:"name"`
			Namespace           string `yaml:"namespace"`
			YamlFile            string `yaml:"yaml_file"`
		} `yaml:"spring_petclinic"`
		Workload struct {
			Namespace string `yaml:"namespace"`
			YamlFile  string `yaml:"yaml_file"`
		} `yaml:"workload"`
	} `yaml:"outerloop"`
	PackageRepository struct {
		Image     string `yaml:"image"`
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"package_repository"`
	TanzunetCredsSecret struct {
		Export    bool   `yaml:"export"`
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
		Password  string `yaml:"password"`
		Registry  string `yaml:"registry"`
		Username  string `yaml:"username"`
	} `yaml:"tanzunet_creds_secret"`
	ImageSecret struct {
		Export    bool   `yaml:"export"`
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
		Password  string `yaml:"password"`
		Registry  string `yaml:"registry"`
		Username  string `yaml:"username"`
	} `yaml:"image_secret"`
	Tap struct {
		Name             string `yaml:"name"`
		Namespace        string `yaml:"namespace"`
		PackageName      string `yaml:"package_name"`
		PollTimeout      string `yaml:"poll_timeout"`
		ValuesSchemaFile string `yaml:"values_schema_file"`
		Version          string `yaml:"version"`
	} `yaml:"tap"`
}{}

var tapValuesSchema = struct {
	Buildservice struct {
		KpDefaultRepository         string `yaml:"kp_default_repository"`
		KpDefaultRepositoryPassword string `yaml:"kp_default_repository_password"`
		KpDefaultRepositoryUsername string `yaml:"kp_default_repository_username"`
		TanzunetPassword            string `yaml:"tanzunet_password"`
		TanzunetUsername            string `yaml:"tanzunet_username"`
	} `yaml:"buildservice"`
	CeipPolicyDisclosed bool `yaml:"ceip_policy_disclosed"`
	Cnrs                struct {
		DomainName interface{} `yaml:"domain_name,omitempty"`
	} `yaml:"cnrs,omitempty"`
	Contour struct {
		Envoy struct {
			Service struct {
				Type string `yaml:"type,omitempty"`
			} `yaml:"service,omitempty"`
		} `yaml:"envoy,omitempty"`
	} `yaml:"contour,omitempty"`
	Grype struct {
		Namespace             string `yaml:"namespace"`
		TargetImagePullSecret string `yaml:"targetImagePullSecret"`
	} `yaml:"grype"`
	Learningcenter struct {
		IngressDomain string `yaml:"ingressDomain"`
	} `yaml:"learningcenter"`
	OotbSupplyChainBasic struct {
		Gitops struct {
			SSHSecret string `yaml:"ssh_secret"`
		} `yaml:"gitops"`
		Registry struct {
			Repository string `yaml:"repository"`
			Server     string `yaml:"server"`
		} `yaml:"registry"`
	} `yaml:"ootb_supply_chain_basic"`
	OotbSupplyChainTesting struct {
		Gitops struct {
			SSHSecret string `yaml:"ssh_secret"`
		} `yaml:"gitops"`
		Registry struct {
			Repository string `yaml:"repository"`
			Server     string `yaml:"server"`
		} `yaml:"registry"`
	} `yaml:"ootb_supply_chain_testing"`
	OotbSupplyChainTestingScanning struct {
		Gitops struct {
			SSHSecret string `yaml:"ssh_secret"`
		} `yaml:"gitops"`
		Registry struct {
			Repository string `yaml:"repository"`
			Server     string `yaml:"server"`
		} `yaml:"registry"`
	} `yaml:"ootb_supply_chain_testing_scanning"`
	Profile     string `yaml:"profile"`
	SupplyChain string `yaml:"supply_chain"`
	TapGui      struct {
		AppConfig struct {
			App struct {
				BaseURL string `yaml:"baseUrl,omitempty"`
				Title   string `yaml:"title,omitempty"`
			} `yaml:"app,omitempty"`
			Backend struct {
				BaseURL string `yaml:"baseUrl,omitempty"`
				Cors    struct {
					Origin string `yaml:"origin,omitempty"`
				} `yaml:"cors,omitempty"`
			} `yaml:"backend,omitempty"`
			Catalog struct {
				Locations []struct {
					Target string `yaml:"target,omitempty"`
					Type   string `yaml:"type,omitempty"`
				} `yaml:"locations,omitempty"`
			} `yaml:"catalog,omitempty"`
		} `yaml:"app_config,omitempty"`
		ServiceType string `yaml:"service_type,omitempty"`
	} `yaml:"tap_gui,omitempty"`
}{}

func TestMain(m *testing.M) {
	logFile, err := SetLogger(filepath.Join(GetFileDir(), "logs"))
	if err != nil {
		log.Fatal(fmt.Errorf("error while setting log file %s: %w", logFile, err))
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(fmt.Errorf("error while getting user home directory: %w", err))
	}
	testenv = env.NewWithKubeConfig(filepath.Join(home, ".kube", "config"))

	configBytes, err := os.ReadFile(filepath.Join(GetFileDir(), "suite-config.yaml"))
	if err != nil {
		log.Fatal(fmt.Errorf("error while reading config file: %w", err))
	}
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		log.Fatal(fmt.Errorf("error while unmarshalling config file: %w", err))
	}

	config.Tap.ValuesSchemaFile = filepath.Join(GetFileDir(), config.Tap.ValuesSchemaFile)

	outerloopResourcesDir := "outerloop-resources"
	config.Outerloop.Mysql.YamlFile = filepath.Join(GetFileDir(), outerloopResourcesDir, config.Outerloop.Mysql.YamlFile)
	config.Outerloop.ScanPolicy.YamlFile = filepath.Join(GetFileDir(), outerloopResourcesDir, config.Outerloop.ScanPolicy.YamlFile)
	config.Outerloop.SpringPetclinic.YamlFile = filepath.Join(GetFileDir(), outerloopResourcesDir, config.Outerloop.SpringPetclinic.YamlFile)
	config.Outerloop.Workload.YamlFile = filepath.Join(GetFileDir(), outerloopResourcesDir, config.Outerloop.Workload.YamlFile)

	tapValuesSchemaBytes, err := os.ReadFile(config.Tap.ValuesSchemaFile)
	if err != nil {
		log.Fatal(fmt.Errorf("error while reading tap values schema file: %w", err))
	}
	err = yaml.Unmarshal(tapValuesSchemaBytes, &tapValuesSchema)
	if err != nil {
		log.Fatal(fmt.Errorf("error while unmarshalling tap values schema file: %w", err))
	}

	testenv.Setup(
		envfuncs.CheckAndDeploy("kapp-controller", []string{"https://github.com/vmware-tanzu/carvel-kapp-controller/releases/latest/download/release.yml"}, "default"),           // temporary, to be replaced by cluster essentials script
		envfuncs.CheckAndDeploy("secretgen-controller", []string{"https://github.com/vmware-tanzu/carvel-secretgen-controller/releases/download/v0.5.0/release.yml"}, "default"), // temporary, to be replaced by cluster essentials script
		envfuncs.CreateNamespaces(config.Namespaces),
		envfuncs.CreateSecret(config.TanzunetCredsSecret.Name, config.TanzunetCredsSecret.Registry, config.TanzunetCredsSecret.Username, config.TanzunetCredsSecret.Password, config.TanzunetCredsSecret.Namespace, config.TanzunetCredsSecret.Export),
		envfuncs.CreateSecret(config.ImageSecret.Name, config.ImageSecret.Registry, config.ImageSecret.Username, config.ImageSecret.Password, config.ImageSecret.Namespace, config.ImageSecret.Export),
		envfuncs.AddPackageRepository(config.PackageRepository.Name, config.PackageRepository.Image, config.PackageRepository.Namespace),
		envfuncs.CheckIfPackageRepositoryReconciled(config.PackageRepository.Name, config.PackageRepository.Namespace, 10),
		envfuncs.InstallPackage(config.Tap.Name, config.Tap.PackageName, config.Tap.Version, config.Tap.Namespace, config.Tap.ValuesSchemaFile, config.Tap.PollTimeout),
		envfuncs.CheckIfPackageInstalled(config.Tap.Name, config.Tap.Namespace, 10),
	)

	testenv.Finish(
		envfuncs.UninstallPackage(config.Tap.Name, config.Tap.Namespace),
		envfuncs.DeletePackageRepository(config.PackageRepository.Name, config.PackageRepository.Namespace),
		envfuncs.DeleteSecret(config.ImageSecret.Name, config.ImageSecret.Namespace),
		envfuncs.DeleteSecret(config.TanzunetCredsSecret.Name, config.TanzunetCredsSecret.Namespace),
		envfuncs.DeleteNamespaces(config.Namespaces),
	)

	os.Exit(testenv.Run(m))
}
