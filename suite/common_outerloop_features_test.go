//go:build all || outerloop || outerloop_basic || outerloop_testing || outerloop_testing_scanning || outerloop_basic_delivery || outerloop_scan_multiple_apps

package suite

import (
	"context"
	"encoding/base64"
	"fmt"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/git"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/github"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectlCmds"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_helpers"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_libs"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubernetes/client"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/misc"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzuCmds"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzu_helpers"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzu_libs"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"log"
	"net"
	"os"
	"path/filepath"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	"strings"
	"testing"
	"time"
)

type outerloopConfiguration struct {
	CatalogInfoYaml string `yaml:"catalog_info_yaml"`
	Mysql           struct {
		YamlFile string `yaml:"yaml_file"`
	} `yaml:"mysql"`
	Namespace string `yaml:"namespace"`
	Project   struct {
		Host                string `yaml:"host"`
		WebpageRelativePath string `yaml:"webpage_relative_path"`
		File                string `yaml:"file"`
		Name                string `yaml:"name"`
		DestName            string `yaml:"dest_name"`
		RepoTemplate        string `yaml:"repo_template"`
		DestRepoTemplate    string `yaml:"dest_repo_template"`
		NewString           string `yaml:"new_string"`
		OriginalString      string `yaml:"original_string"`
		CommitMessage       string `yaml:"commit_message"`
		Repository          string `yaml:"repository"`
		Username            string `yaml:"username"`
		Email               string `yaml:"email"`
		AccessToken         string `yaml:"access_token"`
	} `yaml:"project"`
	ScanPolicy struct {
		YamlFile string `yaml:"yaml_file"`
	} `yaml:"scan_policy"`
	SpringPetclinicPipeline struct {
		YamlFile string `yaml:"yaml_file"`
	} `yaml:"spring_petclinic_pipeline"`
	Workload struct {
		Name                 string `yaml:"name"`
		YamlFile             string `yaml:"yaml_file"`
		TestYamlFile         string `yaml:"test_yaml_file"`
		BuildNameSuffix      string `yaml:"build_name_suffix"`
		PipelineName         string `yaml:"pipeline_name"`
		TaskRunInfix         string `yaml:"taskrun_name_infix"`
		TaskRunTestSuffix    string `yaml:"taskrun_test_suffix"`
		ServiceBindingSuffix string `yaml:"service_binding_suffix"`
		GitopsYamlFile       string `yaml:"gitops_yaml_file"`
		GitSSHSecretYamlFile string `yaml:"gitssh_secret_yaml_file"`
	} `yaml:"workload"`
	BuildPacks struct {
		ScanPolicy string `yaml:"scan_policy"`
		Workloads  []struct {
			Name                string `yaml:"name"`
			GitRepository       string `yaml:"git_repository"`
			GitBranch           string `yaml:"git_branch"`
			WebpageRelativePath string `yaml:"webpage_relative_path"`
			ContainsConventions bool   `yanl:"contains_conventions"`
		} `yaml:"workloads"`
	} `yaml:"buildpacks"`
}

var outerloopResourcesDir = filepath.Join(utils.GetFileDir(), "resources", "outerloop")

func getOuterloopConfig() (outerloopConfiguration, error) {
	log.Printf("getting outerloop config")

	outerloopConfig := outerloopConfiguration{}
	file := filepath.Join(outerloopResourcesDir, "outerloop-config.yaml")

	// read file
	outerloopConfigBytes, err := os.ReadFile(file)
	if err != nil {
		log.Printf("error while reading outerloop config file %s", file)
		log.Printf("error: %s", err)
		return outerloopConfig, err
	} else {
		log.Printf("read outerloop config file %s", file)
	}

	// unmarshall
	err = yaml.Unmarshal(outerloopConfigBytes, &outerloopConfig)
	if err != nil {
		log.Printf("error while unmarshalling outerloop config file %s", file)
		log.Printf("error: %s", err)
		return outerloopConfig, err
	} else {
		log.Printf("unmarshalled file %s", file)
	}

	// update outerloop config for full file paths
	outerloopConfig.Mysql.YamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.Mysql.YamlFile)
	outerloopConfig.ScanPolicy.YamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.ScanPolicy.YamlFile)
	outerloopConfig.SpringPetclinicPipeline.YamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.SpringPetclinicPipeline.YamlFile)
	outerloopConfig.Workload.YamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.Workload.YamlFile)
	outerloopConfig.Workload.TestYamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.Workload.TestYamlFile)
	outerloopConfig.Workload.GitopsYamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.Workload.GitopsYamlFile)
	outerloopConfig.Workload.GitSSHSecretYamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.Workload.GitSSHSecretYamlFile)
	outerloopConfig.BuildPacks.ScanPolicy = filepath.Join(outerloopResourcesDir, outerloopConfig.BuildPacks.ScanPolicy)
	return outerloopConfig, nil
}

var outerloopConfig, _ = getOuterloopConfig()

var deployMysqldbService = features.New("deploy-mysqldb-service-via-yaml").
	Assess("deploy-mysqldb-service", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("deploying mysqldb service")

		// deploy app
		err := kubectlCmds.KubectlApplyConfiguration(outerloopConfig.Mysql.YamlFile, outerloopConfig.Namespace)
		if err != nil {
			t.Error("error while deploying mysqldb service")
			t.FailNow()
		} else {
			t.Log("deployed mysqldb service")
		}

		return ctx
	}).
	Feature()

var deploySpringpetclinicPipeline = features.New("deploy-pipeline-app-via-yaml-configurations").
	Assess("deploy-springpetclinic-pipeline", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("deploying springpetclinic-pipeline")

		// deploy app
		err := kubectlCmds.KubectlApplyConfiguration(outerloopConfig.SpringPetclinicPipeline.YamlFile, outerloopConfig.Namespace)
		if err != nil {
			t.Error("error while deploying springpetclinic-pipeline")
			t.FailNow()
		} else {
			t.Log("deployed springpetclinic-pipeline")
		}

		return ctx
	}).
	Feature()

var deploySpringpetclinicPipelineWithGitops = features.New("deploy-pipeline-app-via-yaml-configurations-with-gitops").
	Assess("deploy-springpetclinic-pipeline", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("deploying springpetclinic-pipeline")

		// deploy app
		err := kubectlCmds.KubectlApplyConfiguration(outerloopConfig.Workload.GitopsYamlFile, outerloopConfig.Namespace)
		if err != nil {
			t.Error("error while deploying springpetclinic-pipeline")
			t.FailNow()
		} else {
			t.Log("deployed springpetclinic-pipeline")
		}

		return ctx
	}).
	Feature()

