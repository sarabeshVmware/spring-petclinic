package kubectl_helpers

import (
	"encoding/base64"
	"encoding/json"
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
	log.Printf("Validating 'Spring Boot' labels, name: %s, namespace: %s", name, namespace)
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
	log.Printf("Get imagerepository status, name: %s, namespace: %s", name, namespace)
	imagerepos := kubectl_lib.GetImageRepositories(name, namespace)
	if len(imagerepos) < 1 {
		log.Println("No images found")
		return "None"
	}
	log.Printf("imagerepository status : %s", imagerepos[len(imagerepos)-1].READY)
	return imagerepos[len(imagerepos)-1].READY
}

func GetLatestBuildStatus(name string, namespace string) string {
	log.Printf("Get build status, name: %s, namespace: %s", name, namespace)
	builds := kubectl_lib.GetBuilds(name, namespace)
	if len(builds) < 1 {
		log.Println("No builds found")
		return "None"
	}
	log.Printf("build status : %s", builds[len(builds)-1].SUCCEEDED)
	return builds[len(builds)-1].SUCCEEDED
}

func GetLatestImageStatus(namespace string) string {
	log.Printf("Get latest image status, namespace: %s", namespace)
	image := kubectl_lib.GetLatestImage(namespace)
	log.Printf("latest image status : %s", image.READY)
	return image.READY
}

func GetKsvcStatus(name string, namespace string) string {
	log.Printf("Get ksvc image status, name: %s, namespace: %s", name, namespace)
	ksvc := kubectl_lib.GetKsvc(name, namespace)
	log.Printf("ksvc image status : %s", ksvc[0].READY)
	return ksvc[0].READY
}

