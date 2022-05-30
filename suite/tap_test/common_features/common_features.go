package common_features

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	exec2 "os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/docker"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/git"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/github"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/imgpkg"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectlCmds"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_helpers"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_libs"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubernetes/client"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/misc"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzuCmds"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzu_helpers"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzu_libs"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/models"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

var tiltprocCmdKey = "tiltprocCmd"
var rootDir = filepath.Join(utils.GetFileDir(), "../../")
var buildName = ""
var ksvcLatestReady = ""
var revisionName = ""

func compile(filepath string) {
	app := "./mvnw"
	arg0 := "compile"
	cmd := exec2.Command(app, arg0)
	cmd.Dir = filepath
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(stdout))
}

func UpdatePackageRepository(t *testing.T, name string, registry string, namespace string) features.Feature {
	return features.New("updating package repository").
		Assess(fmt.Sprintf("updating-packaging-repository-%s", name), func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			log.Printf("updating pacakage repository %s", name)
			tanzu_libs.TanzuUpdatePackageRepository(name, registry, namespace)
			updated := tanzu_helpers.CheckIfPackageRepositoryReconciled(name, namespace, 5, 30)
			if updated {
				t.Logf("Updated repository : %s, image: %s successfully", name, registry)
			} else {
				t.Error(fmt.Errorf("update FAILED for repository : %s, image: %s", name, registry))
				t.Fail()
			}
			return ctx
		}).Feature()
}

func AddPackageRepository(t *testing.T, name string, registry string, version string, namespace string) features.Feature {
	return features.New("adding package repository").
		Assess(fmt.Sprintf("adding-packaging-repository-%s", name), func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			log.Printf("adding package repository %s (%s) in namespace %s", name, registry, namespace)
			registryWithVersion := fmt.Sprintf("%s:%s", registry, version)
			// add repo
			err := tanzuCmds.TanzuAddPackageRepository(name, registryWithVersion, namespace)
			if err == nil {
				t.Logf("Installed repository : %s, image: %s:%s successfully", name, registry, version)
			} else {
				t.Error(fmt.Errorf("install FAILED for repository : %s, image: %s:%s", name, registry, version))
				t.Fail()
			}
			updated := tanzu_helpers.CheckIfPackageRepositoryReconciled(name, namespace, 5, 30)
			if updated {
				t.Logf("Updated repository : %s, image: %s:%s successfully", name, registry, version)
			} else {
				t.Error(fmt.Errorf("update FAILED for repository : %s, image: %s:%s", name, registry, version))
				t.Fail()
			}
			return ctx
		}).Feature()
}

func DeletePackageRepository(t *testing.T, name string, namespace string) features.Feature {
	return features.New("deleting package repository").
		Assess(fmt.Sprintf("deleting-packaging-repository-%s", name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			log.Printf("deleting package repository %s in namespace %s", name, namespace)

			// delete repo
			err := tanzuCmds.TanzuDeletePackageRepository(name, namespace)
			if err != nil {
				t.Errorf("error while deleting package repository %s in namespace %s", name, namespace)
				t.Fail()
			}
			return ctx
		}).Feature()
}

func InstallPackage(t *testing.T, name string, packageRepository string, version string, namespace string, valuesFile string, pollTimeout string) features.Feature {
	return features.New("installing package").
		Assess(fmt.Sprintf("installing-package-%s", name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			log.Printf("installing tap package  %s (%s)", name, packageRepository)
			err := tanzu_libs.InstallPackage(name, packageRepository, version, namespace, valuesFile, pollTimeout)
			if err != nil {
				// if error, check via kubectl, not tanzu-cli
				pass := kubectl_helpers.ValidateTAPInstallation(name, namespace, 10, 60)
				if !pass {
					kubectl_helpers.LogFailedResourcesDetails(namespace)
					log.Printf("error while installing package %s (%s)", name, packageRepository)
					return ctx
				} else {
					return ctx
				}
			}
			return ctx
		}).
		Feature()

}

