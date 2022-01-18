package kubectl_helpers

import (
	"fmt"
	"log"
	kubectl_lib "pkg/kubectl/kubectl_libs"
	"pkg/utils/linux_util"
	"strings"
)

func ValidateAppLiveViewLabels(name string, namespace string) bool {
	log.Println("Validating 'App Live View' labels")
	raw := kubectl_lib.GetPodintentJson(name, namespace)
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

func GetImageRepositoryStatus(name string, namespace string) string {
	log.Println("Get imagerepository status...")
	imagerepos := kubectl_lib.GetImageRepositories(name, namespace)
	if len(imagerepos) > 1 {
		log.Println("Multiple images found. Returning status of first image.")
	}
	log.Printf("imagerepository status : %s", imagerepos[0].READY)
	return imagerepos[0].READY
}

func GetBuildStatus(name string, namespace string) string {
	log.Println("Get build status...")
	builds := kubectl_lib.GetBuilds(name, namespace)
	if len(builds) > 1 {
		log.Println("Multiple builds found. Returning status of first build.")
	}
	log.Printf("build status : %s", builds[0].SUCCEEDED)
	return builds[0].SUCCEEDED
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
	log.Printf("ksvc image status : %s", ksvc.READY)
	return ksvc.READY
}

func ValidateImageScans(name string, namespace string) bool {
	log.Println("Validating image scans")
	imageScan := kubectl_lib.GetImageScan(name, namespace)
	if imageScan.PHASE == "Completed" && imageScan.CRITICAL == "" && imageScan.HIGH == "" && imageScan.MEDIUM == "" && imageScan.LOW == "" && imageScan.UNKNOWN == "" && imageScan.CVETOTAL == "" {
		return true
	}
	return false
}

func ValidateSourceScans(name string, namespace string) bool {
	log.Println("Validating source scans")
	srcScan := kubectl_lib.GetSourceScan(name, namespace)
	if srcScan.PHASE == "Completed" && srcScan.CRITICAL == "" && srcScan.HIGH == "" && srcScan.MEDIUM == "" && srcScan.LOW == "" && srcScan.UNKNOWN == "" && srcScan.CVETOTAL == "" {
		return true
	}
	return false
}

func ValidatePipelineExists(name string, namespace string) bool {
	log.Println("Validating pipeline exists")
	pipeline := kubectl_lib.GetPipeline(name, namespace)
	return (kubectl_lib.GetPipelineOutput{}) != pipeline
}

func ValidatePipelineRuns(name string, namespace string) bool {
	log.Println("Validating pipeline runs")
	prs := kubectl_lib.GetPipelineRuns(name, namespace)
	if prs.SUCCEEDED == "True" && prs.REASON == "Succeeded" {
		return true
	}
	return false
}

func ValidateServiceBindings(name string, namespace string) bool {
	log.Println("Validating service bindings")
	svcBindings := kubectl_lib.GetServiceBindings(name, namespace)
	if svcBindings.READY == "True" && svcBindings.REASON == "Ready" {
		return true
	}
	return false
}

func ValidateTrainingPortalStatus(name string) bool {
	log.Println("Validating training portals")
	tp := kubectl_lib.GetTrainingPortals(name)
	return tp.STATUS == "Running"
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
