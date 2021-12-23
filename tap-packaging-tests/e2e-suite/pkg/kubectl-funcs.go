// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"fmt"
	"log"
)

func CreateNamespace(namespace string) (string, error) {
	cmd := fmt.Sprintf("kubectl create ns %s", namespace)
	log.Printf("command executed: %s", cmd)
	output, err := RunCommand(cmd)

	return output, err
}

func DeleteNamespace(namespace string) (string, error) {
	cmd := fmt.Sprintf("kubectl delete ns %s", namespace)
	log.Printf("command executed: %s", cmd)
	output, err := RunCommand(cmd)

	return output, err
}
