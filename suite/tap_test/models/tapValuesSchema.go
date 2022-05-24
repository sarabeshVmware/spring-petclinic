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
		ServiceType    string `yaml:"service_type"`
		IngressEnabled bool   `yaml:"ingressEnabled"`
		AppConfig      struct {
			App struct {
				BaseURL string `yaml:"baseUrl"`
			} `yaml:"app"`
			Catalog struct {
				Locations []struct {
					Type   string `yaml:"type"`
					Target string `yaml:"target"`
				} `yaml:"locations"`
			} `yaml:"catalog"`
			Backend struct {
				BaseURL string `yaml:"baseUrl"`
				Cors    struct {
					Origin string `yaml:"origin"`
				} `yaml:"cors"`
			} `yaml:"backend"`
			Kubernetes struct {
				ServiceLocatorMethod struct {
					Type string `yaml:"type"`
				} `yaml:"serviceLocatorMethod"`
				ClusterLocatorMethods []struct {
					Type     string `yaml:"type"`
					Clusters []struct {
						URL                 string `yaml:"url"`
						Name                string `yaml:"name"`
						AuthProvider        string `yaml:"authProvider"`
						ServiceAccountToken string `yaml:"serviceAccountToken"`
						SkipTLSVerify       bool   `yaml:"skipTLSVerify"`
					} `yaml:"clusters"`
				} `yaml:"clusterLocatorMethods"`
			} `yaml:"kubernetes"`
		} `yaml:"app_config"`
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

func GetProfileTapValuesSchema(profile string) (TapValuesSchema, error) {
	log.Printf("getting tap values schema")

	suiteConfig := GetSuiteConfig()
	tapValuesSchema := TapValuesSchema{}
	var file string
	if profile == "build" {
		file = suiteConfig.Multicluster.BuildTapValuesFile
	} else if profile == "run" {
		file = suiteConfig.Multicluster.RunTapValuesFile
	} else if profile == "view" {
		file = suiteConfig.Multicluster.ViewTapValuesFile
	} else {
		file = suiteConfig.Tap.ValuesSchemaFile
	}

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
