// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tanzuCmds

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

func TanzuInstallPackage(name string, packageName string, version string, namespace string, valuesFile string, pollTimeout string) error {
	log.Printf("installing package %s (%s) in namespace %s", name, packageName, namespace)

	// execute cmd
	cmd := fmt.Sprintf("tanzu package install %s -p %s -v %s -n %s", name, packageName, version, namespace)
	if valuesFile != "" {
		cmd += fmt.Sprintf(" -f %s", valuesFile)
	}
	if pollTimeout != "" {
		cmd += fmt.Sprintf(" --poll-timeout %s", pollTimeout)
	}
	output, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while installing package %s (%s) in namespace %s", name, packageName, namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("package %s (%s) installed in namespace %s", name, packageName, namespace)
		log.Printf("output: %s", output)
	}

	return err
}

func TanzuUninstallPackage(name string, namespace string) error {
	log.Printf("uninstalling package %s from namespace %s", name, namespace)

	// execute cmd
	cmd := fmt.Sprintf("tanzu package installed delete %s -n %s -y", name, namespace)
	output, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while uninstalling package %s from namespace %s", name, namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("package %s uninstalled from namespace %s", name, namespace)
		log.Printf("output: %s", output)
	}

	return err
}

func TanzuUpdatePackage(name string, packageName string, version string, namespace string, valuesFile string) error {
	log.Printf("updating package %s in namespace %s", name, namespace)

	// execute cmd
	cmd := fmt.Sprintf("tanzu package installed update %s -p %s -v %s -n %s -f %s", name, packageName, version, namespace, valuesFile)
	output, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while updating package %s in namespace %s", name, namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("package %s updated in namespace %s", name, namespace)
		log.Printf("output: %s", output)
	}

	return err
}

func TanzuListInstalledPackages(namespace string) error {
	log.Printf("listing packages in namespace %s", namespace)

	// execute cmd
	cmd := fmt.Sprintf("tanzu package installed list -n %s", namespace)
	output, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while listing packages in namespace %s", namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("packages updated in namespace %s", namespace)
		log.Printf("output: %s", output)
	}

	return err
}

func TanzuGetPackageInstalledStatus(name string, namespace string) (string, error) {
	log.Printf("getting package %s installation status in namespace %s", name, namespace)

	// get installation status
	cmd := fmt.Sprintf("tanzu package installed get %s -n %s -o json", name, namespace)
	output, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while getting package %s installation status in namespace %s", name, namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
		return "", err
	} else {
		log.Printf("package %s in namespace %s installation status retrieved", name, namespace)
		log.Printf("output: %s", output)
	}

	packageInstalled := []struct {
		Status string `json:"status"`
	}{}

	// unmarshall
	if strings.HasPrefix(output, "[") {
		err = json.Unmarshal([]byte(output), &packageInstalled)
	} else {

		outputArray := strings.SplitN(output, "\n", 2)
		strippedOutput := outputArray[1]
		err = json.Unmarshal([]byte(strippedOutput), &packageInstalled)
	}
	if err != nil {
		log.Printf("error while unmarshalling output %s", output)
		log.Printf("error: %s", err)
		return "", err
	} else {
		log.Printf("unmarshalled output %s", output)
	}

	// check len
	if len(packageInstalled) <= 0 {
		err = fmt.Errorf("list empty for package installed status for package %s", name)
		log.Printf("error while checking length of packages installed")
		log.Printf("error: %s", err)
		return "", err
	}

	return packageInstalled[0].Status, nil
}

func TanzuAddPackageRepository(name string, image string, namespace string) error {
	log.Printf("adding package repository %s (%s) in namespace %s", name, image, namespace)

	// execute cmd
	cmd := fmt.Sprintf("tanzu package repository add %s --url %s -n %s", name, image, namespace)
	output, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while adding package repository %s (%s) in namespace %s", name, image, namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("package repository %s (%s) added in namespace %s", name, image, namespace)
		log.Printf("output: %s", output)
	}

	return err
}

