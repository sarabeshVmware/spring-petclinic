// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package envfuncs

import (
	"context"
	"fmt"
	"log"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/client"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/exec"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

func CheckAndDeploy(name string, files []string, namespace string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		isDeployed, err := client.CheckDeploymentExists("kapp-controller", cfg.Client().RESTConfig())
		if err != nil {
			return ctx, fmt.Errorf("error while checking for deployment %s: %w", name, err)
		}
		if isDeployed {
			log.Printf("deployment %s exists", name)
			return ctx, nil
		}
		cmd, output, err := exec.KappDeployAppInNamespace(name, files, namespace)
		log.Printf("command executed: %s", cmd)
		if err != nil {
			return ctx, fmt.Errorf("error while deploying %s in namespace %s: %w: %s", name, namespace, err, output)
		}
		log.Printf("deployment %s successful in namespace %s", name, namespace)
		return ctx, nil
	}
}
