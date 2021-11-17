// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/buger/jsonparser"
)

type PackageRepository struct {
	Name      string `yaml:"name"`
	Image     string `yaml:"image"`
	Namespace string `yaml:"namespace"`
}

type PackageRepoOutput struct {
	Details    string `json:"details"`
	Name       string `json:"name"`
	Repository string `json:"repository"`
	Status     string `json:"status"`
}

func AddPackageRepository(packageRepository PackageRepository) {
	log.Printf("Adding package repository CR: %s", packageRepository.Name)
	Run(fmt.Sprintf("tanzu package repository add %s --url %s -n %s", packageRepository.Name, packageRepository.Image, packageRepository.Namespace))
	CheckPackageRepositoryStatus(packageRepository)
}

func ListPackageRepositories(namespace string) []PackageRepoOutput {
	var addedPkgrs []PackageRepoOutput
	log.Printf("Retriving Package repository in namespace: %s", namespace)
	repoList, _ := Run(fmt.Sprintf("tanzu package repository list -n %s -o json", namespace))
	err := json.Unmarshal(repoList, &addedPkgrs)
	CheckError(err)
	return addedPkgrs
}

func DeletePackageRepository(namespace string) {
	addedPkgr := ListPackageRepositories(namespace)
	if len(addedPkgr) != 0 {
		log.Printf("Deleting package repository: %s", addedPkgr[0].Name)
		Run(fmt.Sprintf("tanzu package repository delete %s -n %s", addedPkgr[0].Name, namespace))
	}
}

func CheckPackageRepositoryStatus(packageRepository PackageRepository) {
	log.Printf("Checking package repository status: %s", packageRepository.Name)
	packageRepositoryStatus, _ := Run(fmt.Sprintf("tanzu package repository get %s -n %s -o json", packageRepository.Name, packageRepository.Namespace))
	jsonparser.ArrayEach(packageRepositoryStatus, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		repository, err := jsonparser.GetString(value, "repository")
		CheckError(err)
		if repository == packageRepository.Image {
			status, err := jsonparser.GetString(value, "status")
			CheckError(err)
			if status == "Reconciling" || status == "" {
				time.Sleep(5 * time.Second)
				CheckPackageRepositoryStatus(packageRepository)
			} else if status == "Reconcile succeeded" {
				log.Printf("Reconcile succeeded for package repository: %s", packageRepository.Name)
			} else {
				log.Fatalf("Reconcile not succeeded for package repository: %s", packageRepository.Name)
			}
		}
	})
}
