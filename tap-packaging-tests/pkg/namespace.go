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
apiVersion: v1
kind: Secret
metadata:
  name: tap-registry
  annotations:
    secretgen.carvel.dev/image-pull-secret: ""
type: kubernetes.io/dockerconfigjson
data:
  .dockerconfigjson: e30K

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: default # use value from "Install Default Supply Chain"
secrets:
  - name: registry-credentials
imagePullSecrets:
  - name: registry-credentials
  - name: tap-registry

---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: kapp-permissions
  annotations:
    kapp.k14s.io/change-group: "role"
rules:
  - apiGroups:
      - servicebinding.io
    resources: ['servicebindings']
    verbs: ['*']
  - apiGroups:
      - serving.knative.dev
    resources: ['services']
    verbs: ['*']
  - apiGroups: [""]
    resources: ['configmaps']
    verbs: ['get', 'watch', 'list', 'create', 'update', 'patch', 'delete']

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: kapp-permissions
  annotations:
    kapp.k14s.io/change-rule: "upsert after upserting role"
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: kapp-permissions
subjects:
  - kind: ServiceAccount
    name: default # use value from "Install Default Supply Chain"
`
	os.WriteFile(tempFile.Name(), []byte(configuration), 0666)
	ApplyConfigurationInNamespace(tempFile.Name(), namespace)
}
