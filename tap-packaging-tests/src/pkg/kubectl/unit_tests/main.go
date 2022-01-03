package main

import (
	"log"
	kubectl_helper "pkg/kubectl/kubectl_helper"
	"strings"
)

func main() {
	kubectl_helper.ValidateAppLiveViewLabels()
	kubectl_helper.GetImageRepositoryStatus("tanzu-java-web-app", "tap-install")
}
