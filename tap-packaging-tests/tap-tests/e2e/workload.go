// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"fmt"
	"log"

	tap "gitlab.eng.vmware.com/tap/tap-packages/tap-packaging-tests/pkg"
)

func DeleteWorkload(workload string, namespace string) {
	log.Printf("Deleting workload: %s", workload)
	tap.Run(fmt.Sprintf("tanzu apps workload delete %s -n %s -y", workload, namespace))
}

func CreateWorkload(workload string, server string, repository string, sourceImage string, localPath string, namespace string) {
	log.Printf("Creating workload: %s", workload)
	tap.Run(fmt.Sprintf("tanzu apps workload create %s --source-image %s/%s/%s --type web --yes --local-path=%s -n %s", workload, server, repository, sourceImage, localPath, namespace))
}

func UpdateWorkload(workload string, server string, repository string, sourceImage string, localPath string, namespace string) {
	log.Printf("Updating workload: %s", workload)
	tap.Run(fmt.Sprintf("tanzu apps workload update %s --source-image %s/%s/%s --type web --yes --local-path=%s -n %s --live-update=true", workload, server, repository, sourceImage, localPath, namespace))
}