func ValidateImageScans(name string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Printf("Validating image scans, name: %s, namespace: %s", name, namespace)
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		imageScan := kubectl_lib.GetImageScan(name, namespace)
		if (imageScan == kubectl_lib.GetImageScanOutput{}) {
			log.Println("Image scan is not started yet")
		} else if imageScan.PHASE == "Completed" && ((imageScan.CRITICAL >= "1") || (imageScan.HIGH >= "1") || (imageScan.UNKNOWN >= "1")) {
			log.Println("Image scan complete, CVE(s) found")
			// TODO: tanzu insight list CVEs
			break
		} else if imageScan.PHASE == "Completed" && (imageScan.CRITICAL == "0" || imageScan.CRITICAL == "") && (imageScan.HIGH == "0" || imageScan.HIGH == "") && (imageScan.UNKNOWN == "0" || imageScan.UNKNOWN == "") {
			log.Println("Image scan complete successfully")
			result = true
			break
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if !result {
		log.Printf("Image scan failed/not completed after %d mins", timeoutInMins)
		_, err := kubectl_lib.DescribeImageScan(name, namespace)
		if err != nil {
			log.Printf("error :%s", err)
		}
	}
	return result
}

func ValidateSourceScans(name string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Printf("Validating source scans, name: %s, namespace: %s", name, namespace)
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		srcScan := kubectl_lib.GetSourceScan(name, namespace)
		if (srcScan == kubectl_lib.GetSourceScanOutput{}) {
			log.Println("Source scan is not started yet")
		} else if srcScan.PHASE == "Completed" && srcScan.CVETOTAL >= "1" {
			log.Println("Source scan complete, CVE(s) found")
			// TODO: tanzu insight list CVEs
			break
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
	log.Printf("Validating pipeline exists, name: %s, namespace: %s", name, namespace)
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

func ValidatePipelineRuns(prefix string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Printf("Validating pipeline runs, prefix: %s, namespace: %s", prefix, namespace)
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		prs := kubectl_lib.GetPipelineRuns("", namespace)
		if (prs == kubectl_lib.GetPipelineRunsOutput{}) {
			log.Println("Pipeline runs not created yet")
		} else if prs.SUCCEEDED == "True" && prs.REASON == "Succeeded" && strings.HasPrefix(prs.NAME, prefix) {
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
	log.Printf("Validating service bindings,name: %s, namespace: %s", name, namespace)
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
	log.Printf("Validating training portals, name: %s", name)
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
	log.Printf("Validating 'Learning Center', name: %s, namespace: %s", name, namespace)
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

func VerifyBuildStatus(buildName string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Printf("Validating build status, name: %s, namespace: %s", buildName, namespace)
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		builds := kubectl_lib.GetBuilds(buildName, namespace)
		if len(builds) < 1 {
			log.Println("Builds are not generated yet")
		} else if builds[0].SUCCEEDED == "True" {
			log.Printf("Build %s status is verified successfully. Status is %s", builds[0].NAME, builds[0].SUCCEEDED)
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
	log.Printf("Get podintents status, name: %s, namespace: %s", name, namespace)
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
	log.Printf("Get workload status, name: %s, namespace: %s", name, namespace)
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

func VerifyKsvcStatus(name string, namespace string, latestReady string, timeoutInMins int, intervalInSeconds int) bool {
	log.Printf("Validating ksvc status  READY == true and LATESTREADY == %s, name: %s, namespace: %s", latestReady, name, namespace)
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		ksvc := kubectl_lib.GetKsvc(name, namespace)
		if len(ksvc) < 1 {
			log.Println("Knative services are not generated yet")
		} else if (ksvc[0].READY == "True") && (ksvc[0].LATESTREADY >= latestReady) {
			log.Printf("Knative %s status is verified successfully. Status is %s. LatestReady is %s.", ksvc[0].NAME, ksvc[0].READY, ksvc[0].LATESTREADY)
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
	log.Printf("Validating Image repository status, name: %s, namespace: %s", name, namespace)
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
	log.Printf("Validating Git repository status, name: %s, namespace: %s", name, namespace)
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

func VerifyTaskrunStatus(taskrunPrefix string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Printf("Validating task run status, taskrunPrefix:%s, namespace:%s", taskrunPrefix, namespace)
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		taskruns := kubectl_lib.GetTaskruns("", namespace)
		if len(taskruns) < 1 {
			log.Println("taskruns are not generated yet")
		} else {
			for _, taskrun := range taskruns {
				if taskrun.SUCCEEDED == "True" && strings.HasPrefix(taskrun.NAME, taskrunPrefix) {
					log.Printf("taskrun %s status is verified successfully, status is %s", taskrun.NAME, taskrun.SUCCEEDED)
					result = true
					return result
				}
			}
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

func VerifyTestTaskrunStatus(taskrunPrefix string, taskrunSuffix string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Printf("Validating task run status, taskrunPrefix: %s, taskrunSuffix: %s, namespace: %s", taskrunPrefix, taskrunSuffix, namespace)
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		taskruns := kubectl_lib.GetTaskruns("", namespace)
		if len(taskruns) < 1 {
			log.Println("taskruns are not generated yet")
		} else {
			for _, taskrun := range taskruns {
				if taskrun.SUCCEEDED == "True" && strings.HasPrefix(taskrun.NAME, taskrunPrefix) && strings.HasSuffix(taskrun.NAME, taskrunSuffix) {
					log.Printf("taskrun %s status is verified successfully, status is %s", taskrun.NAME, taskrun.SUCCEEDED)
					result = true
					return result
				}
			}
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
	log.Printf("Validating podintent status, name: %s, namespace: %s", name, namespace)
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
	log.Printf("Validating TAP installation status. pkgname: %s, namespace: %s", pkgName, namespace)
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

func ValidateLatestImageStatus(namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Printf("Validating image scans in namespace: %s", namespace)
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		status := GetLatestImageStatus(namespace)
		if status == "True" {
			log.Println("Latest image validated successfully")
			result = true
			break
		} else {
			log.Printf("Latest image status: %s", status)
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if !result {
		log.Printf("Latest image not successfull after %d mins", timeoutInMins)
	}
	return result
}

func ValidateDeliverables(name string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Printf("Validating deliverable %s in namespace %s", name, namespace)
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		deliverables := kubectl_lib.GetDeliverables(name, namespace)
		if len(deliverables) < 1 {
			log.Println("Deliverable is ready yet")
		} else if deliverables[0].READY == "True" && deliverables[0].REASON == "Ready" {
			log.Printf("Deliverable %s is ready", deliverables[0].NAME)
			result = true
			break
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if !result {
		log.Printf("Deliverable is not ready after %d mins", timeoutInMins)
	}
	return result
}

func ValidateBuildClusterDeliverableStatus(name string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Printf("Validating deliverable %s in namespace %s", name, namespace)
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		deliverables := kubectl_lib.GetDeliverables(name, namespace)
		if len(deliverables) < 1 {
			log.Println("Deliverable is not ready yet")
		} else if deliverables[0].READY == "False" && deliverables[0].REASON == "DeliveryNotFound" {
			log.Printf("Deliverable %s is %s", deliverables[0].NAME, deliverables[0].READY)
			result = true
			break
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if !result {
		log.Printf("Deliverable is not ready after %d mins", timeoutInMins)
	}
	return result
}

type Secrets struct {
	Name string `json:"name"`
}

func PatchServiceAccountWithNewSecret(sa string, namespace string, newSecret string) bool {

	log.Println("Patching the default service account")
	raw := kubectl_lib.GetServiceAccountJson(sa, namespace)
	if raw == nil {
		return false
	}

	var secret Secrets
	secret.Name = newSecret
	var secrets = append(raw.Secrets, secret)

	var secretPatch, err = json.Marshal(secrets)
	if err != nil {
		log.Printf("error unmarshaling: %v", err)
	}
	log.Printf("Patch to be added : %s", string(secretPatch))
	res := kubectl_lib.PatchServiceAccount(sa, namespace, "'{\"secrets\":"+string(secretPatch)+"}'")
	if !res {
		log.Println("Error while patching")
		return false
	}
	log.Println("Patch completed")
	return true
}

func ValidateWorkloadStatus(name string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Println("Validating workloads ")
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		status := GetWorkloadStatus(name, namespace)
		if status == "True" {
			log.Println("Workload validated successfully")
			result = true
			break
		} else {
			log.Printf("workload ready status: %s", status)
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if !result {
		log.Printf("Workload validation not successfull after %d mins", timeoutInMins)
	}
	return result
}

func GetLatestRevision(config_name string, namespace string, timeoutInMins int, intervalInSeconds int) string {
	log.Printf("Get revisions for config: %s in namespace: %s", config_name, namespace)

	time.Sleep(time.Duration(60) * time.Second)

	finalTimeout := (timeoutInMins - 1) * 60
	revisionName := ""
	for finalTimeout >= 0 {
		revs := kubectl_lib.GetRevisions("", namespace)
		for i := len(revs) - 1; i >= 0; i-- {
			if revs[i].CONFIG_NAME == config_name {
				revisionName = revs[i].NAME
				log.Printf("Latest revision is %s", revisionName)
				break
			}
		}
		if revisionName != "" {
			log.Printf("Found latest revision: %s", revisionName)
			break
		}
		log.Printf("%s not found, Waiting for %d seconds before retry", config_name, intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if revisionName == "" {
		log.Printf("Revision is not ready after %d mins", timeoutInMins)
	}

	return revisionName
}

func ValidateRevisionStatus(revision_name, config string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Printf("Validating revision %s having config %s in namespace %s", revision_name, config, namespace)
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		revision := kubectl_lib.GetRevisions(revision_name, namespace)
		if len(revision) < 1 {
			log.Printf("Revision: %s is not created yet", revision_name)
		} else if revision[0].READY == "True" && revision[0].CONFIG_NAME == config {
			log.Printf("Revision %s is ready", revision[0].NAME)
			result = true
			break
		} else {
			log.Printf("revision status: %s, revision config: %s", revision[0].READY, revision[0].CONFIG_NAME)
		}

		log.Printf("%s not found, Waiting for %d seconds before retry", config, intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if !result {
		log.Printf("Revision is not ready after %d mins", timeoutInMins)
	}
	return result
}

func VerifyNewerBuildStatus(oldBuildName string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Printf("Validating newer build status. oldBuildName: %s", oldBuildName)
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		builds := kubectl_lib.GetBuilds("", namespace)
		if len(builds) < 1 {
			log.Println("Builds are not generated yet")
		} else if builds[len(builds)-1].SUCCEEDED == "True" && builds[len(builds)-1].NAME > oldBuildName {
			log.Printf("Build %s status is verified successfully. Status is %s", builds[len(builds)-1].NAME, builds[len(builds)-1].SUCCEEDED)
			result = true
			break
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if !result {
		log.Printf("Newer build compared to %s is not ready after %d mins", oldBuildName, timeoutInMins)
	}
	return result
}

func VerifyNewerKsvcStatus(name string, namespace string, revision string, timeoutInMins int, intervalInSeconds int) bool {
	log.Printf("Validating ksvc status READY == true and revision >= : %s", revision)
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		ksvc := kubectl_lib.GetKsvc(name, namespace)
		if len(ksvc) < 1 {
			log.Println("Knative services are not generated yet")
		} else if (ksvc[len(ksvc)-1].READY == "True") && (ksvc[len(ksvc)-1].LATESTREADY >= revision) {
			log.Printf("Knative %s status is verified successfully. Status is %s. LatestCreated is %s.", ksvc[len(ksvc)-1].NAME, ksvc[len(ksvc)-1].READY, ksvc[len(ksvc)-1].LATESTCREATED)
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

func GetNewerRevision(old_revision_name, config_name string, namespace string, timeoutInMins int, intervalInSeconds int) string {
	log.Printf("Get revisions newer than %s for config: %s in namespace: %s", old_revision_name, config_name, namespace)

	finalTimeout := timeoutInMins * 60
	revisionName := ""
	for finalTimeout > 0 {
		rev := GetLatestRevision(config_name, namespace, 1, 30)
		if rev > old_revision_name {
			revisionName = rev
			break
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if revisionName == "" {
		log.Printf("No new revisions found for config: %s", config_name)
	}
	return revisionName
}

func GetImageDigest(imageName string, namespace string, timeoutInMins int, intervalInSeconds int) string {
	log.Printf("Fetch image digest for image %s in namespace: %s", imageName, namespace)
	finalTimeout := timeoutInMins * 60
	imageDigest := ""
	for finalTimeout > 0 {
		images := kubectl_lib.GetImages(imageName, namespace) // has to be in ready state
		if len(images) < 1 {
			log.Println("Images are not generated yet")
		} else if images[0].READY == "True" {
			imageDigest = strings.Split(images[0].LATESTIMAGE, "@")[1]
			log.Printf("imageDigests %s :", imageDigest)
			break
		} else {
			log.Printf("Image: %s, ready: %s", imageName, images[0].READY)
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	return imageDigest
}

func CheckDeploymentExists(name string, namespace string, timeoutInMins int, intervalInSeconds int) bool {
	log.Printf("Validating deployment exists, name: %s, namespace: %s", name, namespace)
	finalTimeout := timeoutInMins * 60
	result := false
	for finalTimeout > 0 {
		deployment := kubectl_lib.GetDeployments("", "")
		if len(deployment) <= 0 {
			log.Println("No deployments exist")
		} else {
			for _, dep := range deployment {
				if dep.NAME == name {
					result = true
					break
				}
			}
		}
		if result {
			log.Println("Deployment exists")
			break
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if !result {
		log.Printf("Deployment does not exist after %d mins", timeoutInMins)
	}
	return result
}

func GetCurrentClusterURL() string {
	log.Printf("Getting current cluster URL")
	output := kubectl_lib.GetCurrentConfigView()
	//log.Printf("output: %v", output)
	log.Printf("output_clusters: %v", output.Clusters[0].Cluster.Server)

	//log.Printf("output_clusters_server: %v", output.Clusters[0].Cluster.Server)
	return output.Clusters[0].Cluster.Server
}

func GetClusterToken(name string, namespace string) string {
	serviceAccount := kubectl_lib.GetServiceAccountJson(name, namespace)
	secretName := serviceAccount.Secrets[0].Name
	getSecrets := kubectl_lib.GetSecrets(secretName, namespace)
	clusterencodedToken := getSecrets.Data.Token
	decodedToken, err := base64.StdEncoding.DecodeString(clusterencodedToken)
	if err != nil {
		log.Printf("error while decoding token")
	}
	return string(decodedToken)
}

func GetServiceExternalIP(service string, namespace string, timeoutInMins int, intervalInSeconds int) string {
	log.Printf("Get external IP of service %s in namespace: %s", service, namespace)

	finalTimeout := timeoutInMins * 60
	externalIP := ""
	for finalTimeout > 0 {
		svc := kubectl_lib.GetServices(service, namespace)
		if len(svc) < 1 {
			log.Printf("%s service not yet created", service)
		} else if svc[0].EXTERNAL_IP == "" || svc[0].EXTERNAL_IP == "<none>" {
			log.Print("External IP is not found")
		} else {
			return svc[0].EXTERNAL_IP
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
	}
	if externalIP == "" {
		log.Printf("External IP not generated for service %s in namespace %s", service, namespace)
	}
	return externalIP
}
