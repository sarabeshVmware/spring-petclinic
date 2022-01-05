package kubectl_helper

import (
	"log"
	kubectl_lib "pkg/kubectl/kubectl_lib"
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
	log.Printf("Status.Template.Metadata.Annotations.ConventionsAppsTanzuVmwareComAppliedConventions", raw.Status.Template.Metadata.Annotations.ConventionsAppsTanzuVmwareComAppliedConventions)
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
	log.Printf("Status.Template.Metadata.Annotations.ConventionsAppsTanzuVmwareComAppliedConventions", raw.Status.Template.Metadata.Annotations.ConventionsAppsTanzuVmwareComAppliedConventions)
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
	log.Printf("Status.Template.Metadata.Annotations.ConventionsAppsTanzuVmwareComAppliedConventions", raw.Status.Template.Metadata.Annotations.ConventionsAppsTanzuVmwareComAppliedConventions)
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

func GetKsvcImageStatus(name string, namespace string) string {
	log.Println("Get ksvc image status...")
	ksvIimage := kubectl_lib.GetKsvcImage(name, namespace)
	log.Printf("ksvc image status : %s", ksvIimage.READY)
	return ksvIimage.READY
}
