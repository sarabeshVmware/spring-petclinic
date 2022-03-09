package tanzu_libs

// Available Commands:
//   get         Get details for an available package or the openAPI schema of a package with a specific version
//   list        List available packages

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

type ListAvailablePackagesOutput struct {
	NAME, DISPLAY_NAME, SHORT_DESCRIPTION, LATEST_VERSION, NAMESPACE string
}

func ListAvailablePackages(namespace string) []ListAvailablePackagesOutput {
	installedPackages := []ListAvailablePackagesOutput{}
	cmd := "tanzu package available list"
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
		var installedPackage ListAvailablePackagesOutput
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
