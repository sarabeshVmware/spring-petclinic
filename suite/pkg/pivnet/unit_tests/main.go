package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"

	pivnet_helpers "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/pivnet/pivnet_helpers"
	pivnet_libs "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/pivnet/pivnet_libs"
)

type ConfigData struct {
	Host           string   `yaml:"host"`
	APIToken       string   `yaml:"api-token"`
	ProductSlug    string   `yaml:"product-slug"`
	ReleaseVersion string   `yaml:"release-version"`
	EulaSlug       string   `yaml:"eula-slug"`
	ReleaseType    string   `yaml:"release-type"`
	ArtifactPath   string   `yaml:"artifact-path"`
	Digest         string   `yaml:"digest"`
	UserGroups     []string `yaml:"user-groups"`
}

func create_release() {

	var config ConfigData

	filename, _ := filepath.Abs("./pkg/pivnet/unit_tests/config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	pivnet_libs.Login(config.Host, config.APIToken)
	pivnet_libs.CreateRelease(config.ProductSlug, config.ReleaseVersion, config.ReleaseType, config.EulaSlug)
	artifact_det := pivnet_libs.CreateArtifactReference(config.ReleaseVersion, config.ProductSlug, config.ArtifactPath, config.Digest)
	pivnet_helpers.WaitTillArtifactReferenceIsReady(config.ProductSlug, artifact_det.ID)
	pivnet_libs.AddArtifactReference(config.ProductSlug, config.ReleaseVersion, artifact_det.ID)
	pivnet_libs.UpdateRelease(config.ProductSlug, config.ReleaseVersion, "selected-user-groups")
	userGroupList := pivnet_libs.ListUserGroups()
	fmt.Println(userGroupList)
	for _, value := range config.UserGroups {
		fmt.Println(value)
		pivnet_libs.AddUserGroup(config.ProductSlug, config.ReleaseVersion, "437")
	}

}

func main() {
	create_release()
}
