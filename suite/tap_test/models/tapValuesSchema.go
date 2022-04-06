package models

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type TapValuesSchema struct {
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
}

func GetTapValuesSchema() (TapValuesSchema, error) {
	log.Printf("getting tap values schema")

	suiteConfig := GetSuiteConfig()
	tapValuesSchema := TapValuesSchema{}
	file := suiteConfig.Tap.ValuesSchemaFile

	// read file
	tapValuesSchemaBytes, err := os.ReadFile(file)
	if err != nil {
		log.Printf("error while reading tap values schema file %s", file)
		log.Printf("error: %s", err)
		return tapValuesSchema, err
	} else {
		log.Printf("read tap values schema file %s", file)
	}

	// unmarshal
	err = yaml.Unmarshal(tapValuesSchemaBytes, &tapValuesSchema)
	if err != nil {
		log.Printf("error while unmarshalling tap values schema file %s", file)
		log.Printf("error: %s", err)
		return tapValuesSchema, err
	} else {
		log.Printf("unmarshalled file %s", file)
	}

	return tapValuesSchema, nil
}
