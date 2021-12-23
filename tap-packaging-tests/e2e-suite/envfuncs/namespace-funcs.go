// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package envfuncs

import (
	"context"
	"fmt"
	"log"

	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"

	e2e "gitlab.eng.vmware.com/tap/tap-packages/tap-packaging-tests/e2e-suite/pkg"
)

func CreateNamespaces(namespaces []string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		for _, namespace := range namespaces {
			log.Printf("creating namespace %s", namespace)
			output, err := e2e.CreateNamespace(namespace)
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
			output, err := e2e.DeleteNamespace(namespace)
			if err != nil {
				return ctx, fmt.Errorf("error while deleting namespace %s:  %w: %s", namespace, err, output)
			}
			log.Printf("namespace %s deleted: %s", namespace, output)
		}

		return ctx, nil
	}
}
