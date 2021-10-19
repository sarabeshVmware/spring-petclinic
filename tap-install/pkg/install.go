// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"log"

	tap "gitlab.eng.vmware.com/tap/tap-packaging-tests/pkg"
)

func Install(preCleanup bool, postCleanup bool) {
	log.Printf("Request: Install packages")
	input := GetInput()

	log.Printf("Request: Cleanup pre-installation (%t)", preCleanup)
	if preCleanup {
		Cleanup()
	}

	tap.CreateNamespace(input.Namespace)
	tap.CreateImagepullSecrets(input.Secrets, input.Namespace)
	tap.AddPackageRepository(input.PackageRepository, input.Namespace)
	tap.CheckPackageRepositoryStatus(input.PackageRepository, input.Namespace)
	tap.ListPackages(input.Namespace)
	// tap.ListValuesSchema(input.Packages, input.Namespace)
	tap.InstallPackages(input.Packages, input.Namespace, input.ValuesDirectory)

	log.Printf("Request: Cleanup post-installation (%t)", postCleanup)
	if postCleanup {
		Cleanup()
	}
}
