//go:build outerloop

package suite

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/client"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/exec"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_helpers"
	tanzu_lib "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzu_libs"
	"gopkg.in/yaml.v3"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
)

type outerloopConfiguration struct {
	CatalogInfoYaml    string `yaml:"catalog_info_yaml"`
	Clusterrolebinding struct {
		Name           string `yaml:"name"`
		Clusterrole    string `yaml:"clusterrole"`
		ServiceAccount string `yaml:"serviceAccount"`
	} `yaml:"clusterrolebinding"`
	Mysql struct {
		Name     string `yaml:"name"`
		YamlFile string `yaml:"yaml_file"`
	} `yaml:"mysql"`
	Namespace string `yaml:"namespace"`
	Project   struct {
		Application         string `yaml:"application"`
		File                string `yaml:"file"`
		Name                string `yaml:"name"`
		NewString           string `yaml:"new_string"`
		WebpageRelativePath string `yaml:"webpage_relative_path"`
		OriginalString      string `yaml:"original_string"`
		Repository          string `yaml:"repository"`
	} `yaml:"project"`
	ScanPolicy struct {
		YamlFile string `yaml:"yaml_file"`
	} `yaml:"scan_policy"`
	SpringPetclinic struct {
		BuildNamePrefix     string `yaml:"build_name_prefix"`
		GitrepositoryName   string `yaml:"gitrepository_name"`
		ImagerepositoryName string `yaml:"imagerepository_name"`
		KsvcName            string `yaml:"ksvc_name"`
		Name                string `yaml:"name"`
		PodintentName       string `yaml:"podintent_name"`
		TaskrunNamePrefix   string `yaml:"taskrun_name_prefix"`
		YamlFile            string `yaml:"yaml_file"`
	} `yaml:"spring_petclinic"`
	Workload struct {
		Name     string `yaml:"name"`
		YamlFile string `yaml:"yaml_file"`
	} `yaml:"workload"`
}

var outerloopResourcesDir = filepath.Join(resourcesDir, "outerloop")

func getOuterloopConfig() (outerloopConfiguration, error) {
	outerloopConfig := outerloopConfiguration{}

	// read file
	outerloopConfigBytes, err := os.ReadFile(filepath.Join(outerloopResourcesDir, "outerloop-config.yaml"))
	if err != nil {
		return outerloopConfig, fmt.Errorf("error while reading outerloop config file: %w", err)
	}
	err = yaml.Unmarshal(outerloopConfigBytes, &outerloopConfig)
	if err != nil {
		return outerloopConfig, fmt.Errorf("error while unmarshalling outerloop config file: %w", err)
	}

	// update outerloop config for full file paths
	outerloopConfig.Mysql.YamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.Mysql.YamlFile)
	outerloopConfig.ScanPolicy.YamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.ScanPolicy.YamlFile)
	outerloopConfig.SpringPetclinic.YamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.SpringPetclinic.YamlFile)
	outerloopConfig.Workload.YamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.Workload.YamlFile)

	return outerloopConfig, nil
}

var outerloopConfig, _ = getOuterloopConfig()

var deployMysqlService = features.New("deploy-mysql-service-app-via-yaml-configurations").
	Assess("deploy-mysql-service", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		name, files, namespace := outerloopConfig.Mysql.Name, []string{outerloopConfig.Mysql.YamlFile}, outerloopConfig.Namespace

		t.Logf("deploying mysql service %s in namespace %s", name, namespace)
		cmd, output, err := exec.KappDeployAppInNamespace(name, files, namespace)
		t.Logf("command executed: %s", cmd)
		if err != nil {
			t.Error(fmt.Errorf("error while deploying mysql service %s in namespace %s: %w: %s", name, namespace, err, output))
			t.FailNow()
		}
		t.Logf("mysql service %s deployed in namespace %s: %s", name, namespace, output)
		return ctx
	}).
	Feature()

var deployWorkload = features.New("deploy-workload").
	Assess("deploy-workload", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		name, file, namespace := outerloopConfig.Workload.Name, outerloopConfig.Workload.YamlFile, outerloopConfig.Namespace

		t.Logf("deploying workload %s in namespace %s", name, namespace)
		cmd, output, err := exec.TanzuDeployWorkload(file, namespace)
		t.Logf("command executed: %s", cmd)
		if err != nil {
			t.Error(fmt.Errorf("error while deploying workload %s using %s in namespace %s: %w: %s", name, file, namespace, err, output))
			t.FailNow()
		}
		t.Logf("workload %s deployed in namespace %s: %s", name, namespace, output)
		return ctx
	}).
	Feature()