var deployScanPolicy = features.New("deploy-scan-policy-via-yaml").
	Assess("deploy-scanpolicy", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("deploying scanpolicy")

		// deploy app
		err := kubectlCmds.KubectlApplyConfiguration(outerloopConfig.ScanPolicy.YamlFile, outerloopConfig.Namespace)
		if err != nil {
			t.Error("error while deploying scanpolicy")
			t.FailNow()
		} else {
			t.Log("deployed scanpolicy")
		}

		return ctx
	}).
	Feature()

var deployLenientScanPolicy = features.New("deploy-scan-policy-via-yaml").
	Assess("deploy-scanpolicy", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("deploying scanpolicy")

		// deploy app
		err := kubectlCmds.KubectlApplyConfiguration(outerloopConfig.BuildPacks.ScanPolicy, outerloopConfig.Namespace)
		if err != nil {
			t.Error("error while deploying scanpolicy")
			t.FailNow()
		} else {
			t.Log("deployed scanpolicy")
		}

		return ctx
	}).
	Feature()

var deployWorkload = features.New("deploy-workload").
	Assess("deploy-workload", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("deploying workload")

		// deploy workload
		err := tanzuCmds.TanzuDeployWorkload(outerloopConfig.Workload.YamlFile, outerloopConfig.Namespace)
		if err != nil {
			t.Error("error while deploying workload")
			t.FailNow()
		} else {
			t.Log("deployed workload")
		}

		return ctx
	}).
	Feature()

var deployWorkloadWithTest = features.New("deploy-workload-with-test").
	Assess("deploy-workload-test", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("deploying workload")

		// deploy workload
		err := tanzuCmds.TanzuDeployWorkload(outerloopConfig.Workload.TestYamlFile, outerloopConfig.Namespace)
		if err != nil {
			t.Error("error while deploying workload")
			t.FailNow()
		} else {
			t.Log("deployed workload")
		}

		return ctx
	}).
	Feature()

var verifyGrypePackageInstalled = features.New("check-grype-package-installed").
	Assess("check-grype", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("checking grype is installed")

		// check
		grypeInstalled := tanzu_helpers.IsGrypeInstalled(suiteConfig.Tap.Namespace)
		if !grypeInstalled {
			t.Error("grype is not installed")
			t.FailNow()
		} else {
			t.Log("grype is installed")
		}

		return ctx
	}).
	Feature()

var verifyScanningPackageInstalled = features.New("check-scanning-package-installed").
	Assess("check-scanning", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("checking scanning is installed")

		// check
		scanningInstalled := tanzu_helpers.IsScanningInstalled(suiteConfig.Tap.Namespace)
		if !scanningInstalled {
			t.Error("scanning is not installed")
			t.FailNow()
		} else {
			t.Log("scanning is installed")
		}

		return ctx
	}).
	Feature()

var verifyPipelineStatus = features.New("verify-pipeline-status").
	Assess("verify-pipeline-installed", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying pipeline status")

		// check
		pipelineInstalled := kubectl_helpers.ValidatePipelineExists(outerloopConfig.Workload.PipelineName, outerloopConfig.Namespace, 5, 30)
		if !pipelineInstalled {
			t.Error("pipeline not installed")
			t.FailNow()
		} else {
			t.Log("pipeline installed")
		}

		return ctx
	}).
	Feature()

var verifySourceScanStatus = features.New("verify-source-scan-status").
	Assess("verify-source-scan-completed", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying source scan status")

		// check
		sourceScanCompleted := kubectl_helpers.ValidateSourceScans(outerloopConfig.Workload.Name, outerloopConfig.Namespace, 5, 30)
		if !sourceScanCompleted {
			t.Error("source scan completed")
			t.FailNow()
		} else {
			t.Log("source scan completed successfully")
		}

		return ctx
	}).
	Feature()

var verifyImageScanStatus = features.New("verify-imagescan-status").
	Assess("verify-imagescan-completed", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying image scan status")

		// check
		imageScanCompleted := kubectl_helpers.ValidateImageScans(outerloopConfig.Workload.Name, outerloopConfig.Namespace, 5, 30)
		if !imageScanCompleted {
			t.Error("image scan completed")
			t.FailNow()
		} else {
			t.Log("image scan completed successfully")
		}

		return ctx
	}).
	Feature()

var verifyPipelineRunStatus = features.New("verify-pipeline-runs-status").
	Assess("verify-pipeline-runs-succeded", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying pipeline runs status")

		// check
		pipelineRunSucceeded := kubectl_helpers.ValidatePipelineRuns(outerloopConfig.Workload.Name, outerloopConfig.Namespace, 5, 30)
		if !pipelineRunSucceeded {
			t.Error("pipeline runs not succeeded")
			t.FailNow()
		} else {
			t.Log("pipeline runs succeeded")
		}

		return ctx
	}).
	Feature()

var verifyImageskpac = features.New("verify-images.kpac-status").
	Assess("verify-images.kpac-true", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Logf("verifying latest image status")

		// check
		if !kubectl_helpers.ValidateLatestImageStatus(outerloopConfig.Namespace, 15, 60) {
			t.Error("image status is not true")
			t.FailNow()
		} else {
			t.Log("image status is true")
		}

		return ctx
	}).
	Feature()

var verifyGitrepoStatus = features.New("verify-gitrepo-status").
	Assess("verify-gitrepo-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying gitrepo ready status")

		// check
		gitrepoReady := kubectl_helpers.VerifyGitRepoStatus(outerloopConfig.Workload.Name, outerloopConfig.Namespace, 5, 30)
		if !gitrepoReady {
			t.Error("gitrepo not ready")
			t.FailNow()
		} else {
			t.Log("gitrepo ready")
		}

		return ctx
	}).
	Feature()

var verifyBuildStatus = features.New("verify-build-status").
	Assess("verify-build-succeeded", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying build succeeded status")

		buildName = fmt.Sprintf("%s%s", outerloopConfig.Workload.Name, outerloopConfig.Workload.BuildNameSuffix)
		buildSucceeded := kubectl_helpers.VerifyBuildStatus(buildName, outerloopConfig.Namespace, 15, 60)
		if !buildSucceeded {
			t.Error("build not succeeded")
			t.FailNow()
		} else {
			t.Log("build succeeded")
		}

		return ctx
	}).
	Feature()

