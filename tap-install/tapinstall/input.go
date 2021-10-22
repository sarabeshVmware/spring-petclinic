// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tapinstall

import (
	"os"
	"path/filepath"

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

func GetInput() Input {
	tapInstallDir := filepath.Dir(tap.GetCurrentDir())
	inputBytes, err := os.ReadFile(filepath.Join(tapInstallDir, "user_input.yaml"))
	tap.CheckError(err)
	input := Input{}
	err = yaml.Unmarshal(inputBytes, &input)
	tap.CheckError(err)
	input.ValuesDirectory = filepath.Join(tapInstallDir, "values")
	return input
}
