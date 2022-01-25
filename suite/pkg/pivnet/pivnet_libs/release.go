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

func CreateRelease(productSlug string, releaseVersion string, releaseType string, eulaSlug string) *CreateReleaseOutput {
	var raw *CreateReleaseOutput
	cmd := fmt.Sprintf("./pivnet-cli create-release --product-slug %s --release-version %s --release-type %s --eula-slug %s --format json", productSlug, releaseVersion, eulaSlug, releaseType)
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return raw
	}
	res := strings.Split(response, "{")
	in := []byte(res[1])

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

func UpdateRelease(productSlug string, releaseVersion string, availability string) *UpdateReleaseOutput {
	// productSlug = "tanzu-application-platform"
	// releaseVersion = "1.0.1-build.test"
	cmd := fmt.Sprintf("./pivnet-cli update-release --product-slug=%s --release-version %s  --availability=%s --format json", productSlug, releaseVersion, availability)

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
