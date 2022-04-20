package models

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	"gopkg.in/yaml.v2"
)

type SuiteConfig struct {
	CreateNamespaces []string `yaml:"create_namespaces"`
	Innerloop        struct {
		Workload struct {
			Name                string `yaml:"name"`
			Namespace           string `yaml:"namespace"`
			URL                 string `yaml:"url"`
			Gitrepository       string `yaml:"gitrepository"`
			YamlFile            string `yaml:"yaml_file"`
			PodintentName       string `yaml:"podintent_name"`
			ApplicationFilePath string `yaml:"application_file_path"`
			NewString           string `yaml:"new_string"`
			OriginalString      string `yaml:"original_string"`
			BuildNameSuffix     string `yaml:"build_name_suffix"`
			ImageDeliverySuffix string `yaml:"image_delivery_suffix"`
		} `yaml:"workload"`
	} `yaml:"innerloop"`
	PackageRepository struct {
		Image     string `yaml:"image"`
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
		Registry  string `yaml:"registry"`
	} `yaml:"package_repository"`
	NonTanzuRepository []struct {
		Repository   string `yaml:"repository"`
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
		PasswordType string `yaml:"passwordType"`
		Server       string `yaml:"server"`
	} `yaml:"non-tanzu-repository"`
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
		TanzunetHost            string `yaml:"tanzunet_host"`
		TanzunetApiToken        string `yaml:"tanzunet_api_token"`
		ProductFileId           int    `yaml:"product_file_id"`
		ReleaseVersion          string `yaml:"release_version"`
		ProductSlug             string `yaml:"product_slug"`
		DownloadBundle          string `yaml:"download_bundle"`
		InstallBundle           string `yaml:"install_bundle"`
		InstallRegistryHostname string `yaml:"install_registry_hostname"`
		InstallRegistryUsername string `yaml:"install_registry_username"`
		InstallRegistryPassword string `yaml:"install_registry_password"`
	} `yaml:"tanzu-cluster-essentials"`
	GitCredentials struct {
		Username string `yaml:"username"`
		Email    string `yaml:"email"`
	} `yaml:"git-credentials"`
	UpgradeVersions struct {
		Image             string `yaml:"image"`
		TapVersion        string `yaml:"tap-version"`
		UpgradeImage      string `yaml:"upgrade-image"`
		UpgradeTapVersion string `yaml:"upgrade-tap-version"`
	} `yaml:"upgrade-versions"`
}

var suiteResourcesDir = filepath.Join(utils.GetFileDir(), "../../resources/suite")
var suiteDir = filepath.Join(utils.GetFileDir(), "../..")

func GetSuiteConfig() SuiteConfig {
	var suiteConfig = SuiteConfig{}
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
	suiteConfig.Innerloop.Workload.YamlFile = filepath.Join(suiteDir, suiteConfig.Innerloop.Workload.YamlFile)

	return suiteConfig
}