func DeletePackage(t *testing.T, name string, namespace string) features.Feature {
	return features.New("deleting-package").
		Assess(fmt.Sprintf("deleting-package-%s", name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			err := tanzu_libs.DeleteInstalledPackage(name, namespace)
			if err != nil {
				t.Error(fmt.Errorf("Uninstallation FAILED for package : %s", name))
				t.Fail()
			}
			return ctx
		}).
		Feature()
}

func CheckIfPackageInstalled(name string, namespace string, recursiveCount int, secondsGap int) features.Feature {
	return features.New("checking package").
		Assess(fmt.Sprintf("checking-package-%s", name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			log.Printf("checking package %s installation status", name)

			for ; recursiveCount >= 0; recursiveCount-- {
				// get status
				packageInstalledStatus, err := tanzuCmds.TanzuGetPackageInstalledStatus(name, namespace)
				if err != nil {
					t.Errorf("error while getting package %s in namespace %s installation status", name, namespace)
					t.Fail()
				}

				// check
				if packageInstalledStatus == "Reconciling" || packageInstalledStatus == "" {
					log.Printf("package %s is getting installed", name)
					log.Printf("sleeping for %d seconds", secondsGap)
					time.Sleep(time.Duration(secondsGap) * time.Second)
				} else if packageInstalledStatus == "Reconcile succeeded" {
					log.Printf("package %s is installed", name)
				} else if packageInstalledStatus == "Reconcile Failed" {
					t.Errorf("package %s installation failed", name)
					t.Fail()
				} else {
					t.Errorf("package %s installation unknown", name)
					t.Fail()
				}
			}

			return ctx
		}).Feature()
}

func UpdateTapVersion(t *testing.T, name string, tapPackageName string, namespace string, valuesFile string, tapVersion string, pollTimeout string) features.Feature {
	return features.New(fmt.Sprintf("updating-tap-version-%s", tapVersion)).
		Assess(fmt.Sprintf("updating-tap-package-%s", tapVersion), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			log.Printf("updating tap- %s ", tapVersion)
			tanzu_libs.UpdateInstalledPackage(name, tapPackageName, tapVersion, namespace, valuesFile, pollTimeout)
			updated := tanzu_helpers.ValidateInstalledPackageVersion(name, namespace, tapVersion, 30, 60)
			if updated {
				t.Logf("Updated tap version: %s successfully", tapVersion)
			} else {
				t.Error(fmt.Errorf("update FAILED for tap version: %s", tapVersion))
				t.Fail()
			}
			availablePkgs := tanzu_libs.ListInstalledPackages(namespace)
			for _, pkg := range availablePkgs {
				installed := tanzu_helpers.ValidateInstalledPackageStatus(pkg.NAME, namespace, 10, 30)
				if installed {
					t.Logf("Installed package : %s, version: %s successfully", pkg.NAME, pkg.PACKAGE_VERSION)
				} else {
					t.Errorf("Installation FAILED for package : %s, version: %s, status: %s", pkg.NAME, pkg.PACKAGE_VERSION, pkg.STATUS)
					t.Fail()
				}
			}
			log.Printf("final packages version after tap update...")
			tanzu_libs.ListInstalledPackages(namespace)
			return ctx
		}).
		Feature()
}

func UpdateTapProfileSupplyChain(t *testing.T, name string, tapPackageName string, tapVersion string, profile string, supplyChain string, namespace string) features.Feature {
	return features.New("update-tap-profile-supplychain").
		Assess("update-package", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("updating tap package")

			// get schema and update values
			tapValuesSchema, err := models.GetTapValuesSchema()
			if err != nil {
				t.Error("error while getting tap values schema")
				t.FailNow()
			}
			tapValuesSchema.Profile = profile
			tapValuesSchema.SupplyChain = supplyChain

			// create temporary file
			t.Log("creating tempfile for tap values schema")
			tempFile, err := ioutil.TempFile("", "tap-values*.yaml")
			if err != nil {
				t.Error("error while creating tempfile for tap values schema")
				t.FailNow()
			} else {
				t.Log("created tempfile")
			}
			defer os.Remove(tempFile.Name())

			// write the updated schema to the temporary file
			err = utils.WriteYAMLFile(tempFile.Name(), tapValuesSchema)
			if err != nil {
				t.Error("error while writing updated tap values schema to YAML file")
				t.FailNow()
			} else {
				t.Log("wrote tap values schema to file")
			}

			// update tap
			err = tanzuCmds.TanzuUpdatePackage(name, tapPackageName, tapVersion, namespace, tempFile.Name())
			if err != nil {
				t.Error("error while updating tap")
				t.FailNow()
			} else {
				t.Log("updated tap")
			}

			return ctx
		}).
		Feature()
}

