// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tapinstall

import (
	"log"

	tap "gitlab.eng.vmware.com/tap/tap-packaging-tests/pkg"
)

func Cleanup(configFile string, valuesDir string) {
	input := GetInput(configFile, valuesDir)
	log.Printf("Request: Cleanup")
	tap.UninstallPackages(input.Namespace)
	tap.DeletePackageRepository(input.Namespace)
	tap.DeleteImagepullSecrets(input.Namespace)
	tap.DeleteNamespace(input.Namespace)
}
