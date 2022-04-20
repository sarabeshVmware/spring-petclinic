package kubectl_libs

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

type GetDeliverablesOutput struct {
	NAME, SOURCE, DELIVERY, READY, REASON, AGE string
}

func GetDeliverables(deliverableName string, namespace string) []GetDeliverablesOutput {
	deliverables := []GetDeliverablesOutput{}
	cmd := "kubectl get deliverable"
	if deliverableName != "" {
		cmd += fmt.Sprintf(" %s", deliverableName)
	}
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += " -A"
	}
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return deliverables
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return deliverables
	}

	header_index := 0
	if strings.HasPrefix(temp[0], "I04") {
		header_index = 1
	}

	if len(temp) <= header_index+1 {
		log.Printf("Output : %s", temp[0])
		return deliverables
	}

	ss := linux_util.FieldIndices(temp[header_index])
	headers := linux_util.GetFields(temp[header_index], ss)
	for _, element := range temp[header_index+1:] {
		words := linux_util.GetFields(element, ss)
		var wl GetDeliverablesOutput
		for index, value := range words {
			reflect.ValueOf(&wl).Elem().FieldByName(headers[index]).SetString(value)
		}
		deliverables = append(deliverables, wl)
	}

	fmt.Printf("deliverables: %+v\n", deliverables)
	return deliverables
}
