package pivnet_helpers

import (
	"log"
	"strings"
	"time"

	pivnet_libs "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/pivnet/pivnet_libs"
)

func WaitTillArtifactReferenceIsReady(productSlug string, artifactReferenceId int, timeoutInMins int, intervalInSeconds int) bool {
	log.Println("Validating artifacts creation status")
	finalTimeout := timeoutInMins * 60
	result := false
	timeSpent := 0
	for finalTimeout > 0 {
		artifact_Details := pivnet_libs.GetArtifactReference(productSlug, artifactReferenceId)
		if artifact_Details.ReplicationStatus == "complete" {
			log.Printf("Artifact created. Total time taken: %d mins %d seconds", timeSpent/60, timeSpent%60)
			result = true
			break
		}
		log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		finalTimeout -= intervalInSeconds
		timeSpent += intervalInSeconds
	}
	if !result {
		log.Printf("Artifacts are not generated after %d mins", timeoutInMins)
	}
	return result
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
