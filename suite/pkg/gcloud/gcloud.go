package gcloud

import (
	"fmt"
	"log"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

func DeleteImageContainer(image string) error {
	log.Printf("deleting image container %s", image)

	// execute cmd
	cmd := fmt.Sprintf("gcloud container images delete %s", image)
	output, err := linux_util.ExecuteCmdInBashMode(cmd)
	if err != nil {
		log.Printf("error while deleting %s", image)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("error while deleting %s", image)
		log.Printf("output: %s", output)
	}

	return err
}
