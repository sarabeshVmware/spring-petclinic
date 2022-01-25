package pivnet_helpers

import (
	"log"
	"time"

	pivnet_libs "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/pivnet/pivnet_libs"
)

func WaitTillArtifactReferenceIsReady(productSlug string, artifactReferenceId string) bool {
	artifact_Details := pivnet_libs.GetArtifactReference("tanzu-application-platform", "27548")
	if artifact_Details.ReplicationStatus != "in_progress" {
		log.Println("Artifact created")
		return true
	} else {
		time.Sleep(300)
		return false
	}

}
