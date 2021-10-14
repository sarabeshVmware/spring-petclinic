// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Input struct {
	Namespace         string            `yaml:"namespace"`
	Secrets           []Secret          `yaml:"secrets"`
	PackageRepository PackageRepository `yaml:"package_repository"`
	Packages          []Package         `yaml:"packages"`
}

func GetInput() Input {
	inputBytes, err := os.ReadFile(filepath.Join(GetProjectDirectory(), "user_input.yaml"))
	CheckError(err)
	input := Input{}
	err = yaml.Unmarshal(inputBytes, &input)
	CheckError(err)
	return input
}

func Cleanup() {
	input := GetInput()
	log.Printf("Request: Cleanup")
	UninstallPackages(input.Namespace)
	DeletePackageRepository(input.Namespace)
	DeleteImagepullSecrets(input.Namespace)
	DeleteNamespace(input.Namespace)
}

func Install(preCleanup bool, postCleanup bool) {
	log.Printf("Request: Install packages")
	input := GetInput()

	log.Printf("Request: Cleanup pre-installation (%t)", preCleanup)
	if preCleanup {
		Cleanup()
	}

	CreateNamespace(input.Namespace)
	CreateImagepullSecrets(input.Secrets, input.Namespace)
	AddPackageRepository(input.PackageRepository, input.Namespace)
	CheckPackageRepositoryStatus(input.PackageRepository, input.Namespace)
	ListPackages(input.Namespace)
	// ListValuesSchema(input.Packages, input.Namespace)
	InstallPackages(input.Packages, input.Namespace)

	log.Printf("Request: Cleanup post-installation (%t)", postCleanup)
	if postCleanup {
		Cleanup()
	}
}
