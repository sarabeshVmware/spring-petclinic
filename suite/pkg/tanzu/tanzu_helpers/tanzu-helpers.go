package tanzu_helpers

import (
	"log"
	"time"

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

func ValidateInstalledPackageStatus(name string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Println("Executing: ValidateInstalledPackageStatus")
	result := false
	finalTimeout := timeoutInMins * 60
	for finalTimeout > 0 {
		pkg := tanzu_libs.GetInstalledPackages(name, namespace)
		if len(pkg) < 1 {
			log.Println("Package installation not started yet")
		} else if pkg[0].Status == "Reconcile succeeded" {
			log.Printf("Package %s installation is verified successfully.", name)
			result = true
			break
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	return result
}

func ValidateWorkloadDeleted(workloadName string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Println("Executing: ValidateWorkloadDeleted")
	result := false
	finalTimeout := timeoutInMins * 60
	for finalTimeout > 0 {
		wl := tanzu_libs.ListAppWorkloads("", namespace)
		found := false
		for _, element := range wl {
			if element.NAME == workloadName {
				found = true
				break
			}
		}
		if !found {
			log.Printf("Workload %s not found. Deleted successfully", workloadName)
			return true
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	return result
}
