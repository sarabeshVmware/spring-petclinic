package models

import (
	"log"
	"os"
	"path/filepath"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	"gopkg.in/yaml.v2"
)

type OuterloopConfig struct {
	CatalogInfoYaml string `yaml:"catalog_info_yaml"`
	Mysql           struct {
		YamlFile string `yaml:"yaml_file"`
	} `yaml:"mysql"`
	Namespace string `yaml:"namespace"`
	Project   struct {
		Host                string `yaml:"host"`
		WebpageRelativePath string `yaml:"webpage_relative_path"`
		File                string `yaml:"file"`
		Name                string `yaml:"name"`
		DestName            string `yaml:"dest_name"`
		RepoTemplate        string `yaml:"repo_template"`
		DestRepoTemplate    string `yaml:"dest_repo_template"`
		NewString           string `yaml:"new_string"`
		OriginalString      string `yaml:"original_string"`
		CommitMessage       string `yaml:"commit_message"`
		Repository          string `yaml:"repository"`
		Username            string `yaml:"username"`
		Email               string `yaml:"email"`
		AccessToken         string `yaml:"access_token"`
	} `yaml:"project"`
	ScanPolicy struct {
		YamlFile string `yaml:"yaml_file"`
	} `yaml:"scan_policy"`
	SpringPetclinicPipeline struct {
		Name     string `yaml:"name"`
		YamlFile string `yaml:"yaml_file"`
	} `yaml:"spring_petclinic_pipeline"`
	TestTargetRepo string `yaml:"test_target_repo"`
	Workload       struct {
		Name                 string `yaml:"name"`
		YamlFile             string `yaml:"yaml_file"`
		TestYamlFile         string `yaml:"test_yaml_file"`
		BuildNameSuffix      string `yaml:"build_name_suffix"`
		PipelineName         string `yaml:"pipeline_name"`
		TaskRunInfix         string `yaml:"taskrun_name_infix"`
		TaskRunTestSuffix    string `yaml:"taskrun_test_suffix"`
		ServiceBindingSuffix string `yaml:"service_binding_suffix"`
		GitopsYamlFile       string `yaml:"gitops_yaml_file"`
		GitSSHSecretYamlFile string `yaml:"gitssh_secret_yaml_file"`
	} `yaml:"workload"`
	BuildPacks struct {
		ScanPolicy       string `yaml:"scan_policy"`
		PipelineYamlFile string `yaml:"pipeline_yaml_file"`
		Workloads        []struct {
			Name                string `yaml:"name"`
			GitRepository       string `yaml:"git_repository"`
			GitBranch           string `yaml:"git_branch"`
			WebpageRelativePath string `yaml:"webpage_relative_path"`
			ContainsConventions bool   `yanl:"contains_conventions"`
		} `yaml:"workloads"`
	} `yaml:"buildpacks"`
	Domain        string `yaml:"domain"`
	MetadataStore struct {
		Domain    string `yaml:"domain"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata_store"`
}

var outerloopResourcesDir = filepath.Join(utils.GetFileDir(), "../../resources/outerloop")

func GetOuterloopConfig() (OuterloopConfig, error) {
	log.Printf("getting outerloop config")

	outerloopConfig := OuterloopConfig{}
	file := filepath.Join(outerloopResourcesDir, "outerloop-config.yaml")

	// read file
	outerloopConfigBytes, err := os.ReadFile(file)
	if err != nil {
		log.Printf("error while reading outerloop config file %s", file)
		log.Printf("error: %s", err)
		return outerloopConfig, err
	} else {
		log.Printf("read outerloop config file %s", file)
	}

	// unmarshall
	err = yaml.Unmarshal(outerloopConfigBytes, &outerloopConfig)
	if err != nil {
		log.Printf("error while unmarshalling outerloop config file %s", file)
		log.Printf("error: %s", err)
		return outerloopConfig, err
	} else {
		log.Printf("unmarshalled file %s", file)
	}

	// update outerloop config for full file paths
	outerloopConfig.Mysql.YamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.Mysql.YamlFile)
	outerloopConfig.ScanPolicy.YamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.ScanPolicy.YamlFile)
	outerloopConfig.SpringPetclinicPipeline.YamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.SpringPetclinicPipeline.YamlFile)
	outerloopConfig.BuildPacks.PipelineYamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.BuildPacks.PipelineYamlFile)
	outerloopConfig.Workload.YamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.Workload.YamlFile)
	outerloopConfig.Workload.TestYamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.Workload.TestYamlFile)
	outerloopConfig.Workload.GitopsYamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.Workload.GitopsYamlFile)
	outerloopConfig.Workload.GitSSHSecretYamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.Workload.GitSSHSecretYamlFile)
	outerloopConfig.BuildPacks.ScanPolicy = filepath.Join(outerloopResourcesDir, outerloopConfig.BuildPacks.ScanPolicy)
	return outerloopConfig, nil
}
