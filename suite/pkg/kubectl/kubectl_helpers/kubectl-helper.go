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

func VerifyBuildStatus(namespace string) bool {
	count := 20
	for count <= 20 {
		if count == 0 {
			log.Println("Builds are not generated after 10 mins")
			return false
		}
		builds := kubectl_lib.GetBuilds("", namespace)
		if len(builds) < 1 {
			log.Println("Builds are not generated yet")
		} else {
			status := builds[len(builds)-1].SUCCEEDED
			build_name := builds[len(builds)-1].NAME
			if status == "Unknown" {
				log.Printf("Build %s status is Unknown", build_name)

			} else if status == "True" {
				log.Printf("Build %s status is verified successfully. Status is %s", build_name, status)
				return true
			}
		}
		log.Printf("Waiting for 30s for builds getting generated ...")
		time.Sleep(30 * time.Second)
		count -= 1
	}
	return false
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

func VerifyKsvcStatus(name string, namespace string) bool {
	count := 10
	for count <= 10 {
		if count == 0 {
			log.Println("Ksvc are not generated after 5 mins")
			return false
		}
		ksvc := kubectl_lib.GetKsvc(name, namespace)
		if len(ksvc) < 1 {
			log.Println("Knative services are not generated yet")
		} else {
			status := ksvc[0].READY
			ksvc_name := ksvc[0].NAME
			if status == "True" {
				log.Printf("Knative %s status is verified successfully. Status is %s", ksvc_name, status)
				return true
			} else {
				log.Printf("Knative %s status is not verified successfully. Status is %s", ksvc_name, status)
			}
		}
		log.Printf("Waiting for 30s for ksvc getting generated ...")
		time.Sleep(30 * time.Second)
		count -= 1
	}
	return false
}

func VerifyImageRepositoryStatus(name string, namespace string) bool {
	count := 20
	for count <= 20 {
		if count == 0 {
			log.Println("Image repositories are not generated after 10 mins")
			return false
		}
		img := kubectl_lib.GetImageRepositories(name, namespace)
		if len(img) < 1 {
			log.Println("Image repository is not generated yet")
		} else {
			status := img[0].READY
			img_name := img[0].NAME
			if status == "True" {
				log.Printf("Image repository %s status is verified successfully. Status is %s", img_name, status)
				return true
			} else {
				log.Printf("Image repository %s status is not verified successfully. Status is %s", img_name, status)

			}
		}
		log.Printf("Waiting for 30s for image repository getting generated ...")
		time.Sleep(30 * time.Second)
		count -= 1
	}
	return false
}

func VerifyGitRepoStatus(name string, namespace string) bool {
	count := 10
	for count <= 10 {
		if count == 0 {
			log.Println("gitrepo is not generated after 5 mins")
			return false
		}
		repo := kubectl_lib.GetGitrepo(name, namespace)
		if len(repo) < 1 {
			log.Println("gitrepo is not generated yet")
		} else {
			status := repo[0].READY
			repoName := repo[0].NAME
			if status == "True" {
				log.Printf("gitrepo %s status is verified successfully, status is %s", repoName, status)
				return true
			} else {
				log.Printf("gitrepo %s status is not verified successfully, status is %s", repoName, status)

			}
		}
		log.Printf("waiting for 30s for Git repo getting generated ...")
		time.Sleep(30 * time.Second)
		count -= 1
	}
	return false
}

func VerifyTaskrunStatus(namespace string) bool {
	count := 20
	for count <= 20 {
		if count == 0 {
			log.Println("taskruns are not generated after 10 mins")
			return false
		}
		taskruns := kubectl_lib.GetTaskruns("", namespace)
		if len(taskruns) < 1 {
			log.Println("taskruns are not generated yet")
		} else {
			status := taskruns[len(taskruns)-1].SUCCEEDED
			taskrunName := taskruns[len(taskruns)-1].NAME
			if status == "Unknown" {
				log.Printf("taskrun %s status is Unknown", taskrunName)

			} else if status == "True" {
				log.Printf("taskrun %s status is verified successfully, status is %s", taskrunName, status)
				return true
			}
		}
		log.Printf("waiting for 30s for taskruns getting generated ...")
		time.Sleep(30 * time.Second)
		count -= 1
	}
	return false
}
