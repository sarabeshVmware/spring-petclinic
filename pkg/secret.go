// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"encoding/json"
	"fmt"
	"log"
)

type Secret struct {
	Name      string `yaml:"name"`
	Registry  string `yaml:"registry"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Namespace string `yaml:"namespace"`
}

type SecretOutput struct {
	Age      string `json:"age"`
	Exported string `json:"exported"`
	Name     string `json:"name"`
	Registry string `json:"registry"`
}

func CreateDockerRegistrySecret(secret Secret) {
	log.Printf("Creating secret: %s", secret.Name)
	Run_DontLogCommand(fmt.Sprintf("kubectl create secret docker-registry %s --docker-server %s --docker-username %s --docker-password %s -n %s",
		secret.Name, secret.Registry, secret.Username, secret.Password, secret.Namespace))
}

func CreateTanzuSecret(secret Secret) {
	log.Printf("Creating secret: %s", secret.Name)
	Run_DontLogCommand(fmt.Sprintf("tanzu secret registry add %s --server %s --username %s --password %s --export-to-all-namespaces -n %s -y",
		secret.Name, secret.Registry, secret.Username, secret.Password, secret.Namespace))
}

func ListTanzuSecrets(namespace string) []SecretOutput {
	var secrets []SecretOutput
	log.Printf("Secrets in namespace: %s", namespace)
	secretsList, _ := Run(fmt.Sprintf("tanzu secret registry list -n %s -o json", namespace))
	err := json.Unmarshal(secretsList, &secrets)
	CheckError(err)
	return secrets
}

func DeleteTanzuSecrets(namespace string) {
	addedSecrets := ListTanzuSecrets(namespace)
	for _, secret := range addedSecrets {
		log.Printf("Deleting secret: %s", secret.Name)
		Run(fmt.Sprintf("tanzu secret registry delete %s -n %s -y", secret.Name, namespace))
	}
}
