// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package innerloop

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	tap "gitlab.eng.vmware.com/tap/tap-packaging-tests/pkg"
	tapInstall "gitlab.eng.vmware.com/tap/tap-packaging-tests/tap-install/tapinstall"
	e2e "gitlab.eng.vmware.com/tap/tap-packaging-tests/tap-tests/e2e"
	"gopkg.in/yaml.v3"
)

func InnerloopSourceBuildDeploy(installPackages bool, preCleanup bool, postCleanup bool) {
	log.Printf("Testing innerloop (source build deploy):")
	configFile := filepath.Join(tap.GetCurrentDir(), "source-build-deploy.yaml")
	valuesDir := tapInstall.GetDefaultValuesDir()

	config := tapInstall.GetConfig(configFile, valuesDir)

	defaultSupplyChainSchema := struct {
		Registry struct {
			Repository string `yaml:"repository"`
			Server     string `yaml:"server"`
		} `yaml:"registry"`
	}{}
	defaultSupplyChainSchemaBytes, err := os.ReadFile(filepath.Join(valuesDir, "default-supply-chain.yaml"))
	tap.CheckError(err)
	err = yaml.Unmarshal(defaultSupplyChainSchemaBytes, &defaultSupplyChainSchema)
	tap.CheckError(err)

	acceleratorProject := "tanzu-java-web-app"
	workload := "tanzu-java-web-app"
	sourceImage := fmt.Sprintf("%s-src", workload)
	oldString := "Greetings from Spring Boot + Tanzu!"
	newString := "Greetings from Spring Boot + TAP!"

	if installPackages {
		tapInstall.Install(configFile, valuesDir, preCleanup, false)
	}

	tap.RunWithBash(`ps aux | grep -i kubectl | grep -v grep | awk {'print $2'} | xargs kill`)
	pidAppAcceleratorPortForward, _ := tap.RunAndDisown("kubectl -n accelerator-system port-forward svc/acc-ui-server 8877:80")
	time.Sleep(5 * time.Second)
	defer tap.KillPID(pidAppAcceleratorPortForward)
	pidAppLiveViewPortForward, _ := tap.RunAndDisown("kubectl -n tap-install port-forward service/application-live-view-5112 5112:5112")
	time.Sleep(5 * time.Second)
	defer tap.KillPID(pidAppLiveViewPortForward)

	e2e.GetAppAcceleratorExternalIP()
	e2e.GetAppLiveViewExternalIP()

	e2e.ListAccelerators()
	e2e.GenerateAcceleratorProject(acceleratorProject, acceleratorProject, defaultSupplyChainSchema.Registry.Server, true)

	e2e.DeleteWorkload(workload, config.Namespace)
	e2e.CreateWorkload(workload, defaultSupplyChainSchema.Registry.Server, defaultSupplyChainSchema.Registry.Repository, sourceImage, acceleratorProject, config.Namespace)

	verify(workload, config.Namespace, oldString, newString, false)

	e2e.UpdateFile("tanzu-java-web-app/src/main/java/com/example/springboot/HelloController.java", oldString, newString)
	e2e.UpdateFile("tanzu-java-web-app/src/test/java/com/example/springboot/HelloControllerTest.java", oldString, newString)
	e2e.UpdateWorkload(workload, defaultSupplyChainSchema.Registry.Server, defaultSupplyChainSchema.Registry.Repository, sourceImage, acceleratorProject, config.Namespace)

	verify(workload, config.Namespace, oldString, newString, true)

	log.Printf("Innerloop (source build deploy) successful.")

	if postCleanup {
		tapInstall.Cleanup(configFile, valuesDir)
	}
}

func verify(workload string, namespace string, oldString string, newString string, testNew bool) {
	e2e.VerifyImageRepositoryReadyStatus(workload, namespace)
	e2e.VerifyBuildStatus(fmt.Sprintf("%s-build-*", workload), namespace, true)
	e2e.VerifyKnativeServiceStatus(workload, namespace)
	e2e.VerifyWorkloadStatus(workload, namespace)
	e2e.VerifyApplicationRunningWithValidationString(e2e.GetEnvoyExternalIP(), "tanzu-java-web-app.tap-install.example.com", oldString, newString, testNew)
}
