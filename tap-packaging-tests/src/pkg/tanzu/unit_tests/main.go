package main

import (
	tanzu_helper "pkg/tanzu/tanzu_helpers"
	tanzu_lib "pkg/tanzu/tanzu_libs"
)

func main() {

	// Testing tanzu helper methods
	tanzu_helper.IsGrypeInstalled("tap-install")
	tanzu_helper.IsScanningInstalled("tap-install")
	tanzu_helper.ValidateInstalledPackageStatus("accelerator", "tap-install", "")

	// Testing tanzu lib methods
	tanzu_lib.ListInstalledPackages("tap-install")
	tanzu_lib.GetInstalledPackages("accelerator", "tap-install")
}
