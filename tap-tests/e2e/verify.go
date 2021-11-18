// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e


import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/pivotal/kpack/pkg/apis/build/v1alpha2"
	tap "gitlab.eng.vmware.com/tap/tap-packaging-tests/pkg"
	//util "gitlab.eng.vmware.com/tap/tap-packaging-tests/pkg/util"
	"context"
	corev1 "k8s.io/api/core/v1"
)

func VerifyImageRepositoryReadyStatus(name string, namespace string) {
	log.Printf("Checking image repository status for: %s", name)
	imageRepositoryStatusBytes, _ := tap.RunWithBash(fmt.Sprintf(`kubectl get imagerepositories -n %s | awk '{if($1=="%s")print $4}'`, namespace, name))
	imageRepositoryStatus := strings.TrimSpace(string(imageRepositoryStatusBytes))
	log.Printf("Image repository %s state: %s", name, imageRepositoryStatus)
	if imageRepositoryStatus == "True" {
		log.Printf("Image repository ready: %s", name)
	} else if imageRepositoryStatus == "False" {
		log.Fatalf("Image repository is not ready: %s", name)
	} else {
		time.Sleep(5 * time.Second)
		VerifyImageRepositoryReadyStatus(name, namespace)
	}
}

// func verifyBuildsAreGenerated(namespace string) {
// 	log.Printf("Checking if builds are generated in namespace: %s", namespace)
// 	buildsBytes, _ := tap.Run(fmt.Sprintf("kubectl get build -n %s", namespace))
// 	builds := strings.TrimSpace(string(buildsBytes))
// 	if builds == fmt.Sprintf("No resources found in %s namespace.", namespace) {
// 		log.Printf("Builds not generated in namespace: %s", namespace)
// 		time.Sleep(5 * time.Second)
// 		verifyBuildsAreGenerated(namespace)
// 	} else {
// 		log.Printf("Builds generated in namespace: %s", namespace)
// 	}
// }

// func verifyBuildIsNotInUnknownState(buildName string, namespace string) {
// 	log.Printf("Checking if build is in unknown state: %s", buildName)
// 	buildStateBytes, _ := tap.RunWithBash(fmt.Sprintf(`kubectl get build -n %s | awk '{if($1=="%s")print $2}'`, namespace, buildName))
// 	buildState := strings.TrimSpace(string(buildStateBytes))
// 	if buildState == "Unknown" {
// 		log.Printf("Build in unknown state: %s", buildName)
// 		time.Sleep(5 * time.Second)
// 		verifyBuildIsNotInUnknownState(buildName, namespace)
// 	} else {
// 		log.Printf("Build not in unknown state: %s", buildName)
// 		log.Printf("Build image: %s", buildState)
// 	}
// }

// func VerifyBuildStatus(pattern string, namespace string, verifyBuildsGeneration bool) {
// 	log.Printf("Checking build status for pattern: %s", pattern)
// 	if verifyBuildsGeneration {
// 		verifyBuildsAreGenerated(namespace)
// 	}

// 	log.Printf("Getting build names matching pattern: %s", pattern)
// 	buildNamesBytes, _ := tap.RunWithBash(fmt.Sprintf(`kubectl get build -n %s | awk '{if($1~"%s")print $1}' | xargs`, namespace, pattern))
// 	buildNames := strings.Split(strings.TrimSpace(string(buildNamesBytes)), " ")
// 	allBuildsFailed := true
// 	for _, buildName := range buildNames {
// 		log.Printf("Checking build status for build: %s", buildName)
// 		verifyBuildIsNotInUnknownState(buildName, namespace)
// 		buildStateBytes, _ := tap.RunWithBash(fmt.Sprintf(`kubectl get build -n %s | awk '{if($1=="%s")print $3}'`, namespace, buildName))
// 		buildState := strings.TrimSpace(string(buildStateBytes))
// 		log.Printf("Build %s status: %s", buildName, buildState)
// 		if buildState == "True" {
// 			log.Printf("Build ready: %s", buildName)
// 			allBuildsFailed = false
// 			break
// 		}
// 	}
// 	if allBuildsFailed {
// 		log.Fatalf("No builds ready for the pattern: %s", pattern)
// 	}
// }

