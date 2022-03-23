package main

import (
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
	log.Println("Logging into tanzunet")
	if !pivnet_libs.Login(config.Host, config.APIToken) {
		log.Fatalln("Unable to login to tanzunet to fetch current release version")
	}

	createRelease := true
	artifacts := pivnet_libs.ListArtifactReferences(config.ProductSlug, "", config.Digest)
	if len(artifacts) != 0 {
		log.Println("Artifacts with given sha already exists.")
		for _, artf := range artifacts {
			if len(artf.ReleaseVersions) != 0 {
				log.Printf("Skipping release creation\n Artifacts with the given sha is already added to the following release(s) %v", artf.ReleaseVersions)
				createRelease = false
				break
			}
		}
	}
	if createRelease {
		log.Println("Create release in tanzunet")
		rel_det := pivnet_libs.CreateRelease(config.ProductSlug, config.ReleaseVersion, config.ReleaseType, config.EulaSlug)
		log.Printf("Release created, id: %d, version: %s\n", rel_det.ID, rel_det.Version)

		log.Println("Add artifact reference in tanzunet")
		artifact_det := pivnet_libs.CreateArtifactReference("tap-package-repo-bundle", config.ProductSlug, config.ArtifactPath, config.Digest)
		log.Printf("Artifact created, id: %d, name: %s\n", artifact_det.ID, artifact_det.Name)

		log.Println("Waiting till artifacts addition is complete")
		artifactsAdded := pivnet_helpers.WaitTillArtifactReferenceIsReady(config.ProductSlug, artifact_det.ID, 40, 60)
		if !artifactsAdded {
			log.Fatal("Exiting as artifact did not get created in tanzunet")
		}
		pivnet_libs.ListArtifactReferences(config.ProductSlug, "", config.Digest)

		log.Println("Add artifact reference to release")
		res := pivnet_libs.AddArtifactReference(config.ProductSlug, config.ReleaseVersion, artifact_det.ID)
		if !res {
			log.Fatal("Exiting as artifact not added to release")
		}

		// add files and filegroups from previous release
		lstInd := strings.LastIndex(config.ReleaseVersion, ".")
		verNo, _ := strconv.Atoi(config.ReleaseVersion[lstInd+1:])
		if verNo == 1 {
			log.Println("Creating the first build in the release, no prior release found to copy product files and file groups from.")
		} else {
			preVerNo := strconv.Itoa(verNo - 1)
			previousVersion := config.ReleaseVersion[:lstInd] + "." + preVerNo

			log.Printf("Fetching 'File Groups' from release: %s", previousVersion)
			filegroups := pivnet_libs.ListFileGroups(config.ProductSlug, previousVersion)
			for _, fg := range filegroups {
				log.Printf("Adding File Group '%s' to the release '%s'", fg.Name, config.ReleaseVersion)
				fgadded := pivnet_libs.AddFileGroup(fg.ID, config.ProductSlug, config.ReleaseVersion)
				if !fgadded {
					log.Printf("Unable to add file group '%s' to release", fg.Name)
				}
			}

			log.Printf("Fetching 'Product Files' from release: %s", previousVersion)
			prodfiles := pivnet_libs.ListProductFiles(config.ProductSlug, previousVersion)
			fileGroupsRel := pivnet_libs.ListFileGroups(config.ProductSlug, config.ReleaseVersion)
			for _, pf := range prodfiles {
				// skip file if it already got added as part of file group
				added := false
				for _, fg := range fileGroupsRel {
					for _, file := range fg.ProductFiles {
						if file.ID == pf.ID {
							added = true
							log.Printf("Skipping Product File '%s' as it is already added as part of file group.", pf.Name)
						}
					}
				}
				if !added {
					log.Printf("Adding Product File '%s' to the release '%s'", pf.Name, config.ReleaseVersion)
					pfadded := pivnet_libs.AddProductFile(pf.ID, config.ProductSlug, config.ReleaseVersion)
					if !pfadded {
						log.Printf("Unable to add product file '%s' to release", pf.Name)
					}
				}
			}

		}

		log.Println("Updating release availability")
		updatedRel := pivnet_libs.UpdateRelease(config.ProductSlug, config.ReleaseVersion, "selected-user-groups")
		log.Printf("Release availability updated to '%s'", updatedRel.Availability)

		log.Println("Get all user groups")
		userGroupList := pivnet_libs.ListUserGroups(config.ProductSlug)
		for _, user := range userGroupList {
			for _, u := range config.UserGroups {
				if user.Name == u {
					log.Println("Add user group to release")
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
