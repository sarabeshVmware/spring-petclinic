package kubectl_libs

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

type GetServiceBindingsOutput struct {
	NAME, READY, REASON, AGE string
}

func GetServiceBindings(name string, namespace string) GetServiceBindingsOutput {
	var svcBindings GetServiceBindingsOutput
	cmd := "kubectl get servicebindings"
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
		return svcBindings
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return svcBindings
	}

	ss := linux_util.FieldIndices(temp[0])
	headers, words := linux_util.GetFields(temp[0], ss), linux_util.GetFields(temp[1], ss)

	for index, value := range words {
		reflect.ValueOf(&svcBindings).Elem().FieldByName(headers[index]).SetString(value)
	}
	fmt.Printf("svcBindings: %+v\n", svcBindings)
	return svcBindings
}
