package tanzu_libs

// Usage:
//   tanzu insight image [command]

// Aliases:
//   image, images, imgs, img

// Available Commands:
//   add             Add an image report
//   get             Get image by digest
//   packages        Get image packages
//   vulnerabilities Get image vulnerabilities

// Flags:
//   -h, --help   help for image

import (
	"fmt"
	"log"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

func ListInsightImagesVulnerabilities(imageDigest string) (string, error) {

	cmd := fmt.Sprintf("tanzu insight images vulnerabilities --digest %s --format text", imageDigest)
	output, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while getting vulnerabilities for %s", imageDigest)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("vulnerabilities output for image %s:", imageDigest)
	}
	return output, err
}

func GetInsightImages(imageDigest string) (string, error) {

	cmd := fmt.Sprintf("tanzu insight images get --digest %s --format text", imageDigest)
	output, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while getting vulnerabilities for %s", imageDigest)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("vulnerabilities output for image %s:", imageDigest)
	}
	return output, err
}
