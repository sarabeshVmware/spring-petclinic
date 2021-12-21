// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func GetInstalledPackages(namespace string) ([]PackageInstalledOutput, error) {
	cmd := fmt.Sprintf("tanzu package installed list -n %s -o json", namespace)
	log.Printf("Getting installed packages in namespace %s: %s", namespace, cmd)
	output, err := RunCommand(cmd)
	if err != nil {
		return []PackageInstalledOutput{}, err
	}

	packages := []PackageInstalledOutput{}
	err = json.Unmarshal([]byte(output), &packages)
	if err != nil {
		return packages, err
	}

	return packages, nil
}

func GetInstalledPackageStatus(installedPackageName string, namespace string) (string, error) {
	cmd := fmt.Sprintf("tanzu package installed get %s -n %s -o json", installedPackageName, namespace)
	log.Printf("Checking package installed status for package %s: %s", installedPackageName, cmd)
	output, err := RunCommand(cmd)
	if err != nil {
		return "", err
	}

	packageInstall := []PackageInstalledOutput{}
	err = json.Unmarshal([]byte(output), &packageInstall)
	if err != nil {
		return "", err
	}

	if len(packageInstall) <= 0 {
		return packageInstall[0].Status, fmt.Errorf("list empty for package installed status for package: %s", installedPackageName)
	}

	return packageInstall[0].Status, nil
}

func GetPackagesList() ([]Package, error) {
	packagesList := []Package{}

	packagesFileBytes, err := os.ReadFile(GetPackagesYamlFilepath())
	if err != nil {
		return packagesList, err
	}

	err = yaml.Unmarshal(packagesFileBytes, &packagesList)
	if err != nil {
		return packagesList, err
	}

	return packagesList, nil
}

func GetDependentPackagesInfo(packageInfo Package) ([]Package, error) {
	dependentPackagesInfo := []Package{}
	for _, packageDependency := range packageInfo.PackageDependencies {
		packagesList, err := GetPackagesList()
		if err != nil {
			return []Package{}, err
		}
		for _, packageInfo := range packagesList {
			if packageInfo.Package == packageDependency {
				dependentPackagesInfo = append(dependentPackagesInfo, packageInfo)
			}
		}
	}

	return dependentPackagesInfo, nil
}

func GetPackageInfoFromName(packageName string) (Package, error) {
	packagesList, err := GetPackagesList()
	if err != nil {
		return Package{}, err
	}

	for _, packageInfo := range packagesList {
		if packageInfo.Name == packageName {
			return packageInfo, nil
		}
	}

	return Package{}, fmt.Errorf("package %s not found in the packages list", packageName)
}
