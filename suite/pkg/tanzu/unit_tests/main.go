package main

import (
	tanzu_lib "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzu_libs"

	tanzu_helper "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzu_helpers"
)

func main() {

	// Testing tanzu helper methods
	tanzu_helper.IsGrypeInstalled("tap-install")
	tanzu_helper.IsScanningInstalled("tap-install")
	tanzu_helper.ValidateInstalledPackageStatus("accelerator", "tap-install", 5, 30)

	// Testing tanzu lib methods
	tanzu_lib.ListInstalledPackages("tap-install")
	tanzu_lib.GetInstalledPackages("accelerator", "tap-install")
}
