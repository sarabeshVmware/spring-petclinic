// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import "log"

func CheckIfNamespaceExists(namespace string) bool {
	log.Printf("Checking if namespace exists: %s", namespace)
	_, err := RunCommand(Command{CommandName: "kubectl", Arguments: []string{"get", "namespaces", namespace}, AllowError: true})
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
		RunCommand(Command{CommandName: "kubectl", Arguments: []string{"delete", "ns", namespace}})
	}
}

func CreateNamespace(namespace string) {
	log.Printf("Creating namespace: %s", namespace)
	RunCommand(Command{CommandName: "kubectl", Arguments: []string{"create", "ns", namespace}})
}
