// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package envfuncs

import (
	"context"
	"fmt"
	"log"
	"time"

	"gitlab.eng.vmware.com/tap/tap-packages/tap-packaging-tests/suite/exec"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

func InstallPackage(name string, packageName string, version string, namespace string, valuesFile string, pollTimeout string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("installing package %s", name)
		cmd, output, err := exec.TanzuInstallPackage(name, packageName, version, namespace, valuesFile, pollTimeout)
		log.Printf("command executed: %s", cmd)
		if err != nil {
			return ctx, fmt.Errorf("error while installing package %s: %w: %s", name, err, output)
		}
		log.Printf("package %s installed: %s", name, output)
		return ctx, nil
	}
}

func UninstallPackage(name string, namespace string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("uninstalling package %s", name)
		cmd, output, err := exec.TanzuUninstallPackage(name, namespace)
		log.Printf("command executed: %s", cmd)
		if err != nil {
			return ctx, fmt.Errorf("error while uninstalling package %s: %w: %s", name, err, output)
		}
		log.Printf("package %s uninstalled: %s", name, output)
		return ctx, nil
	}
}

func CheckIfPackageInstalled(name string, namespace string, recursiveCount int) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("checking package %s installation status", name)
		log.Printf("getting package %s installation status", name)
		cmd, output, err := exec.TanzuGetPackageInstalledStatus(name, namespace)
		log.Printf("command executed: %s", cmd)
		if err != nil {
			return ctx, fmt.Errorf("error while getting package %s installation status: %w: %s", name, err, output)
		}
		for recursiveCount > 0 {
			if output == "Reconciling" || output == "" {
				log.Printf("package %s is getting installed: %s", name, output)
				log.Printf("sleeping: 60 seconds")
				time.Sleep(1 * time.Minute)
				recursiveCount -= 1
			} else if output == "Reconcile succeeded" {
				log.Printf("package %s is installed: %s", name, output)
				return ctx, nil
			} else if output == "Reconcile Failed" {
				return ctx, fmt.Errorf("package %s installation failed: %s", name, output)
			} else {
				return ctx, fmt.Errorf("package %s installation unknown: %s", name, output)
			}
		}
		return ctx, fmt.Errorf(`package %s is not getting in "Reconcile succeeded" state after %d iterations (%s)`, name, recursiveCount, output)
	}
}
