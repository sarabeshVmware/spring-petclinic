package tanzu_libs

// Usage:
//   tanzu package install INSTALLED_PACKAGE_NAME --package-name PACKAGE_NAME --version VERSION [flags]

import (
	"fmt"
	"log"
	"strings"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

func InstallePackage(installedPackageName string, packageName string, version string, namespace string, valuesFile string, pollTimeout string) {

	cmd := fmt.Sprintf("tanzu package install %s --package-name %s --version %s --namespace %s", installedPackageName, packageName, version, namespace)
	if valuesFile != "" {
		cmd += fmt.Sprintf(" --values-file %s", valuesFile)
	}
	if pollTimeout != "" {
		cmd += fmt.Sprintf(" --poll-timeout %s", pollTimeout)
	}
	res, err := linux_util.ExecuteCmd(cmd)
	if err != nil && !strings.Contains(res, "Uninstalled package") {
		log.Printf("Error while deleting the package %s. Error %v,  Output %s", packageName, err, res)
	}

}
