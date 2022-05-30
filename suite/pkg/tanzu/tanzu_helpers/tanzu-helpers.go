package tanzu_helpers

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzuCmds"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzu_libs"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/models"
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

func ValidateInstalledPackageVersion(name string, namespace string, expectedVersion string, timeoutInMins int, intervalInSeconds int) bool {
	log.Println("Executing: ValidateInstalledPackageStatus")
	result := false
	finalTimeout := timeoutInMins * 60
	for finalTimeout > 0 {
		pkg := tanzu_libs.GetInstalledPackages(name, namespace)
		if len(pkg) < 1 {
			log.Println("Package installation not started yet")
		} else if pkg[0].Status == "Reconcile succeeded" {
			if pkg[0].PackageVersion == expectedVersion {
				log.Printf("Package %s version %s installation is verified successfully.", name, expectedVersion)
				result = true
				break
			} else {
				log.Printf("Package %s version %s is installed.", name, pkg[0].PackageVersion)
				log.Printf("expected version : %s", expectedVersion)
				result = false
				break
			}
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

func CheckIfPackageRepositoryReconciled(name string, namespace string, recursiveCount int, secondsGap int) bool {
	log.Printf("checking package repository %s status", name)

	for ; recursiveCount >= 0; recursiveCount-- {
		// get status
		packageRepositoryStatus, err := tanzuCmds.TanzuGetPackageRepositoryStatus(name, namespace)
		if err != nil {
			log.Printf("error while getting package repository %s in namespace %s status", name, namespace)
			return false
		}

		// check
		if packageRepositoryStatus == "Reconciling" || packageRepositoryStatus == "" {
			log.Printf("package repository %s is getting reconciled", name)
			log.Printf("sleeping for %d seconds", secondsGap)
			time.Sleep(time.Duration(secondsGap) * time.Second)
		} else if packageRepositoryStatus == "Reconcile succeeded" {
			log.Printf("package repository %s reconcilation succeeded", name)
			return true
		} else if packageRepositoryStatus == "Reconcile Failed" {
			log.Printf("package repository %s reconcilation failed", name)
			return false
		} else {
			log.Printf("package repository %s reconcilation unknown", name)
			return false
		}
	}
	log.Printf(`package repository %s is not getting in "Reconcile succeeded" state`, name)
	return false
}

func UpdateTapValues(tapValuesSchema models.TapValuesSchema, tapName string, tapPackageName string, tapVersion string, tapNamespace string) error {
	log.Print("creating tempfile for tap values schema")
	tempFile, err := ioutil.TempFile("", "tap-values*.yaml")
	if err != nil {
		log.Printf("error while creating tempfile for tap values schema")
		return err
	} else {
		log.Print("created tempfile")
	}
	defer os.Remove(tempFile.Name())

	// write the updated schema to the temporary file
	err = utils.WriteYAMLFile(tempFile.Name(), tapValuesSchema)
	if err != nil {
		log.Print("error while writing updated tap values schema to YAML file")
		return err
	} else {
		log.Print("wrote tap values schema to file")
	}

	// update tap
	err = tanzu_libs.UpdateInstalledPackage(tapName, tapPackageName, tapVersion, tapNamespace, tempFile.Name(), "")
	if err != nil {
		log.Print("error while updating tap")
		return err
	} else {
		log.Print("updated tap")
	}
	return err
}
