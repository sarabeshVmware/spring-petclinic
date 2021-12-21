// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package envfuncs

import (
	"context"
	"fmt"
	"log"
	"time"

	e2e "gitlab.eng.vmware.com/tap/tap-packages/tap-packaging-tests/e2e-suite/pkg"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

func CheckIfPackageInstalled(packageName string, namespace string, recursiveCount int) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("checking if package is installed: %s", packageName)

		installedPackages, err := e2e.GetInstalledPackages(namespace)
		if err != nil {
			return ctx, fmt.Errorf("error while getting installed packages: %w", err)
		}
		log.Printf("installed packages: \n%s", installedPackages)

		for _, installedPackage := range installedPackages {
			if packageName == installedPackage.PackageName {
				log.Printf("checking package status, iteration: %d", recursiveCount)
				installedPackageStatus, err := e2e.GetInstalledPackageStatus(installedPackage.Name, namespace)
				if err != nil {
					return ctx, fmt.Errorf("error while getting installed package %s: %w", installedPackage.Name, err)
				}
				for recursiveCount > 0 {
					if installedPackageStatus == "Reconciling" {
						log.Printf("package %s is getting installed: %s", packageName, installedPackage.Status)
						log.Printf("sleep: 60 seconds")
						time.Sleep(1 * time.Minute)
						recursiveCount -= 1
					} else if installedPackageStatus == "Reconcile succeeded" {
						log.Printf("package %s is installed: %s", packageName, installedPackage.Status)
						return ctx, nil
					} else if installedPackageStatus == "Reconcile Failed" {
						return ctx, fmt.Errorf("package %s installation failed: %s", packageName, installedPackage.Status)
					} else {
						return ctx, fmt.Errorf("package %s installation unknown: %s", packageName, installedPackage.Status)
					}
				}
				return ctx, fmt.Errorf(`package %s is not getting in "Reconcile succeeded" state after %d iterations: %s`, packageName, recursiveCount, installedPackage.Status)
			}
		}
		return ctx, fmt.Errorf("package %s not found in the list of installed packages", packageName)
	}
}
