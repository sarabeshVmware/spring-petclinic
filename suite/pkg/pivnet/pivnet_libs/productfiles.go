package pivnet_libs

import (
	"encoding/json"
	"fmt"
	"strings"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

type ListProductFilesOutput []struct {
	ID           int    `json:"id"`
	AwsObjectKey string `json:"aws_object_key"`
	FileType     string `json:"file_type"`
	FileVersion  string `json:"file_version"`
	Sha256       string `json:"sha256,omitempty"`
	Name         string `json:"name"`
	Links        struct {
		Download struct {
			Href string `json:"href"`
		} `json:"download"`
	} `json:"_links"`
}

func ListProductFiles(productSlug string, releaseVersion string) ListProductFilesOutput {
	fmt.Println("Executing ListProductFiles")
	var raw ListProductFilesOutput
	cmd := fmt.Sprintf("pivnet-cli product-files --product-slug %s --limit %s --format json", productSlug, releaseVersion)
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

func AddProductFile(productFileId int, productSlug string, releaseVersion string) bool {
	fmt.Println("Executing AddProductFile")
	cmd := fmt.Sprintf("pivnet-cli add-product-file --product-file-id %d --product-slug %s  --release-version %s", productFileId, productSlug, releaseVersion)
	response, err := linux_util.ExecuteCmd(cmd)
	if strings.Contains(response, "Release already contains this product file") {
		return true
	}
	if err != nil && response != "" {
		return false
	}
	return true
}
