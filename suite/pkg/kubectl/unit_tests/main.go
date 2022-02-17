package main

import (
	kubectl_helper "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_helpers"
	kubectl_lib "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_libs"
)

func main() {

	// Testing kubectl helper methods

	kubectl_helper.ValidateAppLiveViewLabels("tanzu-java-web-app-git", "tap-install")
	kubectl_helper.ValidateAppLiveViewConventions("tanzu-java-web-app-git", "tap-install")
	kubectl_helper.ValidateSpringBootLabels("tanzu-java-web-app-git", "tap-install")
	kubectl_helper.ValidateSpringBootConventions("tanzu-java-web-app-git", "tap-install")
	kubectl_helper.ValidateDeveloperConventions("tanzu-java-web-app-git", "tap-install")
	kubectl_helper.GetLatestImageRepositoryStatus("spring-petclinic-delivery", "tap-install")
	kubectl_helper.GetLatestBuildStatus("spring-petclinic-build-1", "tap-install")
	kubectl_helper.GetLatestImageRepositoryStatus("", "tap-install")
	kubectl_helper.GetLatestBuildStatus("", "tap-install")
	kubectl_helper.GetLatestImageStatus("tap-install")
	kubectl_helper.GetKsvcStatus("spring-petclinic", "tap-install")
	kubectl_helper.ValidateImageScans("spring-petclinic", "tap-install", 5, 30)
	kubectl_helper.ValidateSourceScans("spring-petclinic", "tap-install", 5, 30)
	kubectl_helper.ValidatePipelineExists("developer-defined-tekton-pipeline", "tap-install", 5, 30)
	kubectl_helper.ValidatePipelineRuns("spring-petclinic-s8jxd", "tap-install", 5, 30)
	kubectl_helper.ValidateServiceBindings("spring-petclinic-spring-petclinic-db", "tap-install", 5, 30)
	kubectl_helper.ValidateTrainingPortalStatus("learning-center-guided", 5, 30)
	kubectl_helper.ValidateLearningCenter("learningcenter-portal", "learning-center-guided-ui")
	kubectl_helper.VerifyImageRepositoryStatus("spring-petclinic-delivery", "tap-install", 15, 30)
	kubectl_helper.VerifyKsvcStatus("spring-petclinic", "tap-install", 5, 30)

	/// Testing kubectl helper methods
	kubectl_lib.GetPodintent("tanzu-java-web-app-git", "tap-install")
	kubectl_lib.GetWorkload("tanzu-java-web-app-git", "tap-install")
	kubectl_lib.GetImageRepositories("tanzu-java-web-app-git-delivery", "tap-install")
	kubectl_lib.GetBuilds("tanzu-java-web-app-git-build-1", "tap-install")
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
