// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package envfuncs

import (
	"context"
	"fmt"
	"log"
	"time"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzuCmds"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

func AddPackageRepository(name string, image string, namespace string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("adding package repository %s (%s) in namespace %s", name, image, namespace)

		// add repo
		err := tanzuCmds.TanzuAddPackageRepository(name, image, namespace)
		if err != nil {
			return ctx, fmt.Errorf("error while adding package repository %s (%s) in namespace %s", name, image, namespace)
		}

		return ctx, nil
	}
}

func DeletePackageRepository(name string, namespace string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("deleting package repository %s in namespace %s", name, namespace)

		// delete repo
		err := tanzuCmds.TanzuDeletePackageRepository(name, namespace)
		if err != nil {
			return ctx, fmt.Errorf("error while deleting package repository %s in namespace %s", name, namespace)
		}

		return ctx, nil
	}
}

func CheckIfPackageRepositoryReconciled(name string, namespace string, recursiveCount int, secondsGap int) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("checking package repository %s status", name)

		for ; recursiveCount >= 0; recursiveCount-- {
			// get status
			packageRepositoryStatus, err := tanzuCmds.TanzuGetPackageRepositoryStatus(name, namespace)
			if err != nil {
				return ctx, fmt.Errorf("error while getting package repository %s in namespace %s status", name, namespace)
			}

			// check
			if packageRepositoryStatus == "Reconciling" || packageRepositoryStatus == "" {
				log.Printf("package repository %s is getting reconciled", name)
				log.Printf("sleeping for %d seconds", secondsGap)
				time.Sleep(time.Duration(secondsGap) * time.Second)
			} else if packageRepositoryStatus == "Reconcile succeeded" {
				log.Printf("package repository %s reconcilation succeeded", name)
				return ctx, nil
			} else if packageRepositoryStatus == "Reconcile Failed" {
				return ctx, fmt.Errorf("package repository %s reconcilation failed", name)
			} else {
				return ctx, fmt.Errorf("package repository %s reconcilation unknown", name)
			}
		}
		return ctx, fmt.Errorf(`package repository %s is not getting in "Reconcile succeeded" state`, name)
	}
}