var verifyPodintents = features.New("verify-podintents-labels-conventions").
	Assess("verify-podintent-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying podintent ready status")

		// check
		if !kubectl_helpers.VerifyPodIntentStatus(outerloopConfig.Workload.Name, outerloopConfig.Namespace, 5, 30) {
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
		alvLabelsPresent := kubectl_helpers.ValidateAppLiveViewLabels(outerloopConfig.Workload.Name, outerloopConfig.Namespace)
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
		springbootconventionsLabelsPresent := kubectl_helpers.ValidateSpringBootLabels(outerloopConfig.Workload.Name, outerloopConfig.Namespace)
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
		appliveviewConventionsPresent := kubectl_helpers.ValidateAppLiveViewConventions(outerloopConfig.Workload.Name, outerloopConfig.Namespace)
		if !appliveviewConventionsPresent {
			t.Error("appliveview conventions absent in podintent")
			t.FailNow()
		} else {
			t.Log("appliveview conventions present in podintent")
		}

		return ctx
	}).
	Assess("verify-podintent-springbootconventions-conventions", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying springbootconventions conventions present in podintent")

		// check
		springbootconventionsConventionsPresent := kubectl_helpers.ValidateSpringBootConventions(outerloopConfig.Workload.Name, outerloopConfig.Namespace)
		if !springbootconventionsConventionsPresent {
			t.Error("springbootconventions conventions absent in podintent")
			t.FailNow()
		} else {
			t.Log("springbootconventions conventions present in podintent")
		}

		return ctx
	}).
	Feature()

var verifyRevisionStatus = features.New("verify-revision-status").
	Assess("verify-revision-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying revision ready status")

		revisionName = kubectl_helpers.GetLatestRevision(outerloopConfig.Workload.Name, outerloopConfig.Namespace, 1, 30)
		revisionReady := kubectl_helpers.ValidateRevisionStatus(revisionName, outerloopConfig.Workload.Name, outerloopConfig.Namespace, 5, 30)
		if !revisionReady {
			t.Error("revision not ready")
			t.FailNow()
		} else {
			t.Log("revision ready")
		}
		return ctx
	}).
	Feature()

var verifyKsvcStatus = features.New("verify-ksvc-status").
	Assess("verify-ksvc-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying ksvc ready status")

		ksvcReady := kubectl_helpers.VerifyKsvcStatus(outerloopConfig.Workload.Name, outerloopConfig.Namespace, revisionName, 5, 30)
		if !ksvcReady {
			t.Error("ksvc not ready")
			t.FailNow()
		} else {
			t.Log("ksvc ready")
		}

		return ctx
	}).
	Feature()

var verifyTaskrunStatus = features.New("verify-taskrun-status").
	Assess("verify-taskrun-succeeded", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying taskrun succeeded status")

		taskRunPrefix := fmt.Sprintf("%s%s", outerloopConfig.Workload.Name, outerloopConfig.Workload.TaskRunInfix)
		taskrunSucceeded := kubectl_helpers.VerifyTaskrunStatus(taskRunPrefix, outerloopConfig.Namespace, 5, 30)
		if !taskrunSucceeded {
			t.Error("taskrun not succeeded")
			t.FailNow()
		} else {
			t.Log("taskrun succeeded")
		}

		return ctx
	}).
	Feature()

var verifyTestTaskrunStatus = features.New("verify-test-taskrun-status").
	Assess("verify-test-taskrun-succeeded", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying test taskrun succeeded status")

		taskrunSucceeded := kubectl_helpers.VerifyTestTaskrunStatus(outerloopConfig.Workload.Name, outerloopConfig.Workload.TaskRunTestSuffix, outerloopConfig.Namespace, 5, 30)
		if !taskrunSucceeded {
			t.Error("taskrun not succeeded")
			t.FailNow()
		} else {
			t.Log("taskrun succeeded")
		}

		return ctx
	}).
	Feature()

var verifyWorkloadStatus = features.New("verify-workload-status").
	Assess("verify-workload-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying workload ready status")

		// check
		workloadStatus := kubectl_helpers.GetWorkloadStatus(outerloopConfig.Workload.Name, outerloopConfig.Namespace)
		if workloadStatus != "True" {
			t.Error("workload not ready")
			t.FailNow()
		} else {
			t.Log("workload ready")
		}

		return ctx
	}).
	Feature()

var verifyWebpageOriginal = features.New("verify-webpage-original").
	Assess("get-externalip-and-check-webpage-for-original-string", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("getting external ip and checking for original string")

		// get external IP
		externalIP, err := client.GetServiceExternalIP("envoy", "tanzu-system-ingress", cfg.Client().RESTConfig(), 2, 30)
		if err != nil {
			t.Error("error while getting external IP")
			t.FailNow()
		} else {
			t.Log("external IP retrieved")
		}

		// set url
		url := fmt.Sprintf("%s/%s", externalIP, outerloopConfig.Project.WebpageRelativePath)
		if !strings.HasPrefix(url, "http://") {
			url = "http://" + url
		}

		webpageContainsString, _ := misc.VerifyWebpageContainsString(outerloopConfig.Project.Host, url, outerloopConfig.Project.OriginalString, 10, 10, 30)
		if !webpageContainsString {
			t.Error("webpage does not contains string")
			t.FailNow()
		} else {
			t.Log("webpage contains string")
		}

		return ctx
	}).
	Feature()

var gitUpdate = features.New("git-update").
	Assess("git-config", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("setting git config")

		// set git config
		err := git.GitConfig(outerloopConfig.Project.Username, outerloopConfig.Project.Email)
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
		err := git.GitClone(utils.GetFileDir(), outerloopConfig.Project.Repository)
		if err != nil {
			t.Error("error while cloning git repo")
			t.FailNow()
		} else {
			t.Log("cloned git repo")
		}

		return ctx
	}).
	Assess("git-seturl", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("setting git remote url")

		// set remote url
		err := git.GitSetRemoteUrl(filepath.Join(utils.GetFileDir(), outerloopConfig.Project.Name), outerloopConfig.Project.AccessToken, outerloopConfig.Project.Repository)
		if err != nil {
			t.Error("error while setting git remote url")
			t.FailNow()
		} else {
			t.Log("set git remote url")
		}

		return ctx
	}).
	Assess("replace-string-in-file", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("replacing string in file")

		// replace string
		err := utils.ReplaceStringInFile(filepath.Join(utils.GetFileDir(), outerloopConfig.Project.Name, outerloopConfig.Project.File), outerloopConfig.Project.OriginalString, outerloopConfig.Project.NewString)
		if err != nil {
			t.Error("error while replacing string in file")
			t.FailNow()
		} else {
			t.Log("replaced string in file")
		}

		return ctx
	}).
	Assess("git-add", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("adding files to git index")

		// add files
		err := git.GitAdd(filepath.Join(utils.GetFileDir(), outerloopConfig.Project.Name), []string{outerloopConfig.Project.File})
		if err != nil {
			t.Error("error while adding files to git index")
			t.FailNow()
		} else {
			t.Log("added files to git index")
		}

		return ctx
	}).
	Assess("git-commit", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("git committing index")

		// commit
		err := git.GitCommit(filepath.Join(utils.GetFileDir(), outerloopConfig.Project.Name), outerloopConfig.Project.CommitMessage)
		if err != nil {
			t.Error("error while committing git index")
			t.FailNow()
		} else {
			t.Log("committed git index")
		}

		return ctx
	}).
	Assess("git-push", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("pushing changes to repo")

		// push
		err := git.GitPush(filepath.Join(utils.GetFileDir(), outerloopConfig.Project.Name), false)
		if err != nil {
			t.Error("error while pushing changes to repo")
			t.Fail() // DON'T DO t.FailNow() AS WE WANT TO CLEAN UP REGARDLESS OF THE STATE OF THE TEST
		} else {
			t.Log("pushed changes to repo")
		}

		return ctx
	}).
	Feature()

