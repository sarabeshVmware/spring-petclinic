// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package envfuncs

import (
	"context"
	"fmt"
	"log"
	"time"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/exec"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

func AddPackageRepository(name string, image string, namespace string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("adding package repository %s", name)
		cmd, output, err := exec.TanzuAddPackageRepository(name, image, namespace)
		log.Printf("command executed: %s", cmd)
		if err != nil {
			return ctx, fmt.Errorf("error while adding package repository %s: %w: %s", name, err, output)
		}
		log.Printf("package repository %s added: %s", name, output)
		return ctx, nil
	}
}

func DeletePackageRepository(name string, namespace string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("deleting package repository %s", name)
		cmd, output, err := exec.TanzuDeletePackageRepository(name, namespace)
		log.Printf("command executed: %s", cmd)
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
		cmd, output, err := exec.TanzuGetPackageRepositoryStatus(name, namespace)
		log.Printf("command executed: %s", cmd)
		if err != nil {
			return ctx, fmt.Errorf("error while getting package repository %s status: %w: %s", name, err, output)
		}
		for recursiveCount > 0 {
			if output == "Reconciling" || output == "" {
				log.Printf("package repository %s is getting reconciled: %s", name, output)
				log.Printf("sleeping: 60 seconds")
				time.Sleep(1 * time.Minute)
				recursiveCount -= 1
			} else if output == "Reconcile succeeded" {
				log.Printf("package repository %s reconcilation succeeded: %s", name, output)
				return ctx, nil
			} else if output == "Reconcile Failed" {
				return ctx, fmt.Errorf("package repository %s reconcilation failed: %s", name, output)
			} else {
				return ctx, fmt.Errorf("package repository %s reconcilation unknown: %s", name, output)
			}
		}
		return ctx, fmt.Errorf(`package repository %s is not getting in "Reconcile succeeded" state after %d iterations: %s`, name, recursiveCount, output)
	}
}
