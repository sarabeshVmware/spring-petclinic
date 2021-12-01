// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package innerloop

import (
	"fmt"
	"log"
	"os"

	tap "gitlab.eng.vmware.com/tap/tap-packages/tap-packaging-tests/pkg"
	e2e "gitlab.eng.vmware.com/tap/tap-packages/tap-packaging-tests/tap-tests/e2e"
	"gopkg.in/yaml.v3"
)

func InnerloopSourceBuildDeploy(tapPackageInfo tap.Package) {
	log.Printf("Testing innerloop (source build deploy):")

	namespace := tapPackageInfo.Namespace
	scRepository, scServer := GetSCRegistryDetails(tapPackageInfo.ValuesFile)
	acceleratorProject := "tanzu-java-web-app"
	workload := "tanzu-java-web-app"
	sourceImage := fmt.Sprintf("%s-src", workload)
	oldString := "Greetings from Spring Boot + Tanzu!"
	newString := "Greetings from Spring Boot + TAP!"

	tap.SetupDeveloperNamespacePostInstallation(namespace)

	// tap.RunWithBash(`ps aux | grep -i kubectl | grep -v grep | awk {'print $2'} | xargs kill`)
	// pidAppAcceleratorPortForward, _ := tap.RunAndDisown("kubectl -n accelerator-system port-forward svc/acc-ui-server 8877:80")
	// time.Sleep(5 * time.Second)
	// defer tap.KillPID(pidAppAcceleratorPortForward)
	// pidAppLiveViewPortForward, _ := tap.RunAndDisown("kubectl -n tap-install port-forward service/application-live-view-5112 5112:5112")
	// time.Sleep(5 * time.Second)
	// defer tap.KillPID(pidAppLiveViewPortForward)

	appAccExternalIP := e2e.GetAppAcceleratorExternalIP()

	// Setting Env variable ACC_SERVER_URL
	tap.RunWithBash(fmt.Sprintf("export ACC_SERVER_URL=http://%s", appAccExternalIP))

	e2e.GetAppLiveViewExternalIP()

	e2e.ListAccelerators()
	e2e.GenerateAcceleratorProject(acceleratorProject, acceleratorProject, scServer, true, appAccExternalIP)

	tap.SetupDeveloperNamespacePostInstallation(namespace)
	e2e.DeleteWorkload(workload, namespace)
	e2e.CreateWorkload(workload, scServer, scRepository, sourceImage, acceleratorProject, namespace)

	verify(workload, namespace, oldString, newString, false)

	e2e.UpdateFile("tanzu-java-web-app/src/main/java/com/example/springboot/HelloController.java", oldString, newString)
	e2e.UpdateFile("tanzu-java-web-app/src/test/java/com/example/springboot/HelloControllerTest.java", oldString, newString)
	e2e.UpdateWorkload(workload, scServer, scRepository, sourceImage, acceleratorProject, namespace)

	verify(workload, namespace, oldString, newString, true)

	log.Printf("Innerloop (source build deploy) successful.")
}

func verify(workload string, namespace string, oldString string, newString string, testNew bool) {
	e2e.VerifyImageRepositoryReadyStatus(workload, namespace)
	// e2e.VerifyBuildStatus(fmt.Sprintf("%s-build-*", workload), namespace, true)
	e2e.VerifyBuildStatus()
	e2e.VerifyKnativeServiceStatus(workload, namespace)
	e2e.VerifyWorkloadStatus(workload, namespace)
	e2e.VerifyApplicationRunningWithValidationString(e2e.GetEnvoyExternalIP(), "tanzu-java-web-app.tap-install.example.com", oldString, newString, testNew)
}

func GetSCRegistryDetails(valuesFile string) (string, string) {
	type supplyChainSchema struct {
		SupplyChain          string `yaml:"supply_chain"`
		OotbSupplyChainBasic struct {
			Registry struct {
				Server     string `yaml:"server"`
				Repository string `yaml:"repository"`
			} `yaml:"registry"`
		} `yaml:"ootb_supply_chain_basic"`
		OotbSupplyChainTesting struct {
			Registry struct {
				Server     string `yaml:"server"`
				Repository string `yaml:"repository"`
			} `yaml:"registry"`
		} `yaml:"ootb_supply_chain_testing"`
		OotbSupplyChainTestingScanning struct {
			Registry struct {
				Server     string `yaml:"server"`
				Repository string `yaml:"repository"`
			} `yaml:"registry"`
		} `yaml:"ootb_supply_chain_testing_scanning"`
	}
	var scSchema supplyChainSchema
	var repository string
	var server string

	supplyChainSchemaBytes, err := os.ReadFile(valuesFile)
	tap.CheckError(err)
	err = yaml.Unmarshal(supplyChainSchemaBytes, &scSchema)
	tap.CheckError(err)
	if scSchema.SupplyChain == "" || scSchema.SupplyChain == "basic" {
		scregistry := scSchema.OotbSupplyChainBasic.Registry
		repository = scregistry.Repository
		server = scregistry.Server
	} else if scSchema.SupplyChain == "testing" {
		scregistry := scSchema.OotbSupplyChainTesting.Registry
		repository = scregistry.Repository
		server = scregistry.Server
	} else if scSchema.SupplyChain == "scanning" {
		scregistry := scSchema.OotbSupplyChainTestingScanning.Registry
		repository = scregistry.Repository
		server = scregistry.Server
	} else {
		log.Println("Invalid Supply chain schema in values.yaml file")
	}

	log.Printf("Supply chain: repository %s, server %s", repository, server)
	return repository, server
}
