// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
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

func ListPackages(namespace string) {
	log.Printf("Available packages in namespace: %s", namespace)
	RunCommand(Command{CommandName: "tanzu", Arguments: []string{"package", "available", "list", "-n", namespace}})
}

func ListValuesSchema(packages []Package, namespace string) {
	for _, packageInfo := range packages {
		log.Printf("Values schemas for package: %s", packageInfo.Name)
		RunCommand(Command{CommandName: "tanzu", Arguments: []string{"package", "available", "get", packageInfo.Name + "/" + packageInfo.Version, "--values-schema", "-n", namespace}})
	}
}

func InstallPackages(packages []Package, namespace string) {
	for _, packageInfo := range packages {
		valuesSchemaFile := filepath.Join(GetValuesDirectory(), "values.yaml")
		if packageInfo.UseValuesFile != "" {
			valuesSchemaFile = filepath.Join(GetValuesDirectory(), packageInfo.UseValuesFile)
		}
		log.Printf("Installing package: %s", packageInfo.Name)
		RunCommand(Command{CommandName: "tanzu", Arguments: []string{"package", "install", packageInfo.InstalledName,
			"-p", packageInfo.Name, "-v", packageInfo.Version, "-n", namespace, "-f", valuesSchemaFile}})
		ValidatePackage(packageInfo, namespace)
	}
}

func ValidatePackage(packageInfo Package, namespace string) {
	log.Printf("Validating package: %s", packageInfo.Name)
	packageInstalled, _ := RunCommand(Command{CommandName: "tanzu", Arguments: []string{"package", "installed", "get", packageInfo.InstalledName, "-n", namespace, "-o", "json"}})
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

func UninstallPackages(packages []Package, namespace string) {
	for _, packageInfo := range packages {
		log.Printf("Uninstalling package: %s", packageInfo.Name)
		RunCommand(Command{CommandName: "tanzu", Arguments: []string{"package", "installed", "delete", packageInfo.InstalledName, "-n", namespace, "-y"}})
	}
}
