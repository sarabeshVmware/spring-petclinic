// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

type Package struct {
	Description         string   `yaml:"description"`
	Name                string   `yaml:"name"`
	Namespace           string   `yaml:"namespace"`
	Package             string   `yaml:"package"`
	ValuesFile          string   `yaml:"values_file,omitempty"`
	Version             string   `yaml:"version"`
	PackageDependencies []string `yaml:"package_dependencies,omitempty"`
}

type PackageInstalledOutput struct {
	Name           string `json:"name"`
	PackageName    string `json:"package-name"`
	PackageVersion string `json:"package-version"`
	Status         string `json:"status"`
}
