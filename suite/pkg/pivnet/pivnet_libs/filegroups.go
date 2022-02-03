package pivnet_libs

import (
	"encoding/json"
	"fmt"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

type ListFileGroupsOutput []struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Product struct {
	} `json:"product"`
	ProductFiles []struct {
		ID           int    `json:"id"`
		AwsObjectKey string `json:"aws_object_key"`
		FileType     string `json:"file_type"`
		FileVersion  string `json:"file_version"`
		Sha256       string `json:"sha256"`
		Name         string `json:"name"`
		Links        struct {
			Download struct {
				Href string `json:"href"`
			} `json:"download"`
		} `json:"_links"`
	} `json:"product_files"`
}

func ListFileGroups(productSlug string, releaseVersion string) ListFileGroupsOutput {
	fmt.Println("Executing ListFileGroups")
	var raw ListFileGroupsOutput
	cmd := fmt.Sprintf("pivnet-cli file-groups --product-slug %s --release-version %s --format json", productSlug, releaseVersion)
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

func AddFileGroup(fileGroupId int, productSlug string, releaseVersion string) bool {
	fmt.Println("Executing AddFileGroup")
	cmd := fmt.Sprintf("pivnet-cli add-file-group --file-group-id %d --product-slug %s  --release-version %s", fileGroupId, productSlug, releaseVersion)
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil && response != "" {
		return false
	}
	return true
}
