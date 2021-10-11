// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"encoding/json"
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
		RunCommand(Command{CommandName: "kubectl", Arguments: []string{"create", "secret", "docker-registry", secret.Name, "-n", namespace,
			"--docker-server", secret.Registry, "--docker-username", secret.Username, "--docker-password", secret.Password}, DontLogCommand: true})
	}
}

func CreateImagepullSecrets(secrets []Secret, namespace string) {
	for _, secret := range secrets {
		log.Printf("Creating secret: %s", secret.Name)
		RunCommand(Command{CommandName: "tanzu", Arguments: []string{"imagepullsecret", "add", secret.Name,
			"--registry", secret.Registry, "--username", secret.Username, "--password", secret.Password,
			"--export-to-all-namespaces", "-n", namespace}, DontLogCommand: true})
	}
}

func ListImagepullSecrets(namespace string) []SecretOutput {
	var secrets []SecretOutput
	log.Printf("Image Pull Secrets in namespace: %s", namespace)
	secretsList, _ := RunCommand(Command{CommandName: "tanzu", Arguments: []string{"imagepullsecret", "list", "-n", namespace, "-ojson"}})
	err := json.Unmarshal(secretsList, &secrets)
	CheckError(err)
	return secrets
}

func DeleteImagepullSecrets(namespace string) {
	addedSecrets := ListImagepullSecrets(namespace)
	for _, secret := range addedSecrets {
		log.Printf("Deleting secret: %s", secret.Name)
		RunCommand(Command{CommandName: "tanzu", Arguments: []string{"imagepullsecret", "delete", secret.Name, "-n", namespace, "-y"}})
	}
}
