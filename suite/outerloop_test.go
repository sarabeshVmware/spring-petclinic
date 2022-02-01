//go:build outerloop

package suite

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type outerloopConfig struct {
	CatalogInfoYaml    string `yaml:"catalog_info_yaml"`
	Clusterrolebinding struct {
		Name           string `yaml:"name"`
		Clusterrole    string `yaml:"clusterrole"`
		ServiceAccount string `yaml:"serviceAccount"`
	} `yaml:"clusterrolebinding"`
	Mysql struct {
		Name     string `yaml:"name"`
		YamlFile string `yaml:"yaml_file"`
	} `yaml:"mysql"`
	Namespace string `yaml:"namespace"`
	Project   struct {
		Application         string `yaml:"application"`
		File                string `yaml:"file"`
		Name                string `yaml:"name"`
		NewString           string `yaml:"new_string"`
		WebpageRelativePath string `yaml:"webpage_relative_path"`
		OriginalString      string `yaml:"original_string"`
		Repository          string `yaml:"repository"`
	} `yaml:"project"`
	ScanPolicy struct {
		YamlFile string `yaml:"yaml_file"`
	} `yaml:"scan_policy"`
	SpringPetclinic struct {
		BuildNamePrefix     string `yaml:"build_name_prefix"`
		GitrepositoryName   string `yaml:"gitrepository_name"`
		ImagerepositoryName string `yaml:"imagerepository_name"`
		KsvcName            string `yaml:"ksvc_name"`
		Name                string `yaml:"name"`
		PodintentName       string `yaml:"podintent_name"`
		TaskrunNamePrefix   string `yaml:"taskrun_name_prefix"`
		YamlFile            string `yaml:"yaml_file"`
	} `yaml:"spring_petclinic"`
	Workload struct {
		Name     string `yaml:"name"`
		YamlFile string `yaml:"yaml_file"`
	} `yaml:"workload"`
}

var outerloopResourcesDir = filepath.Join(resourcesDir, "outerloop")

func getOuterloopConfig() (outerloopConfig, error) {
	outerloopConfig := outerloopConfig{}

	// read file
	outerloopConfigBytes, err := os.ReadFile(filepath.Join(outerloopResourcesDir, "outerloop-config.yaml"))
	if err != nil {
		return outerloopConfig, fmt.Errorf("error while reading outerloop config file: %w", err)
	}
	err = yaml.Unmarshal(outerloopConfigBytes, &outerloopConfig)
	if err != nil {
		return outerloopConfig, fmt.Errorf("error while unmarshalling outerloop config file: %w", err)
	}

	// update outerloop config for full file paths
	outerloopConfig.Mysql.YamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.Mysql.YamlFile)
	outerloopConfig.ScanPolicy.YamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.ScanPolicy.YamlFile)
	outerloopConfig.SpringPetclinic.YamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.SpringPetclinic.YamlFile)
	outerloopConfig.Workload.YamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.Workload.YamlFile)

	return outerloopConfig, nil
}
