// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package envfuncs

import (
	"context"
	"fmt"
	"log"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_helpers"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_libs"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

func InstallClusterEssentials(tanzunetHost string, tanzunetApiToken string, productFileId int, releaseVersion string, productSlug string, downloadBundle string, installBundle string, installRegistryHostname string, InstallRegistryUsername string, installRegistryPassword string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		kappControllerDeployed := kubectl_helpers.CheckDeploymentExists("kapp-controller", "", 1, 30)
		secGenControllerDeployed := kubectl_helpers.CheckDeploymentExists("secretgen-controller", "", 1, 30)
		if kappControllerDeployed || secGenControllerDeployed {
			log.Println("kapp-controller or secretgen-controller deployment exists.")
			return ctx, nil
		}
		// log.Println("Download artifacts from Tanzu Network")
		// log.Println("Logging into tanzunet")
		// if !pivnet_libs.Login(tanzunetHost, tanzunetApiToken) {
		// 	log.Fatalln("Unable to login to tanzunet")
		// }
		// if !pivnet_libs.DownloadProductFile(productFileId, productSlug, releaseVersion) {
		// 	log.Fatalln("Unable to download product file")
		// }
		// extract_cluster_essentials_cmd := fmt.Sprintf("mkdir ./tanzu-cluster-essentials; tar -xvf %s -C ./tanzu-cluster-essentials", downloadBundle)
		// response, err := linux_util.ExecuteCmdInBashMode(extract_cluster_essentials_cmd)
		// if err != nil {
		// 	return ctx, fmt.Errorf("error while deploying cluster-essentials: %w: %s", err, response)
		// }

		// log.Println("Installing Tanzu cluster Essentials...")
		// log.Println("Setting up required environment variables for installing Tanzu Cluster Essentials.")
		// os.Setenv("INSTALL_BUNDLE", installBundle)
		// log.Printf("INSTALL_BUNDLE env set to: %s", os.Getenv("INSTALL_BUNDLE"))
		// os.Setenv("INSTALL_REGISTRY_HOSTNAME", installRegistryHostname)
		// log.Printf("INSTALL_REGISTRY_HOSTNAME env set to: %s", os.Getenv("INSTALL_REGISTRY_HOSTNAME"))
		// os.Setenv("INSTALL_REGISTRY_USERNAME", InstallRegistryUsername)
		// log.Printf("INSTALL_REGISTRY_USERNAME env set to: %s", os.Getenv("INSTALL_REGISTRY_USERNAME"))
		// os.Setenv("INSTALL_REGISTRY_PASSWORD", installRegistryPassword)
		// log.Println("INSTALL_REGISTRY_PASSWORD env set.")

		err := kubectl_libs.KubectlApplyConfiguration("resources/suite/kapp-controller.yml", "")
		if err != nil {
			return ctx, fmt.Errorf("error while setting up developer namespace %s", "kapp-controller")
		}

		err1 := kubectl_libs.KubectlApplyConfiguration("resources/suite/secretgen-controller.yml", "")
		if err1 != nil {
			return ctx, fmt.Errorf("error while setting up developer namespace %s", "secretgen-controller")
		}

		log.Printf("cluster-essentials deployment successful")
		return ctx, nil
	}
}
