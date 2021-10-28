// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	tap "gitlab.eng.vmware.com/tap/tap-packaging-tests/pkg"
	tapInstall "gitlab.eng.vmware.com/tap/tap-packaging-tests/tap-install/tapinstall"
	innerloop "gitlab.eng.vmware.com/tap/tap-packaging-tests/tap-tests/e2e/innerloop"
)

func setLogger() {
	os.MkdirAll("logs", 0755)
	logFilePath := filepath.Join("logs", fmt.Sprintf("log_%s.log", time.Now().Format(time.RFC3339Nano)))
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	tap.CheckError(err)
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func main() {
	setLogger()
	rootCmd := &cobra.Command{
		Use:   "tap-tests",
		Short: "TAP packaging tests CLI",
	}
	rootCmd.AddCommand(installCommand())
	rootCmd.AddCommand(cleanupCommand())
	rootCmd.AddCommand(e2eCommand())
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

func e2eCommand() *cobra.Command {
	var innerloopSourceBuildDeploy, installPackages, preCleanup, postCleanup bool
	e2eCmd := &cobra.Command{
		Use:   "e2e",
		Short: "End-to-end testing",
		Run: func(cmd *cobra.Command, args []string) {
			if innerloopSourceBuildDeploy {
				innerloop.InnerloopSourceBuildDeploy(installPackages, preCleanup, postCleanup)
			}
		},
	}
	e2eCmd.Flags().BoolVar(&innerloopSourceBuildDeploy, "innerloop-source-build-deploy", true, "Test innerloop: source build deploy.")
	e2eCmd.Flags().BoolVar(&installPackages, "install", false, "Install packages pre-testing.")
	e2eCmd.Flags().BoolVar(&preCleanup, "pre-cleanup", false, "Cleanup namespace, secrets, repository and packages before installation.")
	e2eCmd.Flags().BoolVar(&postCleanup, "post-cleanup", false, "Cleanup namespace, secrets, repository and packages after testing.")
	return e2eCmd
}
