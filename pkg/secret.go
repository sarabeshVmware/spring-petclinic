// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"encoding/json"
	"fmt"
	"log"
)

type Secret struct {
	Name     string `yaml:"name"`
	Registry string `yaml:"registry"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type SecretOutput struct {
	Age      string `json:"age"`
	Exported string `json:"exported"`
	Name     string `json:"name"`
	Registry string `json:"registry"`
}

func CreateDockerRegistrySecrets(secrets []Secret, namespace string) {
	for _, secret := range secrets {
		log.Printf("Creating secret: %s", secret.Name)
		Run_DontLogCommand(fmt.Sprintf("kubectl create secret docker-registry %s --docker-server %s --docker-username %s --docker-password %s -n %s",
			secret.Name, secret.Registry, secret.Username, secret.Password, namespace))
	}
}

func CreateImagepullSecrets(secrets []Secret, namespace string) {
	for _, secret := range secrets {
		log.Printf("Creating secret: %s", secret.Name)
		Run_DontLogCommand(fmt.Sprintf("tanzu imagepullsecret add %s --registry %s --username %s --password %s --export-to-all-namespaces -n %s",
			secret.Name, secret.Registry, secret.Username, secret.Password, namespace))
	}
}

func ListImagepullSecrets(namespace string) []SecretOutput {
	var secrets []SecretOutput
	log.Printf("Image Pull Secrets in namespace: %s", namespace)
	secretsList, _ := Run(fmt.Sprintf("tanzu imagepullsecret list -n %s -o json", namespace))
	err := json.Unmarshal(secretsList, &secrets)
	CheckError(err)
	return secrets
}

func DeleteImagepullSecrets(namespace string) {
	addedSecrets := ListImagepullSecrets(namespace)
	for _, secret := range addedSecrets {
		log.Printf("Deleting secret: %s", secret.Name)
		Run(fmt.Sprintf("tanzu imagepullsecret delete %s -n %s -y", secret.Name, namespace))
	}
}
