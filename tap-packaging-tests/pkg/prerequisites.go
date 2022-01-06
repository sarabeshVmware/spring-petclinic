// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"context"
	"log"
	"time"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CheckPrerequisites() {
	log.Printf("Checking prerequisites:")

	var retries int = 10
	for retries > 0 {
		deployments, err := GetClientset().AppsV1().Deployments(apiv1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
		if CheckErrorWithoutExit(err) != true {
			log.Println("Status code is :", deployments.StatusCode)
			break
		} else {
			retries -= 1
			log.Printf("Retry after 30 seconds")
			time.Sleep(30 * time.Second)
		}
	}
	deployments, err := GetClientset().AppsV1().Deployments(apiv1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	CheckError(err)
	kappControllerInstalled, secretgenControllerInstalled := false, false
	for _, item := range deployments.Items {
		if item.Namespace == "kapp-controller" {
			kappControllerInstalled = true
		}
		if item.Namespace == "secretgen-controller" {
			secretgenControllerInstalled = true
		}
	}

	if kappControllerInstalled {
		log.Printf("kapp-controller already deployed.")
	} else {
		log.Printf("Deploying kapp-controller:")
		DeployApp("kapp-controller", []string{"https://github.com/vmware-tanzu/carvel-kapp-controller/releases/latest/download/release.yml"})
	}

	if secretgenControllerInstalled {
		log.Printf("secretgen-controller already deployed.")
	} else {
		log.Printf("Deploying secretgen-controller:")
		DeployApp("secretgen-controller", []string{"https://github.com/vmware-tanzu/carvel-secretgen-controller/releases/download/v0.5.0/release.yml"})
	}
}
