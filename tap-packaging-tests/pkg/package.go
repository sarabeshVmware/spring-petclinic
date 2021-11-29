// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/buger/jsonparser"
	"gopkg.in/yaml.v3"
)

type Package struct {
	Description         string   `yaml:"description"`
	Name                string   `yaml:"name"`
	Namespace           string   `yaml:"namespace"`
	Package             string   `yaml:"package"`
	ValuesFile          string   `yaml:"values_file,omitempty"`
	Version             string   `yaml:"version"`
	PackageDependencies []string `yaml:"package_dependencies,omitempty"`
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

func ListValuesSchema(packageInfo Package) {
	log.Printf("Values schemas for package: %s", packageInfo.Package)
	Run(fmt.Sprintf("tanzu package available get %s/%s --values-schema -n %s", packageInfo.Package, packageInfo.Version, packageInfo.Namespace))
}

func GetDependentPackagesInfo(parentPackage Package, packagesList []Package) []Package {
	log.Printf("Checking for package dependencies: %s", parentPackage.Package)
	dependentPackagesInfo := []Package{}
	for _, packageDependency := range parentPackage.PackageDependencies {
		for _, packageInfo := range packagesList {
			if packageInfo.Package == packageDependency {
				log.Printf("Dependency for package %s: %s", packageInfo.Package, packageDependency)
				dependentPackagesInfo = append(dependentPackagesInfo, packageInfo)
			}
		}
	}
	return dependentPackagesInfo
}

func CheckIfPackageInstalled(packageInfo Package) bool {
	log.Printf("Checking if package is installed: %s", packageInfo.Package)
	packageInstalled := false
	for _, installedPackage := range ListInstalledPackages(packageInfo.Namespace) {
		if packageInfo.Package == installedPackage.PackageName {
			count := 10
			for count >= 0 {
				if count == 0 {
					log.Printf("Package is not in Reconcile Succeded state after 10 mins. Uninstalling it...")
					UninstallPackage(packageInfo.Namespace, installedPackage.Name)
					packageInstalled = false
					break
				}
				installedPackageStatus := GetInstalledPackageStatus(installedPackage.Name, packageInfo.Namespace)
				if installedPackageStatus == "Reconcile succeeded" {
					log.Printf("Package %s is installed. Status is :%s", packageInfo.Package, installedPackage.Status)
					packageInstalled = true
					break

				} else if installedPackageStatus == "Reconciling" || installedPackageStatus == "Reconcile Failed" {
					log.Printf("Package %s is installed. Status is :%s", packageInfo.Package, installedPackage.Status)
					//packageInstalled = true
					log.Printf("Wating for 1 minute ...")
					time.Sleep(1 * time.Minute)
					count -= 1
				}

			}
		}
	}
	return packageInstalled
}

func GetInstalledPackageStatus(installedPackageName string, namespace string) string {
	var packageInstall []PackageInstalledOutput
	status := ""
	log.Printf("Checking packageInstall status for package: %s", installedPackageName)
	pkgi, _ := Run(fmt.Sprintf("tanzu package installed get %s -n %s -o json", installedPackageName, namespace))
	err := json.Unmarshal(pkgi, &packageInstall)
	CheckError(err)
	if len(packageInstall) > 0 {
		log.Printf("packageInstall status for package: %s is: %s", installedPackageName, packageInstall[0].Status)
		return packageInstall[0].Status
	}
	return status
}


func GetPackageInfoFromName(packageName string, packagesList []Package) Package {
	for _, packageInfo := range packagesList {
		if packageInfo.Name == packageName {
			return packageInfo
		}
	}
	log.Fatalf("Package not found in the provided list: %s", packageName)
	return Package{}
}

func InstallPackageByInfo(packageInfo Package, packagesList []Package) {
	log.Printf("Installing package: %s", packageInfo.Package)

	if CheckIfPackageInstalled(packageInfo) {
		log.Printf("Package already installed: %s", packageInfo.Package)
		return
	}

	// install package dependencies:
	for _, dependentPackageInfo := range GetDependentPackagesInfo(packageInfo, packagesList) {
		log.Printf("Installing package dependency: %s", dependentPackageInfo.Package)
		InstallPackageByInfo(dependentPackageInfo, packagesList)
	}

	// pre-requisites for packages:
	if packageInfo.Package == "scanning.apps.tanzu.vmware.com" {
		log.Printf("Handling pre-requisites for scan-controller:")
		HandleScanControllerPreRequisites(packageInfo)
	}

	// install:
	installCmd := fmt.Sprintf("tanzu package install %s -p %s -v %s -n %s --poll-timeout 30m", packageInfo.Name, packageInfo.Package, packageInfo.Version, packageInfo.Namespace)
	if packageInfo.ValuesFile != "" {
		installCmd += fmt.Sprintf(" -f %s", packageInfo.ValuesFile)
	}
	Run(installCmd)

	ValidatePackage(packageInfo)

	// handle post-installation:
	if packageInfo.Package == "cnrs.tanzu.vmware.com" {
		log.Printf("Handling post-installation for cloud-native-runtimes:")
		HandleCloudNativeRuntimesPostInstallation()
	}
	if packageInfo.Package == "image-policy-webhook.signing.run.tanzu.vmware.com" {
		log.Printf("Handling post-installation for image-policy-webhook:")
		HandleImagePolicyWebhookPostInstallation()
	}
}

func InstallPackageByName(packageName string, packagesList []Package) {
	packageInfo := GetPackageInfoFromName(packageName, packagesList)
	InstallPackageByInfo(packageInfo, packagesList)
}

func ValidatePackage(packageInfo Package) {
	log.Printf("Validating package: %s", packageInfo.Package)
	packageInstalled, _ := Run(fmt.Sprintf("tanzu package installed get %s -n %s -o json", packageInfo.Name, packageInfo.Namespace))
	status, err := jsonparser.GetString(packageInstalled, "[0]", "status")
	CheckError(err)
	if status == "Reconciling" || status == "" {
		time.Sleep(5 * time.Second)
		ValidatePackage(packageInfo)
	} else if status == "Reconcile succeeded" {
		log.Printf("Reconcile succeeded for package install: %s", packageInfo.Package)
	} else {
		log.Fatalf("Reconcile not succeeded for package install: %s", packageInfo.Package)
	}
}

func UninstallPackages(namespace string) {
	installedpackages := ListInstalledPackages(namespace)
	for _, each := range installedpackages {
		log.Printf("Uninstalling package: %s", each.Name)
		Run(fmt.Sprintf("tanzu package installed delete %s -n %s -y", each.Name, namespace))
	}
}

func UninstallPackage(namespace string, InstalledPackageName string) {
	log.Printf("Uninstalling package: %s", InstalledPackageName)
	Run(fmt.Sprintf("tanzu package installed delete %s -n %s -y", InstalledPackageName, namespace))
}
func HandleScanControllerPreRequisites(packageInfo Package) {
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
	scanControllerBytes, err := os.ReadFile(packageInfo.ValuesFile)
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
	err = os.WriteFile(packageInfo.ValuesFile, scanControllerBytes, 0666)
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
