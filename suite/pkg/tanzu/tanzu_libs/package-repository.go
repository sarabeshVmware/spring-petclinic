package tanzu_libs

// Usage:
//   tanzu package repository [command]

// Available Commands:
//   add         Add a package repository
//   delete      Delete a package repository
//   get         Get details for a package repository
//   list        List package repositories
//   update      Update a package repository

import (
	"fmt"
	"log"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

func TanzuUpdatePackageRepository(name string, repositoryUrl string, namespace string) error {
	log.Printf("updating package %s in namespace %s", name, namespace)

	// execute cmd
	cmd := fmt.Sprintf("tanzu package repository update %s --url %s -n %s", name, repositoryUrl, namespace)
	output, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while updating repository %s in namespace %s", name, namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("package %s updated in namespace %s", name, namespace)
		log.Printf("output: %s", output)
	}

	return err
}
