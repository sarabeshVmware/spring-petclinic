package kubectl_libs

import (
	"fmt"
	"log"
	linux_util "pkg/utils/linux_util"
	"reflect"
	"strings"
)

type GetPipelineRunsOutput struct {
	NAME, SUCCEEDED, REASON, STARTTIME, COMPLETIONTIME string
}

func GetPipelineRuns(name string, namespace string) GetPipelineRunsOutput {
	var pipeline GetPipelineRunsOutput
	cmd := "kubectl get prs"
	if name != "" {
		cmd += fmt.Sprintf(" %s", name)
	}
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += " -A"
	}
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return pipeline
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return pipeline
	}

	ss := linux_util.FieldIndices(temp[0])
	headers, words := linux_util.GetFields(temp[0], ss), linux_util.GetFields(temp[1], ss)

	for index, value := range words {
		reflect.ValueOf(&pipeline).Elem().FieldByName(headers[index]).SetString(value)
	}
	fmt.Printf("pipeline: %+v\n", pipeline)
	return pipeline
}
