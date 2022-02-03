package pivnet_libs

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

type CreateArtifactReferenceOutput struct {
	ID                int    `json:"id"`
	ArtifactPath      string `json:"artifact_path"`
	Digest            string `json:"digest"`
	Name              string `json:"name"`
	ReplicationStatus string `json:"replication_status"`
}

func CreateArtifactReference(name string, productSlug string, artifactPath string, digest string) *CreateArtifactReferenceOutput {
	log.Println("Executing CreateArtifactReference")
	var raw *CreateArtifactReferenceOutput
	cmd := fmt.Sprintf("pivnet-cli create-artifact-reference --name %s --product-slug=%s --artifact-path=%s --digest=%s --format json", name, productSlug, artifactPath, digest)
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return raw
	}
	in := []byte(response)

	if err := json.Unmarshal(in, &raw); err != nil {
		panic(err)
	}
	return raw
}

type GetArtifactReferenceOutput struct {
	ID                int    `json:"id"`
	ArtifactPath      string `json:"artifact_path"`
	Digest            string `json:"digest"`
	Name              string `json:"name"`
	ReplicationStatus string `json:"replication_status"`
}

func GetArtifactReference(productSlug string, artifactReferenceId int) *GetArtifactReferenceOutput {
	log.Println("Executing GetArtifactReference")
	var raw *GetArtifactReferenceOutput
	cmd := fmt.Sprintf("pivnet-cli artifact-reference --product-slug=%s --artifact-reference-id %d --format json", productSlug, artifactReferenceId)
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return raw
	}
	in := []byte(response)
	if err := json.Unmarshal(in, &raw); err != nil {
		panic(err)
	}
	return raw
}

func AddArtifactReference(productSlug string, releaseVersion string, artifactReferenceId int) bool {
	log.Println("Executing AddArtifactReference")
	count := 5
	for count <= 5 {
		if count == 0 {
			log.Println("Unable to add artifacts after 5 attempts")
			return false
		}
		cmd := fmt.Sprintf("pivnet-cli add-artifact-reference --product-slug=%s --release-version %s --artifact-reference-id=%d --format json", productSlug, releaseVersion, artifactReferenceId)
		response, err := linux_util.ExecuteCmd(cmd)
		if err == nil && response == "" {
			return true
		}
		log.Printf("Waiting for 30s to retry")
		time.Sleep(30 * time.Second)
		count -= 1
	}
	return false
}

type ListArtifactReferencesOutput []struct {
	ID              int      `json:"id"`
	ArtifactPath    string   `json:"artifact_path"`
	Digest          string   `json:"digest"`
	Name            string   `json:"name"`
	ReleaseVersions []string `json:"release_versions"`
}

func ListArtifactReferences(productSlug string, releaseVersion string, digest string) ListArtifactReferencesOutput {
	log.Println("Executing ListArtifactReferences")
	var raw ListArtifactReferencesOutput
	cmd := fmt.Sprintf("pivnet-cli artifact-references --product-slug=%s", productSlug)
	if releaseVersion != "" {
		cmd += fmt.Sprintf(" --release-version %s", releaseVersion)
	}
	if digest != "" {
		cmd += fmt.Sprintf(" --digest %s", digest)
	}
	cmd += " --format json"
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return raw
	}
	in := []byte(response)
	if err := json.Unmarshal(in, &raw); err != nil {
		panic(err)
	}
	return raw
}