func TanzuDeployWorkload(t *testing.T, workloadYamlFile string, namespace string) features.Feature {
	return features.New("deploy-tanzu-workload-via-yaml").
		Assess(fmt.Sprintf("deploy-workload-from-%s", workloadYamlFile), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("deploying workload yaml %s", workloadYamlFile)
			//workloadFilePath := filepath.Join(rootDir, workloadYamlFile)
			// deploy app
			err := tanzu_libs.TanzuApplyWorkload(namespace, workloadYamlFile)
			if err != nil {
				t.Errorf("error while deploying workload yaml %s", workloadYamlFile)
				t.FailNow()
			} else {
				t.Logf("deployed workload yaml %s", workloadYamlFile)
			}

			return ctx
		}).
		Feature()
}

func TanzuDeleteWorkload(t *testing.T, name string, namespace string) features.Feature {
	return features.New("delete-tanzu-workload").
		Assess(fmt.Sprintf("delete-workload-from-%s", name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("deploying workload %s", name)
			// deploy app
			err := tanzu_libs.DeleteWorkload(name, namespace)
			if err != nil {
				t.Errorf("error while deleting workload %s", name)
				t.FailNow()
			} else {
				t.Logf("deleting workload %s", name)
			}

			return ctx
		}).
		Feature()
}

func GitClone(t *testing.T, gitUsername string, gitEmail string, gitRepository string) features.Feature {
	return features.New("git-update").
		Assess("git-config", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("setting git config")

			// set git config
			err := git.GitConfig(gitUsername, gitEmail)
			if err != nil {
				t.Error("error while setting git config")
				t.FailNow()
			} else {
				t.Log("set git config")
			}

			return ctx
		}).
		Assess("git-clone", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("cloning git repo")

			// clone
			err := git.GitClone(rootDir, gitRepository)
			if err != nil {
				t.Error("error while cloning git repo")
				t.FailNow()
			} else {
				t.Log("cloned git repo")
			}

			return ctx

		}).
		Feature()
}

func ReplaceStringInFile(t *testing.T, originalString string, newString string, filePath string, workload string) features.Feature {
	return features.New("replace-string-in-file").
		Assess("replace-tanzu-to-tap ", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			oldString := originalString
			newString := newString
			filePath := filepath.Join(rootDir, filePath)
			t.Logf("Replace from string %s to string %s in file %s", oldString, newString, filePath)
			err := utils.ReplaceStringInFile(filePath, oldString, newString)
			t.Logf("Compiling and building app %s", workload)
			compile(filepath.Join(rootDir, workload))
			if err != nil {
				t.Error(fmt.Errorf("error while replacing string in file %s : %w", filePath, err))
				t.FailNow()
			}
			return ctx
		}).
		Feature()
}

func VerifyTanzuWorkloadStatus(t *testing.T, name string, namespace string) features.Feature {
	return features.New(fmt.Sprintf("verify-%sworkload-status", name)).
		Assess(fmt.Sprintf("verify-tanzu-%s-ready", name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verifying tanzu %s ready status", name)

			// check
			status := kubectl_helpers.ValidateWorkloadStatus(name, namespace, 10, 10)
			t.Logf("workload %s validation status : %v", name, status)
			if !status {
				t.Error(fmt.Errorf("workload %s is not ready.", name))
				t.Fail()
			}
			return ctx
		}).
		Feature()
}

func CreateGithubRepo(t *testing.T, name string, repoTemplate string, accessToken string) features.Feature {
	return features.New("create-github-repo").
		Assess("create-github-repo", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("creating github repo")

			// create repo
			err := github.CreateGithubRepo(name, repoTemplate, accessToken)
			if err != nil {
				t.Error("error while creating repo ")
				t.FailNow()
			} else {
				t.Log("created repo")
			}
			return ctx
		}).
		Feature()
}

