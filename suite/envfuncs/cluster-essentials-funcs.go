// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package envfuncs

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubernetes/client"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

func InstallClusterEssentials(bundle string, registry string, username string, password string, filename string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		kappControllerDeployed, err := client.CheckDeploymentExists("kapp-controller", cfg.Client().RESTConfig())
		if err != nil {
			return ctx, fmt.Errorf("error while checking for kapp-controller deployment: %w", err)
		}
		secGenControllerDeployed, err := client.CheckDeploymentExists("secretgen-controller", cfg.Client().RESTConfig())
		if err != nil {
			return ctx, fmt.Errorf("error while checking for secretgen-controller deployment: %w", err)
		}
		if kappControllerDeployed || secGenControllerDeployed {
			log.Println("kapp-controller or secretgen-controller deployment exists.")
			return ctx, nil
		}
		log.Println("Installing Tanzu cluster Essentials...")
		log.Println("Setting up required environment variables for installing Tanzu Cluster Essentials.")
		os.Setenv("INSTALL_BUNDLE", bundle)
		log.Printf("INSTALL_BUNDLE env set to: %s", os.Getenv("INSTALL_BUNDLE"))
		os.Setenv("INSTALL_REGISTRY_HOSTNAME", registry)
		log.Printf("INSTALL_REGISTRY_HOSTNAME env set to: %s", os.Getenv("INSTALL_REGISTRY_HOSTNAME"))
		os.Setenv("INSTALL_REGISTRY_USERNAME", username)
		log.Printf("INSTALL_REGISTRY_USERNAME env set to: %s", os.Getenv("INSTALL_REGISTRY_USERNAME"))
		os.Setenv("INSTALL_REGISTRY_PASSWORD", password)
		log.Println("INSTALL_REGISTRY_PASSWORD env set.")
		//executefrom := filedir
		wd, _ := os.Getwd()
		//file := path.Join(wd, filedir, "install.sh")
		file := path.Join(wd, filename)
		output, err := linux_util.RunBashFile(file, "")
		log.Printf("File %s executed successfully", file)
		if err != nil {
			return ctx, fmt.Errorf("error while deploying cluster-essentials: %w: %s", err, output)
		}
		log.Printf("cluster-essentials deployment successful")
		return ctx, nil
	}
}
