// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"encoding/json"
	"fmt"
	"log"
)

type packageInstalledOutput struct {
	Name           string `json:"name"`
	PackageName    string `json:"package-name"`
	PackageVersion string `json:"package-version"`
	Status         string `json:"status"`
}

type packageRepositoryOutput []struct {
	Name       string `json:"name"`
	Reason     string `json:"reason"`
	Repository string `json:"repository"`
	Status     string `json:"status"`
	Tag        string `json:"tag"`
	Version    string `json:"version"`
}

func InstallPackage(name string, packageName string, version string, namespace string, valuesFile string, pollTimeout string) (string, error) {
	cmd := fmt.Sprintf("tanzu package install %s -p %s -v %s -n %s", name, packageName, version, namespace)
	if valuesFile != "" {
		cmd += fmt.Sprintf(" -f %s", valuesFile)
	}
	if pollTimeout != "" {
		cmd += fmt.Sprintf(" --poll-timeout %s", pollTimeout)
	}
	log.Printf("command executed: %s", cmd)
	output, err := RunCommand(cmd)
	if err != nil {
		return "", err
	}

	return output, err
}

func UninstallPackage(name string, namespace string) (string, error) {
	cmd := fmt.Sprintf("tanzu package installed delete %s -n %s -y", name, namespace)
	log.Printf("command executed: %s", cmd)
	output, err := RunCommand(cmd)
	if err != nil {
		return "", err
	}

	return output, err
}

func GetPackageInstalledStatus(name string, namespace string) (string, error) {
	cmd := fmt.Sprintf("tanzu package installed get %s -n %s -o json", name, namespace)
	log.Printf("command executed: %s", cmd)
	output, err := RunCommand(cmd)
	if err != nil {
		return "", err
	}

	packageInstall := []packageInstalledOutput{}
	err = json.Unmarshal([]byte(output), &packageInstall)
	if err != nil {
		return "", err
	}

	if len(packageInstall) <= 0 {
		return "", fmt.Errorf("list empty for package installed status for package %s", name)
	}

	return packageInstall[0].Status, nil
}

func AddPackageRepository(name string, image string, namespace string) (string, error) {
	cmd := fmt.Sprintf("tanzu package repository add %s --url %s -n %s", name, image, namespace)
	log.Printf("command executed: %s", cmd)
	output, err := RunCommand(cmd)

	return output, err
}

func DeletePackageRepository(name string, namespace string) (string, error) {
	cmd := fmt.Sprintf("tanzu package repository delete %s -n %s -y", name, namespace)
	log.Printf("command executed: %s", cmd)
	output, err := RunCommand(cmd)

	return output, err
}

func GetPackageRepositoryStatus(name string, namespace string) (string, error) {
	cmd := fmt.Sprintf("tanzu package repository get %s -n %s -o json", name, namespace)
	log.Printf("command executed: %s", cmd)
	output, err := RunCommand(cmd)
	if err != nil {
		return "", err
	}

	packageRepository := packageRepositoryOutput{}
	err = json.Unmarshal([]byte(output), &packageRepository)
	if err != nil {
		return "", err
	}

	if len(packageRepository) <= 0 {
		return "", fmt.Errorf("list empty for package repository status for repository %s", name)
	}

	return packageRepository[0].Status, nil
}

func CreateSecret(name string, registry string, username string, password string, namespace string, export bool) (string, error) {
	cmd := fmt.Sprintf("tanzu secret registry add %s --server %s --username %s --password %s -n %s -y", name, registry, username, password, namespace)
	if export {
		cmd += " --export-to-all-namespaces"
	}
	// don't log command as it contains password
	output, err := RunCommand(cmd)

	return output, err
}

func DeleteSecret(name string, namespace string) (string, error) {
	cmd := fmt.Sprintf("tanzu secret registry delete %s -n %s -y", name, namespace)
	log.Printf("command executed: %s", cmd)
	output, err := RunCommand(cmd)

	return output, err
}
