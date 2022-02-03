package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"

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
	if !pivnet_libs.Login(config.Host, config.APIToken) {
		log.Fatalln("Unable to login to tanzunet to fetch current release version")
	}

	createRelease := true
	artifacts := pivnet_libs.ListArtifactReferences(config.ProductSlug, "", config.Digest)
	if len(artifacts) != 0 {
		fmt.Println("Artifacts with given sha already exists.")
		for _, artf := range artifacts {
			if len(artf.ReleaseVersions) != 0 {
				fmt.Println("Artifacts with the given sha is already added to atleast 1 release. Skipping release creation.")
				createRelease = false
				break
			}
		}
	}
	if createRelease {
		fmt.Println("Create release in tanzunet")
		rel_det := pivnet_libs.CreateRelease(config.ProductSlug, config.ReleaseVersion, config.ReleaseType, config.EulaSlug)
		fmt.Printf("Release created, id: %d, version: %s\n", rel_det.ID, rel_det.Version)

		fmt.Println("Add artifact reference in tanzunet")
		artifact_det := pivnet_libs.CreateArtifactReference("tap-package-repo-bundle", config.ProductSlug, config.ArtifactPath, config.Digest)
		fmt.Printf("Artifact created, id: %d, name: %s\n", artifact_det.ID, artifact_det.Name)

		fmt.Println("Waiting till artifacts addition is complete")
		pivnet_helpers.WaitTillArtifactReferenceIsReady(config.ProductSlug, artifact_det.ID)
		pivnet_libs.ListArtifactReferences(config.ProductSlug, "", config.Digest)

		fmt.Println("Add artifact reference to release")
		res := pivnet_libs.AddArtifactReference(config.ProductSlug, config.ReleaseVersion, artifact_det.ID)
		if !res {
			log.Fatal("Exiting as artifact not added to release")
		}

		// add files and filegroups from previous release
		lstInd := strings.LastIndex(config.ReleaseVersion, ".")
		verNo, _ := strconv.Atoi(config.ReleaseVersion[lstInd+1:])
		if verNo == 1 {
			fmt.Println("Creating the first build in the release, no prior release found to copy product files and file groups from.")
		} else {
			preVerNo := strconv.Itoa(verNo - 1)
			previousVersion := config.ReleaseVersion[:lstInd] + "." + preVerNo

			fmt.Println("Fetching 'File Groups' from release: %s", previousVersion)
			filegroups := pivnet_libs.ListFileGroups(config.ProductSlug, previousVersion)
			for _, fg := range filegroups {
				fmt.Println("Adding File Group '%s' to the release '%s'", fg.Name, config.ReleaseVersion)
				fgadded := pivnet_libs.AddFileGroup(fg.ID, config.ProductSlug, config.ReleaseVersion)
				if !fgadded {
					log.Println("Unable to add file group '%s' to release", fg.Name)
				}
			}

			fmt.Println("Fetching 'Product Files' from release: %s", previousVersion)
			prodfiles := pivnet_libs.ListProductFiles(config.ProductSlug, previousVersion)
			for _, pf := range prodfiles {
				fmt.Println("Adding Product File '%s' to the release '%s'", pf.Name, config.ReleaseVersion)
				fgadded := pivnet_libs.AddFileGroup(pf.ID, config.ProductSlug, config.ReleaseVersion)
				if !fgadded {
					log.Println("Unable to add product file '%s' to release", pf.Name)
				}
			}

		}

		fmt.Println("Updating release availability")
		updatedRel := pivnet_libs.UpdateRelease(config.ProductSlug, config.ReleaseVersion, "selected-user-groups")
		fmt.Printf("Release availability updated to '%s'", updatedRel.Availability)

		fmt.Println("Get all user groups")
		userGroupList := pivnet_libs.ListUserGroups(config.ProductSlug)
		for _, user := range userGroupList {
			for _, u := range config.UserGroups {
				if user.Name == u {
					fmt.Println("Add user group to release")
					userAdded := pivnet_libs.AddUserGroup(config.ProductSlug, config.ReleaseVersion, user.ID)
					if !userAdded {
						log.Println("Unable to assign user group '%s' to release", u)
					}
					break
				}
			}
		}
	}
}