var verifyGitrepoStatus = features.New("verify-gitrepo-status").
	Assess("verify-gitrepo-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		if !kubectl_helpers.VerifyGitRepoStatus(outerloopConfig.SpringPetclinic.PodintentName, outerloopConfig.Namespace) {
			t.Error(fmt.Errorf("gitrepo %s is not ready.", outerloopConfig.SpringPetclinic.GitrepositoryName))
			t.FailNow()
		}
		return ctx
	}).
	Feature()

var verifyBuildStatus = features.New("verify-build-status").
	Assess("verify-build-succeeded", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		if !kubectl_helpers.VerifyBuildStatus(outerloopConfig.Namespace) {
			t.Error(fmt.Errorf("build is not in succeeded status for namespace %s", outerloopConfig.Namespace))
			t.FailNow()
		}
		return ctx

	}).
	Feature()

var verifyPodintents = features.New("verify-pod-intents-annotations-labels").
	Assess("verify-podintent-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		if kubectl_helpers.GetPodIntentStatus(outerloopConfig.SpringPetclinic.PodintentName, outerloopConfig.Namespace) != "True" {
			t.Error(fmt.Errorf("podintent %s is not ready.", outerloopConfig.SpringPetclinic.PodintentName))
			t.FailNow()
		}
		return ctx
	}).
	Assess("verify-podintent-alv-lables", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		if !kubectl_helpers.ValidateAppLiveViewLabels(outerloopConfig.SpringPetclinic.PodintentName, outerloopConfig.Namespace) {
			t.Error(fmt.Errorf("app live view lables are not added to the podintent"))
			t.FailNow()
		}
		return ctx
	}).
	Assess("verify-podintent-springbootconventions-lables", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		if !kubectl_helpers.ValidateSpringBootLabels(outerloopConfig.SpringPetclinic.PodintentName, outerloopConfig.Namespace) {
			t.Error(fmt.Errorf("spring boot conventions lables are not added to the podintent"))
			t.FailNow()
		}
		return ctx
	}).
	Assess("verify-podintent-alv-annotations", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		if !kubectl_helpers.ValidateAppLiveViewConventions(outerloopConfig.SpringPetclinic.PodintentName, outerloopConfig.Namespace) {
			t.Error(fmt.Errorf("app live view annotations are not added to the podintent"))
			t.FailNow()
		}
		return ctx
	}).
	Assess("verify-podintent-springbootconventions-annotations", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		if !kubectl_helpers.ValidateSpringBootConventions(outerloopConfig.SpringPetclinic.PodintentName, outerloopConfig.Namespace) {
			t.Error(fmt.Errorf("spring-boot-conventions annotations are not added to the podintent"))
			t.FailNow()
		}
		return ctx
	}).
	Feature()

var verifyKsvcStatus = features.New("verify-ksvc-status").
	Assess("verify-ksvc-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		if kubectl_helpers.GetKsvcStatus(outerloopConfig.SpringPetclinic.KsvcName, outerloopConfig.Namespace) != "True" {
			t.Error(fmt.Errorf("ksvc %s is not ready", outerloopConfig.SpringPetclinic.KsvcName))
			t.FailNow()
		}
		return ctx

		// return stepfuncs.VerifyKsvcReady(ctx, t, cfg, outerloopConfig.SpringPetclinic.KsvcName, outerloopConfig.Namespace)
	}).
	Feature()

var verifyTaskrunStatus = features.New("verify-taskrun-status").
	Assess("verify-taskrun-succeeded", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		if !kubectl_helpers.VerifyTaskrunStatus(outerloopConfig.Namespace) {
			t.Error(fmt.Errorf("taskrun is not in succeeded status for namespace %s", outerloopConfig.Namespace))
			t.FailNow()
		}
		return ctx

		// return stepfuncs.VerifyTaskrunSucceeded(ctx, t, cfg, outerloopConfig.SpringPetclinic.TaskrunNamePrefix, outerloopConfig.Namespace)
	}).
	Feature()

