// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package envfuncs

import (
	"context"
	"fmt"
	"log"

	e2e "gitlab.eng.vmware.com/tap/tap-packages/tap-packaging-tests/suite/pkg"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

func CreateSecret(name string, registry string, username string, password string, namespace string, export bool) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("creating secret %s", name)
		output, err := e2e.CreateSecret(name, registry, username, password, namespace, export)
		if err != nil {
			return ctx, fmt.Errorf("error while creating secret %s: %w: %s", name, err, output)
		}
		log.Printf("secret %s created: %s", name, output)

		return ctx, nil
	}
}

func DeleteSecret(name string, registry string, username string, password string, namespace string, export bool) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("deleting secret %s", name)
		output, err := e2e.DeleteSecret(name, namespace)
		if err != nil {
			return ctx, fmt.Errorf("error while deleting secret %s: %w: %s", name, err, output)
		}
		log.Printf("secret %s deleted: %s", name, output)

		return ctx, nil
	}
}
