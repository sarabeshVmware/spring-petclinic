package kubectl_libs

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

type GetPkgiOutput struct {
	NAME, PACKAGE_NAME, PACKAGE_VERSION, DESCRIPTION, AGE string
}

func GetPkgi(pkgName string, namespace string) []GetPkgiOutput {
	pkgsi := []GetPkgiOutput{}
	cmd := "kubectl get pkgi"
	if pkgName != "" {
		cmd += fmt.Sprintf(" %s", pkgName)
	}
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += " -A"
	}
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return pkgsi
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return pkgsi
	}

	header_index := 0
	if strings.HasPrefix(temp[0], "I04") {
		header_index = 1
	}

	if len(temp) <= header_index+1 {
		log.Printf("Output : %s", temp[0])
		return pkgsi
	}
	ss := linux_util.FieldIndicesWithSingleSpace(temp[header_index])
	headers := linux_util.GetFields(temp[header_index], ss)
	for index, ele := range headers {
		headers[index] = strings.ReplaceAll(ele, " ", "_")
	}
	for _, element := range temp[header_index+1:] {
		words := linux_util.GetFields(element, ss)
		var pkg GetPkgiOutput
		for index, value := range words {
			reflect.ValueOf(&pkg).Elem().FieldByName(headers[index]).SetString(value)
		}
		pkgsi = append(pkgsi, pkg)
	}

	fmt.Printf("pkgsi: %+v\n", pkgsi)
	return pkgsi
}

func DescribePkgi(pkgName string, namespace string) {
	cmd := fmt.Sprintf("kubectl describe pkgi %s -n %s", pkgName, namespace)
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return
	}
	log.Println(response)
}
