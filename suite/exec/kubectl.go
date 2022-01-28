// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
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

func KubectlApplyConfiguration(file string, namespace string) (string, string, error) {
	cmd := fmt.Sprintf("kubectl apply -n %s -f %s", namespace, file)
	output, err := RunCommand(cmd)
	return cmd, output, err
}

func KubectlDeleteConfiguration(file string, namespace string) (string, string, error) {
	cmd := fmt.Sprintf("kubectl delete -n %s -f %s", namespace, file)
	output, err := RunCommand(cmd)
	return cmd, output, err
}

func matchStatusAndType(cmd string, name string, namespace string, conditionType string, conditionStatus string, checkPrefix bool) (bool, error) {
	output, err := RunCommand(cmd)
	if err != nil {
		return false, err
	}
	repositories := struct {
		Items []struct {
			Metadata struct {
				Name string `json:"name"`
			} `json:"metadata"`
			Status struct {
				Conditions []struct {
					Status string `json:"status"`
					Type   string `json:"type"`
				} `json:"conditions"`
			} `json:"status"`
		} `json:"items"`
	}{}
	err = json.Unmarshal([]byte(output), &repositories)
	if err != nil {
		return false, err
	}
	if len(repositories.Items) <= 0 {
		return false, fmt.Errorf("list empty for namespace %s", namespace)
	}
	for _, item := range repositories.Items {
		if item.Metadata.Name == name || (checkPrefix && strings.HasPrefix(item.Metadata.Name, name)) {
			for _, condition := range item.Status.Conditions {
				if condition.Type == conditionType {
					return condition.Status == conditionStatus, nil
				}
			}
			return false, fmt.Errorf("type %s not found in the list of conditions", conditionType)
		}
	}
	return false, fmt.Errorf("%s not found in the list for namespace %s", name, namespace)
}

func KubectlIsImageRepositoryReady(name string, namespace string) (string, bool, error) {
	cmd := fmt.Sprintf("kubectl get imagerepositories -n %s -o json", namespace)
	ready, err := matchStatusAndType(cmd, name, namespace, "Ready", "True", false)
	return cmd, ready, err
}

func KubectlIsGitRepositoryReady(name string, namespace string) (string, bool, error) {
	cmd := fmt.Sprintf("kubectl get gitrepo -n %s -o json", namespace)
	ready, err := matchStatusAndType(cmd, name, namespace, "Ready", "True", false)
	return cmd, ready, err
}

func KubectlIsBuildSucceeded(prefix string, namespace string) (string, bool, error) {
	cmd := fmt.Sprintf("kubectl get build -n %s -o json", namespace)
	succeeded, err := matchStatusAndType(cmd, prefix, namespace, "Succeeded", "True", true)
	return cmd, succeeded, err
}

func KubectlIsPodintentReady(name string, namespace string) (string, bool, error) {
	cmd := fmt.Sprintf("kubectl get podintents -n %s -o json", namespace)
	ready, err := matchStatusAndType(cmd, name, namespace, "Ready", "True", false)
	return cmd, ready, err
}

func KubectlIsKsvcReady(name string, namespace string) (string, bool, error) {
	cmd := fmt.Sprintf("kubectl get ksvc -n %s -o json", namespace)
	ready, err := matchStatusAndType(cmd, name, namespace, "Ready", "True", false)
	return cmd, ready, err
}

func KubectlIsTaskrunSucceeded(prefix string, namespace string) (string, bool, error) {
	cmd := fmt.Sprintf("kubectl get taskruns -n %s -o json", namespace)
	succeeded, err := matchStatusAndType(cmd, prefix, namespace, "Succeeded", "True", true)
	return cmd, succeeded, err
}

func KubectlIsPodintentAnnotationExists(annotationKey string, annotationValue string, checkOnlyKey bool, podintent string, namespace string) (string, bool, error) {
	cmd := fmt.Sprintf("kubectl get podintents -n %s -o json", namespace)
	output, err := RunCommand(cmd)
	if err != nil {
		return cmd, false, err
	}
	exists := false
	_, err = jsonparser.ArrayEach([]byte(output), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		name, _ := jsonparser.GetString(value, "metadata", "name")
		if name == podintent {
			jsonparser.ObjectEach(value, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
				if string(key) == annotationKey && (checkOnlyKey || string(value) == annotationValue) {
					exists = true
				}
				return nil
			}, "status", "template", "metadata", "annotations")
		}
	}, "items")
	if err != nil {
		return cmd, false, err
	}
	return cmd, exists, nil
}

func KubectlIsPodintentLabelExists(labelKey string, labelValue string, checkOnlyKey bool, podintent string, namespace string) (string, bool, error) {
	cmd := fmt.Sprintf("kubectl get podintents -n %s -o json", namespace)
	output, err := RunCommand(cmd)
	if err != nil {
		return cmd, false, err
	}
	exists := false
	_, err = jsonparser.ArrayEach([]byte(output), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		name, _ := jsonparser.GetString(value, "metadata", "name")
		if name == podintent {
			jsonparser.ObjectEach(value, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
				if string(key) == labelKey && (checkOnlyKey || string(value) == labelValue) {
					exists = true
				}
				return nil
			}, "status", "template", "metadata", "labels")
		}
	}, "items")
	if err != nil {
		return cmd, false, err
	}
	return cmd, exists, nil
}
