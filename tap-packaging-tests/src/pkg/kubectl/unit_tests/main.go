package main

import (
	kubectl_helper "pkg/kubectl/kubectl_helper"
)

func main() {
	kubectl_helper.ValidateAppLiveViewLabels()
	kubectl_helper.GetImageRepositoryStatus("tanzu-java-web-app", "tap-install")
}
