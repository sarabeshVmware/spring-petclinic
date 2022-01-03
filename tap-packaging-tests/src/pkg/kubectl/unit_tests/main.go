package main

import (
	"fmt"
	kubectl_helper "pkg/kubectl/kubectl_helper"
)

func main() {
	kubectl_helper.ValidateAppLiveViewLabels()
	kubectl_helper.GetImageRepositoryStatus("tanzu-java-web-app", "tap-install")

	// Examples of kubectl linux cmd ouput parsing

	a := kubectl_helper.GetLatestImage("tap-install")
	fmt.Printf("%+v\n", a)

	b := kubectl_helper.GetPodintent("tanzu-java-web-app", "tap-install")
	fmt.Printf("podIntents: %+v\n", b)

	c := kubectl_helper.GetWorkload("tanzu-java-web-app", "tap-install")
	fmt.Printf("%+v\n", c)

	d := kubectl_helper.GetImageRepositories("tanzu-java-web-app", "tap-install")
	fmt.Printf("%+v\n", d)

	e := kubectl_helper.GetBuilds("tanzu-java-web-app-build-1", "tap-install")
	fmt.Printf("%+v\n", e)

	// Examples of kubectl json output parsing

	f := kubectl_helper.ValidateAppLiveViewLabels()
	fmt.Println(f)

}
