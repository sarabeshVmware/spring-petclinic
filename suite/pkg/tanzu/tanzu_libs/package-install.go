package tanzu_libs

// Usage:
//   tanzu package install INSTALLED_PACKAGE_NAME --package-name PACKAGE_NAME --version VERSION [flags]

import (
	"fmt"
	"log"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

func InstallPackage(installedPackageName string, packageName string, version string, namespace string, valuesFile string, pollTimeout string) error {

	cmd := fmt.Sprintf("tanzu package install %s --package-name %s --version %s --namespace %s", installedPackageName, packageName, version, namespace)
	if valuesFile != "" {
		cmd += fmt.Sprintf(" --values-file %s", valuesFile)
	}
	if pollTimeout != "" {
		cmd += fmt.Sprintf(" --poll-timeout %s", pollTimeout)
	}
	res, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while installing package %s (%s) in namespace %s", installedPackageName, packageName, namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", res)
	} else {
		log.Printf("package %s (%s) installed in namespace %s", installedPackageName, packageName, namespace)
		log.Printf("output: %s", res)
	}
	return err
}
