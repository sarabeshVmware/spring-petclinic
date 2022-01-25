package tanzu_libs

// Usage:
//   tanzu apps workload [command]

// Aliases:
//   workload, workloads

// Available Commands:
//   apply       Apply configuration to a new or existing workload
//   create      Create a workload with specified configuration
//   delete      Delete workload(s)
//   get         Get details from a workload
//   list        Table listing of workloads
//   tail        Watch workload related logs
//   update      Update configuration of an existing workload
import (
	"fmt"
	"log"
	"strings"
	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

func DeleteWorkload(name string, namespace string) {

	cmd := fmt.Sprintf("tanzu apps workload delete -y %s", name)
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += " -A"
	}	
	res, err := linux_util.ExecuteCmd(cmd)
	if err != nil && strings.Contains(res, "Deleted workload"){
		log.Println("Error while deleting the workload %s. Error %w, Output %s", name, err, res)
	}

}
