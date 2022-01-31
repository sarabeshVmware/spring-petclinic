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

func main() {
	var config ConfigData

	filename, _ := filepath.Abs("./pkg/pivnet/scripts/config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}
	fmt.Println("Logging into tanzunet")
	pivnet_libs.Login(config.Host, config.APIToken)
	fmt.Println("Create release in tanzunet")
	pivnet_libs.CreateRelease(config.ProductSlug, config.ReleaseVersion, config.ReleaseType, config.EulaSlug)
	fmt.Println("Add artifact reference in tanzunet")
	artifact_det := pivnet_libs.CreateArtifactReference(config.ReleaseVersion, config.ProductSlug, config.ArtifactPath, config.Digest)
	fmt.Println("Waiting till artifacts addition is complete")
	pivnet_helpers.WaitTillArtifactReferenceIsReady(config.ProductSlug, artifact_det.ID)
	fmt.Println("Add artifact reference to release")
	pivnet_libs.AddArtifactReference(config.ProductSlug, config.ReleaseVersion, artifact_det.ID)
	fmt.Println("Updating release availability")
	pivnet_libs.UpdateRelease(config.ProductSlug, config.ReleaseVersion, "selected-user-groups")
	fmt.Println("Get all user groups")
	userGroupList := pivnet_libs.ListUserGroups(config.ProductSlug)
	for _, user := range userGroupList {
		for _, u := range config.UserGroups {
			if user.Name == u {
				fmt.Println("Add user group to release")
				pivnet_libs.AddUserGroup(config.ProductSlug, config.ReleaseVersion, user.ID)
				break
			}
		}
	}
}
