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

func TanzuDeployWorkloadByCommand(workload string, namespace string, gitRepository string, gitBranch string, workloadType string, hasTests string) error {
	log.Printf("deploying workload %s in namespace %s", workload, namespace)

	// execute cmd
	cmd := fmt.Sprintf("tanzu apps workload create %s --git-repo %s --git-branch %s --label \"apps.kubernetes.io/name=%s\" --label \"app.kubernetes.io/part-of=%s\" --label \"apps.tanzu.vmware.com/workload-type=%s\" --label \"apps.tanzu.vmware.com/has-tests=%s\" --label \"tanzu.app.live.view.application.actuator.port=8080\"-y -n %s", workload, gitRepository, gitBranch, workload, workload, workloadType, hasTests, namespace)
	output, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while deploying workload %s in namespace %s", workload, namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("workload %s deployed in namespace %s", workload, namespace)
		log.Printf("output: %s", output)
	}

	return err
}

func DeleteWorkload(name string, namespace string) error {

	cmd := fmt.Sprintf("tanzu apps workload delete -y %s", name)
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += " -A"
	}
	res, err := linux_util.ExecuteCmd(cmd)
	if err != nil && strings.Contains(res, "Deleted workload") {
		log.Println("Error while deleting the workload %s. Error %w, Output %s", name, err, res)
	}
	return err

}

func DeleteAllWorkload(namespace string) error {

	cmd := "tanzu apps workload delete --all -y"
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += " -A"
	}
	res, err := linux_util.ExecuteCmd(cmd)
	if err != nil && strings.Contains(res, "Deleted workloads") {
		log.Println("Error while deleting the workloads. Error %w, Output %s", err, res)
	}
	return err

}
