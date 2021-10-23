// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/buger/jsonparser"
	"gopkg.in/yaml.v3"
)

type Package struct {
	Name          string `yaml:"name"`
	InstalledName string `yaml:"installed_name"`
	Version       string `yaml:"version"`
	UseValuesFile string `yaml:"use_values_file"`
}

type PackageInstalledOutput struct {
	Name           string `json:"name"`
	PackageName    string `json:"package-name"`
	PackageVersion string `json:"package-version"`
	Status         string `json:"status"`
}

func ListPackages(namespace string) {
	log.Printf("Available packages in namespace: %s", namespace)
	Run(fmt.Sprintf("tanzu package available list -n %s", namespace))
}

func ListInstalledPackages(namespace string) []PackageInstalledOutput {
	var packages []PackageInstalledOutput
	log.Printf("Installed packages in namespace: %s", namespace)
	packagesList, _ := Run(fmt.Sprintf("tanzu package installed list -n %s -o json", namespace))
	err := json.Unmarshal(packagesList, &packages)
	CheckError(err)
	return packages
}

func ListValuesSchema(packages []Package, namespace string) {
	for _, packageInfo := range packages {
		log.Printf("Values schemas for package: %s", packageInfo.Name)
		Run(fmt.Sprintf("tanzu package available get %s/%s --values-schema -n %s", packageInfo.Name, packageInfo.Version, namespace))
	}
}

func InstallPackages(packages []Package, namespace string, ValuesDirectory string) {
	for _, packageInfo := range packages {
		log.Printf("Installing package: %s", packageInfo.Name)

		// pre-requisites for packages:
		if packageInfo.Name == "appliveview.tanzu.vmware.com" {
			log.Printf("Handling pre-requisites for appliveview:")
			HandleAppLiveViewPreRequisites(packageInfo, ValuesDirectory)
		}
		if packageInfo.Name == "scanning.apps.tanzu.vmware.com" {
			log.Printf("Handling pre-requisites for scan-controller:")
			HandleScanControllerPreRequisites(packageInfo, ValuesDirectory)
		}

		// install:
		installCmd := fmt.Sprintf("tanzu package install %s -p %s -v %s -n %s", packageInfo.InstalledName, packageInfo.Name, packageInfo.Version, namespace)
		if packageInfo.UseValuesFile != "" {
			installCmd += fmt.Sprintf(" -f %s", filepath.Join(ValuesDirectory, packageInfo.UseValuesFile))
		}
		if packageInfo.Name == "cnrs.tanzu.vmware.com" || packageInfo.Name == "buildservice.tanzu.vmware.com" {
			installCmd += " --poll-timeout 30m"
		}
		Run(installCmd)

		ValidatePackage(packageInfo, namespace)

		// handle post-installation:
		if packageInfo.Name == "cnrs.tanzu.vmware.com" {
			log.Printf("Handling post-installation for cloud-native-runtimes:")
			HandleCloudNativeRuntimesPostInstallation()
		}
		if packageInfo.Name == "image-policy-webhook.signing.run.tanzu.vmware.com" {
			log.Printf("Handling post-installation for image-policy-webhook:")
			HandleImagePolicyWebhookPostInstallation()
		}
	}
}

func ValidatePackage(packageInfo Package, namespace string) {
	log.Printf("Validating package: %s", packageInfo.Name)
	packageInstalled, _ := Run(fmt.Sprintf("tanzu package installed get %s -n %s -o json", packageInfo.InstalledName, namespace))
	status, err := jsonparser.GetString(packageInstalled, "[0]", "status")
	CheckError(err)
	if status == "Reconciling" {
		time.Sleep(5 * time.Second)
		ValidatePackage(packageInfo, namespace)
	} else if status == "Reconcile succeeded" {
		log.Printf("Reconcile succeeded for package install: %s", packageInfo.Name)
	} else {
		log.Fatalf("Reconcile not succeeded for package install: %s", packageInfo.Name)
	}
}

func UninstallPackages(namespace string) {
	installedpackages := ListInstalledPackages(namespace)
	for _, each := range installedpackages {
		log.Printf("Uninstalling package: %s", each.Name)
		Run(fmt.Sprintf("tanzu package installed delete %s -n %s -y", each.Name, namespace))
	}
}

func HandleAppLiveViewPreRequisites(packageInfo Package, ValuesDirectory string) {
	valuesSchemaFile := filepath.Join(ValuesDirectory, packageInfo.UseValuesFile)
	appliveviewSchemaBytes, err := os.ReadFile(valuesSchemaFile)
	CheckError(err)
	appliveviewSchema := struct {
		ServerNamespace string `yaml:"server_namespace"`
	}{}
	err = yaml.Unmarshal([]byte(appliveviewSchemaBytes), &appliveviewSchema)
	CheckError(err)
	if appliveviewSchema.ServerNamespace != "" {
		CreateNamespace(appliveviewSchema.ServerNamespace)
	} else {
		CreateNamespace("app-live-view")
	}
}