func DeleteGithubRepo(t *testing.T, name string, accessToken string) features.Feature {
	return features.New("delete-github-repo").
		Assess("delete-github-repo", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("deleting github repo")

			// create repo
			err := github.DeleteGithubRepo(name, accessToken)
			if err != nil {
				t.Error("error while deleting repo ")
				t.FailNow()
			} else {
				t.Log("deleted repo")
			}
			return ctx
		}).
		Feature()
}

func ApplyKubectlConfigurationFile(t *testing.T, configurationFile string, namespace string) features.Feature {
	return features.New("deploy-yaml").
		Assess(fmt.Sprintf("deploy-%s", configurationFile), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("deploying %s", configurationFile)

			// deploy app
			err := kubectlCmds.KubectlApplyConfiguration(configurationFile, namespace)
			if err != nil {
				t.Errorf("error while deploying %s", configurationFile)
				t.FailNow()
			} else {
				t.Logf("deployed %s", configurationFile)
			}

			return ctx
		}).
		Feature()
}

func VerifyWorkloadResponse(t *testing.T, workloadUrl string, verificationString string, relativePath string) features.Feature {
	return features.New("verify-workload-response").
		Assess("get-externalip-and-check-webpage-for-string", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("getting external ip and checking for string")

			// get external IP
			url := kubectl_helpers.GetServiceExternalIP("envoy", "tanzu-system-ingress", 2, 30)
			if relativePath != "" {
				url = fmt.Sprintf("%s/%s", url, relativePath)
			}
			// set url
			if !strings.HasPrefix(url, "http://") {
				url = "http://" + url
			}

			webpageContainsString, _ := misc.VerifyWebpageContainsString(workloadUrl, url, verificationString, 20, 10, 30)
			if !webpageContainsString {
				t.Error("webpage does not contains string")
				t.Fail() // DON'T DO t.FailNow() AS WE WANT TO CLEAN UP REGARDLESS OF THE STATE OF THE TEST
			} else {
				t.Log("webpage contains string")
			}

			return ctx
		}).
		Feature()
}

func ImgPkgCopyToRepo(t *testing.T, sourceBundle string, targetRepo string) features.Feature {
	return features.New("imgpkg copy").
		Assess(fmt.Sprintf("copying image bundles from %s to %s", sourceBundle, targetRepo), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("copying image bundles from %s to %s", sourceBundle, targetRepo)

			// deploy app
			err := imgpkg.ImgpkgCopy(sourceBundle, targetRepo)
			if err != nil {
				t.Errorf("error while copying image bundles from %s to %s", sourceBundle, targetRepo)
				t.FailNow()
			} else {
				t.Logf("copied image bundles from %s to %s", sourceBundle, targetRepo)
			}

			return ctx
		}).
		Feature()
}

func CreateSecret(t *testing.T, name string, registry string, username string, password string, passwordType string, namespace string, export bool) features.Feature {
	return features.New(fmt.Sprintf("creating secret %s", name)).
		Assess(fmt.Sprintf("creating secret %s", name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			log.Printf("creating secret %s (registry %s, username %s) in namespace %s", name, registry, username, namespace)

			// create secret
			err := tanzuCmds.TanzuCreateSecret(name, registry, username, password, passwordType, namespace, export)
			if err != nil {
				t.Errorf("error while creating secret %s", name)
				t.FailNow()
			} else {
				t.Logf("created secret %s", name)
			}

			return ctx
		}).Feature()
}

func DeleteSecret(t *testing.T, name string, namespace string) features.Feature {
	return features.New(fmt.Sprintf("deleting secret %s", name)).
		Assess(fmt.Sprintf("deleting secret %s", name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			log.Printf("deleting secret %s in namespace %s", name, namespace)

			// create secret
			err := tanzuCmds.TanzuDeleteSecret(name, namespace)
			if err != nil {
				t.Errorf("error while deleting secret %s", name)
				t.FailNow()
			} else {
				t.Logf("deleted secret %s", name)
			}

			return ctx
		}).Feature()
}

