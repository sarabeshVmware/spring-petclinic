package pivnet_libs

import (
	"fmt"
	"log"
	"strings"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

func Login(host string, apiToken string) bool {
	cmd := fmt.Sprintf("pivnet-cli login --host %s --api-token %s", host, apiToken)

	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return false
	}

	if !strings.Contains(response, "Logged-in successfully") {
		log.Println("Login failed")
		return false
	}
	return true

}
