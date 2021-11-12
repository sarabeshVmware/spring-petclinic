// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tapinstall

import (
	"os"
	"path/filepath"

	tap "gitlab.eng.vmware.com/tap/tap-packaging-tests/pkg"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Namespace         string                `yaml:"namespace"`
	Secrets           []tap.Secret          `yaml:"secrets"`
	PackageRepository tap.PackageRepository `yaml:"package_repository"`
	Packages          []tap.Package         `yaml:"packages"`
	ValuesDirectory   string
}

func GetConfig(configFile string, valuesDir string) Config {
	configBytes, err := os.ReadFile(configFile)
	tap.CheckError(err)
	config := Config{}
	err = yaml.Unmarshal(configBytes, &config)
	tap.CheckError(err)
	config.ValuesDirectory = valuesDir
	return config
}

func GetDefaultValuesDir() string {
	return filepath.Join(filepath.Dir(tap.GetCurrentDir()), "values")
}