func HandleScanControllerPreRequisites(packageInfo Package, ValuesDirectory string) {
	tempFile, err := ioutil.TempFile("", "configuration*.yaml")
	CheckError(err)
	defer os.Remove(tempFile.Name())

	configuration := `
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: metadata-store-read-write
  namespace: metadata-store
rules:
- resources: ["all"]
  verbs: ["get", "create", "update"]
  apiGroups: [ "metadata-store/v1" ]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: metadata-store-read-write
  namespace: metadata-store
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: metadata-store-read-write
subjects:
- kind: ServiceAccount
  name: metadata-store-read-write-client
  namespace: metadata-store

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: metadata-store-read-write-client
  namespace: metadata-store
automountServiceAccountToken: false
`
	os.WriteFile(tempFile.Name(), []byte(configuration), 0666)
	ApplyConfiguration(tempFile.Name())

	// update values schema
	metadataStoreUrl, _ := RunWithBash(`kubectl -n metadata-store get service -o name | grep app | xargs kubectl -n metadata-store get -o jsonpath='{.spec.ports[].name}{"://"}{.metadata.name}{"."}{.metadata.namespace}{".svc.cluster.local:"}{.spec.ports[].port}'`)
	metadataStoreCa, _ := RunWithBash(`kubectl get secret app-tls-cert -n metadata-store -o json | jq -r '.data."ca.crt"' | base64 -d`)
	metadataStoreToken, _ := RunWithBash(`kubectl get secret $(kubectl get sa -n metadata-store metadata-store-read-write-client -o json | jq -r '.secrets[0].name') -n metadata-store -o json | jq -r '.data.token' | base64 -d`)
	valuesSchemaFile := filepath.Join(ValuesDirectory, packageInfo.UseValuesFile)
	scanControllerBytes, err := os.ReadFile(valuesSchemaFile)
	CheckError(err)
	scanControllerSchema := struct {
		MetadataStoreUrl         string `yaml:"metadataStoreUrl"`
		MetadataStoreCa          string `yaml:"metadataStoreCa"`
		MetadataStoreTokenSecret string `yaml:"metadataStoreTokenSecret"`
	}{}
	err = yaml.Unmarshal([]byte(scanControllerBytes), &scanControllerSchema)
	CheckError(err)
	scanControllerSchema.MetadataStoreUrl = string(metadataStoreUrl)
	scanControllerSchema.MetadataStoreCa = string(metadataStoreCa)
	scanControllerBytes, err = yaml.Marshal(&scanControllerSchema)
	CheckError(err)
	err = os.WriteFile(valuesSchemaFile, scanControllerBytes, 0666)
	CheckError(err)
	log.Printf("Updated values schema for scan-controller: \n%s", string(scanControllerBytes))

	CreateNamespace("scan-link-system")

	configuration = fmt.Sprintf(`
---
apiVersion: v1
kind: Secret
metadata:
  name: %s
  namespace: scan-link-system
type: kubernetes.io/opaque
stringData:
  token: %s
`, scanControllerSchema.MetadataStoreTokenSecret, metadataStoreToken)
	os.WriteFile(tempFile.Name(), []byte(configuration), 0666)
	ApplyConfiguration(tempFile.Name())
}

func HandleCloudNativeRuntimesPostInstallation() {
	log.Printf("Creating an empty secret:")
	Run_AllowError("kubectl create secret generic pull-secret --from-literal=.dockerconfigjson={} --type=kubernetes.io/dockerconfigjson")
	log.Printf("Annotating the empty secret as a target of the secretgen controller:")
	Run_AllowError(`kubectl annotate secret pull-secret secretgen.carvel.dev/image-pull-secret=""`)
	log.Printf("Adding the secret to the service account:")
	RunWithBash(`kubectl patch serviceaccount default -p '{"imagePullSecrets": [{"name": "pull-secret"}]}'`)
	Run("kubectl describe serviceaccount default")
}

func HandleImagePolicyWebhookPostInstallation() {
	tempFile, err := ioutil.TempFile("", "configuration*.yaml")
	CheckError(err)
	defer os.Remove(tempFile.Name())

	configuration := `
apiVersion: signing.run.tanzu.vmware.com/v1alpha1
kind: ClusterImagePolicy
metadata:
 name: image-policy
spec:
 verification:
   exclude:
     resources:
       namespaces:
       - kube-system
   keys:
   - name: cosign-key
     publicKey: |
       -----BEGIN PUBLIC KEY-----
       MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEhyQCx0E9wQWSFI9ULGwy3BuRklnt
       IqozONbbdbqz11hlRJy9c7SG+hdcFl9jE9uE/dwtuwU2MqU9T/cN0YkWww==
       -----END PUBLIC KEY-----
   images:
   - namePattern: gcr.io/projectsigstore/cosign*
     keys:
     - name: cosign-key
`
	os.WriteFile(tempFile.Name(), []byte(configuration), 0666)
	ApplyConfiguration(tempFile.Name())
}
