package kubectl_helpers

import (
	"fmt"
	"log"
	"strings"
	"time"

	kubectl_lib "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_libs"
	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

func ValidateAppLiveViewLabels(name string, namespace string) bool {
	log.Println("Validating 'App Live View' labels")
	raw := kubectl_lib.GetPodintentJson(name, namespace)
	if raw == nil {
		return false
	}
	log.Printf("Status.Template.Metadata.Labels.TanzuAppLiveView --> Expected : 'true', Observed: '%s'", raw.Status.Template.Metadata.Labels.TanzuAppLiveView)
	log.Printf("Status.Template.Metadata.Labels.TanzuAppLiveViewApplicationFlavours --> Expected : 'spring-boot', Observed: '%s'", raw.Status.Template.Metadata.Labels.TanzuAppLiveViewApplicationFlavours)
	if (raw.Status.Template.Metadata.Labels.TanzuAppLiveView == "true") && (raw.Status.Template.Metadata.Labels.TanzuAppLiveViewApplicationFlavours == "spring-boot") {
		log.Println("Validation passed")
		return true
	} else {
		log.Println("Validation failed")
		return false
	}
}

func ValidateAppLiveViewConventions(name string, namespace string) bool {
	validateConventions := [3]string{"appliveview-sample/app-live-view-connector", "appliveview-sample/app-live-view-appflavours", "appliveview-sample/app-live-view-systemproperties"}
	log.Println("Validating 'App Live View' conventions")
	raw := kubectl_lib.GetPodintentJson(name, namespace)
	if raw == nil {
		return false
	}
	log.Printf("Status.Template.Metadata.Annotations.ConventionsAppsTanzuVmwareComAppliedConventions: %s", raw.Status.Template.Metadata.Annotations.ConventionsAppsTanzuVmwareComAppliedConventions)
	for _, value := range validateConventions {
		if !(strings.Contains(raw.Status.Template.Metadata.Annotations.ConventionsAppsTanzuVmwareComAppliedConventions, value)) {
			log.Println("Validation failed")
			return false
		}
	}
	log.Println("Validation passed")
	return true
}

func ValidateSpringBootLabels(name string, namespace string) bool {
	log.Println("Validating 'Spring Boot' labels")
	raw := kubectl_lib.GetPodintentJson(name, namespace)
	if raw == nil {
		return false
	}
	log.Printf("Status.Template.Metadata.Labels.ConventionsAppsTanzuVmwareComFramework  --> Expected : 'spring-boot', Observed: '%s'", raw.Status.Template.Metadata.Labels.ConventionsAppsTanzuVmwareComFramework)
	if raw.Status.Template.Metadata.Labels.ConventionsAppsTanzuVmwareComFramework == "spring-boot" {
		log.Println("Validation passed")
		return true
	} else {
		log.Println("Validation failed")
		return false
	}
}

func ValidateSpringBootConventions(name string, namespace string) bool {
	validateConventions := [4]string{"spring-boot-convention/spring-boot", "spring-boot-convention/spring-boot-graceful-shutdown", "spring-boot-convention/spring-boot-web", "spring-boot-convention/spring-boot-actuator"}
	log.Println("Validating 'Spring Boot' conventions")
	raw := kubectl_lib.GetPodintentJson(name, namespace)
	if raw == nil {
		return false
	}
	log.Printf("Status.Template.Metadata.Annotations.ConventionsAppsTanzuVmwareComAppliedConventions: %s", raw.Status.Template.Metadata.Annotations.ConventionsAppsTanzuVmwareComAppliedConventions)
	for _, value := range validateConventions {
		if !(strings.Contains(raw.Status.Template.Metadata.Annotations.ConventionsAppsTanzuVmwareComAppliedConventions, value)) {
			log.Println("Validation failed")
			return false
		}
	}
	log.Println("Validation passed")
	return true
}

func ValidateDeveloperConventions(name string, namespace string) bool {
	validateConventions := [2]string{"developer-conventions/live-update-convention", "developer-conventions/add-source-image-label"}
	log.Println("Validating 'Developer' conventions")
	raw := kubectl_lib.GetPodintentJson(name, namespace)
	if raw == nil {
		return false
	}
	log.Printf("Status.Template.Metadata.Annotations.ConventionsAppsTanzuVmwareComAppliedConventions: %s", raw.Status.Template.Metadata.Annotations.ConventionsAppsTanzuVmwareComAppliedConventions)
	for _, value := range validateConventions {
		if !(strings.Contains(raw.Status.Template.Metadata.Annotations.ConventionsAppsTanzuVmwareComAppliedConventions, value)) {
			log.Println("Validation failed")
			return false
		}
	}
	log.Println("Validation passed")
	return true
}

func GetLatestImageRepositoryStatus(name string, namespace string) string {
	log.Println("Get imagerepository status...")
	imagerepos := kubectl_lib.GetImageRepositories(name, namespace)
	if len(imagerepos) < 1 {
		log.Println("No images found")
		return "None"
	}
	log.Printf("imagerepository status : %s", imagerepos[len(imagerepos)-1].READY)
	return imagerepos[len(imagerepos)-1].READY
}

func GetLatestBuildStatus(name string, namespace string) string {
	log.Println("Get build status...")
	builds := kubectl_lib.GetBuilds(name, namespace)
	if len(builds) < 1 {
		log.Println("No builds found")
		return "None"
	}
	log.Printf("build status : %s", builds[len(builds)-1].SUCCEEDED)
	return builds[len(builds)-1].SUCCEEDED
}

func GetLatestImageStatus(namespace string) string {
	log.Println("Get latest image status...")
	image := kubectl_lib.GetLatestImage(namespace)
	log.Printf("latest image status : %s", image.READY)
	return image.READY
}

func GetKsvcStatus(name string, namespace string) string {
	log.Println("Get ksvc image status...")
	ksvc := kubectl_lib.GetKsvc(name, namespace)
	log.Printf("ksvc image status : %s", ksvc[0].READY)
	return ksvc[0].READY
}

func ValidateImageScans(name string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Println("Validating image scans")
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		imageScan := kubectl_lib.GetImageScan(name, namespace)
		if (imageScan == kubectl_lib.GetImageScanOutput{}) {
			log.Println("Image scan is not started yet")
		} else if imageScan.PHASE == "Completed" && imageScan.CRITICAL == "" && imageScan.HIGH == "" && imageScan.MEDIUM == "" && imageScan.LOW == "" && imageScan.UNKNOWN == "" && imageScan.CVETOTAL == "" {
			log.Println("Image scan complete successfully")
			result = true
			break
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if !result {
		log.Printf("Image scan not completed successfully after %d mins", timeoutInMins)
	}
	return result
}

func ValidateSourceScans(name string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Println("Validating source scans")
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		srcScan := kubectl_lib.GetSourceScan(name, namespace)
		if (srcScan == kubectl_lib.GetSourceScanOutput{}) {
			log.Println("Source scan is not started yet")
		} else if srcScan.PHASE == "Completed" && srcScan.CRITICAL == "" && srcScan.HIGH == "" && srcScan.MEDIUM == "" && srcScan.LOW == "" && srcScan.UNKNOWN == "" && srcScan.CVETOTAL == "" {
			log.Println("Source scan complete successfully")
			result = true
			break
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if !result {
		log.Printf("Source scan not completed successfully after %d mins", timeoutInMins)
	}
	return result
}

func ValidatePipelineExists(name string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Println("Validating pipeline exists")
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		pipeline := kubectl_lib.GetPipeline(name, namespace)
		if (pipeline == kubectl_lib.GetPipelineOutput{}) {
			log.Println("Pipeline not created yet")
		} else if pipeline.NAME == name {
			log.Println("Pipeline created successfully")
			result = true
			break
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if !result {
		log.Printf("Pipeline not created after %d mins", timeoutInMins)
	}
	return result
}

func ValidatePipelineRuns(name string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Println("Validating pipeline runs")
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		prs := kubectl_lib.GetPipelineRuns(name, namespace)
		if (prs == kubectl_lib.GetPipelineRunsOutput{}) {
			log.Println("Pipeline runs not created yet")
		} else if prs.SUCCEEDED == "True" && prs.REASON == "Succeeded" {
			log.Println("Pipeline runs created successfully")
			result = true
			break
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if !result {
		log.Printf("Pipeline runs not created after %d mins", timeoutInMins)
	}
	return result
}

func ValidateServiceBindings(name string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Println("Validating service bindings")
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		svcBindings := kubectl_lib.GetServiceBindings(name, namespace)
		if (svcBindings == kubectl_lib.GetServiceBindingsOutput{}) {
			log.Println("Service bindings not ready yet")
		} else if svcBindings.READY == "True" && svcBindings.REASON == "Ready" {
			log.Printf("Service bindings %s is ready", svcBindings.NAME)
			result = true
			break
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if !result {
		log.Printf("Service binding is not ready after %d mins", timeoutInMins)
	}
	return result
}

func ValidateTrainingPortalStatus(name string, timeoutInMins int, intervalInSeconds int) bool {
	log.Println("Validating training portals")
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		tp := kubectl_lib.GetTrainingPortals(name)
		if (tp == kubectl_lib.GetTrainingPortalsOutput{}) {
			log.Println("Training portal is not ready yet")
		} else if tp.STATUS == "Running" {
			log.Printf("Training portal %s is ready", tp.NAME)
			result = true
			break
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if !result {
		log.Printf("Training portal is not ready after %d mins", timeoutInMins)
	}
	return result
}

func ValidateLearningCenter(name string, namespace string) bool {
	log.Println("Validating 'Learning Center'")
	img := kubectl_lib.GetIngress(name, namespace)
	cmd1 := fmt.Sprintf("echo '%s %s' >> /etc/hosts", img.ADDRESS, img.HOSTS)
	linux_util.ExecuteCmd(cmd1)
	cmd2 := fmt.Sprintf("curl -i %s", img.HOSTS)
	res, err := linux_util.ExecuteCmd(cmd2)
	if err != nil {
		log.Println("error")
	}
	return strings.Contains(res, "HTTP/1.1 302 Found")
}

func VerifyBuildStatus(namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Println("Validating build status")
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		builds := kubectl_lib.GetBuilds("", namespace)
		if len(builds) < 1 {
			log.Println("Builds are not generated yet")
		} else if builds[len(builds)-1].SUCCEEDED == "True" {
			log.Printf("Build %s status is verified successfully. Status is %s", builds[len(builds)-1].NAME, builds[len(builds)-1].SUCCEEDED)
			result = true
			break
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if !result {
		log.Printf("Build is not ready after %d mins", timeoutInMins)
	}
	return result
}

func GetPodIntentStatus(name string, namespace string) string {
	log.Println("Get podintents status...")
	podintents := kubectl_lib.GetPodintent(name, namespace)
	if len(podintents) > 0 {
		log.Println("Found podintents")
		log.Printf("podintents status : %s", podintents[0].READY)
		return podintents[0].READY
	} else {
		log.Printf("podintent not found")
		return "None"
	}

}

func GetWorkloadStatus(name string, namespace string) string {
	log.Println("Get workload status...")
	workloads := kubectl_lib.GetWorkload(name, namespace)
	if len(workloads) > 0 {
		log.Println("Found workloads")
		log.Printf("workloads status : %s", workloads[0].READY)
		return workloads[0].READY
	} else {
		log.Printf("workload not found")
		return "None"
	}

}

func VerifyKsvcStatus(name string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Println("Validating ksvc status")
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		ksvc := kubectl_lib.GetKsvc(name, namespace)
		if len(ksvc) < 1 {
			log.Println("Knative services are not generated yet")
		} else if ksvc[0].READY == "True" {
			log.Printf("Knative %s status is verified successfully. Status is %s", ksvc[0].NAME, ksvc[0].READY)
			result = true
			break
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if !result {
		log.Printf("Ksvc is not generated after %d mins", timeoutInMins)
	}
	return result
}

func VerifyImageRepositoryStatus(name string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Println("Validating Image repository status")
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		img := kubectl_lib.GetImageRepositories(name, namespace)
		if len(img) < 1 {
			log.Println("Image repository is not generated yet")
		} else if img[0].READY == "True" {
			log.Printf("Image repository %s status is verified successfully. Status is %s", img[0].NAME, img[0].READY)
			result = true
			break
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if !result {
		log.Printf("Image repository is not generated after %d mins", timeoutInMins)
	}
	return result
}

func VerifyGitRepoStatus(name string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Println("Validating Git repository status")
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		repo := kubectl_lib.GetGitrepo(name, namespace)
		if len(repo) < 1 {
			log.Println("gitrepo is not generated yet")
		} else if repo[0].READY == "True" {
			log.Printf("gitrepo %s status is verified successfully, status is %s", repo[0].NAME, repo[0].READY)
			result = true
			break
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if !result {
		log.Printf("Git repository is not generated after %d mins", timeoutInMins)
	}
	return result
}

func VerifyTaskrunStatus(namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Println("Validating task run status")
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		taskruns := kubectl_lib.GetTaskruns("", namespace)
		if len(taskruns) < 1 {
			log.Println("taskruns are not generated yet")
		} else if taskruns[len(taskruns)-1].SUCCEEDED == "True" {
			log.Printf("taskrun %s status is verified successfully, status is %s", taskruns[len(taskruns)-1].NAME, taskruns[len(taskruns)-1].SUCCEEDED)
			result = true
			break
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if !result {
		log.Printf("task run is not generated after %d mins", timeoutInMins)
	}
	return result
}

func VerifyPodIntentStatus(name string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Println("Validating podintent status")
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		podintentStatus := GetPodIntentStatus(name, namespace)
		if podintentStatus == "True" {
			log.Printf("podintent %s status is %s", name, podintentStatus)
			result = true
			break
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if !result {
		log.Printf("Podintent is not ready after %d mins", timeoutInMins)
	}
	return result
}

func LogFailedResourcesDetails(namespace string) {
	pkgiList := kubectl_lib.GetPkgi("", namespace)
	for _, value := range pkgiList {
		if value.DESCRIPTION != "Reconcile succeeded" {
			log.Printf("Describe pkgi %s", value.NAME)
			kubectl_lib.DescribePkgi(value.NAME, namespace)
		}
	}

}

func ValidateTAPInstallation(pkgName string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Println("Validating TAP installation status")
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		pkg := kubectl_lib.GetPkgi(pkgName, namespace)
		if len(pkg) < 1 {
			log.Println("tap reconcile is not successed yet")
		} else if pkg[len(pkg)-1].DESCRIPTION == "Reconcile succeeded" {
			log.Printf("package %s status is verified successfully, status is %s", pkg[len(pkg)-1].NAME, pkg[len(pkg)-1].DESCRIPTION)
			result = true
			break
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if !result {
		log.Printf("TAP install did not reconcile after %d mins", timeoutInMins)
	}
	return result
}
