package pivnet_helpers

import (
	"log"
	"strings"
	"time"

	pivnet_libs "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/pivnet/pivnet_libs"
)

func WaitTillArtifactReferenceIsReady(productSlug string, artifactReferenceId int) bool {
	count := 40
	for count <= 40 {
		if count == 0 {
			log.Fatalf("Artifacts are not generated after 20 mins")
			return false
		}
		artifact_Details := pivnet_libs.GetArtifactReference(productSlug, artifactReferenceId)
		if artifact_Details.ReplicationStatus == "complete" {
			log.Println("Artifact created")
			return true
		}
		log.Printf("Waiting for 30s for artifacts getting generated ...")
		time.Sleep(30 * time.Second)
		count -= 1
	}
	return false
}

func GetLatestRelease(productSlug string, versionPrefix string) string {
	releases := pivnet_libs.ListReleases(productSlug, 10)
	latestRel := ""
	for _, rel := range releases {
		if strings.HasPrefix(rel.Version, versionPrefix) {
			latestRel = rel.Version
			log.Printf("Latest release: %s", rel.Version)
			break
		}
	}
	return latestRel
}