var verifyWebpageNew = features.New("verify-webpage-new").
	Assess("get-externalip-and-check-webpage-for-new-string", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("getting external ip and checking for new string")

		// get external IP
		externalIP, err := client.GetServiceExternalIP("envoy", "tanzu-system-ingress", cfg.Client().RESTConfig(), 2, 30)
		if err != nil {
			t.Error("error while getting external IP")
			t.Fail() // DON'T DO t.FailNow() AS WE WANT TO CLEAN UP REGARDLESS OF THE STATE OF THE TEST
		} else {
			t.Log("external IP retrieved")
		}

		// set url
		url := fmt.Sprintf("%s/%s", externalIP, outerloopConfig.Project.WebpageRelativePath)
		if !strings.HasPrefix(url, "http://") {
			url = "http://" + url
		}

		webpageContainsString, _ := misc.VerifyWebpageContainsString(outerloopConfig.Project.Host, url, outerloopConfig.Project.NewString, 10, 10, 30)
		if !webpageContainsString {
			t.Error("webpage does not contains string")
			t.Fail() // DON'T DO t.FailNow() AS WE WANT TO CLEAN UP REGARDLESS OF THE STATE OF THE TEST
		} else {
			t.Log("webpage contains string")
		}

		return ctx
	}).
	Feature()

var gitReset = features.New("git-reset").
	Assess("git-reset", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("resetting git repo")

		// reset
		err := git.GitResetFromHead(filepath.Join(utils.GetFileDir(), outerloopConfig.Project.Name), 1)
		if err != nil {
			t.Error("error while resetting git repo")
			t.Fail() // DON'T DO t.FailNow() AS WE WANT TO CLEAN UP REGARDLESS OF THE STATE OF THE TEST
		} else {
			t.Log("resetted git repo")
		}

		return ctx
	}).
	Assess("git-push-force", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("force pushing changes to repo")

		// force push
		err := git.GitPush(filepath.Join(utils.GetFileDir(), outerloopConfig.Project.Name), true)
		if err != nil {
			t.Error("error while force pushing changes to repo")
			t.Fail() // DON'T DO t.FailNow() AS WE WANT TO CLEAN UP REGARDLESS OF THE STATE OF THE TEST
		} else {
			t.Log("force pushed changes to repo")
		}

		return ctx
	}).
	Feature()

var removeProjectDir = features.New("remove-project-dir").
	Assess("remove-dir", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("removing project directory")

		// remove
		err := utils.RemoveDirectory(outerloopConfig.Project.Name)
		if err != nil {
			t.Error("error while removing directory")
			t.Fail() // DON'T DO t.FailNow() AS WE WANT TO CLEAN UP REGARDLESS OF THE STATE OF THE TEST
		} else {
			t.Log("removed directory")
		}

		return ctx
	}).
	Feature()

var deleteWorkload = features.New("delete-workload").
	Assess("delete-workload", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("deleting workload")

		// delete workload
		err := tanzuCmds.TanzuDeleteWorkload(outerloopConfig.Workload.YamlFile, outerloopConfig.Namespace)
		if err != nil {
			t.Error("error while deleting workload")
			t.Fail() // DON'T DO t.FailNow() AS WE WANT TO CLEAN UP REGARDLESS OF THE STATE OF THE TEST
		} else {
			t.Log("deleted workload")
		}
		return ctx
	}).
	Assess(fmt.Sprintf("validate-%s-deletion", outerloopConfig.Workload.Name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		workloadDeleted := tanzu_helpers.ValidateWorkloadDeleted(outerloopConfig.Workload.Name, outerloopConfig.Namespace, 5, 30)
		if !workloadDeleted {
			t.Errorf("error while validating workload %s deletion", outerloopConfig.Workload.Name)
			t.Fail()
		} else {
			t.Logf("validated workload %s deletion", outerloopConfig.Workload.Name)
		}
		// workaround for kapp-controller issue: https://github.com/vmware-tanzu/carvel-kapp-controller/issues/416
		t.Logf("Waiting for 2 mins after workload deletion to avoid ns getting stuck at deletion")
		time.Sleep(time.Duration(120) * time.Second)
		return ctx
	}).
	Feature()

var createGithubRepo = features.New("create-github-repo").
	Assess("create-github-repo", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("creating github repo")

		// create repo
		err := github.CreateGithubRepo(outerloopConfig.Project.Name, outerloopConfig.Project.RepoTemplate, outerloopConfig.Project.AccessToken)
		if err != nil {
			t.Error("error while creating repo ")
			t.FailNow()
		} else {
			t.Log("created repo")
		}
		return ctx
	}).
	Feature()

var deleteGithubRepo = features.New("delete-github-repo").
	Assess("delete-github-repo", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("deleting github repo")

		// create repo
		err := github.DeleteGithubRepo(outerloopConfig.Project.Name, outerloopConfig.Project.AccessToken)
		if err != nil {
			t.Error("error while deleting repo ")
			t.FailNow()
		} else {
			t.Log("deleted repo")
		}
		return ctx
	}).
	Feature()

var createGitSSHSecret = features.New("create-git-ssh-secret").
	Assess("create-secret", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("Creating git-ssh secret")
		err := kubectlCmds.KubectlApplyConfiguration(outerloopConfig.Workload.GitSSHSecretYamlFile, outerloopConfig.Namespace)
		if err != nil {
			t.Error("error while creating git-ssh secret")
			t.Fail()
		} else {
			t.Log("created git-ssh secret")
		}
		return ctx
	}).
	Feature()

var patchServiceAccountSecrets = features.New("patch-sa-secret").
	Assess("patch-sa-secret", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("Patching default sa secrets")
		res := kubectl_helpers.PatchServiceAccountWithNewSecret("default", outerloopConfig.Namespace, "git-ssh")
		if !res {
			t.Error("error while patching sa secret")
			t.Fail()
		} else {
			t.Log("patched sa secret")
		}
		return ctx
	}).
	Feature()

var verifyDeliverables = features.New("verify-deliverables").
	Assess("verify-deliverables-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying deliverables ready status")

		// check
		if !kubectl_helpers.ValidateDeliverables(outerloopConfig.Workload.Name, outerloopConfig.Namespace, 5, 30) {
			t.Error("deliverables not ready")
			t.FailNow()
		} else {
			t.Log("deliverables ready")
		}

		return ctx
	}).
	Feature()