func DockerLogin(t *testing.T, server string, username string, password string) features.Feature {
	return features.New(fmt.Sprintf("Docker login server")).
		Assess(fmt.Sprintf("Logging in server %s", server), func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {

			err := docker.DockerLogin(server, username, password)
			if err != nil {
				t.Errorf("error while loging in server %s", server)
				t.FailNow()
			} else {
				t.Logf("docker login success for %s", server)
			}

			return ctx
		}).Feature()
}

func ChangeContext(t *testing.T, clusterContext string) features.Feature {
	return features.New("changing cluster context").
		Assess(fmt.Sprintf("changing cluster context to %s", clusterContext), func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {

			_, err := kubectl_libs.UseContext(clusterContext)
			if err != nil {
				t.Errorf("error while changing context to %s", clusterContext)
				t.FailNow()
			} else {
				t.Logf("context changed to %s", clusterContext)
			}
			return ctx
		}).Feature()
}

func VerifyTanzuJavaWebAppImageRepository(t *testing.T, name string, namespace string) features.Feature {
	return features.New("verify-image-repositories").
		Assess("verify-image-repositories", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verify image-repositories status")
			status := kubectl_helpers.VerifyImageRepositoryStatus(name, namespace, 10, 30)
			t.Logf("ImageRepository %s status is : %t", name, status)
			if !status {
				t.Error(fmt.Errorf("ImageRepository %s is not ready.", name))
				t.Fail()
			}
			return ctx
		}).Feature()
}

func GenerateAcceleratorProject(t *testing.T, namespace string) features.Feature {
	return features.New("generate-acc-project-and-unzip").
		Assess("generate-project", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			service, accNamespace := "acc-server", "accelerator-system"
			t.Logf("getting external ip for %s (namespace %s)", service, accNamespace)
			serviceExternalIp := kubectl_helpers.GetServiceExternalIP(service, accNamespace, 2, 30)
			// if err != nil {
			// 	t.Error(fmt.Errorf("error while getting external ip for %s (namespace %s): %w", service, accNamespace, err))
			// 	t.FailNow()
			// }
			t.Logf("external ip for %s (namespace %s): %s", "server", accNamespace, serviceExternalIp)
			t.Logf("sleeping for 1 minute before generating project")
			t.Logf("generating tanzu java web app accelerator project")
			tapValuesSchema, err := models.GetProfileTapValuesSchema("iterate")
			if err != nil {
				t.Error(fmt.Errorf("error while getting tap values schema: %w", err))
			}
			// generate project
			repositoryPrefix := tapValuesSchema.OotbSupplyChainBasic.Registry.Server + "/" + tapValuesSchema.OotbSupplyChainBasic.Registry.Repository
			err = tanzuCmds.TanzuGenerateAccelerator("tanzu-java-web-app", "tanzu-java-web-app", repositoryPrefix, serviceExternalIp, namespace, 4, 30)
			if err != nil {
				t.Error("error while generating accelerator project")
				t.FailNow()
			} else {
				t.Log("accelerator project generated")
			}

			return ctx
		}).
		Assess("unzip-project", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			project := filepath.Join(rootDir, "tanzu-java-web-app")
			zipFile := "tanzu-java-web-app" + ".zip"
			t.Logf("listing existing project files if exists")
			output, err := linux_util.ExecuteCmd(fmt.Sprintf("ls -lt %s", project))
			t.Logf("command executed: ls -lt %s. output %s", project, output)
			if err == nil {
				t.Logf("deleting %s folder", project)
				output, err := linux_util.ExecuteCmd(fmt.Sprintf("rm -rf %s", project))
				t.Logf("command executed: rm -rf %s. output %s", project, output)
				if err != nil {
					t.Error(fmt.Errorf("error while deleting project files %s: %w: %s", project, err, output))
					t.FailNow()
				}
			}

			t.Logf("unzipping file %s", zipFile)
			output, err = linux_util.ExecuteCmd(fmt.Sprintf("unzip %s -d %s", zipFile, rootDir))
			t.Logf("command executed: unzip %s. output %s", zipFile, output)
			if err != nil {
				t.Error(fmt.Errorf("error while unzip accelerator project zip file %s: %w: %s", zipFile, err, output))
				t.FailNow()
			}
			t.Logf("Accelerator project zip files %s unzipped successfully", zipFile)

			return ctx
		}).
		Feature()
}

