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

	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
	tap "gitlab.eng.vmware.com/tap/tap-packages/tap-packaging-tests/pkg"
	innerloop "gitlab.eng.vmware.com/tap/tap-packages/tap-packaging-tests/tap-tests/e2e/innerloop"
)

func setLogger() {
	os.MkdirAll("logs", 0755)
	logFilePath := filepath.Join("logs", fmt.Sprintf("log_%s.log", time.Now().Format(time.RFC3339Nano)))
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	tap.CheckError(err)
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func getPackagesList() []tap.Package {
	packagesFileBytes, err := os.ReadFile(filepath.Join(tap.GetCurrentDir(), "packages.yaml"))
	tap.CheckError(err)
	packagesList := []tap.Package{}
	err = yaml.Unmarshal(packagesFileBytes, &packagesList)
	tap.CheckError(err)
	return packagesList
}

func main() {
	setLogger()
	rootCmd := &cobra.Command{
		Use:   "tap-tests",
		Short: "TAP packaging tests CLI",
	}
	rootCmd.AddCommand(createCommand())
	rootCmd.AddCommand(installCommand())
	rootCmd.AddCommand(cleanupCommand())
	rootCmd.AddCommand(e2eCommand())
	err := rootCmd.Execute()
	tap.CheckError(err)
}

func createCommand() *cobra.Command {
	var createResourcesFile string
	var installPrerequisites bool
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create resources",
		Run: func(cmd *cobra.Command, args []string) {
			if installPrerequisites {
				tap.HandlePrerequisites()
			}

			log.Printf("Creating resources from file: %s", createResourcesFile)

			createResourcesBytes, err := os.ReadFile(createResourcesFile)
			tap.CheckError(err)
			createResources := struct {
				Namespaces          []string                `yaml:"namespaces"`
				Secrets             []tap.Secret            `yaml:"secrets"`
				PackageRepositories []tap.PackageRepository `yaml:"package_repositories"`
			}{}
			err = yaml.Unmarshal(createResourcesBytes, &createResources)
			tap.CheckError(err)

			for _, namespace := range createResources.Namespaces {
				tap.CreateNamespace(namespace)
			}
			for _, secret := range createResources.Secrets {
				tap.CreateTanzuSecret(secret)
			}
			for _, packageRepository := range createResources.PackageRepositories {
				tap.AddPackageRepository(packageRepository)
			}
		},
	}
	createCmd.Flags().StringVarP(&createResourcesFile, "create-resources-file", "f", filepath.Join(tap.GetCurrentDir(), "create-resources.yaml"), "Create resources YAML file.")
	createCmd.Flags().BoolVar(&installPrerequisites, "install-prerequisites", false, "Install prerequisites such as kapp-controller, secretgen-controller, cert-manager, flux-system, etc.")
	return createCmd
}

func installCommand() *cobra.Command {
	packagesList := getPackagesList()

	packagesNames := []string{}
	for _, packageInfo := range packagesList {
		packagesNames = append(packagesNames, packageInfo.Name)
	}

	/*
		    NOTE:
			- Defining sub-commands (via a for-loop) for installing different packages doesn't work because the cobra.Command variable passed to the
			  AddCommand() function has the same address for all the packages and hence, it registers the same InstallPackage() call for all the sub-commands.
			- That is, all sub-commands will install the last package in packages.yaml (workshops.learningcenter in our case).
			- Passing the cobra.Command address without using a variable also didn't work, and so didn't creating a map for storing different cobra.Command variables. Just how golang works!
			- We are thus using args instead of separate sub-commands. It also allows us to do multiple packages installation at once.
	*/

	var valuesFile, developerNamespace string
	installCmd := &cobra.Command{
		Use:       "install package1 [package2, ..]",
		Short:     "Install package",
		Args:      cobra.MinimumNArgs(1),
		ValidArgs: packagesNames, // NOTE: Requires https://github.com/spf13/cobra/pull/841 to be merged to function properly
		Run: func(cmd *cobra.Command, args []string) {
			for _, packageName := range args {
				packageInfo := tap.GetPackageInfoFromName(packageName, packagesList)
				if valuesFile != "" {
					packageInfo.ValuesFile = valuesFile
				}
				tap.InstallPackageByInfo(packageInfo, packagesList)
			}
			if developerNamespace != "" {
				tap.SetupDeveloperNamespacePostInstallation(developerNamespace)
			}
		},
	}
	installCmd.Flags().StringVarP(&valuesFile, "values-file", "f", "", "Values schema YAML file.")
	installCmd.Flags().StringVar(&developerNamespace, "developer-namespace", "", "Setup developer namespace.")
	return installCmd
}

func cleanupCommand() *cobra.Command {
	cleanupCmd := &cobra.Command{
		Use:   "clean namespace1 [namespace2, ..]",
		Short: "Clean packages, secrets, package repositories etc. from namespace(s)",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, namespace := range args {
				log.Printf("Cleanup requested from namespace: %s", namespace)
				tap.UninstallPackages(namespace)
				tap.DeletePackageRepository(namespace)
				tap.DeleteTanzuSecrets(namespace)
				tap.DeleteNamespace(namespace)
			}
		},
	}
	return cleanupCmd
}

func e2eCommand() *cobra.Command {
	packagesList := getPackagesList()

	e2eCmd := &cobra.Command{
		Use:   "e2e",
		Short: "End-to-end testing",
	}

	var install bool
	innerloopCmd := &cobra.Command{
		Use:   "innerloop",
		Short: "End-to-end innerloop test",
		Run: func(cmd *cobra.Command, args []string) {
			if install {
				tap.InstallPackageByName("tap-dev-light", packagesList)
				tap.InstallPackageByName("app-accelerator", packagesList)
			}
			tapPackageInfo := tap.GetPackageInfoFromName("tap-dev-light", packagesList)
			innerloop.InnerloopSourceBuildDeploy(tapPackageInfo)
		},
	}
	innerloopCmd.Flags().BoolVar(&install, "install", false, "Install packages pre-testing.")

	e2eCmd.AddCommand(innerloopCmd)
	return e2eCmd
}
