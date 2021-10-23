// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"fmt"
	"log"
	"os"
)

func CreateClusterRoleBinding(name string, clusterRole string, serviceAccount string) {
	log.Printf("Creating cluster role binding: %s", name)
	Run(fmt.Sprintf("kubectl create clusterrolebinding %s --clusterrole=%s --serviceaccount=%s", name, clusterRole, serviceAccount))
}

func DeployApp(name string, files []string, namespace string) {
	log.Printf("Deploying app: %s", name)
	cmd := fmt.Sprintf("kapp deploy -a %s -n %s -y", name, namespace)
	for _, file := range files {
		cmd += fmt.Sprintf(" -f %s", file)
	}
	Run(cmd)
}

func ApplyConfiguration(file string) {
	configuration, err := os.ReadFile(file)
	CheckError(err)
	log.Printf("Applying configuration: \n%s", string(configuration))
	Run(fmt.Sprintf("kubectl apply -f %s", file))
}

func ApplyConfigurationInNamespace(file string, namespace string) {
	configuration, err := os.ReadFile(file)
	CheckError(err)
	log.Printf("Applying configuration in namespace %s: \n%s", namespace, string(configuration))
	Run(fmt.Sprintf("kubectl apply -n %s -f %s", namespace, file))
}
