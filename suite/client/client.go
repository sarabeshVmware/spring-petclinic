// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetServiceExternalIP(service string, namespace string, c *rest.Config) (string, error) {
	clientset, err := kubernetes.NewForConfig(c)
	if err != nil {
		return "", err
	}
	svcList, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.name=" + service})
	if err != nil {
		return "", err
	}
	for _, svc := range svcList.Items {
		if len(svc.Status.LoadBalancer.Ingress) != 0 {
			externalIP := svc.Status.LoadBalancer.Ingress[0].IP
			if externalIP != "" {
				return externalIP, nil
			}
			externalIP = svc.Status.LoadBalancer.Ingress[0].Hostname
			if externalIP != "" {
				return externalIP, nil
			}
		}
	}
	return "", fmt.Errorf("external IP not found for service %s in namespace %s", service, namespace)
}

func GetServicePort(service string, namespace string, c *rest.Config) (int, error) {
	clientset, err := kubernetes.NewForConfig(c)
	if err != nil {
		return -1, err
	}
	svcList, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.name=" + service})
	if err != nil {
		return -1, err
	}
	for _, svc := range svcList.Items {
		if len(svc.Spec.Ports) != 0 {
			return int(svc.Spec.Ports[0].Port), nil
		}
	}
	return -1, fmt.Errorf("port not found for service %s in namespace %s", service, namespace)
}

func CheckDeploymentExists(name string, c *rest.Config) (bool, error) {
	clientset, err := kubernetes.NewForConfig(c)
	if err != nil {
		return false, err
	}
	deployments, err := clientset.AppsV1().Deployments(apiv1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, err
	}
	for _, item := range deployments.Items {
		if item.Name == name {
			return true, nil
		}
	}
	return false, nil
}