var verifyServiceBindings = features.New("verify-service-bindings").
	Assess("verify-service-bindings-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying service bindings ready status")

		sbname := fmt.Sprintf("%[1]s-%[1]s%[2]s", outerloopConfig.Workload.Name, outerloopConfig.Workload.ServiceBindingSuffix)
		if !kubectl_helpers.ValidateServiceBindings(sbname, outerloopConfig.Namespace, 5, 30) {
			t.Error("service bindings not ready")
			t.FailNow()
		} else {
			t.Log("service bindings ready")
		}

		return ctx
	}).
	Feature()

var verifyBuildStatusAfterUpdate = features.New("verify-build-status").
	Assess("verify-build-succeeded", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying build succeeded status")

		buildSucceeded := kubectl_helpers.VerifyNewerBuildStatus(buildName, outerloopConfig.Namespace, 15, 60)
		if !buildSucceeded {
			t.Error("build not succeeded")
			t.FailNow()
		} else {
			t.Log("build succeeded")
		}
		return ctx
	}).
	Feature()

var verifyKsvcStatusAfterUpdate = features.New("verify-ksvc-status").
	Assess("verify-ksvc-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying ksvc ready status")

		ksvcReady := kubectl_helpers.VerifyNewerKsvcStatus(outerloopConfig.Workload.Name, outerloopConfig.Namespace, revisionName, 5, 30)
		if !ksvcReady {
			t.Error("ksvc not ready")
			t.FailNow()
		} else {
			t.Log("ksvc ready")
		}

		return ctx
	}).
	Feature()

var verifyRevisionStatusAfterUpdate = features.New("verify-revision-status").
	Assess("verify-revision-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying revision ready status")

		revisionName = kubectl_helpers.GetNewerRevision(revisionName, outerloopConfig.Workload.Name, outerloopConfig.Namespace, 5, 30)
		revisionReady := kubectl_helpers.ValidateRevisionStatus(revisionName, outerloopConfig.Workload.Name, outerloopConfig.Namespace, 5, 30)
		if !revisionReady {
			t.Error("revision not ready")
			t.FailNow()
		} else {
			t.Log("revision ready")
		}
		return ctx
	}).
	Feature()

var createDestGithubRepo = features.New("create-github-repo").
	Assess("create-github-repo", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("creating github repo")

		// create repo
		err := github.CreateGithubRepo(outerloopConfig.Project.DestName, outerloopConfig.Project.DestRepoTemplate, outerloopConfig.Project.AccessToken)
		if err != nil {
			t.Error("error while creating repo ")
			t.FailNow()
		} else {
			t.Log("created repo")
		}
		return ctx
	}).
	Feature()

var deleteDestGithubRepo = features.New("delete-github-repo").
	Assess("delete-github-repo", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("deleting github repo")

		// create repo
		err := github.DeleteGithubRepo(outerloopConfig.Project.DestName, outerloopConfig.Project.AccessToken)
		if err != nil {
			t.Error("error while deleting repo ")
			t.FailNow()
		} else {
			t.Log("deleted repo")
		}
		return ctx
	}).
	Feature()

var deployBuildPacksPipeline = features.New("deploy-pipeline-app-via-yaml-configurations").
	Assess("deploy-buildpacks-pipeline", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("deploying buildpacks-test-pipeline")

		// deploy app
		err := kubectlCmds.KubectlApplyConfiguration(outerloopConfig.SpringPetclinicPipeline.YamlFile, outerloopConfig.Namespace)
		if err != nil {
			t.Error("error while deploying buildpacks-test-pipeline")
			t.FailNow()
		} else {
			t.Log("deployed buildpacks-test-pipeline")
		}

		return ctx
	}).
	Feature()

var deployBuildPackWorkloads = features.New("deploy-buildpack-workloads").
	Assess("deploying-buildpack-workloads-test", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("deploying workloads")

		for _, workload := range outerloopConfig.BuildPacks.Workloads {
			deployWorkload := features.New(fmt.Sprintf("deploying-workload-%s", workload.Name)).
				Assess(fmt.Sprintf("%s", workload.Name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
					var branch string
					if workload.GitBranch == "" {
						branch = "main"
					} else {
						branch = workload.GitBranch
					}
					err := tanzu_libs.TanzuDeployWorkloadByCommand(workload.Name, outerloopConfig.Namespace, workload.GitRepository, branch, "web", "true")
					if err != nil {
						t.Errorf("error while deploying %s", workload.Name)
						t.Fail()
					} else {
						t.Logf("deployed workload %s", workload.Name)
					}
					return ctx
				}).
				Feature()
			testenv.Test(t, deployWorkload)
		}
		return ctx
	}).
	Feature()

var verifyBuildPackWorkloadsGitrepoStatus = features.New("verify-buildpacks-gitrepo-status").
	Assess("verify-gitrepo-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying gitrepo ready status")

		// check
		for _, workload := range outerloopConfig.BuildPacks.Workloads {
			verifyGitrepo := features.New(fmt.Sprintf("verifying-gitrepo-%s", workload.Name)).
				Assess(fmt.Sprintf("%s", workload.Name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
					gitrepoReady := kubectl_helpers.VerifyGitRepoStatus(workload.Name, outerloopConfig.Namespace, 5, 30)
					if !gitrepoReady {
						t.Errorf("%s gitrepo not ready", workload.Name)
						t.Fail()
					} else {
						t.Logf("deployed workload %s", workload.Name)
					}
					return ctx
				}).
				Feature()
			testenv.Test(t, verifyGitrepo)
		}
		return ctx
	}).
	Feature()

var deleteBuildPackWorkloads = features.New("delete-buildpacks-workloads").
	Assess("delete-buildpack-workloads", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("deleting workloads")

		// delete workload
		for _, workload := range outerloopConfig.BuildPacks.Workloads {
			deleteWorkload := features.New(fmt.Sprintf("deleting-workload-%s", workload.Name)).
				Assess(fmt.Sprintf("%s", workload.Name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
					err := tanzu_libs.DeleteWorkload(workload.Name, outerloopConfig.Namespace)
					if err != nil {
						t.Errorf("error while deleting workload %s", workload.Name)
						t.Fail() // DON'T DO t.FailNow() AS WE WANT TO CLEAN UP REGARDLESS OF THE STATE OF THE TEST
					} else {
						t.Logf("deleted workload %s", workload.Name)
					}
					return ctx
				}).
				Assess(fmt.Sprintf("validate-%s-deletion", workload.Name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
					workloadDeleted := tanzu_helpers.ValidateWorkloadDeleted(workload.Name, outerloopConfig.Namespace, 5, 30)
					if !workloadDeleted {
						t.Errorf("error while validating workload %s deletion", workload.Name)
						t.Fail()
					} else {
						t.Logf("validated workload %s deletion", workload.Name)
					}
					return ctx
				}).
				Feature()
			testenv.Test(t, deleteWorkload)
		}
		// workaround for kapp-controller issue: https://github.com/vmware-tanzu/carvel-kapp-controller/issues/416
		t.Logf("Waiting for 2 mins after workload deletion to avoid ns getting stuck at deletion")
		time.Sleep(time.Duration(120) * time.Second)
		return ctx
	}).
	Feature()

