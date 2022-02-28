//go:build all || outerloop || outerloop_basic || outerloop_testing || outerloop_testing_scanning

package suite

import (
	"context"
	"fmt"
	"log"
	"os"
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
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	"gopkg.in/yaml.v3"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

type outerloopConfiguration struct {
	CatalogInfoYaml string `yaml:"catalog_info_yaml"`
	Mysql           struct {
		Name     string `yaml:"name"`
		YamlFile string `yaml:"yaml_file"`
	} `yaml:"mysql"`
	Namespace string `yaml:"namespace"`
	Project   struct {
		Application         string `yaml:"application"`
		Host                string `yaml:"host"`
		WebpageRelativePath string `yaml:"webpage_relative_path"`
		File                string `yaml:"file"`
		Name                string `yaml:"name"`
		RepoTemplate        string `yaml:"repo_template"`
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
	SpringPetclinicPipeline struct {
		YamlFile string `yaml:"yaml_file"`
	} `yaml:"spring_petclinic_pipeline"`
	Workload struct {
		Name                string `yaml:"name"`
		YamlFile            string `yaml:"yaml_file"`
		TestYamlFile        string `yaml:"test_yaml_file"`
		BuildNamePrefix     string `yaml:"build_name_prefix"`
		GitrepositoryName   string `yaml:"gitrepository_name"`
		ImagerepositoryName string `yaml:"imagerepository_name"`
		KsvcName            string `yaml:"ksvc_name"`
		PodintentName       string `yaml:"podintent_name"`
		TaskrunNamePrefix   string `yaml:"taskrun_name_prefix"`
		ImageScanName       string `yaml:"imagescan_name"`
		SourceScanName      string `yaml:"sourcescan_name"`
		PipelineName        string `yaml:"pipeline_name"`
		AppName				string `yaml:"app_name"`
	} `yaml:"workload"`
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
	outerloopConfig.SpringPetclinic.YamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.SpringPetclinic.YamlFile)
	outerloopConfig.SpringPetclinicPipeline.YamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.SpringPetclinicPipeline.YamlFile)
	outerloopConfig.Workload.YamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.Workload.YamlFile)
	outerloopConfig.Workload.TestYamlFile = filepath.Join(outerloopResourcesDir, outerloopConfig.Workload.TestYamlFile)

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
		sourceScanCompleted := kubectl_helpers.ValidateSourceScans(outerloopConfig.Workload.SourceScanName, outerloopConfig.Namespace, 5, 30)
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
		imageScanCompleted := kubectl_helpers.ValidateImageScans(outerloopConfig.Workload.ImageScanName, outerloopConfig.Namespace, 5, 30)
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
		pipelineRunSucceeded := kubectl_helpers.ValidatePipelineRuns("", outerloopConfig.Namespace, 5, 30)
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
		if !kubectl_helpers.ValidateLatestImageStatus(suiteConfig.Innerloop.Workload.Namespace, 10, 30) {
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
		gitrepoReady := kubectl_helpers.VerifyGitRepoStatus(outerloopConfig.Workload.GitrepositoryName, outerloopConfig.Namespace, 5, 30)
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

		// check
		buildSucceeded := kubectl_helpers.VerifyBuildStatus(outerloopConfig.Namespace, 15, 60)
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
		if !kubectl_helpers.VerifyPodIntentStatus(outerloopConfig.SpringPetclinic.PodintentName, outerloopConfig.Namespace, 5, 30) {
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
		alvLabelsPresent := kubectl_helpers.ValidateAppLiveViewLabels(outerloopConfig.Workload.PodintentName, outerloopConfig.Namespace)
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
		springbootconventionsLabelsPresent := kubectl_helpers.ValidateSpringBootLabels(outerloopConfig.Workload.PodintentName, outerloopConfig.Namespace)
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
		appliveviewConventionsPresent := kubectl_helpers.ValidateAppLiveViewConventions(outerloopConfig.Workload.PodintentName, outerloopConfig.Namespace)
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
		springbootconventionsConventionsPresent := kubectl_helpers.ValidateSpringBootConventions(outerloopConfig.Workload.PodintentName, outerloopConfig.Namespace)
		if !springbootconventionsConventionsPresent {
			t.Error("springbootconventions conventions absent in podintent")
			t.FailNow()
		} else {
			t.Log("springbootconventions conventions present in podintent")
		}

		return ctx
	}).
	Feature()

var verifyKsvcStatus = features.New("verify-ksvc-status").
	Assess("verify-ksvc-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying ksvc ready status")

		// check
		ksvcReady := kubectl_helpers.VerifyKsvcStatus(outerloopConfig.Workload.KsvcName, outerloopConfig.Namespace, 5, 30)
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

		// check
		taskrunSucceeded := kubectl_helpers.VerifyTaskrunStatus(outerloopConfig.Namespace, 5, 30)
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
		externalIP, err := client.GetServiceExternalIP("envoy", "tanzu-system-ingress", cfg.Client().RESTConfig())
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
		externalIP, err := client.GetServiceExternalIP("envoy", "tanzu-system-ingress", cfg.Client().RESTConfig())
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

		// deploy workload
		err := tanzuCmds.TanzuDeleteWorkload(outerloopConfig.Workload.YamlFile, outerloopConfig.Namespace)
		if err != nil {
			t.Error("error while deleting workload")
			t.Fail() // DON'T DO t.FailNow() AS WE WANT TO CLEAN UP REGARDLESS OF THE STATE OF THE TEST
		} else {
			t.Log("deleted workload")
		}

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

var verifyDeliverables = features.New("verify-deliverables").
	Assess("verify-deliverables-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying deliverables ready status")

		// check
		if !kubectl_helpers.ValidateDeliverables(outerloopConfig.Workload.AppName, outerloopConfig.Namespace, 5, 30) {
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

		// check
		sbname := fmt.Sprintf("%[1]s-%[1]s-db", outerloopConfig.Workload.AppName)
		if !kubectl_helpers.ValidateServiceBindings(sbname, outerloopConfig.Namespace, 5, 30) {
			t.Error("service bindings not ready")
			t.FailNow()
		} else {
			t.Log("service bindings ready")
		}

		return ctx
	}).
	Feature()
