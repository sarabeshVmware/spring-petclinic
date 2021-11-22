// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"context"
	"log"
	"os"
	"path/filepath"

	tap "gitlab.eng.vmware.com/tap/tap-packaging-tests/pkg"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func GetAppAcceleratorExternalIP() string {
	log.Printf("Getting app accelerator external IP:")
	// appAccExternalIPBytes, _ := tap.RunWithBash(`kubectl get svc -A | awk '{if($2=="acc-ui-server")print $5}'`)
	// appAccExternalIP := strings.TrimSpace(string(appAccExternalIPBytes))
	appAccExternalIP := GetServiceExternalIp("acc-ui-server", "accelerator-system")
	log.Printf("App Accelerator external IP: %s", appAccExternalIP)
	return appAccExternalIP
}

func GetAppLiveViewExternalIP() string {
	log.Printf("Getting app live view external IP:")
	// appLiveViewExternalIPBytes, _ := tap.RunWithBash(`kubectl get svc -A | awk '{if($2=="application-live-view-5112")print $5}'`)
	// appLiveViewExternalIP := strings.TrimSpace(string(appLiveViewExternalIPBytes))
	appLiveViewExternalIP := GetServiceExternalIp("application-live-view-5112", "app-live-view")
	log.Printf("App Live View external IP: %s", appLiveViewExternalIP)
	return appLiveViewExternalIP
}

func GetEnvoyExternalIP() string {
	log.Printf("Getting envoy external IP:")
	// envoyExternalIPBytes, _ := tap.RunWithBash(`kubectl get svc -n contour-external | awk '{if($1=="envoy")print $4}'`)
	// envoyExternalIP := strings.TrimSpace(string(envoyExternalIPBytes))
	envoyExternalIP := GetServiceExternalIp("envoy", "contour-external")
	log.Printf("Envoy external IP: %s", envoyExternalIP)
	return envoyExternalIP
}

func GetServiceExternalIp(serviceName string, namespace string) string {
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	var externalIp string
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	tap.CheckError(err)
	clientset, err := kubernetes.NewForConfig(config)
	tap.CheckError(err)
	svcList, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: "metadata.name=" + serviceName,
	})
	tap.CheckError(err)
	for _, svc := range svcList.Items {
		externalIp = svc.Status.LoadBalancer.Ingress[0].IP
		if externalIp == "" { // if no IP, check for hostname
			externalIp = svc.Status.LoadBalancer.Ingress[0].Hostname
		} else {
			break
		}
	}
	return externalIp
}
