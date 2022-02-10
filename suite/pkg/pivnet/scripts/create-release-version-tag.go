package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	versionPrefix := os.Args[1]
	var config ConfigData
	newVersion := ""
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
	version := pivnet_helpers.GetLatestRelease(config.ProductSlug, versionPrefix)
	if version == "" {
		log.Println("No release with give version prefix found. Creating the very first build version string.")
		newVersion = fmt.Sprintf("%s.1", versionPrefix)
	} else {
		lstInd := strings.LastIndex(version, ".")
		verNo, _ := strconv.Atoi(version[lstInd+1:])
		newVerNo := strconv.Itoa(verNo + 1)
		newVersion = version[:lstInd] + "." + newVerNo
	}
	log.Printf("New version:\n%s", newVersion)
}
