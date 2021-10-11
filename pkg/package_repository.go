// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"encoding/json"
	"log"
	"time"

	"github.com/buger/jsonparser"
)

type PackageRepository struct {
	Name  string `yaml:"name"`
	Image string `yaml:"image"`
}

type PackageRepoOutput struct {
	Details    string `json:"details"`
	Name       string `json:"name"`
	Repository string `json:"repository"`
	Status     string `json:"status"`
}

func AddPackageRepository(packageRepository PackageRepository, namespace string) {
	log.Printf("Adding package repository CR: %s", packageRepository.Name)
	RunCommand(Command{CommandName: "tanzu", Arguments: []string{"package", "repository", "add", packageRepository.Name, "--url", packageRepository.Image, "-n", namespace}})
}

func ListPackageRepositories(namespace string) []PackageRepoOutput {
	var addedPkgrs []PackageRepoOutput
	log.Printf("Retriving Package repository in namespace: %s", namespace)
	repoList, _ := RunCommand(Command{CommandName: "tanzu", Arguments: []string{"package", "repository", "list", "-n", namespace, "-ojson"}})
	err := json.Unmarshal(repoList, &addedPkgrs)
	CheckError(err)
	return addedPkgrs
}
func DeletePackageRepository(namespace string) {
	addedPkgr := ListPackageRepositories(namespace)[0]
	log.Printf("Deleting package repository: %s", addedPkgr.Name)
	RunCommand(Command{CommandName: "tanzu", Arguments: []string{"package", "repository", "delete", addedPkgr.Name, "-n", namespace}})
}

func CheckPackageRepositoryStatus(packageRepository PackageRepository, namespace string) {
	log.Printf("Checking package repository status: %s", packageRepository.Name)
	packageRepositoryStatus, _ := RunCommand(Command{CommandName: "tanzu", Arguments: []string{"package", "repository", "get", packageRepository.Name, "-n", namespace, "-o", "json"}})
	jsonparser.ArrayEach(packageRepositoryStatus, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		repository, err := jsonparser.GetString(value, "repository")
		CheckError(err)
		if repository == packageRepository.Image {
			status, err := jsonparser.GetString(value, "status")
			CheckError(err)
			if status == "Reconciling" {
				time.Sleep(5 * time.Second)
				CheckPackageRepositoryStatus(packageRepository, namespace)
			} else if status == "Reconcile succeeded" {
				log.Printf("Reconcile succeeded for package repository: %s", packageRepository.Name)
			} else {
				log.Fatalf("Reconcile not succeeded for package repository: %s", packageRepository.Name)
			}
		}
	})
}
