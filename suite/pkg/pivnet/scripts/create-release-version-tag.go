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
	fmt.Println("Logging into tanzunet")
	pivnet_libs.Login(config.Host, config.APIToken)
	version := pivnet_helpers.GetLatestRelease(config.ProductSlug, versionPrefix)
	if version == "" {
		log.Println("No release with give version prefix found. Creating the very first build version string.")
		newVersion = fmt.Sprintf("%s-build.1", versionPrefix)
	} else {
		lstInd := strings.LastIndex(version, ".")
		verNo, _ := strconv.Atoi(version[lstInd+1:])
		newVerNo := strconv.Itoa(verNo + 1)
		newVersion = version[:lstInd] + newVerNo
	}
	log.Printf("New version: %s", newVersion)
}
