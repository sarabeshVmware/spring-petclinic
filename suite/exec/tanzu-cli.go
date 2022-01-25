// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"encoding/json"
	"fmt"
)

func TanzuInstallPackage(name string, packageName string, version string, namespace string, valuesFile string, pollTimeout string) (string, string, error) {
	cmd := fmt.Sprintf("tanzu package install %s -p %s -v %s -n %s", name, packageName, version, namespace)
	if valuesFile != "" {
		cmd += fmt.Sprintf(" -f %s", valuesFile)
	}
	if pollTimeout != "" {
		cmd += fmt.Sprintf(" --poll-timeout %s", pollTimeout)
	}
	output, err := RunCommand(cmd)
	return cmd, output, err
}

func TanzuUpdatePackage(name string, packageName string, version string, namespace string, valuesFile string) (string, string, error) {
	cmd := fmt.Sprintf("tanzu package installed update %s -p %s -v %s -n %s -f %s", name, packageName, version, namespace, valuesFile)
	output, err := RunCommand(cmd)
	return cmd, output, err
}

func TanzuUninstallPackage(name string, namespace string) (string, string, error) {
	cmd := fmt.Sprintf("tanzu package installed delete %s -n %s -y", name, namespace)
	output, err := RunCommand(cmd)
	return cmd, output, err
}

func TanzuGetPackageInstalledStatus(name string, namespace string) (string, string, error) {
	cmd := fmt.Sprintf("tanzu package installed get %s -n %s -o json", name, namespace)
	output, err := RunCommand(cmd)
	if err != nil {
		return cmd, "", err
	}
	packageInstalled := []struct {
		Status string `json:"status"`
	}{}
	err = json.Unmarshal([]byte(output), &packageInstalled)
	if err != nil {
		return cmd, "", err
	}
	if len(packageInstalled) <= 0 {
		return cmd, "", fmt.Errorf("list empty for package installed status for package %s", name)
	}
	return cmd, packageInstalled[0].Status, nil
}

func TanzuAddPackageRepository(name string, image string, namespace string) (string, string, error) {
	cmd := fmt.Sprintf("tanzu package repository add %s --url %s -n %s", name, image, namespace)
	output, err := RunCommand(cmd)
	return cmd, output, err
}

func TanzuDeletePackageRepository(name string, namespace string) (string, string, error) {
	cmd := fmt.Sprintf("tanzu package repository delete %s -n %s -y", name, namespace)
	output, err := RunCommand(cmd)
	return cmd, output, err
}

func TanzuGetPackageRepositoryStatus(name string, namespace string) (string, string, error) {
	cmd := fmt.Sprintf("tanzu package repository get %s -n %s -o json", name, namespace)
	output, err := RunCommand(cmd)
	if err != nil {
		return cmd, "", err
	}
	packageRepository := []struct {
		Status string `json:"status"`
	}{}
	err = json.Unmarshal([]byte(output), &packageRepository)
	if err != nil {
		return cmd, "", err
	}
	if len(packageRepository) <= 0 {
		return cmd, "", fmt.Errorf("list empty for package repository status for repository %s", name)
	}
	return cmd, packageRepository[0].Status, nil
}

func TanzuCreateSecret(name string, registry string, username string, password string, namespace string, export bool) (string, string, error) {
	cmd := fmt.Sprintf("tanzu secret registry add %s --server %s --username %s --password %s -n %s -y", name, registry, username, password, namespace)
	if export {
		cmd += " --export-to-all-namespaces"
	}
	output, err := RunCommand(cmd)
	return cmd, output, err
}

func TanzuDeleteSecret(name string, namespace string) (string, string, error) {
	cmd := fmt.Sprintf("tanzu secret registry delete %s -n %s -y", name, namespace)
	output, err := RunCommand(cmd)
	return cmd, output, err
}

func TanzuDeployWorkload(workloadFile string, namespace string) (string, string, error) {
	cmd := fmt.Sprintf("tanzu apps workload apply -f %s -n %s -y", workloadFile, namespace)
	output, err := RunCommand(cmd)
	return cmd, output, err
}

func TanzuGenerateAccelerator(acceleratorName string, projectName string, repositoryPrefix string, serverIP string, namespace string) (string, string, error){
	cmd:= fmt.Sprintf(`tanzu accelerator generate %s --options '{"projectName":"%s", "repositoryPrefix":"%s", "includeKubernetes": true}' --server-url http://%s`, acceleratorName, projectName, repositoryPrefix, serverIP)
	output, err := RunCommandInBashMode(cmd)
	return cmd, output, err
}