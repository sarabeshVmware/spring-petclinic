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
	name = "1.0.1-build.test"
	productSlug = "tanzu-application-platform"
	artifactPath = "tap-packages:1.0.1-build.ci.24-01-2022-09-06-31"
	digest = "sha256:66424580e6d86d77eea90ccf7aab7659bbc1880732fdadc062f14e64178b3845"
	cmd := fmt.Sprintf("./pivnet-cli create-artifact-reference --name %s --product-slug=%s--artifact-path=%s --digest=%s --format json", name, productSlug, artifactPath, digest)
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

func GetArtifactReference(productSlug string, artifactReferenceId string) *GetArtifactReferenceOutput {
	productSlug = "tanzu-application-platform"
	artifactReferenceId = "27502"
	cmd := fmt.Sprintf("./pivnet-cli artifact-reference --product-slug=t%s --artifact-reference-id %s --format json", productSlug, artifactReferenceId)
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

type AddArtifactReferenceOutput struct {
	ID                int    `json:"id"`
	ArtifactPath      string `json:"artifact_path"`
	Digest            string `json:"digest"`
	Name              string `json:"name"`
	ReplicationStatus string `json:"replication_status"`
}

func AddArtifactReference(productSlug string, releaseVersion string, artifactReferenceId string) *AddArtifactReferenceOutput {
	cmd := fmt.Sprintf("./pivnet-cli add-artifact-reference --product-slug=%s --release-version %s --artifact-reference-id=%s --format json", productSlug, releaseVersion, artifactReferenceId)
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Println("something bad happened")
	}
	in := []byte(response)
	var raw *AddArtifactReferenceOutput
	if err := json.Unmarshal(in, &raw); err != nil {
		panic(err)
	}
	return raw
}
