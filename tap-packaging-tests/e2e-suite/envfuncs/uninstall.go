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

func UninstallPackage(namespace string, installedPackageName string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		cmd := fmt.Sprintf("tanzu package installed delete %s -n %s -y", installedPackageName, namespace)
		log.Printf("uninstalling package %s (namespace %s): %s", installedPackageName, namespace, cmd)
		output, err := e2e.RunCommand(cmd)
		if err != nil {
			return ctx, fmt.Errorf("error while uninstalling package %s: %w", installedPackageName, err)
		}
		log.Printf("uninstall status: %s", output)

		return ctx, nil
	}
}

// func UninstallPackages(namespace string) env.Func {
// 	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
// 		installedpackages, err := e2e.GetInstalledPackages(namespace)
// 		if err != nil {
// 			return ctx, fmt.Errorf("error while uninstalling packages in namespace %s: %w", namespace, err)
// 		}
// 		// TODO: how to call UninstallPackage?
// 		for _, each := range installedpackages {
// 			log.Printf("uninstalling package: %s", each.Name)
// 			DefaultRun(fmt.Sprintf("tanzu package installed delete %s -n %s -y", each.Name, namespace))
// 		}
// 	}
// }
