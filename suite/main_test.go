package suite

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/envfuncs"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	"gopkg.in/yaml.v3"
	"sigs.k8s.io/e2e-framework/pkg/env"
)

var testenv env.Environment

var suiteConfig = struct {
	CreateNamespaces []string `yaml:"create_namespaces"`
	Innerloop        struct {
		Workload struct {
			Name      string `yaml:"name"`
			Namespace string `yaml:"namespace"`
			Name string `yaml:"name"`
			URL string `yaml:"url"`
		} `yaml:"workload"`
	} `yaml:"innerloop"`
	PackageRepository struct {
		Image     string `yaml:"image"`
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"package_repository"`
	TapRegistrySecret struct {
		Export    bool   `yaml:"export"`
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
		Password  string `yaml:"password"`
		Registry  string `yaml:"registry"`
		Username  string `yaml:"username"`
	} `yaml:"tap_registry_secret"`
	RegistryCredentialsSecret struct {
		Export    bool   `yaml:"export"`
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
		Password  string `yaml:"password"`
		Registry  string `yaml:"registry"`
		Username  string `yaml:"username"`
	} `yaml:"registry_credentials_secret"`
	Tap struct {
		Name             string `yaml:"name"`
		Namespace        string `yaml:"namespace"`
		PackageName      string `yaml:"package_name"`
		PollTimeout      string `yaml:"poll_timeout"`
		ValuesSchemaFile string `yaml:"values_schema_file"`
		Version          string `yaml:"version"`
	} `yaml:"tap"`
	TanzuClusterEssentials struct {
		Bundle    string `yaml:"bundle"`
		Registry  string `yaml:"registry"`
		FileName  string `yaml:"filename"`
	} `yaml:"tanzu-cluster-essentials"`
}{}

var tapValuesSchema = struct {
	Accelerator struct {
		Server struct {
			ServiceType string `yaml:"service_type"`
		} `yaml:"server"`
	} `yaml:"accelerator"`
	Buildservice struct {
		KpDefaultRepository         string `yaml:"kp_default_repository"`
		KpDefaultRepositoryPassword string `yaml:"kp_default_repository_password"`
		KpDefaultRepositoryUsername string `yaml:"kp_default_repository_username"`
		TanzunetPassword            string `yaml:"tanzunet_password"`
		TanzunetUsername            string `yaml:"tanzunet_username"`
	} `yaml:"buildservice"`
	CeipPolicyDisclosed bool `yaml:"ceip_policy_disclosed"`
	Contour             struct {
		Envoy struct {
			Service struct {
				Type string `yaml:"type"`
			} `yaml:"service"`
		} `yaml:"envoy"`
	} `yaml:"contour"`
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
		ServiceType string `yaml:"service_type"`
	} `yaml:"tap_gui"`
}{}

var resourcesDir = filepath.Join(utils.GetFileDir(), "resources")
var suiteResourcesDir = filepath.Join(resourcesDir, "suite")

func TestMain(m *testing.M) {
	// set logger
	logFile, err := utils.SetLogger(filepath.Join(utils.GetFileDir(), "logs"))
	if err != nil {
		log.Fatal(fmt.Errorf("error while setting log file %s: %w", logFile, err))
	}

	// get kubeconfig
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(fmt.Errorf("error while getting user home directory: %w", err))
	}
	testenv = env.NewWithKubeConfig(filepath.Join(home, ".kube", "config"))

	// read suite config
	suiteConfigBytes, err := os.ReadFile(filepath.Join(suiteResourcesDir, "suite-config.yaml"))
	if err != nil {
		log.Fatal(fmt.Errorf("error while reading suite config file: %w", err))
	}
	err = yaml.Unmarshal(suiteConfigBytes, &suiteConfig)
	if err != nil {
		log.Fatal(fmt.Errorf("error while unmarshalling suite config file: %w", err))
	}

	// update suite config for full path for values schema
	suiteConfig.Tap.ValuesSchemaFile = filepath.Join(suiteResourcesDir, suiteConfig.Tap.ValuesSchemaFile)

	// read tap values schema
	tapValuesSchemaBytes, err := os.ReadFile(suiteConfig.Tap.ValuesSchemaFile)
	if err != nil {
		log.Fatal(fmt.Errorf("error while reading tap values schema file: %w", err))
	}
	err = yaml.Unmarshal(tapValuesSchemaBytes, &tapValuesSchema)
	if err != nil {
		log.Fatal(fmt.Errorf("error while unmarshalling tap values schema file: %w", err))
	}

	// setup
	developerNamespaceFile := filepath.Join(suiteResourcesDir, "developer-namespace.yaml")
	testenv.Setup(
		envfuncs.InstallClusterEssentials(suiteConfig.TanzuClusterEssentials.Bundle, 
			suiteConfig.TanzuClusterEssentials.Registry ,
			suiteConfig.TapRegistrySecret.Username, 
			suiteConfig.TapRegistrySecret.Password, 
			suiteConfig.TanzuClusterEssentials.FileName),
		// envfuncs.CheckAndDeploy("kapp-controller", []string{"https://github.com/vmware-tanzu/carvel-kapp-controller/releases/latest/download/release.yml"}, "default"),           // temporary, to be replaced by cluster essentials script
		// envfuncs.CheckAndDeploy("secretgen-controller", []string{"https://github.com/vmware-tanzu/carvel-secretgen-controller/releases/download/v0.5.0/release.yml"}, "default"), // temporary, to be replaced by cluster essentials script
		envfuncs.CreateNamespaces(suiteConfig.CreateNamespaces),
		envfuncs.CreateSecret(suiteConfig.TapRegistrySecret.Name, suiteConfig.TapRegistrySecret.Registry, suiteConfig.TapRegistrySecret.Username, suiteConfig.TapRegistrySecret.Password, suiteConfig.TapRegistrySecret.Namespace, suiteConfig.TapRegistrySecret.Export),
		envfuncs.CreateSecret(suiteConfig.RegistryCredentialsSecret.Name, suiteConfig.RegistryCredentialsSecret.Registry, suiteConfig.RegistryCredentialsSecret.Username, suiteConfig.RegistryCredentialsSecret.Password, suiteConfig.RegistryCredentialsSecret.Namespace, suiteConfig.RegistryCredentialsSecret.Export),
		envfuncs.AddPackageRepository(suiteConfig.PackageRepository.Name, suiteConfig.PackageRepository.Image, suiteConfig.PackageRepository.Namespace),
		envfuncs.CheckIfPackageRepositoryReconciled(suiteConfig.PackageRepository.Name, suiteConfig.PackageRepository.Namespace, 10),
		envfuncs.InstallPackage(suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, suiteConfig.Tap.Namespace, suiteConfig.Tap.ValuesSchemaFile, suiteConfig.Tap.PollTimeout),
		envfuncs.CheckIfPackageInstalled(suiteConfig.Tap.Name, suiteConfig.Tap.Namespace, 10),
		envfuncs.SetupDeveloperNamespace(developerNamespaceFile, suiteConfig.CreateNamespaces[0]),
		envfuncs.SetupDeveloperNamespace(developerNamespaceFile, suiteConfig.CreateNamespaces[1]),
	)

	// finish
	testenv.Finish(
		envfuncs.DeleteDeveloperNamespace(developerNamespaceFile, suiteConfig.CreateNamespaces[1]),
		envfuncs.DeleteDeveloperNamespace(developerNamespaceFile, suiteConfig.CreateNamespaces[0]),
		envfuncs.UninstallPackage(suiteConfig.Tap.Name, suiteConfig.Tap.Namespace),
		envfuncs.DeletePackageRepository(suiteConfig.PackageRepository.Name, suiteConfig.PackageRepository.Namespace),
		envfuncs.DeleteSecret(suiteConfig.RegistryCredentialsSecret.Name, suiteConfig.RegistryCredentialsSecret.Namespace),
		envfuncs.DeleteSecret(suiteConfig.TapRegistrySecret.Name, suiteConfig.TapRegistrySecret.Namespace),
		envfuncs.DeleteNamespaces(suiteConfig.CreateNamespaces),
	)

	os.Exit(testenv.Run(m))
}