var verifyBuildPackWorkloadsSourceScanStatus = features.New("verify-buildpacks-source-scan-status").
	Assess("verify-source-scan-completed", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying source scan status")

		// check
		for _, workload := range outerloopConfig.BuildPacks.Workloads {
			verifySourceScan := features.New(fmt.Sprintf("verify-sourcescan-%s", workload.Name)).
				Assess(fmt.Sprintf("%s", workload.Name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
					sourceScanCompleted := kubectl_helpers.ValidateSourceScans(workload.Name, outerloopConfig.Namespace, 5, 30)
					if !sourceScanCompleted {
						t.Errorf("source scan %s completed", workload.Name)
						t.Fail()
					} else {
						t.Logf("source scan %s completed successfully", workload.Name)
					}
					return ctx
				}).
				Feature()
			testenv.Test(t, verifySourceScan)
		}
		return ctx
	}).
	Feature()

var verifyBuildPackWorkloadsBuildStatus = features.New("verify-buildpacks-build-status").
	Assess("verify-build-succeeded", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying build succeeded status")
		for _, workload := range outerloopConfig.BuildPacks.Workloads {
			verifyBuildStatus := features.New(fmt.Sprintf("verify-build-status-%s", workload.Name)).
				Assess(fmt.Sprintf("%s", workload.Name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
					buildName = fmt.Sprintf("%s%s", workload.Name, outerloopConfig.Workload.BuildNameSuffix)
					buildSucceeded := kubectl_helpers.VerifyBuildStatus(buildName, outerloopConfig.Namespace, 15, 60)
					if !buildSucceeded {
						t.Errorf("build %s not succeeded", buildName)
						t.Fail()
					} else {
						t.Logf("build %s succeeded", buildName)
					}
					return ctx
				}).
				Feature()
			testenv.Test(t, verifyBuildStatus)
		}
		return ctx
	}).
	Feature()

var verifyBuildPackWorkloadsPodintents = features.New("verify-buildpacks-podintents-labels-conventions").
	Assess("verify-podintent-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying podintent ready status")

		// check
		for _, workload := range outerloopConfig.BuildPacks.Workloads {
			verifyPodIntent := features.New(fmt.Sprintf("verify-podintent-%s", workload.Name)).
				Assess(fmt.Sprintf("%s", workload.Name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
					if !kubectl_helpers.VerifyPodIntentStatus(workload.Name, outerloopConfig.Namespace, 5, 30) {
						t.Errorf("podintent %s not ready", workload.Name)
						t.Fail()
					} else {
						t.Logf("podintent %s ready", workload.Name)
					}
					return ctx
				}).
				Feature()
			testenv.Test(t, verifyPodIntent)
		}

		return ctx
	}).
	Assess("verify-podintent-alv-lables", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying appliveview labels present in podintent")

		// check
		for _, workload := range outerloopConfig.BuildPacks.Workloads {
			verifyAlvLabels := features.New(fmt.Sprintf("verify-alv-labels-%s", workload.Name)).
				Assess(fmt.Sprintf("%s", workload.Name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
					alvLabelsPresent := kubectl_helpers.ValidateAppLiveViewLabels(workload.Name, outerloopConfig.Namespace)
					if alvLabelsPresent && workload.ContainsConventions {
						t.Logf("appliveview labels present in podintent %s", workload.Name)
					} else if !workload.ContainsConventions {
						t.Logf("appliveview lables absent in podintent %s", workload.Name)
					} else {
						t.Errorf("appliveview lables absent in podintent %s", workload.Name)
						t.Fail()
					}
					return ctx
				}).
				Feature()
			testenv.Test(t, verifyAlvLabels)
		}

		return ctx
	}).
	Assess("verify-podintent-springbootconventions-lables", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying springbootconventions labels present in podintent")

		// check
		for _, workload := range outerloopConfig.BuildPacks.Workloads {
			verifyPodIntent := features.New(fmt.Sprintf("verify-podintent-sprintbootconventions-%s", workload.Name)).
				Assess(fmt.Sprintf("%s", workload.Name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
					springbootconventionsLabelsPresent := kubectl_helpers.ValidateSpringBootLabels(workload.Name, outerloopConfig.Namespace)
					if springbootconventionsLabelsPresent && workload.ContainsConventions {
						t.Logf("springbootconventions labels present in podintent %s", workload.Name)
					} else if !workload.ContainsConventions {
						t.Logf("springbootconventions labels absent in podintent %s", workload.Name)
					} else {
						t.Errorf("springbootconventions lables absent in podintent %s", workload.Name)
						t.Fail()
					}
					return ctx
				}).
				Feature()
			testenv.Test(t, verifyPodIntent)
		}

		return ctx
	}).
	Assess("verify-podintent-alv-conventions", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying appliveview conventions present in podintent")

		// check
		for _, workload := range outerloopConfig.BuildPacks.Workloads {
			verifyPodIntent := features.New(fmt.Sprintf("verify-podintent-alv-conventions-%s", workload.Name)).
				Assess(fmt.Sprintf("%s", workload.Name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
					appliveviewConventionsPresent := kubectl_helpers.ValidateAppLiveViewConventions(workload.Name, outerloopConfig.Namespace)
					if appliveviewConventionsPresent && workload.ContainsConventions {
						t.Logf("appliveview conventions present in podintent %s", workload.Name)
					} else if !workload.ContainsConventions {
						t.Logf("appliveview conventions absent in podintent %s", workload.Name)
					} else {
						t.Errorf("appliveview conventions absent in podintent %s", workload.Name)
						t.Fail()
					}
					return ctx
				}).
				Feature()
			testenv.Test(t, verifyPodIntent)
		}

		return ctx
	}).
	Assess("verify-podintent-springbootconventions-conventions", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying springbootconventions conventions present in podintent")

		// check
		for _, workload := range outerloopConfig.BuildPacks.Workloads {
			verifyPodIntent := features.New(fmt.Sprintf("verify-podintent-springbootconventions-conventions-%s", workload.Name)).
				Assess(fmt.Sprintf("%s", workload.Name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
					springbootconventionsConventionsPresent := kubectl_helpers.ValidateSpringBootConventions(workload.Name, outerloopConfig.Namespace)
					if springbootconventionsConventionsPresent && workload.ContainsConventions {
						t.Logf("springbootconventions conventions present in podintent %s", workload.Name)
					} else if !workload.ContainsConventions {
						t.Logf("springbootconventions conventions absent in podintent %s", workload.Name)
					} else {
						t.Errorf("springbootconventions conventions absent in podintent %s", workload.Name)
						t.Fail()
					}
					return ctx
				}).
				Feature()
			testenv.Test(t, verifyPodIntent)
		}

		return ctx
	}).
	Feature()

