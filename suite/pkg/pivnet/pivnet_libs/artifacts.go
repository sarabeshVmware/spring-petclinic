package pivnet_libs

import (
	"encoding/json"
	"fmt"
	"log"

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
	cmd := fmt.Sprintf("pivnet-cli create-artifact-reference --name %s --product-slug=%s --artifact-path=%s --digest=%s --format json", name, productSlug, artifactPath, digest)
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Println("something bad happened")
	}
	in := []byte(response)
	var raw *CreateArtifactReferenceOutput
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
	cmd := fmt.Sprintf("pivnet-cli artifact-reference --product-slug=%s --artifact-reference-id %d --format json", productSlug, artifactReferenceId)
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Println("something bad happened")
	}
	in := []byte(response)
	var raw *GetArtifactReferenceOutput
	if err := json.Unmarshal(in, &raw); err != nil {
		panic(err)
	}
	return raw
}

func AddArtifactReference(productSlug string, releaseVersion string, artifactReferenceId int) {
	log.Println("Executing AddArtifactReference")
	cmd := fmt.Sprintf("pivnet-cli add-artifact-reference --product-slug=%s --release-version %s --artifact-reference-id=%d --format json", productSlug, releaseVersion, artifactReferenceId)
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil && response != "" {
		log.Println("something bad happened")
	}
}