func VerifyKnativeServiceStatus(name string, namespace string) {
	log.Printf("Checking knative service status for: %s", name)
	ksvcStatusBytes, _ := tap.RunWithBash(fmt.Sprintf(`kubectl get ksvc -n %s | awk '{if($1=="%s")print $5}'`, namespace, name))
	ksvcStatus := strings.TrimSpace(string(ksvcStatusBytes))
	log.Printf("Knative service %s status: %s", name, ksvcStatus)
	if ksvcStatus == "True" {
		log.Printf("Knative service is ready: %s", name)
	} else if ksvcStatus == "False" {
		log.Fatalf("Knative service is not ready: %s", name)
	} else {
		time.Sleep(5 * time.Second)
		VerifyKnativeServiceStatus(name, namespace)
	}
}

func VerifyWorkloadStatus(name string, namespace string) {
	log.Printf("Checking workload status: %s", name)
	workloadStatusBytes, _ := tap.RunWithBash(fmt.Sprintf(`tanzu apps workload get %s -n %s | awk '{if($1=="%s")print $2}'`, name, namespace, name))
	workloadStatus := strings.TrimSpace(string(workloadStatusBytes))
	if workloadStatus == "Ready" {
		log.Printf("Workload is ready: %s", name)
	} else if workloadStatus == "Unknown" {
		log.Printf("Workload status is unknown: %s", name)
		time.Sleep(5 * time.Second)
		VerifyWorkloadStatus(name, namespace)
	} else {
		log.Fatalf("Workload is not ready: %s", name)
	}
}

func VerifyApplicationRunningWithValidationString(envoyExternalIP string, host string, oldString string, newString string, testNew bool) {
	validationString := oldString
	if testNew {
		validationString = newString
	}
	log.Printf("Checking application %s for result: %s", host, validationString)

	if !strings.HasPrefix(envoyExternalIP, "http://") {
		envoyExternalIP = "http://" + envoyExternalIP
	}
	req, err := http.NewRequest("GET", envoyExternalIP, nil)
	tap.CheckError(err)
	req.Host = host
	resp, err := http.DefaultClient.Do(req)
	tap.CheckError(err)
	defer resp.Body.Close()
	resultStringBytes, _ := ioutil.ReadAll(resp.Body)
	resultString := string(resultStringBytes)

	if resultString == validationString {
		log.Printf("Application %s validated, got result: %s", host, validationString)
	} else if testNew && resultString == oldString {
		time.Sleep(5 * time.Second)
		VerifyApplicationRunningWithValidationString(envoyExternalIP, host, oldString, newString, testNew)
	} else {
		log.Fatalf("Application %s not validated, expected: %s, got: %s", host, validationString, resultString)
	}
}

func VerifyBuildStatus(){
	count := 30
	for count <= 30{
		if count == 0{
			log.Fatalf("Builds are not generated after 5 mins")
			break
		}
		result := GetBuilds()
		if len(result.Items) != 0{
			latestBuildIndex :=  len(result.Items) - 1 
			lastConditionIndex := len(result.Items[0].Status.Status.Conditions) - 1
			//Expect(result.Items[0].Status.Status.Conditions[lastConditionIndex].Status).To(Equal(corev1.ConditionTrue))
			if (result.Items[latestBuildIndex].Status.Status.Conditions[lastConditionIndex].Status) == corev1.ConditionUnknown{
				log.Printf("Build %s status is Unknown", result.Items[latestBuildIndex].ObjectMeta.Name)
				//time.Sleep(5 * time.Second)
				
			} else if (result.Items[latestBuildIndex].Status.Status.Conditions[lastConditionIndex].Status) == corev1.ConditionTrue{
				log.Printf("Build %s status is verified successfully. Status is %s", result.Items[latestBuildIndex].ObjectMeta.Name, result.Items[latestBuildIndex].Status.Status.Conditions[lastConditionIndex].Status)
				break
			}
		} else{
			log.Println("Builds are not generated yet")
		}
		log.Printf("Waiting for 10s for builds getting generated ...")
		time.Sleep(10 * time.Second)
		count -= 1
	}
	
	
	

		
}
func GetBuilds() v1alpha2.BuildList{
	var restClient = tap.GetRestClient()
	result := v1alpha2.BuildList{}
	request := restClient.Get()
	request.AbsPath("/apis/kpack.io/v1alpha2/builds").Do(context.TODO()).Into(&result)
	return result
}
