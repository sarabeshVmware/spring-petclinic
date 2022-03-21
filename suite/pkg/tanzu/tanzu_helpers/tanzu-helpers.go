package tanzu_helpers

import (
	"log"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzu_libs"
)

func IsGrypeInstalled(namespace string) bool {
	log.Println("Executing: IsGrypeInstalled")
	grypeInstalled := false
	packages := tanzu_libs.ListInstalledPackages(namespace)
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
	packages := tanzu_libs.ListInstalledPackages(namespace)
	for _, element := range packages {
		if element.NAME == "scanning" {
			scanningInstalled = true
			break
		}
	}
	return scanningInstalled
}