func VerifyTanzuJavaWebAppBuildStatus(t *testing.T, name string, buildNameSuffix string, namespace string) features.Feature {
	return features.New("verify-builds").
		Assess("verify-build-status", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verify build status")
			buildName = fmt.Sprintf("%s%s", name, buildNameSuffix)
			status := kubectl_helpers.VerifyBuildStatus(buildName, namespace, 10, 30)
			t.Logf("Build status is : %t", status)
			if !status {
				t.Error(fmt.Errorf("Build is not ready."))
				t.Fail()
			}
			return ctx
		}).
		Feature()
}

func VerifyTanzuJavaWebAppImagesKpacStatus(t *testing.T, namespace string) features.Feature {
	return features.New("verify-images.kpac").
		Assess("verify-images.kpac-status", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verify latest image status")
			status := kubectl_helpers.GetLatestImageStatus(namespace)
			t.Logf("Image status is: %s", status)
			if status != "True" {
				t.Error(fmt.Errorf("Image is not built/ready."))
				t.Fail()
			}
			return ctx
		}).
		Feature()
}

func VerifyTanzuJavaWebAppPodIntentStatus(t *testing.T, name string, namespace string) features.Feature {
	return features.New("verify-podintents-labels-conventions").
		Assess("verify-podintent-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying podintent ready status")

			// check
			if !kubectl_helpers.VerifyPodIntentStatus(name, namespace, 5, 30) {
				t.Error("podintent not ready")
				t.FailNow()
			} else {
				t.Log("podintent ready")
			}
			return ctx
		}).
		Assess("verify-podintent-alv-lables", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying appliveview labels present in podintent")

			// check
			alvLabelsPresent := kubectl_helpers.ValidateAppLiveViewLabels(name, namespace)
			if !alvLabelsPresent {
				t.Error("appliveview lables absent in podintent")
				t.FailNow()
			} else {
				t.Log("appliveview labels present in podintent")
			}
			return ctx
		}).
		Assess("verify-podintent-springbootconventions-lables", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying springbootconventions labels present in podintent")

			// check
			springbootconventionsLabelsPresent := kubectl_helpers.ValidateSpringBootLabels(name, namespace)
			if !springbootconventionsLabelsPresent {
				t.Error("springbootconventions lables absent in podintent")
				t.FailNow()
			} else {
				t.Log("springbootconventions labels present in podintent")
			}
			return ctx
		}).
		Assess("verify-podintent-alv-conventions", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying appliveview conventions present in podintent")

			// check
			appliveviewConventionsPresent := kubectl_helpers.ValidateAppLiveViewConventions(name, namespace)
			if !appliveviewConventionsPresent {
				t.Error("appliveview conventions absent in podintent")
				t.FailNow()
			} else {
				t.Log("appliveview conventions present in podintent")
			}
			return ctx
		}).
		Assess("verify-pod-intent-devloper-conventions-annotations", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verify if devloper-conventions annotations are added to podintent")
			status := kubectl_helpers.ValidateDeveloperConventions(name, namespace)
			t.Logf("devloper-conventions annotations status is : %t", status)
			if !status {
				t.Error(fmt.Errorf("devloper-conventions annotations are not added to the podintent"))
				t.FailNow()
			}
			return ctx
		}).
		Assess("verify-podintent-springbootconventions-conventions", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying springbootconventions conventions present in podintent")

			// check
			springbootconventionsConventionsPresent := kubectl_helpers.ValidateSpringBootConventions(name, namespace)
			if !springbootconventionsConventionsPresent {
				t.Error("springbootconventions conventions absent in podintent")
				t.FailNow()
			} else {
				t.Log("springbootconventions conventions present in podintent")
			}
			return ctx
		}).
		Feature()
}

