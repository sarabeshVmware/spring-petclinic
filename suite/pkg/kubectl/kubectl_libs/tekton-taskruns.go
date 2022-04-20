package kubectl_libs

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

type GetTaskrunsOutput struct {
	NAME, SUCCEEDED, REASON, STARTTIME, COMPLETIONTIME string
}

func GetTaskruns(taskrunName string, namespace string) []GetTaskrunsOutput {
	taskruns := []GetTaskrunsOutput{}
	cmd := "kubectl get taskruns"
	if taskrunName != "" {
		cmd += fmt.Sprintf(" %s", taskrunName)
	}
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += " -A"
	}
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return taskruns
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return taskruns
	}

	header_index := 0
	if strings.HasPrefix(temp[0], "I04") {
		header_index = 1
	}

	if len(temp) <= header_index+1 {
		log.Printf("Output : %s", temp[0])
		return taskruns
	}

	ss := linux_util.FieldIndices(temp[header_index])
	headers := linux_util.GetFields(temp[header_index], ss)
	for _, element := range temp[header_index+1:] {
		words := linux_util.GetFields(element, ss)
		var taskrun GetTaskrunsOutput
		for index, value := range words {
			reflect.ValueOf(&taskrun).Elem().FieldByName(headers[index]).SetString(value)
		}
		taskruns = append(taskruns, taskrun)
	}

	fmt.Printf("taskruns: %+v\n", taskruns)
	return taskruns
}
