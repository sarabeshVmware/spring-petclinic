package main

import (
	kubectl_helper "pkg/kubectl/kubectl_helper"
	kubectl_lib "pkg/kubectl/kubectl_lib"
)

func main() {

	// Testing kubectl helper methods

	kubectl_helper.ValidateAppLiveViewLabels("tanzu-java-web-app-git", "tap-install")
	kubectl_helper.ValidateAppLiveViewConventions("tanzu-java-web-app-git", "tap-install")
	kubectl_helper.ValidateSpringBootLabels("tanzu-java-web-app-git", "tap-install")
	kubectl_helper.ValidateSpringBootConventions("tanzu-java-web-app-git", "tap-install")
	kubectl_helper.ValidateDeveloperConventions("tanzu-java-web-app-git", "tap-install")
	kubectl_helper.GetImageRepositoryStatus("spring-petclinic-delivery", "tap-install")
	kubectl_helper.GetBuildStatus("spring-petclinic-build-1", "tap-instal")
	kubectl_helper.GetLatestImageStatus("tap-install")
	kubectl_helper.GetKsvcStatus("spring-petclinic", "tap-install")
	kubectl_helper.ValidateImageScans("spring-petclinic", "tap-install")
	kubectl_helper.ValidateSourceScans("spring-petclinic", "tap-install")
	kubectl_helper.ValidatePipelineExists("developer-defined-tekton-pipeline", "tap-install")
	kubectl_helper.ValidatePipelineRuns("spring-petclinic-s8jxd", "tap-install")
	kubectl_helper.ValidateServiceBindings("spring-petclinic-spring-petclinic-db", "tap-install")
	kubectl_helper.ValidateTrainingPortalStatus("learning-center-guided")
	kubectl_helper.ValidateLearningCenter("learningcenter-portal", "learning-center-guided-ui")

	/// Testing kubectl helper methods
	kubectl_lib.GetPodintent("tanzu-java-web-app-git", "tap-install")
	kubectl_lib.GetWorkload("tanzu-java-web-app-git", "tap-install")
	kubectl_lib.GetImageRepositories("tanzu-java-web-app-git", "tap-install")
	kubectl_lib.GetBuilds("tanzu-java-web-app-git", "tap-install")
	kubectl_lib.GetLatestImage("tap-install")
	kubectl_lib.GetKsvc("tanzu-java-web-app-git", "tap-install")
	kubectl_lib.GetSourceScan("tanzu-java-web-app-git", "tap-install")
	kubectl_lib.GetPipeline("developer-defined-tekton-pipeline", "tap-install")
	kubectl_lib.GetPipelineRuns("spring-petclinic-s8jxd", "tap-install")
	kubectl_lib.GetImageScan("spring-petclinic", "tap-install")
	kubectl_lib.GetServiceBindings("spring-petclinic-spring-petclinic-db", "tap-install")
	kubectl_lib.GetTrainingPortals("learning-center-guided")
	kubectl_lib.GetIngress("learningcenter-portal", "learning-center-guided-ui")
	kubectl_lib.GetPodintentJson("tanzu-java-web-app-git", "tap-install")
	kubectl_lib.GetImageRepositoriesJson("spring-petclinic-delivery", "tap-install")
	kubectl_lib.GetRunnablesJson("tanzu-java-web-app-git", "tap-install")
}
