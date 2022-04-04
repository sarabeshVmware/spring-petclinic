// Usage:
//   tanzu package installed [command]

// Available Commands:
//   create      Install a package
//   delete      Delete an installed package
//   get         Get details for an installed package
//   list        List installed packages
//   update      Update an installed package

package tanzu_libs

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

type ListInstalledPackagesOutput struct {
	NAME, PACKAGE_NAME, PACKAGE_VERSION, STATUS, NAMESPACE string
}

func ListInstalledPackages(namespace string) []ListInstalledPackagesOutput {
	installedPackages := []ListInstalledPackagesOutput{}
	cmd := "tanzu package installed list"
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += " -A"
	}
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return installedPackages
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")

	header_index := 0

	if strings.HasPrefix(temp[1], " ") {
		header_index = 1
	} else {
		header_index = 2
	}

	if len(temp) <= header_index+1 {
		log.Printf("Output : %s", temp[0])
		return installedPackages
	}

	indexSpans := linux_util.FieldIndices(temp[header_index])
	headers := linux_util.GetFields(temp[header_index], indexSpans)

	for index, ele := range headers {
		headers[index] = strings.ReplaceAll(ele, "-", "_")
	}

	for _, element := range temp[header_index+1:] {
		words := linux_util.GetFields(element, indexSpans)
		var installedPackage ListInstalledPackagesOutput
		for index, value := range words {
			reflect.ValueOf(&installedPackage).Elem().FieldByName(headers[index]).SetString(value)
			if namespace != "" {
				installedPackage.NAMESPACE = namespace
			}
		}
		installedPackages = append(installedPackages, installedPackage)
	}
	fmt.Printf("installedPackages: %+v\n", installedPackages)
	return installedPackages
}

type GetInstalledPackagesOutput []struct {
	Conditions         string `json:"conditions"`
	Name               string `json:"name"`
	PackageName        string `json:"package-name"`
	PackageVersion     string `json:"package-version"`
	Status             string `json:"status"`
	UsefulErrorMessage string `json:"useful-error-message"`
}

func GetInstalledPackages(name string, namespace string) GetInstalledPackagesOutput {
	var raw GetInstalledPackagesOutput
	cmd := fmt.Sprintf("tanzu package installed get %s", name)
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += " -A"
	}
	cmd += " -o json"

	output, err := linux_util.ExecuteCmd(cmd)

	if err != nil {
		return raw
	}

	if strings.HasPrefix(output, "[") {
		err = json.Unmarshal([]byte(output), &raw)
	} else {
		outputArray := strings.SplitN(output, "\n", 2)
		strippedOutput := outputArray[1]
		err = json.Unmarshal([]byte(strippedOutput), &raw)
	}

	if err != nil {
		panic(err)
	}
	return raw
}

func DeleteInstalledPackage(name string, namespace string) error {

	cmd := fmt.Sprintf("tanzu package installed delete %s --namespace %s --yes", name, namespace)
	res, err := linux_util.ExecuteCmd(cmd)
	if err != nil && !strings.Contains(res, "Uninstalled package") {
		log.Printf("Error while deleting the package %s. Error %v,  Output %s", name, err, res)
	}
	return err
}

func UpdateInstalledPackage(name string, packageName string, version string, namespace string, valuesFile string, pollTimeout string) error {

	cmd := fmt.Sprintf("tanzu package installed update %s --package-name %s --version %s --namespace %s", name, packageName, version, namespace)
	if valuesFile != "" {
		cmd += fmt.Sprintf(" --values-file %s", valuesFile)
	}
	if pollTimeout != "" {
		cmd += fmt.Sprintf(" --poll-timeout %s", pollTimeout)
	}
	res, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while updating package %s (%s) in namespace %s", name, packageName, namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", res)
	} else {
		log.Printf("package %s (%s) updated in namespace %s", name, packageName, namespace)
		log.Printf("output: %s", res)
	}
	return err
}