func TanzuDeletePackageRepository(name string, namespace string) error {
	log.Printf("deleting package repository %s from namespace %s", name, namespace)

	// execute cmd
	cmd := fmt.Sprintf("tanzu package repository delete %s -n %s -y", name, namespace)
	output, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while deleting package repository %s from namespace %s", name, namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("package repository %s deleted from namespace %s", name, namespace)
		log.Printf("output: %s", output)
	}

	return err
}

func TanzuGetPackageRepositoryStatus(name string, namespace string) (string, error) {
	log.Printf("getting package repository %s status in namespace %s", name, namespace)

	// get repo status
	cmd := fmt.Sprintf("tanzu package repository get %s -n %s -o json", name, namespace)
	output, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while getting package repository %s status in namespace %s", name, namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
		return "", err
	} else {
		log.Printf("package repository %s in namespace %s status retrieved", name, namespace)
		log.Printf("output: %s", output)
	}

	packageRepository := []struct {
		Status string `json:"status"`
	}{}

	// unmarshall
	if strings.HasPrefix(output, "[") {
		err = json.Unmarshal([]byte(output), &packageRepository)
	} else {
		outputArray := strings.SplitN(output, "\n", 2)
		strippedOutput := outputArray[1]
		err = json.Unmarshal([]byte(strippedOutput), &packageRepository)
	}
	if err != nil {
		log.Printf("error while unmarshalling output %s", output)
		log.Printf("error: %s", err)
		return "", err
	} else {
		log.Printf("unmarshalled output %s", output)
	}

	// check len
	if len(packageRepository) <= 0 {
		err = fmt.Errorf("list empty for package repository status for package %s", name)
		log.Printf("error while checking length of package repository status")
		log.Printf("error: %s", err)
		return "", err
	}

	return packageRepository[0].Status, nil
}

func TanzuCreateSecret(name string, registry string, username string, password string, passwordType string, namespace string, export bool) error {
	log.Printf("creating secret %s (registry %s, username %s, export %t) in namespace %s", name, registry, username, export, namespace)

	// execute cmd
	cmd := fmt.Sprintf("tanzu secret registry add %s --server %s --username %s -n %s -y", name, registry, username, namespace)
	if passwordType == "key" {
		// create temporary file for cert
		log.Printf("creating tempfile for key")
		tempFile, err := ioutil.TempFile("", "password*.json")
		if err != nil {
			log.Print("error while creating tempfile")
		} else {
			log.Print("created tempfile")
		}
		defer os.Remove(tempFile.Name())
		err = os.WriteFile(tempFile.Name(), []byte(password), 0677)
		if err != nil {
			log.Printf("error while writing to file %s", tempFile.Name())
			log.Printf("error: %s", err)
		} else {
			log.Printf("file %s written", tempFile.Name())
		}

		cmd += fmt.Sprintf(" --password-file %s", tempFile.Name())
	} else {
		cmd += fmt.Sprintf(" --password %s", password)
	}
	if export {
		cmd += " --export-to-all-namespaces"
	}
	output, err := linux_util.ExecuteCmdNoLog(cmd)
	if err != nil {
		log.Printf("error while creating secret %s (registry %s, username %s) in namespace %s", name, registry, username, namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("secret %s (registry %s, username %s) created in namespace %s", name, registry, username, namespace)
		log.Printf("output: %s", output)
	}

	return err
}

func TanzuCreateSecretWithJson(name string, registry string, username string, password string, namespace string, export bool) error {
	log.Printf("creating secret %s (registry %s, username %s, export %t) in namespace %s", name, registry, username, export, namespace)

	// create temporary file for cert
	log.Printf("creating tempfile for key")
	tempFile, err := ioutil.TempFile("", "password*.json")
	if err != nil {
		log.Print("error while creating tempfile")
	} else {
		log.Print("created tempfile")
	}
	defer os.Remove(tempFile.Name())
	err = os.WriteFile(tempFile.Name(), []byte(password), 0677)
	if err != nil {
		log.Printf("error while writing to file %s", tempFile.Name())
		log.Printf("error: %s", err)
	} else {
		log.Printf("file %s written", tempFile.Name())
	}

	// execute cmd
	cmd := fmt.Sprintf("tanzu secret registry add %s --server %s --username %s --password-file %s -n %s -y", name, registry, username, tempFile.Name(), namespace)
	if export {
		cmd += " --export-to-all-namespaces"
	}
	output, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while creating secret %s (registry %s, username %s) in namespace %s", name, registry, username, namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("secret %s (registry %s, username %s) created in namespace %s", name, registry, username, namespace)
		log.Printf("output: %s", output)
	}

	return err
}

