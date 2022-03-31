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

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/git"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/github"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectlCmds"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_helpers"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubernetes/client"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/misc"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzuCmds"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzu_helpers"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzu_libs"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/models"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

var tiltprocCmdKey = "tiltprocCmd"
var rootDir = filepath.Join(utils.GetFileDir(), "../../")

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

func UpdatePackageRepository(t *testing.T, name string, registry string, version string, namespace string) features.Feature {
	return features.New("updating package repository").
		Assess(fmt.Sprintf("updating-packaging-repository-%s", name), func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			log.Printf("updating pacakage repository %s", name)
			tanzu_libs.TanzuUpdatePackageRepository(name, registry, version, namespace)
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

func UpdateTapVersion(t *testing.T, name string, tapPackageName string, namespace string, tapVersion string, pollTimeout string) features.Feature {
	return features.New(fmt.Sprintf("updating-tap-version-%s", tapVersion)).
		Assess(fmt.Sprintf("updating-tap-package-%s", tapVersion), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			log.Printf("updating tap- %s ", tapVersion)
			tanzu_libs.UpdateInstalledPackage(name, tapPackageName, tapVersion, namespace, "", pollTimeout)
			updated := tanzu_helpers.ValidateInstalledPackageVersion(name, namespace, tapVersion, 5, 30)
			if updated {
				t.Logf("Updated tap version: %s successfully", tapVersion)
			} else {
				t.Error(fmt.Errorf("update FAILED for tap version: %s", tapVersion))
				t.Fail()
			}
			availablePkgs := tanzu_libs.ListInstalledPackages(namespace)
			for _, pkg := range availablePkgs {
				installed := tanzu_helpers.ValidateInstalledPackageStatus(pkg.NAME, namespace, 5, 30)
				if installed {
					t.Logf("Installed package : %s, version: %s successfully", pkg.NAME, pkg.PACKAGE_VERSION)
				} else {
					t.Error(fmt.Errorf("Installation FAILED for package : %s, version: %s", pkg.NAME, pkg.PACKAGE_VERSION))
					t.Fail()
				}
			}
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
			url, err := client.GetServiceExternalIP("envoy", "tanzu-system-ingress", cfg.Client().RESTConfig(), 2, 30)
			if err != nil {
				t.Error("error while getting external IP")
				t.Fail() // DON'T DO t.FailNow() AS WE WANT TO CLEAN UP REGARDLESS OF THE STATE OF THE TEST
			} else {
				t.Log("external IP retrieved")
			}

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
