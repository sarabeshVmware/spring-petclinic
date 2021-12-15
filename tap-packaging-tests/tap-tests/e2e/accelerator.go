// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"log"
	"time"

	tap "gitlab.eng.vmware.com/tap/tap-packages/tap-packaging-tests/pkg"
)

func ListAccelerators() []byte {
	log.Printf("Listing accelerators:")
	out, _ := tap.Run("tanzu accelerator list")
	return out
}

func CheckAccelerators() {
	count := 5
	for count >= 0 {
		if count == 0 {
			log.Printf("No accelerators found after 150 seconds ")
			break
		}
		accout := ListAccelerators()
		strout := string(accout)
		if strout == "No accelerators found." {
			log.Println("No accelerators found waiting for 30 secs. Output is :", strout)
			time.Sleep(30 * time.Second)
		} else {
			log.Println("Accelerators found. Output is :\n", strout)
			break
		}
	}
}
func GenerateAcceleratorProject(label string, projectName string, repositoryPrefix string, unzip bool, serverIP string) {
	log.Printf("Generating accelerator project: %s", label)
	tap.RunWithBash(fmt.Sprintf(`tanzu accelerator generate %s --options '{"projectName":"%s", "repositoryPrefix":"%s", "includeKubernetes": true}' --server-url http://%s`, label, projectName, repositoryPrefix, serverIP))
	if unzip {
		out, err := tap.Run_AllowError(fmt.Sprintf("ls -lt %s", projectName))
		log.Printf("Output: \n%s", string(out))
		if err == nil {
			tap.Run(fmt.Sprintf("rm -r %s", projectName))
		}
		tap.Run(fmt.Sprintf("unzip %s.zip", projectName))
	}
}
