package main

import (
	pivnet_helpers "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/pivnet/pivnet_helpers"
	pivnet_libs "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/pivnet/pivnet_libs"
)

func main() {

	pivnet_libs.Login("host", "token")
	pivnet_libs.CreateRelease("tanzu-application-platform", "1.0.1-build.test", "vmware-prerelease-eula", "Beta Release")
	pivnet_libs.CreateArtifactReference("1.0.1-build.test", "tanzu-application-platform", "tap-packages:1.0.1-build.ci.24-01-2022-09-06-31", "sha256:66424580e6d86d77eea90ccf7aab7659bbc1880732fdadc062f14e64178b3845")
	pivnet_libs.GetArtifactReference("tanzu-application-platform", "27548")
	pivnet_helpers.WaitTillArtifactReferenceIsReady("tanzu-application-platform", "27548")
	pivnet_libs.AddArtifactReference("tanzu-application-platform", "1.0.1-build.test1", "27548")
	pivnet_libs.UpdateRelease("tanzu-application-platform ", "1.0.1-build.test")
	pivnet_libs.ListUserGroups()
	pivnet_libs.AddUserGroup("tanzu-application-platform", "1.0.1-build.test", "437")

}
