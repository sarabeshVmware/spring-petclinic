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

// TODO
// func InstallDependencies(packageInfo e2e.Package) error {
// }

func InstallPackageByInfo(packageInfo e2e.Package) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		log.Printf("installing package: %s", packageInfo.Package)

		// TODO how to install dependencies? : how to do recursive: do topological sort?

		// TODO how to install main package? : separate function?

		// install package dependencies:
		dependentPackages, err := e2e.GetDependentPackagesInfo(packageInfo)
		if err != nil {
			return ctx, fmt.Errorf("error while getting dependent packages for %s: %w", packageInfo.Package, err)
		}
		log.Printf("dependencies for package %s: %s", packageInfo.Package, dependentPackages)

		for _, dependentPackageInfo := range dependentPackages {
			log.Printf("installing package dependency: %s", dependentPackageInfo.Package)
			InstallPackageByInfo(dependentPackageInfo)
		}

		// install:
		cmd := fmt.Sprintf("tanzu package install %s -p %s -v %s -n %s --poll-timeout 30m", packageInfo.Name, packageInfo.Package, packageInfo.Version, packageInfo.Namespace)
		if packageInfo.ValuesFile != "" {
			cmd += fmt.Sprintf(" -f %s", packageInfo.ValuesFile)
		}
		log.Printf("installing package %s: %s", packageInfo.Package, cmd)
		output, err := e2e.RunCommand(cmd)
		if err != nil {
			return ctx, fmt.Errorf("error while installing package %s: %w", packageInfo.Package, err)
		}
		log.Printf("install status: %s", output)

		return ctx, nil
	}
}

func InstallPackageByName(packageName string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		packageInfo, err := e2e.GetPackageInfoFromName(packageName)
		if err != nil {
			return ctx, fmt.Errorf("error while installing package by name: %w", err)
		}
		// TOOD: how to call?
		InstallPackageByInfo(packageInfo)

		return ctx, nil
	}
}
