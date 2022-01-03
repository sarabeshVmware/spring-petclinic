package main

import (
	"log"
	kubectl_helper "pkg/kubectl/kubectl_helper"
	"strings"
)

func main() {
	kubectl_helper.ValidateAppLiveViewLabels()
}
