// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func CheckIfNamespaceExists(namespace string) bool {
	log.Printf("Checking if namespace exists: %s", namespace)
	_, err := Run_AllowError(fmt.Sprintf("kubectl get namespaces %s", namespace))
	if err != nil {
		log.Printf("Namespace not found: %s", namespace)
		return false
	}
	log.Printf("Namespace found: %s", namespace)
	return true
}

func DeleteNamespace(namespace string) {
	if CheckIfNamespaceExists(namespace) {
		log.Printf("Deleting namespace: %s", namespace)
		Run(fmt.Sprintf("kubectl delete ns %s", namespace))
	}
}

func CreateNamespace(namespace string) {
	if !CheckIfNamespaceExists(namespace) {
		log.Printf("Creating namespace: %s", namespace)
		Run(fmt.Sprintf("kubectl create ns %s", namespace))
	}
}

func SetupDeveloperNamespacePostInstallation(namespace string) {
	log.Printf("Setting up developer namespace: %s", namespace)
	tempFile, err := ioutil.TempFile("", "configuration*.yaml")
	CheckError(err)
	defer os.Remove(tempFile.Name())

	configuration := `
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: default
secrets:
  - name: image-secret
imagePullSecrets:
  - name: image-secret

---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: default
rules:
- apiGroups: [source.toolkit.fluxcd.io]
  resources: [gitrepositories]
  verbs: ['*']
- apiGroups: [source.apps.tanzu.vmware.com]
  resources: [imagerepositories]
  verbs: ['*']
- apiGroups: [carto.run]
  resources: [deliverables, runnables]
  verbs: ['*']
- apiGroups: [kpack.io]
  resources: [images]
  verbs: ['*']
- apiGroups: [conventions.apps.tanzu.vmware.com]
  resources: [podintents]
  verbs: ['*']
- apiGroups: [""]
  resources: ['configmaps']
  verbs: ['*']
- apiGroups: [""]
  resources: ['pods']
  verbs: ['list']
- apiGroups: [tekton.dev]
  resources: [taskruns, pipelineruns]
  verbs: ['*']
- apiGroups: [tekton.dev]
  resources: [pipelines]
  verbs: ['list']
- apiGroups: [kappctrl.k14s.io]
  resources: [apps]
  verbs: ['*']
- apiGroups: [serving.knative.dev]
  resources: ['services']
  verbs: ['*']
- apiGroups: [servicebinding.io]
  resources: ['servicebindings']
  verbs: ['*']
- apiGroups: [services.apps.tanzu.vmware.com]
  resources: ['resourceclaims']
  verbs: ['*']
- apiGroups: [scanning.apps.tanzu.vmware.com]
  resources: ['imagescans', 'sourcescans']
  verbs: ['*']

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: kapp-permissions
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: default
subjects:
  - kind: ServiceAccount
    name: default
`
	os.WriteFile(tempFile.Name(), []byte(configuration), 0666)
	ApplyConfigurationInNamespace(tempFile.Name(), namespace)
}
