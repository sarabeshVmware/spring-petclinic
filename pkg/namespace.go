// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"fmt"
	"log"
)

func CheckIfNamespaceExists(namespace string) bool {
	log.Printf("Checking if namespace exists: %s", namespace)
	_, err := Run_AllowError(fmt.Sprintf("kubectl get namespaces %s", namespace))
	if err != nil {
		log.Printf("Namespace not found: %s", namespace)
		return false
	}
	log.Printf("Namespace found: %s", namespace)
	return true
}

func DeleteNamespace(namespace string) {
	if CheckIfNamespaceExists(namespace) {
		log.Printf("Deleting namespace: %s", namespace)
		Run(fmt.Sprintf("kubectl delete ns %s", namespace))
	}
}

func CreateNamespace(namespace string) {
	log.Printf("Creating namespace: %s", namespace)
	Run(fmt.Sprintf("kubectl create ns %s", namespace))
}
