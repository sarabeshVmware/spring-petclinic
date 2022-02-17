// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package envfuncs

import (
	"context"
	"fmt"
	"log"
	"time"

	kubectl_helper "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_helpers"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzuCmds"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

func InstallPackage(name string, packageName string, version string, namespace string, valuesFile string, pollTimeout string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("installing package %s (%s)", name, packageName)

		// install package
		err := tanzuCmds.TanzuInstallPackage(name, packageName, version, namespace, valuesFile, pollTimeout)
		if err != nil {
			// if error, check via kubectl, not tanzu-cli
			pass := kubectl_helper.ValidateTAPInstallation(name, namespace, 10, 60)
			if !pass {
				kubectl_helper.LogFailedResourcesDetails(namespace)
				return ctx, fmt.Errorf("error while installing package %s (%s)", name, packageName)
			} else {
				return ctx, nil
			}
		}

		return ctx, nil
	}
}

func UninstallPackage(name string, namespace string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("uninstalling package %s", name)

		// uninstall package
		err := tanzuCmds.TanzuUninstallPackage(name, namespace)
		if err != nil {
			return ctx, fmt.Errorf("error while installing package %s", name)
		}

		return ctx, nil
	}
}

func ListInstalledPackages(namespace string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("listing packages in namespace %s", namespace)

		// list installed packages
		err := tanzuCmds.TanzuListInstalledPackages(namespace)
		if err != nil {
			return ctx, fmt.Errorf("error while listing packages in namespace %s", namespace)
		}

		return ctx, nil
	}
}

func CheckIfPackageInstalled(name string, namespace string, recursiveCount int, secondsGap int) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("checking package %s installation status", name)

		for ; recursiveCount >= 0; recursiveCount-- {
			// get status
			packageInstalledStatus, err := tanzuCmds.TanzuGetPackageInstalledStatus(name, namespace)
			if err != nil {
				return ctx, fmt.Errorf("error while getting package %s in namespace %s installation status", name, namespace)
			}

			// check
			if packageInstalledStatus == "Reconciling" || packageInstalledStatus == "" {
				log.Printf("package %s is getting installed", name)
				log.Printf("sleeping for %d seconds", secondsGap)
				time.Sleep(time.Duration(secondsGap) * time.Second)
			} else if packageInstalledStatus == "Reconcile succeeded" {
				log.Printf("package %s is installed", name)
				return ctx, nil
			} else if packageInstalledStatus == "Reconcile Failed" {
				return ctx, fmt.Errorf("package %s installation failed", name)
			} else {
				return ctx, fmt.Errorf("package %s installation unknown", name)
			}
		}

		return ctx, fmt.Errorf(`package %s is not getting in "Reconcile succeeded" state`, name)
	}
}
