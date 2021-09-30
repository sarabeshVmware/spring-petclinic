// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import "log"

type Secret struct {
	Name     string `yaml:"name"`
	Registry string `yaml:"registry"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
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

func DeleteImagepullSecrets(secrets []Secret, namespace string) {
	for _, secret := range secrets {
		log.Printf("Deleting secret: %s", secret.Name)
		RunCommand(Command{CommandName: "tanzu", Arguments: []string{"imagepullsecret", "delete", secret.Name, "-n", namespace, "-y"}})
	}
}
