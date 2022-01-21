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

func ValidateInstalledPackageStatus(name string, namespace string, status string) bool {
	log.Println("Executing: ValidateInstalledPackageStatus")
	if status == "" {
		status = "Reconcile succeeded" //Default validation status
	}
	pkg := tanzu_libs.GetInstalledPackages(name, namespace)
	return pkg[0].Status == status
}
