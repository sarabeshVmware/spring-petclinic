package kubectl_libs

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

type GetWorkloadOutput struct {
	NAME, SOURCE, SUPPLYCHAIN, READY, REASON, AGE string
}

func GetWorkload(workloadName string, namespace string) []GetWorkloadOutput {
	workloads := []GetWorkloadOutput{}
	cmd := "kubectl get workload"
	if workloadName != "" {
		cmd += fmt.Sprintf(" %s", workloadName)
	}
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += " -A"
	}
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return workloads
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return workloads
	}

	header_index := 0
	if strings.HasPrefix(temp[0], "I04") {
		header_index = 1
	}

	if len(temp) <= header_index+1 {
		log.Printf("Output : %s", temp[0])
		return workloads
	}

	ss := linux_util.FieldIndices(temp[header_index])
	headers := linux_util.GetFields(temp[header_index], ss)
	for _, element := range temp[header_index+1:] {
		words := linux_util.GetFields(element, ss)
		var wl GetWorkloadOutput
		for index, value := range words {
			reflect.ValueOf(&wl).Elem().FieldByName(headers[index]).SetString(value)
		}
		workloads = append(workloads, wl)
	}

	fmt.Printf("workloads: %+v\n", workloads)
	return workloads
}
