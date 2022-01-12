// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"encoding/json"
	"fmt"
)

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
	output, err := RunCommandInBashMode(cmd)
	return cmd, output, err
}

func KubectlCreateClusterRoleBinding(name string, clusterRole string, serviceAccount string) (string, string, error) {
	cmd := fmt.Sprintf("kubectl create clusterrolebinding %s --clusterrole=%s --serviceaccount=%s", name, clusterRole, serviceAccount)
	output, err := RunCommand(cmd)
	return cmd, output, err
}

func KubectlGetImageRepositoryStatus(name string, namespace string) (string, string, error) {
	cmd := fmt.Sprintf("kubectl get imagerepositories -n %s -o json", namespace)
	output, err := RunCommand(cmd)
	if err != nil {
		return cmd, "", err
	}
	imageRepositories := struct {
		Items []struct {
			Metadata struct {
				Name string `json:"name"`
			} `json:"metadata"`
			Status struct {
				Conditions []struct {
					Type string `json:"type"`
				} `json:"conditions"`
			} `json:"status"`
		} `json:"items"`
	}{}
	err = json.Unmarshal([]byte(output), &imageRepositories)
	if err != nil {
		return cmd, "", err
	}
	if len(imageRepositories.Items) <= 0 {
		return cmd, "", fmt.Errorf("list empty for image repositories for namespace %s", namespace)
	}
	for _, item := range imageRepositories.Items {
		if item.Metadata.Name == name {
			return cmd, item.Status.Conditions[len(item.Status.Conditions)-1].Type, nil
		}
	}
	return cmd, "", fmt.Errorf("%s not found in list of image repositories in namespace %s", name, namespace)
}
