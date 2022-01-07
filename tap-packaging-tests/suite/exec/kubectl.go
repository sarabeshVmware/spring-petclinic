// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import "fmt"

func KubectlCreateNamespace(namespace string) (string, string, error) {
	cmd := fmt.Sprintf("kubectl create ns %s", namespace)
	output, err := RunCommand(cmd)
	return cmd, output, err
}

func KubectlDeleteNamespace(namespace string) (string, string, error) {
	cmd := fmt.Sprintf("kubectl delete ns %s", namespace)
	output, err := RunCommand(cmd)
	return cmd, output, err
}

func KubectlPatchServiceAccount(serviceAccount string, namespace string, patch string) (string, string, error) {
	cmd := fmt.Sprintf("kubectl patch serviceaccount %s -n %s -p %s", serviceAccount, namespace, patch)
	output, err := RunCommand(cmd)
	return cmd, output, err
}

func KubectlCreateClusterRoleBinding(name string, clusterRole string, serviceAccount string) (string, string, error) {
	cmd := fmt.Sprintf("kubectl create clusterrolebinding %s --clusterrole=%s --serviceaccount=%s", name, clusterRole, serviceAccount)
	output, err := RunCommand(cmd)
	return cmd, output, err
}
