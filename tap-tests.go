// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"path/filepath"

	"github.com/spf13/cobra"
	tap "gitlab.eng.vmware.com/tap/tap-packaging-tests/pkg"
	tapInstall "gitlab.eng.vmware.com/tap/tap-packaging-tests/tap-install/tapinstall"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "tap-tests",
		Short: "TAP packaging tests CLI",
	}
	rootCmd.AddCommand(installCommand())
	rootCmd.AddCommand(cleanupCommand())
	err := rootCmd.Execute()
	tap.CheckError(err)
}

func installCommand() *cobra.Command {
	var preCleanup, postCleanup bool
	var configFile, valuesDir string
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install packages",
		Run: func(cmd *cobra.Command, args []string) {
			tapInstall.Install(configFile, valuesDir, preCleanup, postCleanup)
		},
	}
	installCmd.Flags().BoolVar(&preCleanup, "pre-cleanup", false, "Cleanup namespace, secrets, repository and packages before installation.")
	installCmd.Flags().BoolVar(&postCleanup, "post-cleanup", false, "Cleanup namespace, secrets, repository and packages after installation.")
	installCmd.Flags().StringVarP(&configFile, "config-file", "f", filepath.Join(tap.GetCurrentDir(), "tap-install", "user-config.yaml"), "User configuration YAML file.")
	installCmd.Flags().StringVarP(&valuesDir, "values-dir", "v", filepath.Join(tap.GetCurrentDir(), "tap-install", "values"), "Directory containing values schemas.")
	return installCmd
}

func cleanupCommand() *cobra.Command {
	var configFile, valuesDir string
	cleanupCmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean packages, secrets, package repositories etc..",
		Run: func(cmd *cobra.Command, args []string) {
			tapInstall.Cleanup(configFile, valuesDir)
		},
	}
	cleanupCmd.Flags().StringVarP(&configFile, "config-file", "f", filepath.Join(tap.GetCurrentDir(), "tap-install", "user-config.yaml"), "User configuration YAML file.")
	cleanupCmd.Flags().StringVarP(&valuesDir, "values-dir", "v", filepath.Join(tap.GetCurrentDir(), "tap-install", "values"), "Directory containing values schemas.")
	return cleanupCmd
}