func TanzuDeleteSecret(name string, namespace string) error {
	log.Printf("deleting secret %s from namespace %s", name, namespace)

	// execute cmd
	cmd := fmt.Sprintf("tanzu secret registry delete %s -n %s -y", name, namespace)
	output, err := linux_util.ExecuteCmdNoLog(cmd)
	if err != nil {
		log.Printf("error while deleting secret %s from namespace %s", name, namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("secret %s deleted from namespace %s", name, namespace)
		log.Printf("output: %s", output)
	}

	return err
}

func TanzuDeployWorkload(workloadFile string, namespace string) error {
	log.Printf("deploying workload file %s in namespace %s", workloadFile, namespace)

	// execute cmd
	cmd := fmt.Sprintf("tanzu apps workload apply -f %s -n %s -y", workloadFile, namespace)
	output, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while deploying workload file %s in namespace %s", workloadFile, namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("workload file %s deployed in namespace %s", workloadFile, namespace)
		log.Printf("output: %s", output)
	}

	return err
}

func TanzuDeleteWorkload(workloadFile string, namespace string) error {
	log.Printf("deleting workload %s in namespace %s", workloadFile, namespace)

	// execute cmd
	cmd := fmt.Sprintf("tanzu apps workload delete -f %s -n %s -y", workloadFile, namespace)
	output, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while deleting workload %s in namespace %s", workloadFile, namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("workload %s deleted in namespace %s", workloadFile, namespace)
		log.Printf("output: %s", output)
	}

	return err
}

func TanzuGenerateAccelerator(acceleratorName string, projectName string, repositoryPrefix string, serverIP string, namespace string, retries int, intervalInSeconds int) error {
	var output string
	var err error

	for iter := 0; iter < retries; iter++ {
		log.Printf(`[iteration %d] generating accelerator project %s (projectName "%s", repositoryPrefix "%s", serverIP "%s") in namespace %s`, iter+1, acceleratorName, projectName, repositoryPrefix, serverIP, namespace)

		// execute cmd
		cmd := fmt.Sprintf(`tanzu accelerator generate %s --options '{"projectName":"%s", "repositoryPrefix":"%s", "includeKubernetes": true}' --server-url http://%s`, acceleratorName, projectName, repositoryPrefix, serverIP)
		output, err = linux_util.ExecuteCmdInBashMode(cmd)
		if err != nil {
			log.Printf(`error while generating accelerator project %s (projectName "%s", repositoryPrefix "%s", serverIP "%s") in namespace %s`, acceleratorName, projectName, repositoryPrefix, serverIP, namespace)
			log.Printf("error: %s", err)
			log.Printf("output: %s", output)
			log.Printf("waiting for %d seconds before retry", intervalInSeconds)
			time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		} else {
			log.Printf(`accelerator project %s (projectName "%s", repositoryPrefix "%s", serverIP "%s") generated in namespace %s`, acceleratorName, projectName, repositoryPrefix, serverIP, namespace)
			log.Printf("output: %s", output)
			return nil
		}
	}

	log.Printf("accelerator project not generated after %d retries", retries)
	return err
}
