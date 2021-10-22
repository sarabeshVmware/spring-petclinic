// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tapinstall

import (
	"os"

	tap "gitlab.eng.vmware.com/tap/tap-packaging-tests/pkg"
	"gopkg.in/yaml.v3"
)

type Input struct {
	Namespace         string                `yaml:"namespace"`
	Secrets           []tap.Secret          `yaml:"secrets"`
	PackageRepository tap.PackageRepository `yaml:"package_repository"`
	Packages          []tap.Package         `yaml:"packages"`
	ValuesDirectory   string
}

func GetInput(configFile string, valuesDir string) Input {
	inputBytes, err := os.ReadFile(configFile)
	tap.CheckError(err)
	input := Input{}
	err = yaml.Unmarshal(inputBytes, &input)
	tap.CheckError(err)
	input.ValuesDirectory = valuesDir
	return input
}
