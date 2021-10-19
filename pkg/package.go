// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/buger/jsonparser"
)

type Package struct {
	Name          string `yaml:"name"`
	InstalledName string `yaml:"installed_name"`
	Version       string `yaml:"version"`
	UseValuesFile string `yaml:"use_values_file"`
}

type PackageInstalledOutput struct {
	Name           string `json:"name"`
	PackageName    string `json:"package-name"`
	PackageVersion string `json:"package-version"`
	Status         string `json:"status"`
}

func ListPackages(namespace string) {
	log.Printf("Available packages in namespace: %s", namespace)
	Run(fmt.Sprintf("tanzu package available list -n %s", namespace))
}

func ListInstalledPackages(namespace string) []PackageInstalledOutput {
	var packages []PackageInstalledOutput
	log.Printf("Installed packages in namespace: %s", namespace)
	packagesList, _ := Run(fmt.Sprintf("tanzu package installed list -n %s -o json", namespace))
	err := json.Unmarshal(packagesList, &packages)
	CheckError(err)
	return packages
}

func ListValuesSchema(packages []Package, namespace string) {
	for _, packageInfo := range packages {
		log.Printf("Values schemas for package: %s", packageInfo.Name)
		Run(fmt.Sprintf("tanzu package available get %s/%s --values-schema -n %s", packageInfo.Name, packageInfo.Version, namespace))
	}
}

func InstallPackages(packages []Package, namespace string, ValuesDirectory string) {
	for _, packageInfo := range packages {
		log.Printf("Installing package: %s", packageInfo.Name)
		if packageInfo.UseValuesFile != "" {
			valuesSchemaFile := filepath.Join(ValuesDirectory, packageInfo.UseValuesFile)
			Run(fmt.Sprintf("tanzu package install %s -p %s -v %s -n %s -f %s", packageInfo.InstalledName, packageInfo.Name, packageInfo.Version, namespace, valuesSchemaFile))
		} else {
			Run(fmt.Sprintf("tanzu package install %s -p %s -v %s -n %s", packageInfo.InstalledName, packageInfo.Name, packageInfo.Version, namespace))
		}
		ValidatePackage(packageInfo, namespace)
	}
}

func ValidatePackage(packageInfo Package, namespace string) {
	log.Printf("Validating package: %s", packageInfo.Name)
	packageInstalled, _ := Run(fmt.Sprintf("tanzu package installed get %s -n %s -o json", packageInfo.InstalledName, namespace))
	status, err := jsonparser.GetString(packageInstalled, "[0]", "status")
	CheckError(err)
	if status == "Reconciling" {
		time.Sleep(5 * time.Second)
		ValidatePackage(packageInfo, namespace)
	} else if status == "Reconcile succeeded" {
		log.Printf("Reconcile succeeded for package install: %s", packageInfo.Name)
	} else {
		log.Fatalf("Reconcile not succeeded for package install: %s", packageInfo.Name)
	}
}

func UninstallPackages(namespace string) {
	installedpackages := ListInstalledPackages(namespace)
	for _, each := range installedpackages {
		log.Printf("Uninstalling package: %s", each.Name)
		Run(fmt.Sprintf("tanzu package installed delete %s -n %s -y", each.Name, namespace))
	}
}
