// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"log"
	"time"

	tap "gitlab.eng.vmware.com/tap/tap-packages/tap-packaging-tests/pkg"
)

func ListAccelerators() {
	log.Printf("Listing accelerators:")
	_, err := tap.Run_AllowError("tanzu accelerator list")
	if err != nil {
		time.Sleep(5 * time.Second)
		ListAccelerators()
	}
}

func GenerateAcceleratorProject(label string, projectName string, repositoryPrefix string, unzip bool, serverIP string) {
	log.Printf("Generating accelerator project: %s", label)
	tap.RunWithBash(fmt.Sprintf(`tanzu accelerator generate %s --options '{"projectName":"%s", "repositoryPrefix":"%s", "includeKubernetes": true}' --server-url http://%s`, label, projectName, repositoryPrefix, serverIP))
	if unzip {
		out, err := tap.Run_AllowError(fmt.Sprintf("ls -lt %s", projectName))
		if err == nil {
			tap.Run(fmt.Sprintf("rm -r %s", projectName))
		}
		log.Printf("Output: \n%s", string(out))
		tap.Run(fmt.Sprintf("unzip %s.zip", projectName))
	}
}
