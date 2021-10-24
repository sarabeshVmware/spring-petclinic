// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tapinstall

import (
	"log"

	tap "gitlab.eng.vmware.com/tap/tap-packaging-tests/pkg"
)

func Install(configFile string, valuesDir string, preCleanup bool, postCleanup bool) {
	log.Printf("Request: Install packages")
	config := GetConfig(configFile, valuesDir)

	log.Printf("Request: Cleanup pre-installation (%t)", preCleanup)
	if preCleanup {
		Cleanup(configFile, valuesDir)
	}

	tap.CreateNamespace(config.Namespace)
	tap.CreateImagepullSecrets(config.Secrets, config.Namespace)
	tap.AddPackageRepository(config.PackageRepository, config.Namespace)
	tap.CheckPackageRepositoryStatus(config.PackageRepository, config.Namespace)
	tap.ListPackages(config.Namespace)
	// tap.ListValuesSchema(config.Packages, config.Namespace)
	tap.InstallPackages(config.Packages, config.Namespace, config.ValuesDirectory)
	tap.SetupDeveloperNamespacePostInstallation(config.Namespace)

	log.Printf("Request: Cleanup post-installation (%t)", postCleanup)
	if postCleanup {
		Cleanup(configFile, valuesDir)
	}
}
