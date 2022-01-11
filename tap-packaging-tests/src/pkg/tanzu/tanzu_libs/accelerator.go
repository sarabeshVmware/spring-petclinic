package tanzu_libs

// Usage:
//   tanzu accelerator [command]

// Aliases:
//   accelerator, acc

// Available Commands:
//   create        Create a new accelerator
//   delete        Delete an accelerator
//   generate      Generate project from accelerator
//   get           Get accelerator info
//   list          List accelerators
//   push          Push local path to source image
//   update        Update an accelerator

import (
	"fmt"
	"log"
	linux_util "pkg/utils/linux_util"
	"reflect"
	"strings"
)

type TestOutput struct {
	NAME, PACKAGE_NAME, PACKAGE_VERSION, STATUS string
}

func Test(namespace string) []TestOutput {
	installedPackages := []TestOutput{}
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
		var installedPackage TestOutput
		for index, value := range words {
			reflect.ValueOf(&installedPackage).Elem().FieldByName(headers[index]).SetString(value)
		}
		installedPackages = append(installedPackages, installedPackage)
	}
	fmt.Printf("installedPackages: %+v\n", installedPackages)
	return installedPackages
}
