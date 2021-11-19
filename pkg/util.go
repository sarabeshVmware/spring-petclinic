// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func CheckError(err error) bool {
	if err != nil {
		log.Fatal(err.Error())
		return true
	}
	return false
}

func GetCurrentDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}

func GetClientset() *kubernetes.Clientset {
	var homeDir, _ = os.UserHomeDir()
	var kubeconfig = path.Join(homeDir, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

func GetRestClient() rest.Interface {
	return GetClientset().RESTClient()
}
