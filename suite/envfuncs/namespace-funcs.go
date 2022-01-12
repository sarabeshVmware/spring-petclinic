// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package envfuncs

import (
	"context"
	"fmt"
	"log"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/exec"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

func CreateNamespaces(namespaces []string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		for _, namespace := range namespaces {
			log.Printf("creating namespace %s", namespace)
			cmd, output, err := exec.KubectlCreateNamespace(namespace)
			log.Printf("command executed: %s", cmd)
			if err != nil {
				return ctx, fmt.Errorf("error while creating namespace %s: %w: %s", namespace, err, output)
			}
			log.Printf("namespace %s created: %s", namespace, output)
		}
		return ctx, nil
	}
}

func DeleteNamespaces(namespaces []string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		for _, namespace := range namespaces {
			log.Printf("deleting namespace %s", namespace)
			cmd, output, err := exec.KubectlDeleteNamespace(namespace)
			log.Printf("command executed: %s", cmd)
			if err != nil {
				return ctx, fmt.Errorf("error while deleting namespace %s:  %w: %s", namespace, err, output)
			}
			log.Printf("namespace %s deleted: %s", namespace, output)
		}
		return ctx, nil
	}
}