var verifyBuildPackWorkloadsImageScanStatus = features.New("verify-buildpacks-imagescan-status").
	Assess("verify-imagescan-completed", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying image scan status")

		// check
		for _, workload := range outerloopConfig.BuildPacks.Workloads {
			verifyImageScan := features.New(fmt.Sprintf("verify-imagescan-%s", workload.Name)).
				Assess(fmt.Sprintf("%s", workload.Name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
					imageScanCompleted := kubectl_helpers.ValidateImageScans(workload.Name, outerloopConfig.Namespace, 5, 30)
					if !imageScanCompleted {
						t.Errorf("image scan %s completed", workload.Name)
						t.Fail()
					} else {
						t.Logf("image scan %s completed successfully", workload.Name)
					}
					return ctx
				}).
				Feature()
			testenv.Test(t, verifyImageScan)
		}
		return ctx
	}).
	Feature()

var verifyBuildPackWorkloadsTaskrunStatus = features.New("verify-buildpacks-taskrun-status").
	Assess("verify-taskrun-succeeded", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying taskrun succeeded status")

		for _, workload := range outerloopConfig.BuildPacks.Workloads {
			verifyTaskRun := features.New(fmt.Sprintf("verify-taskrun-%s", workload.Name)).
				Assess(fmt.Sprintf("%s", workload.Name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
					taskRunPrefix := fmt.Sprintf("%s%s", workload.Name, outerloopConfig.Workload.TaskRunInfix)
					taskrunSucceeded := kubectl_helpers.VerifyTaskrunStatus(taskRunPrefix, outerloopConfig.Namespace, 5, 30)
					if !taskrunSucceeded {
						t.Errorf("taskrun %s not succeeded", workload.Name)
						t.Fail()
					} else {
						t.Logf("taskrun %s succeeded", workload.Name)
					}
					return ctx
				}).
				Feature()
			testenv.Test(t, verifyTaskRun)
		}
		return ctx
	}).
	Feature()

var verifyBuildPackWorkloadsTestTaskrunStatus = features.New("verify-buildpacks-test-taskrun-status").
	Assess("verify-test-taskrun-succeeded", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying test taskrun succeeded status")

		for _, workload := range outerloopConfig.BuildPacks.Workloads {
			verifyTaskRun := features.New(fmt.Sprintf("verify-test-taskrun-%s", workload.Name)).
				Assess(fmt.Sprintf("%s", workload.Name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
					taskrunSucceeded := kubectl_helpers.VerifyTestTaskrunStatus(workload.Name, outerloopConfig.Workload.TaskRunTestSuffix, outerloopConfig.Namespace, 5, 30)
					if !taskrunSucceeded {
						t.Errorf("taskrun %s not succeeded", workload.Name)
						t.Fail()
					} else {
						t.Logf("taskrun %s succeeded", workload.Name)
					}
					return ctx
				}).
				Feature()
			testenv.Test(t, verifyTaskRun)
		}

		return ctx
	}).
	Feature()

var verifyBuildPackWorkloadsWorkloadStatus = features.New("verify-buildpacks-workload-status").
	Assess("verify-workload-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying workload ready status")

		// check
		for _, workload := range outerloopConfig.BuildPacks.Workloads {
			verifyWorkload := features.New(fmt.Sprintf("verify-workload-ready-%s", workload.Name)).
				Assess(fmt.Sprintf("%s", workload.Name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
					workloadStatus := kubectl_helpers.GetWorkloadStatus(workload.Name, outerloopConfig.Namespace)
					if workloadStatus != "True" {
						t.Errorf("workload %s not ready", workload.Name)
						t.Fail()
					} else {
						t.Logf("workload %s ready", workload.Name)
					}
					return ctx
				}).
				Feature()
			testenv.Test(t, verifyWorkload)
		}
		return ctx
	}).
	Feature()

var verifyBuildPackWorkloadsRevisionStatus = features.New("verify-buildpacks-revision-status").
	Assess("verify-revision-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying revision ready status")

		for _, workload := range outerloopConfig.BuildPacks.Workloads {
			verifyRevision := features.New(fmt.Sprintf("verify-revision-ready-%s", workload.Name)).
				Assess(fmt.Sprintf("%s", workload.Name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
					revisionName = kubectl_helpers.GetLatestRevision(workload.Name, outerloopConfig.Namespace, 5, 30)
					revisionReady := kubectl_helpers.ValidateRevisionStatus(revisionName, workload.Name, outerloopConfig.Namespace, 10, 30)
					if !revisionReady {
						t.Errorf("revision %s not ready", revisionName)
						t.Fail()
					} else {
						t.Logf("revision %s ready", revisionName)
					}
					return ctx
				}).
				Feature()
			testenv.Test(t, verifyRevision)
		}
		return ctx
	}).
	Feature()

var verifyBuildPackWorkloadsKsvcStatus = features.New("verify-buildpacks-ksvc-status").
	Assess("verify-ksvc-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying ksvc ready status")

		for _, workload := range outerloopConfig.BuildPacks.Workloads {
			verifyKsvc := features.New(fmt.Sprintf("verify-ksvc-status-%s", workload.Name)).
				Assess(fmt.Sprintf("%s", workload.Name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
					revisionName = kubectl_helpers.GetLatestRevision(workload.Name, outerloopConfig.Namespace, 5, 30)
					ksvcReady := kubectl_helpers.VerifyKsvcStatus(workload.Name, outerloopConfig.Namespace, revisionName, 5, 30)
					if !ksvcReady {
						t.Errorf("ksvc %s not ready", revisionName)
						t.Fail()
					} else {
						t.Logf("ksvc %s ready", revisionName)
					}
					return ctx
				}).
				Feature()
			testenv.Test(t, verifyKsvc)
		}
		return ctx
	}).
	Feature()

var verifyBuildPackWorkloadsReachability = features.New("verify-buildpacks-webpage-reachability").
	Assess("get-externalip-and-check-webpage-reachability", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("getting external ip and checking reachability")

		// get external IP
		externalIP, err := client.GetServiceExternalIP("envoy", "tanzu-system-ingress", cfg.Client().RESTConfig(), 2, 30)
		if err != nil {
			t.Error("error while getting external IP")
			t.Fail()
		} else {
			t.Log("external IP retrieved")
		}

		// set url
		for _, workload := range outerloopConfig.BuildPacks.Workloads {
			verifyWorkload := features.New(fmt.Sprintf("verify-workload-reachability-%s", workload.Name)).
				Assess(fmt.Sprintf("%s", workload.Name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
					url := fmt.Sprintf("%s/%s", externalIP, workload.WebpageRelativePath)
					if !strings.HasPrefix(url, "http://") {
						url = "http://" + url
					}
					host := fmt.Sprintf("%s.%s.example.com", workload.Name, outerloopConfig.Namespace)
					t.Logf("sending GET request host: %s, url: %s", host, url)
					isWebpageReachable, _ := misc.VerifyWebpageReachable(host, url, 10, 30)
					if !isWebpageReachable {
						t.Errorf("webpage %s is not reachable", workload.Name)
						t.Fail()
					} else {
						t.Logf("webpage %s is reachable", workload.Name)
					}
					return ctx
				}).
				Feature()
			testenv.Test(t, verifyWorkload)
		}
		return ctx
	}).
	Feature()

