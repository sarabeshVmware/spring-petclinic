// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package envfuncs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_libs"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

func AddFinalizersToKappControllerClusterRole() env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		clusterConfig := kubectl_libs.GetCurrentContext()
		if !strings.Contains(clusterConfig, "aroapp") {
			log.Printf("Not ARO, skipping workaround")
		} else {
			cmd := "kubectl get clusterrole kapp-controller-cluster-role -o yaml | yq '.rules[3].resources[2]=\"packageinstalls/finalizers\"' | kubectl apply -f -"
			output, err := linux_util.ExecuteCmdInBashMode(cmd)
			if err != nil {
				return ctx, fmt.Errorf("error while editing kapp-controller: %w: %s", err, output)
			}
			log.Printf("Modifying the clusterrole kapp-controller-cluster-role to add packageinstalls/finalizers in the packaging.carvel.dev apiGroup successful")
		}
		return ctx, nil
	}
}

func CreateClusterRoleBinding() env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		clusterConfig := kubectl_libs.GetCurrentContext()
		if !strings.Contains(clusterConfig, "aroapp") {
			log.Printf("Not ARO, skipping workaround")
		} else {
			cmd := "kubectl create clusterrolebinding apps-admin --clusterrole=cluster-admin --serviceaccount=my-apps:default"
			output, err := linux_util.ExecuteCmdInBashMode(cmd)
			if err != nil {
				return ctx, fmt.Errorf("error while creating cluster role binding: %w: %s", err, output)
			}
			log.Printf("Creating cluster role binding successful")
		}
		return ctx, nil
	}
}
