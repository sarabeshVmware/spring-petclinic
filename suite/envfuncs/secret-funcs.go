// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package envfuncs

import (
	"context"
	"fmt"
	"log"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzuCmds"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

func CreateSecret(name string, registry string, username string, password string, namespace string, export bool) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("creating secret %s (registry %s, username %s) in namespace %s", name, registry, username, namespace)

		// create secret
		err := tanzuCmds.TanzuCreateSecret(name, registry, username, password, "string", namespace, export)
		if err != nil {
			return ctx, fmt.Errorf("error while creating secret %s (registry %s, username %s) in namespace %s", name, registry, username, namespace)
		}

		return ctx, nil
	}
}

func DeleteSecret(name string, namespace string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("deleting secret %s from namespace %s", name, namespace)

		// delete secret
		err := tanzuCmds.TanzuDeleteSecret(name, namespace)
		if err != nil {
			return ctx, fmt.Errorf("error while deleting secret %s from namespace %s", name, namespace)
		}

		return ctx, nil
	}
}