func listVulnerabilities(workloadName string, t *testing.T) {
	verifyVulnerability := features.New(fmt.Sprintf("list-cve-%s", workloadName)).
		Assess(fmt.Sprintf("%s", workloadName), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			imageDigest := kubectl_helpers.GetImageDigest(workloadName, outerloopConfig.Namespace, 2, 30)
			log.Printf("imageDigest: %s", imageDigest)
			_, err := tanzu_libs.ListInsightImagesVulnerabilities(imageDigest)
			if err != nil {
				t.Errorf("error while getting vulnerabilities for %s", workloadName)
				t.Fail()
			}
			return ctx
		}).
		Feature()
	testenv.Test(t, verifyVulnerability)
}

var listBuildPackWorkloadsVulnerabilities = features.New("list-buildpacks-vulnerabilities").
	Assess("setup insight plugin configs", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("Setup metadata store and  insight config")
		setupInsightPluginConfig(t, cfg)
		return ctx
	}).
	Assess("list vulnerabilities", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("listing vulnerabilities")
		for _, workload := range outerloopConfig.BuildPacks.Workloads {
			listVulnerabilities(workload.Name, t)
		}
		return ctx
	}).
	Feature()

var verifyBuildPackWorkloadsDataExistInMetadata = features.New("verify-buildpacks-metadata").
	Assess("verify metadata", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verify metadata")

		for _, workload := range outerloopConfig.BuildPacks.Workloads {
			verifyMetadata := features.New(fmt.Sprintf("verify-metadata-%s", workload.Name)).
				Assess(fmt.Sprintf("%s", workload.Name), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {

					//getting image sha from kpack image for the workload
					log.Printf("getting metadata for workload image %s ", workload.Name)

					images := kubectl_libs.GetImages(workload.Name, outerloopConfig.Namespace)
					log.Printf("images: %v", images)
					imageDigest := strings.Split(images[0].LATESTIMAGE, "@")[1]
					log.Printf("imageDigests %s :", imageDigest)

					//getting insight image metadata
					status, err := tanzu_libs.GetInsightImages(imageDigest)
					if err != nil {
						t.Errorf("error while getting metadata for %s", workload.Name)
						t.Fail()
					}
					if status == "" {
						t.Errorf("metadata not available for %s", workload.Name)
						t.Fail()
					}
					return ctx
				}).
				Feature()
			testenv.Test(t, verifyMetadata)
		}
		return ctx
	}).
	Feature()

var RestartScanLinkController = features.New("restaring-scan-link-controller").
	Assess("restart workaround", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verify scan-controller restart")
		_, err := kubectl_libs.RestartScanLinkController()
		if err != nil {
			t.Errorf("error while restarting scan-link controller")
			t.Fail()
		}
		return ctx
	}).
	Feature()

func setupInsightPluginConfig(t *testing.T, cfg *envconf.Config) {
	//getting metadata store app access token
	serviceAccount := kubectl_libs.GetServiceAccountJson("metadata-store-read-write-client", "metadata-store")
	secretName := serviceAccount.Secrets[0].Name
	secret := kubectl_libs.GetSecrets(secretName, "metadata-store")
	encodedToken := string(secret.Data.Token)
	decodedToken, err := base64.StdEncoding.DecodeString(encodedToken)
	if err != nil {
		t.Error("error while decoding token")
		t.FailNow()
	}

	//getting metadata store app loadbalancer external ip
	externalIP, err := client.GetServiceExternalIP("metadata-store-app", "metadata-store", cfg.Client().RESTConfig(), 2, 30)
	if err != nil {
		t.Error("error while getting external IP")
		t.FailNow()
	} else {
		t.Log("external IP retrieved")
	}
	//Check for valid ip or not
	if net.ParseIP(externalIP) == nil {
		fmt.Printf("Invalid IP Address format: %s\n. Fetching IP address via dig command", externalIP)
		cmd := fmt.Sprintf("dig +short '%s' | tail -1", externalIP)
		externalIP, err = linux_util.ExecuteCmdInBashMode(cmd)
		externalIP = strings.TrimSpace(externalIP)
	} else {
		fmt.Printf("IP Address: %s - Valid\n", externalIP)
	}

	//appending ip mapping for metadata service to /etc/hosts
	cmd := fmt.Sprintf("echo '%s %s' >> /etc/hosts", externalIP, "metadata-store-app.metadata-store.svc.cluster.local")
	res, err := linux_util.ExecuteCmdInBashMode(cmd)
	if err != nil {
		log.Println("error")
	}
	t.Logf("%s", res)
	catCmd := "cat /etc/hosts"
	_, err2 := linux_util.ExecuteCmd(catCmd)
	if err2 != nil {
		log.Println("error")
	}

	//getting ca cert from app-tls-cert secret
	caSecret := kubectl_libs.GetSecrets("app-tls-cert", "metadata-store")
	caEncodedToken := string(caSecret.Data.CaCrt)
	caDecodedSecret, err := base64.StdEncoding.DecodeString(caEncodedToken)
	if err != nil {
		t.Error("error while decoding token")
		t.FailNow()
	}

	// create temporary file for cert
	t.Log("creating tempfile for cert")
	tempFile, err := ioutil.TempFile("", "ca*.crt")
	if err != nil {
		t.Error("error while creating tempfile for tap values schema")
		t.FailNow()
	} else {
		t.Log("created tempfile")
	}
	defer os.Remove(tempFile.Name())
	err = os.WriteFile(tempFile.Name(), caDecodedSecret, 0677)
	if err != nil {
		log.Printf("error while writing to file %s", tempFile.Name())
		log.Printf("error: %s", err)
	} else {
		log.Printf("file %s written", tempFile.Name())
	}

	//configure tanzu insight config set-target command
	err = tanzu_libs.TanzuConfigureInsight(tempFile.Name(), string(decodedToken))
	if err != nil {
		t.FailNow()
	}

}

var listSpringPetclinicVulnerabilities = features.New("list-vulnerabilities").
	Assess("setup insight plugin configs", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("Setup metadata store and  insight config")
		setupInsightPluginConfig(t, cfg)
		return ctx
	}).
	Assess("list vulnerabilities", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("listing vulnerabilities")
		listVulnerabilities(outerloopConfig.Workload.Name, t)
		return ctx
	}).
	Feature()
