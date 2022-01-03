package main

import (
	"fmt"
	kubectl_helper "pkg/kubectl/kubectl_helper"
)

func main() {

	// Examples of kubectl linux cmd ouput parsing

	a := kubectl_helper.GetLatestImageStatus("tap-install")
	fmt.Printf("%+v\n", a)

	// Examples of kubectl json output parsing

	f := kubectl_helper.ValidateAppLiveViewLabels()
	fmt.Println(f)

}