var verifyWorkloadStatus = features.New("verify-workload-status").
	Assess("verify-workload-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		if kubectl_helpers.GetWorkloadStatus(outerloopConfig.Workload.Name, outerloopConfig.Namespace) != "True" {
			t.Error(fmt.Errorf("workload %s is not ready", outerloopConfig.Workload.Name))
			t.Fail()
		}
		return ctx
	}).
	Feature()

var ingressEnvoyExternalIpKey = "ingressEnvoyExternalIp"

var getEnvoyExternalIP = features.New("get-ingress-envoy-externalip-port").
	Assess("get-externalip", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		service, namespace := "envoy", "tanzu-system-ingress"

		t.Logf("getting external ip for %s (namespace %s)", service, namespace)
		serviceExternalIp, err := client.GetServiceExternalIP(service, namespace, cfg.Client().RESTConfig())
		if err != nil {
			t.Error(fmt.Errorf("error while getting external ip for %s (namespace %s): %w", service, namespace, err))
			t.FailNow()
		}
		t.Logf("external ip for %s (namespace %s): %s", "server", namespace, serviceExternalIp)
		return context.WithValue(ctx, ingressEnvoyExternalIpKey, serviceExternalIp)

		// ctx, serviceExternalIp := stepfuncs.GetServiceExternalIp(ctx, t, cfg, "envoy", "tanzu-system-ingress")
		// return context.WithValue(ctx, ingressEnvoyExternalIpKey, serviceExternalIp)
	}).
	Feature()

var verifyApplicationRunningOriginal = features.New("verify-application-running").
	Assess("check-for-original-string", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		url, host, webpageRelativePath, validationString := ctx.Value(ingressEnvoyExternalIpKey).(string), outerloopConfig.Project.Application, outerloopConfig.Project.WebpageRelativePath, outerloopConfig.Project.OriginalString

		t.Logf("checking application %s for result: %s", host, validationString)
		validated := false
		iter := 10
		for i := 1; i <= iter; i++ {
			url := fmt.Sprintf("%s/%s", url, webpageRelativePath)
			if !strings.HasPrefix(url, "http://") {
				url = "http://" + url
			}
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Error(fmt.Errorf("error while giving http request: %w", err))
				t.FailNow()
			}
			req.Host = host

			var retries int = 10
			for retries > 0 {
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					retries -= 1
					t.Logf("didn't get response")
					t.Logf("sleeping for 30 seconds")
					time.Sleep(30 * time.Second)
				} else {
					t.Logf("status code is: %d", resp.StatusCode)
					break
				}
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(fmt.Errorf("error while giving http response: %w", err))
				t.FailNow()
			}
			if resp.StatusCode != http.StatusOK {
				t.Error(fmt.Errorf("bad HTTP Response: %s", resp.Status))
				t.FailNow()
			}
			defer resp.Body.Close()
			resultStringBytes, _ := ioutil.ReadAll(resp.Body)
			resultString := string(resultStringBytes)
			t.Logf(resultString)
			if strings.Contains(resultString, validationString) {
				t.Logf("application %s validated, got result: %s", host, validationString)
				validated = true
				break
			} else {
				t.Logf("getting string %s", resultString)
				t.Logf("sleeping for 1 minute")
				time.Sleep(1 * time.Minute)
			}
		}

		if !validated {
			t.Errorf(`application %s not validated %d iterations`, host, iter)
			t.FailNow()
		}
		return ctx

		// return stepfuncs.VerifyApplicationRunningWithValidationString(ctx, t, cfg, ctx.Value(ingressEnvoyExternalIpKey).(string), outerloopConfig.Project.Application, outerloopConfig.Project.OriginalString)
	}).
	Feature()

