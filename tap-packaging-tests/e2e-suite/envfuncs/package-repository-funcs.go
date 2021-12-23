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

func AddPackageRepository(name string, image string, namespace string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("adding package repository %s", name)
		output, err := e2e.AddPackageRepository(name, image, namespace)
		if err != nil {
			return ctx, fmt.Errorf("error while adding package repository %s: %w: %s", name, err, output)
		}
		log.Printf("package repository %s added: %s", name, output)

		return ctx, nil
	}
}

func DeletePackageRepository(name string, image string, namespace string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("deleting package repository %s", name)
		output, err := e2e.DeletePackageRepository(name, namespace)
		if err != nil {
			return ctx, fmt.Errorf("error while deleting package repository %s: %w: %s", name, err, output)
		}
		log.Printf("package repository %s deleted: %s", name, output)

		return ctx, nil
	}
}

func CheckIfPackageRepositoryReconciled(name string, namespace string, recursiveCount int) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("checking package repository %s status", name)

		log.Printf("getting package repository %s status", name)
		packageRepositoryStatus, err := e2e.GetPackageRepositoryStatus(name, namespace)
		if err != nil {
			return ctx, fmt.Errorf("error while getting package repository %s status: %w: %s", name, err, packageRepositoryStatus)
		}
		for recursiveCount > 0 {
			if packageRepositoryStatus == "Reconciling" || packageRepositoryStatus == "" {
				log.Printf("package repository %s is getting reconciled: %s", name, packageRepositoryStatus)
				log.Printf("sleeping: 60 seconds")
				time.Sleep(1 * time.Minute)
				recursiveCount -= 1
			} else if packageRepositoryStatus == "Reconcile succeeded" {
				log.Printf("package repository %s reconcilation succeeded: %s", name, packageRepositoryStatus)
				return ctx, nil
			} else if packageRepositoryStatus == "Reconcile Failed" {
				return ctx, fmt.Errorf("package repository %s reconcilation failed: %s", name, packageRepositoryStatus)
			} else {
				return ctx, fmt.Errorf("package repository %s reconcilation unknown: %s", name, packageRepositoryStatus)
			}
		}

		return ctx, fmt.Errorf(`package repository %s is not getting in "Reconcile succeeded" state after %d iterations: %s`, name, recursiveCount, packageRepositoryStatus)
	}
}
