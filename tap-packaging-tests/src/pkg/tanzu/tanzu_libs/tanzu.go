package tanzu_libs

import (
	"fmt"
	"log"
	linux_util "pkg/utils/linux_util"
	"reflect"
	"strings"
)

type GetAllInstalledPackagesOutput struct {
	NAME, PACKAGE_NAME, PACKAGE_VERSION, STATUS string
}

func ListInstalledPackages(namespace string) []GetAllInstalledPackagesOutput {
	installedPackages := []GetAllInstalledPackagesOutput{}
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
		var installedPackage GetAllInstalledPackagesOutput
		for index, value := range words {
			reflect.ValueOf(&installedPackage).Elem().FieldByName(headers[index]).SetString(value)
		}
		installedPackages = append(installedPackages, installedPackage)
	}
	fmt.Printf("installedPackages: %+v\n", installedPackages)
	return installedPackages
}