func VerifyTanzuJavaWebAppImageRepositoryDelivery(t *testing.T, name string, imageDeliverySuffix string, namespace string) features.Feature {
	return features.New("verify-image-repository-delivery").
		Assess("verify-image-repositories", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verify image-repositories-delivery status")
			imageRepo := name + imageDeliverySuffix
			status := kubectl_helpers.VerifyImageRepositoryStatus(imageRepo, namespace, 10, 30)
			t.Logf("ImageRepository %s status is : %t", imageRepo, status)
			if !status {
				t.Error(fmt.Errorf("ImageRepository %s is not ready.", imageRepo))
				t.Fail()
			}
			return ctx
		}).
		Feature()
}

func VerifyTanzuJavaWebAppRevisionStatus(t *testing.T, name string, namespace string) features.Feature {
	return features.New("verify-revision-status").
		Assess("verify-revision-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying revision ready status")

			revisionName = kubectl_helpers.GetLatestRevision(name, namespace, 1, 30)
			t.Logf("latestRevision set to %s", revisionName)
			revisionReady := kubectl_helpers.ValidateRevisionStatus(revisionName, name, namespace, 5, 30)

			if !revisionReady {
				t.Error("revision not ready")
				t.FailNow()
			} else {
				t.Log("revision ready")
			}
			return ctx
		}).
		Feature()
}

func VerifyTanzuJavaWebAppKsvcStatus(t *testing.T, name string, namespace string) features.Feature {
	return features.New("verify-ksvc-status").
		Assess("verify-ksvc-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verifying ksvc ready status %s", ksvcLatestReady)

			ksvcReady := kubectl_helpers.VerifyKsvcStatus(name, namespace, revisionName, 5, 30)
			if !ksvcReady {
				t.Error("ksvc not ready")
				t.FailNow()
			} else {
				t.Log("ksvc ready")
			}

			return ctx
		}).
		Feature()
}

func VerifyTanzuJavaWebAppResponseBeforeChange(t *testing.T, workloadUrl string, originalString string, namespace string) features.Feature {
	return features.New("verify-app-response").
		Assess("get-externalip-and-check-webpage-for-string", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("getting external ip and checking for string")

			// get external IP
			url, err := client.GetServiceExternalIP("envoy", "tanzu-system-ingress", cfg.Client().RESTConfig(), 2, 30)
			if err != nil {
				t.Error("error while getting external IP")
				t.Fail() // DON'T DO t.FailNow() AS WE WANT TO CLEAN UP REGARDLESS OF THE STATE OF THE TEST
			} else {
				t.Log("external IP retrieved")
			}

			// set url
			if !strings.HasPrefix(url, "http://") {
				url = "http://" + url
			}

			webpageContainsString, _ := misc.VerifyWebpageContainsString(workloadUrl, url, originalString, 10, 10, 30)
			if !webpageContainsString {
				t.Error("webpage does not contains string")
				t.Fail() // DON'T DO t.FailNow() AS WE WANT TO CLEAN UP REGARDLESS OF THE STATE OF THE TEST
			} else {
				t.Log("webpage contains string")
			}

			return ctx
		}).
		Feature()
}

func VerifyTanzuJavaWebAppDeliverable(t *testing.T, name string, namespace string) features.Feature {
	return features.New("verify-deliverables").
		Assess("verify-deliverables-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying deliverables ready status")
			if !kubectl_helpers.ValidateDeliverables(name, namespace, 5, 30) {
				t.Error("deliverables not ready")
				t.FailNow()
			} else {
				t.Log("deliverables ready")
			}
			return ctx
		}).
		Feature()
}

func VerifyTanzuJavaWebAppGitRepository(t *testing.T, name string, namespace string) features.Feature {
	return features.New("verify-image-repositories").
		Assess("verify-image-repositories", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verify image-repositories status")
			status := kubectl_helpers.VerifyGitRepoStatus(name, namespace, 10, 30)
			t.Logf("ImageRepository %s status is : %t", name, status)
			if !status {
				t.Error(fmt.Errorf("ImageRepository %s is not ready.", name))
				t.Fail()
			}
			return ctx
		}).
		Feature()
}
