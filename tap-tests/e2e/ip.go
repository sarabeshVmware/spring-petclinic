// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"log"
	"strings"

	tap "gitlab.eng.vmware.com/tap/tap-packaging-tests/pkg"
)

func GetAppAcceleratorExternalIP() string {
	log.Printf("Getting app accelerator external IP:")
	appAccExternalIPBytes, _ := tap.RunWithBash(`kubectl get svc -A | awk '{if($2=="acc-ui-server")print $5}'`)
	appAccExternalIP := strings.TrimSpace(string(appAccExternalIPBytes))
	log.Printf("App Accelerator external IP: %s", appAccExternalIP)
	return appAccExternalIP
}

func GetAppLiveViewExternalIP() string {
	log.Printf("Getting app live view external IP:")
	appLiveViewExternalIPBytes, _ := tap.RunWithBash(`kubectl get svc -A | awk '{if($2=="application-live-view-5112")print $5}'`)
	appLiveViewExternalIP := strings.TrimSpace(string(appLiveViewExternalIPBytes))
	log.Printf("App Live View external IP: %s", appLiveViewExternalIP)
	return appLiveViewExternalIP
}

func GetEnvoyExternalIP() string {
	log.Printf("Getting envoy external IP:")
	envoyExternalIPBytes, _ := tap.RunWithBash(`kubectl get svc -n contour-external | awk '{if($1=="envoy")print $4}'`)
	envoyExternalIP := strings.TrimSpace(string(envoyExternalIPBytes))
	log.Printf("Envoy external IP: %s", envoyExternalIP)
	return envoyExternalIP
}
