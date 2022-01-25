package pivnet_libs

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

type CreateReleaseOutput struct {
	ID           int    `json:"id"`
	Availability string `json:"availability"`
	Eula         struct {
		Slug  string `json:"slug"`
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Links struct {
		} `json:"_links"`
	} `json:"eula"`
	ReleaseDate string `json:"release_date"`
	ReleaseType string `json:"release_type"`
	Version     string `json:"version"`
	Links       struct {
		ProductFiles struct {
			Href string `json:"href"`
		} `json:"product_files"`
		EulaAcceptance struct {
			Href string `json:"href"`
		} `json:"eula_acceptance"`
	} `json:"_links"`
	UpdatedAt              time.Time `json:"updated_at"`
	SoftwareFilesUpdatedAt time.Time `json:"software_files_updated_at"`
}

func CreateRelease(productSlug string, releaseVersion string, eulaSlug string, releaseType string) *CreateReleaseOutput {
	// productSlug = "tanzu-application-platform"
	// releaseVersion = "1.0.1-build.test"
	// eulaSlug = "vmware-prerelease-eula"
	// releaseType = "Beta Release"
	cmd := fmt.Sprintf("./pivnet-cli create-release --product-slug %s --release-version %s --eula-slug %s --release-type %s --format json", productSlug, releaseVersion, eulaSlug, releaseType)

	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Println("something bad happened")
	}
	res := strings.Split(response, "{")
	in := []byte(res[1])
	var raw *CreateReleaseOutput
	if err := json.Unmarshal(in, &raw); err != nil {
		panic(err)
	}
	return raw
}

type UpdateReleaseOutput struct {
	ID           int    `json:"id"`
	Availability string `json:"availability"`
	Eula         struct {
		Slug  string `json:"slug"`
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Links struct {
		} `json:"_links"`
	} `json:"eula"`
	ReleaseDate string `json:"release_date"`
	ReleaseType string `json:"release_type"`
	Version     string `json:"version"`
	Links       struct {
		ProductFiles struct {
			Href string `json:"href"`
		} `json:"product_files"`
		EulaAcceptance struct {
			Href string `json:"href"`
		} `json:"eula_acceptance"`
	} `json:"_links"`
	UpdatedAt              time.Time `json:"updated_at"`
	SoftwareFilesUpdatedAt time.Time `json:"software_files_updated_at"`
}

func UpdateRelease(productSlug string, releaseVersion string) *UpdateReleaseOutput {
	// productSlug = "tanzu-application-platform"
	// releaseVersion = "1.0.1-build.test"
	cmd := fmt.Sprintf("./pivnet-cli update-release --availability=selected-user-groups  --product-slug=%s --release-version %s --format json", productSlug, releaseVersion)

	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Println("something bad happened")
	}
	res := strings.Split(response, "{")
	in := []byte(res[1])
	var raw *UpdateReleaseOutput
	if err := json.Unmarshal(in, &raw); err != nil {
		panic(err)
	}
	return raw
}
