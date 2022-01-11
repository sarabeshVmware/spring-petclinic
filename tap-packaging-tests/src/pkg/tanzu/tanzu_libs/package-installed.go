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
	linux_util "pkg/utils/linux_util"
	"reflect"
	"strings"
)

type ListInstalledPackagesOutput struct {
	NAME, PACKAGE_NAME, PACKAGE_VERSION, STATUS string
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
	if len(temp) <= 2 {
		log.Printf("Output : %s", temp[0])
		return installedPackages
	}

	indexSpans := linux_util.FieldIndices(temp[1])
	headers := linux_util.GetFields(temp[1], indexSpans)
	for index, ele := range headers {
		headers[index] = strings.ReplaceAll(ele, "-", "_")
	}

	for _, element := range temp[2:] {
		words := linux_util.GetFields(element, indexSpans)
		var installedPackage ListInstalledPackagesOutput
		for index, value := range words {
			reflect.ValueOf(&installedPackage).Elem().FieldByName(headers[index]).SetString(value)
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

func GetInstalledPackages(name string, namespace string) *GetInstalledPackagesOutput {

	cmd := fmt.Sprintf("tanzu package installed get %s", name)
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += " -A"
	}
	cmd += " -o json"

	res1, err1 := linux_util.ExecuteCmd(cmd)
	if err1 != nil {
		log.Println("something bad happened")
	}
	in := []byte(res1)
	var raw *GetInstalledPackagesOutput
	if err := json.Unmarshal(in, &raw); err != nil {
		panic(err)
	}
	return raw
}
