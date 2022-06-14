package kubectl_libs

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

type GetPipelineOutput struct {
	NAME, AGE string
}

func GetPipeline(name string, namespace string) GetPipelineOutput {
	var pipeline GetPipelineOutput
	cmd := "kubectl get pipeline"
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

func DeletePipeline(name string, namespace string) (string, error) {
	cmd := fmt.Sprintf("kubectl delete pipeline %s -n %s", name, namespace)
	res, err := linux_util.ExecuteCmd(cmd)
	return res, err
}
