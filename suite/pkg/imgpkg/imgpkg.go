package imgpkg

import (
	"fmt"
	"log"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

func ImgpkgCopy(sourceBundle string, targetRepo string) error {
	log.Printf("copying images from %s to %s", sourceBundle, targetRepo)

	// execute cmd
	cmd := fmt.Sprintf("imgpkg copy -b %s --to-repo %s", sourceBundle, targetRepo)
	output, err := linux_util.ExecuteCmdInBashMode(cmd)
	if err != nil {
		log.Printf("error while copying images from %s to %s", sourceBundle, targetRepo)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("copied images from %s to %s", sourceBundle, targetRepo)
		log.Printf("output: %s", output)
	}

	return err
}