var gitUpdate = features.New("git-update").
	Assess("git-clone", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		repo, path := outerloopConfig.Project.Repository, utils.GetFileDir()

		t.Logf("cloning repository %s at %s", repo, path)
		cmd, output, err := exec.GitClone(path, repo)
		t.Logf("command executed: %s", cmd)
		if err != nil {
			t.Error(fmt.Errorf("error while cloning repository %s at %s: %w: %s", repo, path, err, output))
			t.FailNow()
		}
		t.Logf("repository %s cloned at %s: %s", repo, path, output)
		return ctx

		// return stepfuncs.GitClone(ctx, t, cfg, GetFileDir(), outerloopConfig.Project.Repository)
	}).
	Assess("modify-file", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		file, originalString, newString := filepath.Join(utils.GetFileDir(), outerloopConfig.Project.Name, outerloopConfig.Project.File), outerloopConfig.Project.OriginalString, outerloopConfig.Project.NewString

		t.Logf("updating file %s", file)
		inputBytes, err := os.ReadFile(file)
		if err != nil {
			t.Error(fmt.Errorf("error while updating file %s: %w", file, err))
			t.FailNow()
		}
		input := strings.ReplaceAll(string(inputBytes), originalString, newString)

		err = os.WriteFile(file, []byte(input), 0677)
		if err != nil {
			t.Error(fmt.Errorf("error while writing file %s: %w", file, err))
			t.FailNow()
		}
		t.Logf("file %s written", file)
		return ctx

		// return stepfuncs.UpdateFileReplaceString(ctx, t, cfg, filepath.Join(GetFileDir(), outerloopConfig.Project.Name, outerloopConfig.Project.File), outerloopConfig.Project.OriginalString, outerloopConfig.Project.NewString)
	}).
	Assess("git-add", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		path, files := filepath.Join(utils.GetFileDir(), outerloopConfig.Project.Name), []string{outerloopConfig.Project.File}

		t.Logf("adding files %s for repository at %s", files, path)
		cmd, output, err := exec.GitAdd(path, files)
		t.Logf("command executed: %s", cmd)
		if err != nil {
			t.Error(fmt.Errorf("error while adding files %s for repository at %s: %w: %s", files, path, err, output))
			t.FailNow()
		}
		t.Logf("files %s added for repository at %s: %s", files, path, output)
		return ctx

		// return stepfuncs.GitAdd(ctx, t, cfg, filepath.Join(GetFileDir(), outerloopConfig.Project.Name), []string{outerloopConfig.Project.File})
	}).
	Assess("git-commit", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		path, message := filepath.Join(utils.GetFileDir(), outerloopConfig.Project.Name), "changes to webpage"

		t.Logf("committing files for repository at %s (message %s)", path, message)
		cmd, output, err := exec.GitCommit(path, message)
		t.Logf("command executed: %s", cmd)
		if err != nil {
			t.Error(fmt.Errorf("error while committing files for repository at %s: %w: %s", path, err, output))
			t.FailNow()
		}
		t.Logf("committed files for repository at %s (message %s): %s", path, message, output)
		return ctx

		// return stepfuncs.GitCommit(ctx, t, cfg, filepath.Join(GetFileDir(), outerloopConfig.Project.Name), "change to webpage")
	}).

	// NOTE: DON'T DO t.FailNow() AS WE WANT TO REVERT CHANGES TO REPO REGARDLESS OF THE STATE OF THE TEST. USE t.Fail() INSTEAD.
	Assess("git-push", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		path, force := filepath.Join(utils.GetFileDir(), outerloopConfig.Project.Name), false

		t.Logf("pushing commits for repository at %s", path)
		cmd, output, err := exec.GitPush(path, force)
		t.Logf("command executed: %s", cmd)
		if err != nil {
			t.Error(fmt.Errorf("error while pushing commits for repository at %s: %w: %s", path, err, output))
			t.Fail()
		}
		t.Logf("pushed commits for repository at %s: %s", path, output)
		return ctx

		// return stepfuncs.GitPush(ctx, t, cfg, filepath.Join(GetFileDir(), outerloopConfig.Project.Name), false)
	}).
	Feature()

