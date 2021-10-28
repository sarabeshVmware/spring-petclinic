// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"log"
	"time"

	tap "gitlab.eng.vmware.com/tap/tap-packaging-tests/pkg"
)

func ListAccelerators() {
	log.Printf("Listing accelerators:")
	_, err := tap.Run_AllowError("tanzu accelerator list")
	if err != nil {
		time.Sleep(5 * time.Second)
		ListAccelerators()
	}
}

func GenerateAcceleratorProject(label string, projectName string, repositoryPrefix string, unzip bool) {
	log.Printf("Generating accelerator project: %s", label)
	tap.RunWithBash(fmt.Sprintf(`tanzu accelerator generate %s --options '{"projectName":"%s", "repositoryPrefix":"%s", "includeKubernetes": true}'`, label, projectName, repositoryPrefix))
	if unzip {
		tap.Run(fmt.Sprintf("rm -r %s", projectName))
		tap.Run(fmt.Sprintf("unzip %s.zip", projectName))
	}
}
