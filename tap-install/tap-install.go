// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/spf13/cobra"
	tap "gitlab.eng.vmware.com/tap/tap-packaging-tests/pkg"
	tapInstall "gitlab.eng.vmware.com/tap/tap-packaging-tests/tap-install/pkg"
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
	var preCleanup bool
	var postCleanup bool
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install packages",
		Run: func(cmd *cobra.Command, args []string) {
			tapInstall.Install(preCleanup, postCleanup)
		},
	}
	installCmd.Flags().BoolVar(&preCleanup, "pre-cleanup", false, "Cleanup namespace, secrets, repository and packages before installation.")
	installCmd.Flags().BoolVar(&postCleanup, "post-cleanup", false, "Cleanup namespace, secrets, repository and packages after installation.")
	return installCmd
}

func cleanupCommand() *cobra.Command {
	cleanupCmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean packages, secrets, package repositories etc..",
		Run: func(cmd *cobra.Command, args []string) {
			tapInstall.Cleanup()
		},
	}
	return cleanupCmd
}
