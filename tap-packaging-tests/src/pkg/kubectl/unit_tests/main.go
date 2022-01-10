package main

import (
	kubectl_helper "pkg/kubectl/kubectl_helper"
	kubectl_lib "pkg/kubectl/kubectl_lib"
)

func main() {

	// Examples of kubectl linux cmd ouput parsing
	kubectl_helper.GetLatestImageStatus("tap-install")
	kubectl_helper.ValidateLearningCenter("learningcenter-portal", "learning-center-guided-ui")

	// Examples of kubectl json output parsing
	kubectl_helper.ValidateAppLiveViewLabels("tanzu-java-web-app-git", "tap-install")

	// Testing kubectl lib methods
	kubectl_lib.GetPodintent("tanzu-java-web-app-git", "tap-install")
	kubectl_lib.GetWorkload("tanzu-java-web-app-git", "tap-install")
	kubectl_lib.GetImageRepositories("tanzu-java-web-app-git", "tap-install")
	kubectl_lib.GetBuilds("tanzu-java-web-app-git", "tap-install")
	kubectl_lib.GetLatestImage("tap-install")
	kubectl_lib.GetKsvc("tanzu-java-web-app-git", "tap-install")
	kubectl_lib.GetSourceScan("tanzu-java-web-app-git", "tap-install")

}
