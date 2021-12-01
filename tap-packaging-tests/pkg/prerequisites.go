// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import "log"

func HandlePrerequisites() {
	log.Printf("Handling prerequisites:")
	DeployApp("kapp-controller", []string{"https://github.com/vmware-tanzu/carvel-kapp-controller/releases/latest/download/release.yml"})
	// DeployApp("cert-manager", []string{"https://github.com/jetstack/cert-manager/releases/download/v1.5.4/cert-manager.yaml"})
	DeployApp("secretgen-controller", []string{"https://github.com/vmware-tanzu/carvel-secretgen-controller/releases/download/v0.5.0/release.yml"})
	// CreateNamespace("flux-system")
	// CreateClusterRoleBinding("default-admin", "cluster-admin", "flux-system:default")
	// DeployAppInNamespace("flux-source-controller", []string{"https://github.com/fluxcd/source-controller/releases/download/v0.15.4/source-controller.crds.yaml", "https://github.com/fluxcd/source-controller/releases/download/v0.15.4/source-controller.deployment.yaml"}, "flux-system")
}
