// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package envfuncs

import (
	"context"
	"fmt"
	"log"

	"gitlab.eng.vmware.com/tap/tap-packages/tap-packaging-tests/suite/exec"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

func CreateSecret(name string, registry string, username string, password string, namespace string, export bool) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("creating secret %s", name)
		_, output, err := exec.TanzuCreateSecret(name, registry, username, password, namespace, export)
		// don't log command as it contains password
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
		cmd, output, err := exec.TanzuDeleteSecret(name, namespace)
		log.Printf("command executed: %s", cmd)
		if err != nil {
			return ctx, fmt.Errorf("error while deleting secret %s: %w: %s", name, err, output)
		}
		log.Printf("secret %s deleted: %s", name, output)
		return ctx, nil
	}
}
