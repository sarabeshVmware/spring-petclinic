// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"
	"log"
	"time"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func CheckDeploymentExists(name string, c *rest.Config) (bool, error) {
	log.Printf("checking if deployment %s exists", name)

	// create new clientset
	clientset, err := kubernetes.NewForConfig(c)
	if err != nil {
		log.Printf("error while creating clientset")
		log.Printf("error: %s", err)
		return false, err
	} else {
		log.Print("created new clientset")
	}

	// get deployments list
	deployments, err := clientset.AppsV1().Deployments(apiv1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("error while getting deployments list for %s", name)
		log.Printf("error: %s", err)
		return false, err
	}

	for _, item := range deployments.Items {
		if item.Name == name {
			log.Printf("%s deployment found", name)
			return true, nil
		}
	}

	log.Printf("%s deployment not found", name)
	return false, nil
}

func GetServiceExternalIP(service string, namespace string, c *rest.Config, timeoutInMins int, intervalInSeconds int) (string, error) {
	log.Printf("getting external IP for service %s in namespace %s", service, namespace)

	externalIP := ""

	// create new clientset
	clientset, err := kubernetes.NewForConfig(c)
	if err != nil {
		log.Printf("error while creating clientset")
		log.Printf("error: %s", err)
		return externalIP, err
	} else {
		log.Print("created new clientset")
	}

	finalTimeout := timeoutInMins * 60
	for finalTimeout > 0 {
		// get services list
		svcList, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.name=" + service})
		if err != nil {
			log.Printf("error while getting service list for service %s in namespace %s", service, namespace)
			log.Printf("error: %s", err)
			return externalIP, err
		}

		// get external IP
		found := false
		for _, svc := range svcList.Items {
			if len(svc.Status.LoadBalancer.Ingress) != 0 {
				externalIP = svc.Status.LoadBalancer.Ingress[0].IP
				if externalIP != "" {
					found = true
					break
				}
				externalIP = svc.Status.LoadBalancer.Ingress[0].Hostname
				if externalIP != "" {
					found = true
					break
				}
			}
		}
		if found {
			log.Printf("found external IP: %s", externalIP)
			return externalIP, nil
		} else {
			err := fmt.Errorf("external IP not found for service %s in namespace %s", service, namespace)
			log.Printf("error: %s", err)
			log.Printf("Waiting for %d seconds before retry", intervalInSeconds)
			time.Sleep(time.Duration(intervalInSeconds) * time.Second)
			finalTimeout -= intervalInSeconds
		}

	}
	log.Printf("error while getting external IP for service %s in namespace %s", service, namespace)
	return externalIP, err
}
