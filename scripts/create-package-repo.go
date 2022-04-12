package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Package struct {
	Name     string   `yaml:"name"`
	Versions []string `yaml:"versions"`
}

type Repository struct {
	Packages []Package `yaml:"packages"`
}

type BundleLock struct {
	APIVersion string    `json:"apiVersion"`
	Kind       string    `json:"kind"`
	Bundle     BundleRef `json:"bundle"` // This generated yaml, but due to lib we need to use `json`
}

type BundleRef struct {
	Image string `json:"image,omitempty"` // This generated yaml, but due to lib we need to use `json`
	Tag   string `json:"tag,omitempty"`   // This generated yaml, but due to lib we need to use `json`
}

func main() {
	var OciRegistry = "dev.registry.tanzu.vmware.com/tanzu-application-platform"
	var PackagesDirectoryPath = filepath.Join("./", "packages")
	var RepoDirectoryPath = filepath.Join("./", "repos")
	var GeneratedRepoDirectoryPath = filepath.Join(RepoDirectoryPath, "generated")
	var repository Repository

	channel := os.Args[1]
	tag := os.Args[2]
	channelToPush := "tap-packages"
	if len(os.Args) > 3 {
		channelToPush = os.Args[3]
	}

	if len(os.Args) > 4 {
		OciRegistry = os.Args[4]
	}

	channelDir := filepath.Join(GeneratedRepoDirectoryPath, channel)
	imgpkgDir := filepath.Join(channelDir, ".imgpkg")
	packagesDir := filepath.Join(channelDir, "packages")

	// Remove any existing generated files
	os.RemoveAll(fmt.Sprint(channelDir))
	err := os.MkdirAll(imgpkgDir, 0755)
	check(err)

	err = os.MkdirAll(packagesDir, 0755)
	check(err)

	targetChannelFilename := filepath.Join(RepoDirectoryPath, channel+".yaml")
	source, err := ioutil.ReadFile(targetChannelFilename)
	check(err)

	err = yaml.Unmarshal(source, &repository)
	check(err)

	var outputPackageYaml = filepath.Join(packagesDir, "packages.yaml")
	outputFile, err := os.Create(outputPackageYaml)
	check(err)

	defer func() {
		if err := outputFile.Close(); err != nil {
			panic(err)
		}
	}()

	exclude_filepath := []string{
		"packages/cert-manager/metadata.yaml",
		"packages/contour/metadata.yaml"}

	for _, p := range repository.Packages {
		metadataFilepath := filepath.Join(PackagesDirectoryPath, p.Name, "metadata.yaml")
		if !excluded_filepath(metadataFilepath, exclude_filepath) {
			copyYaml(metadataFilepath, outputFile)
		}

		for _, version := range p.Versions {
			packageFilepath := filepath.Join(PackagesDirectoryPath, p.Name, fmt.Sprintf("%s.yaml", version))
			copyYaml(packageFilepath, outputFile)
		}
	}

	imagesLockFile := filepath.Join(imgpkgDir, "images.yml")
	execCommand("kbld", []string{"--file", packagesDir, "--imgpkg-lock-output", imagesLockFile})

	bundleLockFilename := "output.yaml"
	//registryPathAndTag := OciRegistry + "/" + channel + ":latest"
	registryPathAndTag := OciRegistry + "/" + channelToPush + ":" + tag
	execCommand("imgpkg", []string{"push", "--tty", "--bundle", registryPathAndTag, "--file", channelDir, "--lock-output", bundleLockFilename})

	bundleLockYamlFile, err := ioutil.ReadFile(bundleLockFilename)
	check(err)

	var bundleLock BundleLock
	err = yaml.Unmarshal(bundleLockYamlFile, &bundleLock)
	check(err)

	fmt.Println("Package Repository pushed to", bundleLock.Bundle.Image)
	os.RemoveAll(bundleLockFilename)

	// fmt.Println("Updating packages version in packages.yaml from ", targetChannelFilename)
	// UpdatePackagesFile(targetChannelFilename)
}

func execCommand(command string, commandArgs []string) {
	fmt.Println("Executing command : ", command, commandArgs)
	output, err := exec.Command(command, commandArgs...).CombinedOutput()
	if err != nil {
		log.Fatal(string(output))
	}
}

func copyYaml(packageFilepath string, outputFile *os.File) {
	source, err := ioutil.ReadFile(packageFilepath)
	check(err)

	var slice = source[0:3]
	if !strings.HasPrefix(string(slice), "---") {
		if _, err := outputFile.WriteString("---\n"); err != nil {
			panic(err)
		}
	}

	_, err = outputFile.Write(source)
	check(err)

	slice = source[len(source)-1:]
	if string(slice) != "\n" {
		if _, err := outputFile.WriteString("\n"); err != nil {
			panic(err)
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func excluded_filepath(filepath string, exclude_filepath []string) bool {
	for _, path := range exclude_filepath {
		if path == filepath {
			return true
		}
	}
	return false
}

func UpdatePackagesFile(betaFilePath string) {
	packagesFilePath := "tap-packaging-tests/packages.yaml"
	fmt.Println("Updating ", packagesFilePath)
	sourceFile, err := os.ReadFile(betaFilePath)
	check(err)
	destFile, err := os.ReadFile(packagesFilePath)
	check(err)

	type BetaPackages struct {
		Packages []struct {
			Name     string   `yaml:"name"`
			Versions []string `yaml:"versions"`
		} `yaml:"packages"`
	}
	type TestPkgs struct {
		Description         string   `yaml:"description"`
		Name                string   `yaml:"name"`
		Namespace           string   `yaml:"namespace"`
		Package             string   `yaml:"package"`
		ValuesFile          string   `yaml:"values_file,omitempty"`
		Version             string   `yaml:"version"`
		PackageDependencies []string `yaml:"package_dependencies,omitempty"`
	}
	beta4pkgs := BetaPackages{}
	err = yaml.Unmarshal(sourceFile, &beta4pkgs)
	check(err)

	testpkg := []TestPkgs{}
	err = yaml.Unmarshal(destFile, &testpkg)

	check(err)
	for i, val := range testpkg {
		for _, betaval := range beta4pkgs.Packages {
			if (val.Name == betaval.Name) || ((val.Name == "tap-full" || val.Name == "tap-dev-light") && betaval.Name == "tap") {
				lastIndex := len(betaval.Versions) - 1
				if val.Version == betaval.Versions[lastIndex] {
					fmt.Printf("Version already updated for package %s \n", val.Name)
				} else {
					fmt.Printf("Changing version for package %s from %s to %s \n", val.Name, val.Version, betaval.Versions[lastIndex])

					testpkg[i].Version = betaval.Versions[lastIndex]
				}
			}
		}
	}
	changedDestFile, err := yaml.Marshal(&testpkg)
	check(err)

	// write to file
	err = os.WriteFile(packagesFilePath, changedDestFile, 0644)
	check(err)
	fmt.Println(packagesFilePath, "Updated successfully.")
}