var verifyApplicationRunningNew = features.New("verify-application-running").
	// NOTE: DON'T DO t.FailNow() AS WE WANT TO REVERT CHANGES TO REPO REGARDLESS OF THE STATE OF THE TEST. USE t.Fail() INSTEAD.
	Assess("check-for-new-string", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		url, host, webpageRelativePath, validationString := ctx.Value(ingressEnvoyExternalIpKey).(string), outerloopConfig.Project.Application, outerloopConfig.Project.WebpageRelativePath, outerloopConfig.Project.NewString

		t.Logf("checking application %s for result: %s", host, validationString)
		validated := false
		iter := 10
		for i := 1; i <= iter; i++ {
			url := fmt.Sprintf("%s/%s", url, webpageRelativePath)
			if !strings.HasPrefix(url, "http://") {
				url = "http://" + url
			}
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Error(fmt.Errorf("error while giving http request: %w", err))
				t.FailNow()
			}
			req.Host = host

			var retries int = 10
			for retries > 0 {
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					retries -= 1
					t.Logf("didn't get response")
					t.Logf("sleeping for 30 seconds")
					time.Sleep(30 * time.Second)
				} else {
					t.Logf("status code is: %d", resp.StatusCode)
					break
				}
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(fmt.Errorf("error while giving http response: %w", err))
				t.FailNow()
			}
			if resp.StatusCode != http.StatusOK {
				t.Error(fmt.Errorf("bad HTTP Response: %s", resp.Status))
				t.FailNow()
			}
			defer resp.Body.Close()
			resultStringBytes, _ := ioutil.ReadAll(resp.Body)
			resultString := string(resultStringBytes)
			t.Logf(resultString)
			if strings.Contains(resultString, validationString) {
				t.Logf("application %s validated, got result: %s", host, validationString)
				validated = true
				break
			} else {
				t.Logf("getting string %s", resultString)
				t.Logf("sleeping for 1 minute")
				time.Sleep(1 * time.Minute)
			}
		}

		if !validated {
			t.Errorf(`application %s not validated %d iterations`, host, iter)
			t.FailNow()
		}
		return ctx

		// return stepfuncs.VerifyApplicationRunningWithValidationString(ctx, t, cfg, ctx.Value(ingressEnvoyExternalIpKey).(string), outerloopConfig.Project.Application, outerloopConfig.Project.NewString)
	}).
	Feature()
var gitReset = features.New("git-reset").
	Assess("git-reset", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		path, count := filepath.Join(utils.GetFileDir(), outerloopConfig.Project.Name), 1

		t.Logf("resetting commits at HEAD~%d for repository at %s", count, path)
		cmd, output, err := exec.GitResetFromHead(path, count)
		t.Logf("command executed: %s", cmd)
		if err != nil {
			t.Error(fmt.Errorf("error while resetting commits at HEAD~%d for repository at %s: %w: %s", count, path, err, output))
			t.FailNow()
		}
		t.Logf("resetted commits at HEAD~%d for repository at %s: %s", count, path, output)
		return ctx

		// return stepfuncs.GitResetFromHead(ctx, t, cfg, filepath.Join(GetFileDir(), outerloopConfig.Project.Name), 1)
	}).
	Assess("git-push-force", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		path, force := filepath.Join(utils.GetFileDir(), outerloopConfig.Project.Name), true

		t.Logf("pushing commits for repository at %s", path)
		cmd, output, err := exec.GitPush(path, force)
		t.Logf("command executed: %s", cmd)
		if err != nil {
			t.Error(fmt.Errorf("error while pushing commits for repository at %s: %w: %s", path, err, output))
			t.FailNow()
		}
		t.Logf("pushed commits for repository at %s: %s", path, output)
		return ctx

		// return stepfuncs.GitPush(ctx, t, cfg, filepath.Join(GetFileDir(), outerloopConfig.Project.Name), true)
	}).
	Feature()

var cleanRemoveProjectDir = features.New("clean-remove-project-dir").
	Assess("remove-dir", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		dir := outerloopConfig.Project.Name

		t.Logf("removing directory %s", dir)
		err := os.RemoveAll(dir)
		if err != nil {
			t.Error(fmt.Errorf("error while removing directory %s: %w", dir, err))
			t.FailNow()
		}
		t.Logf("directory %s removed", dir)
		return ctx

		// return stepfuncs.RemoveDirectory(ctx, t, cfg, filepath.Join(GetFileDir(), outerloopConfig.Project.Name))
	}).
	Feature()

var deleteWorkload = features.New("delete-workload").
	Assess("delete-workload", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		name, file, namespace := outerloopConfig.Workload.Name, outerloopConfig.Workload.YamlFile, outerloopConfig.Namespace

		t.Logf("deleting workload %s from namespace %s", file, namespace)
		tanzu_lib.DeleteWorkload(name, namespace)
		return ctx
	}).
	Feature()
