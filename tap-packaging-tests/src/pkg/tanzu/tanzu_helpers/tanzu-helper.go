package tanzu_helper

import (
	"log"
	tanzu_lib "pkg/tanzu/tanzu_libs"
)

func IsGrypeInstalled(namespace string) bool {
	log.Println("Executing: IsGrypeInstalled")
	grypeInstalled := false
	packages := tanzu_lib.ListInstalledPackages(namespace)
	for _, element := range packages {
		if element.NAME == "grype" {
			grypeInstalled = true
			break
		}
	}
	return grypeInstalled
}

func IsScanningInstalled(namespace string) bool {
	log.Println("Executing: IsScanningInstalled")
	scanningInstalled := false
	packages := tanzu_lib.ListInstalledPackages(namespace)
	for _, element := range packages {
		if element.NAME == "scanning" {
			scanningInstalled = true
			break
		}
	}
	return scanningInstalled
}
