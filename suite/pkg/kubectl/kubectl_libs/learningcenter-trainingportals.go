package kubectl_libs

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

type GetTrainingPortalsOutput struct {
	NAME, URL, ADMINUSERNAME, ADMINPASSWORD, STATUS string
}

func GetTrainingPortals(name string) GetTrainingPortalsOutput {
	var tps GetTrainingPortalsOutput
	cmd := "kubectl get trainingportals"
	if name != "" {
		cmd += fmt.Sprintf(" %s", name)
	}
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return tps
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return tps
	}

	ss := linux_util.FieldIndices(temp[0])
	headers, words := linux_util.GetFields(temp[0], ss), linux_util.GetFields(temp[1], ss)

	for index, value := range words {
		reflect.ValueOf(&tps).Elem().FieldByName(headers[index]).SetString(value)
	}
	fmt.Printf("trainingPortals: %+v\n", tps)
	return tps
}
