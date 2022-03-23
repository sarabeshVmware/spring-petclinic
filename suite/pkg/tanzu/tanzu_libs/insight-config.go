package tanzu_libs

// Usage:
//   tanzu insight config [command]

// Available Commands:
//   set-target  Set metadata store endpoint

// Flags:
//   -h, --help   help for config

import (
	"fmt"
	"log"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

func TanzuConfigureInsight(caFilePath string, accessToken string) error {
	log.Print("setting insight cli config")

	//configuring tanzu
	cmd := fmt.Sprintf("tanzu insight config set-target https://metadata-store-app.metadata-store.svc.cluster.local:8443 --ca-cert %s --access-token %s", caFilePath, accessToken)
	output, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while configuring insight ")
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Print("insight configured ")
		return err
	}
	return err
}
