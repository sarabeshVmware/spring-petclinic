package tanzu_libs

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)


type ListServiceInstancesOutput struct {
	NAME, KIND, SERVICE_TYPE, AGE, SERVICE_REF string
}

func ListServiceInstances(namespace string) []ListServiceInstancesOutput {
	serviceInstances := []ListServiceInstancesOutput{}
	cmd := "tanzu service instance list -owide"
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += " -A"
	}
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return serviceInstances
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 3 {
		log.Printf("Output : %s", temp[0])
		return serviceInstances
	}

	indexSpans := linux_util.FieldIndicesWithSingleSpace(temp[2])
	headers := linux_util.GetFields(temp[2], indexSpans)

	for index, ele := range headers {
		headers[index] = strings.ReplaceAll(ele, " ", "_")
	}

	for _, element := range temp[3:] {
		words := linux_util.GetFields(element, indexSpans)
		var serviceInstance ListServiceInstancesOutput
		for index, value := range words {
			reflect.ValueOf(&serviceInstance).Elem().FieldByName(headers[index]).SetString(value)
		}
		serviceInstances = append(serviceInstances, serviceInstance)
	}
	fmt.Printf("serviceInstances: %+v\n", serviceInstances)
	return serviceInstances


}
