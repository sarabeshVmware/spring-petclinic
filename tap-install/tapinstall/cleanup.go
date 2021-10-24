// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tapinstall

import (
	"log"

	tap "gitlab.eng.vmware.com/tap/tap-packaging-tests/pkg"
)

func Cleanup(configFile string, valuesDir string) {
	config := GetConfig(configFile, valuesDir)
	log.Printf("Request: Cleanup")
	tap.UninstallPackages(config.Namespace)
	tap.DeletePackageRepository(config.Namespace)
	tap.DeleteImagepullSecrets(config.Namespace)
	tap.DeleteNamespace(config.Namespace)
}
