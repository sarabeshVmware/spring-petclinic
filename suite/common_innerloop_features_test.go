//go:build all || innerloop || innerloop_basic || innerloop_basic_git_source

package suite

import (
	"context"
	"fmt"
	"os"
	exec2 "os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/git"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_helpers"
	kubectl_helper "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_helpers"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubernetes/client"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/misc"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzuCmds"
	tanzu_lib "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzu_libs"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

const tiltApp = "tanzu-java-web-app"
const tiltFile = tiltApp + "/Tiltfile"

var tiltprocCmdKey = "tiltprocCmd"

var deployTanzuJavaWebApp = features.New("deploy-tanzu-java-web-app-via-yaml").
	Assess("deploy-tanzu-java-web-app", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("deploying tanzu java web app")

		// deploy app
		err := tanzuCmds.TanzuDeployWorkload(suiteConfig.Innerloop.Workload.YamlFile, suiteConfig.Innerloop.Workload.Namespace)
		if err != nil {
			t.Error("error while deploying tanzu java web app")
			t.FailNow()
		} else {
			t.Log("deployed tanzu java web app")
		}

		return ctx
	}).
	Feature()

var verifyTanzuJavaWebAppStatus = features.New("verify-tanzu-java-web-app-status").
	Assess("verify-tanzu-java-web-app-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying tanzu java web app ready status")

		// check
		workloadStatus := kubectl_helpers.GetWorkloadStatus(suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace)
		if workloadStatus != "True" {
			t.Error("workload not ready")
			t.FailNow()
		} else {
			t.Log("workload ready")
		}

		return ctx
	}).
	Feature()

var gitCloneTanzuJavaWebApp = features.New("git-update").
	Assess("git-config", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("setting git config")

		// set git config
		err := git.GitConfig(suiteConfig.GitCredentials.Username, suiteConfig.GitCredentials.Email)
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
		err := git.GitClone(utils.GetFileDir(), suiteConfig.Innerloop.Workload.Gitrepository)
		if err != nil {
			t.Error("error while cloning git repo")
			t.FailNow()
		} else {
			t.Log("cloned git repo")
		}

		return ctx
	}).
	Feature()

var updateTanzuJavaWebAppTiltFile = features.New("update-allow-context-tilt").
	Assess("update-tilt-file", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		newLine := "allow_k8s_contexts(k8s_context())"
		t.Logf("Appending Line %s in tilt file %s", newLine, tiltFile)
		file, err := os.OpenFile(tiltFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			t.Error(fmt.Errorf("error while opening tilt file: %w", err))
			t.FailNow()
		}
		defer file.Close()
		_, err = file.WriteString(newLine)
		if err != nil {
			t.Error(fmt.Errorf("error while updating tilt file: %w", err))
			t.FailNow()
		}
		return ctx
	}).
	Assess("update-source-image", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Logf("updating source image in tiltfile")
		tapValuesSchema, err := getTapValuesSchema()
		if err != nil {
			t.Error("error while updating tilt file")
			t.FailNow()
		}
		source_image := fmt.Sprintf("%s/%s/%s-source", tapValuesSchema.OotbSupplyChainBasic.Registry.Server, tapValuesSchema.OotbSupplyChainBasic.Registry.Repository, suiteConfig.Innerloop.Workload.Name)
		err = utils.ReplaceStringInFile(tiltFile, "<SOURCE_IMAGE>", source_image)
		if err != nil {
			t.Error(fmt.Errorf("Error while editing tiltfile: %w", err))
			t.FailNow()
		}
		return ctx
	}).
	Assess("update-namespace", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Logf("updating namespace image in tiltfile")
		err := utils.ReplaceStringInFile(tiltFile, "<DEVELOPMENT_NAMESPACE>", suiteConfig.Innerloop.Workload.Namespace)
		if err != nil {
			t.Error(fmt.Errorf("Error while editing tiltfile: %w", err))
			t.FailNow()
		}
		return ctx
	}).
	Feature()

var updateWorkloadTiltUp = features.New("create-workload-tilt-up").
	Assess("tilting-up", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		tiltFile := filepath.Join(utils.GetFileDir(), suiteConfig.Innerloop.Workload.Name, "Tiltfile")
		t.Logf("Setting NAMESPACE environment variable to %s", suiteConfig.Innerloop.Workload.Namespace)
		os.Setenv("NAMESPACE", suiteConfig.Innerloop.Workload.Namespace)
		tiltCmd := fmt.Sprintf("tilt up --file %s --port 11223", tiltFile)
		t.Logf("Running tilt command %s", tiltCmd)
		proc, err := linux_util.RunCommandWithOutWait(tiltCmd)
		t.Logf("command executed: %s", tiltCmd)
		if err != nil {
			t.Error(fmt.Errorf("error while tilting-up : %w", err))
			t.FailNow()
		}
		t.Logf("sleeping for 1 minute")
		time.Sleep(1 * time.Minute)
		return context.WithValue(ctx, tiltprocCmdKey, proc)
	}).
	Feature()

var verifyTanzuJavaWebAppGitRepository = features.New("verify-image-repositories").
	Assess("verify-image-repositories", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Logf("verify image-repositories status")
		status := kubectl_helper.VerifyGitRepoStatus(suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace, 10, 30)
		t.Logf("ImageRepository %s status is : %t", suiteConfig.Innerloop.Workload.Name, status)
		if !status {
			t.Error(fmt.Errorf("ImageRepository %s is not ready.", suiteConfig.Innerloop.Workload.Name))
			t.Fail()
		}
		return ctx
	}).
	Feature()

var verifyTanzuJavaWebAppImageRepository = features.New("verify-image-repositories").
	Assess("verify-image-repositories", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Logf("verify image-repositories status")
		status := kubectl_helper.VerifyImageRepositoryStatus(suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace, 10, 30)
		t.Logf("ImageRepository %s status is : %t", suiteConfig.Innerloop.Workload.Name, status)
		if !status {
			t.Error(fmt.Errorf("ImageRepository %s is not ready.", suiteConfig.Innerloop.Workload.Name))
			t.Fail()
		}
		return ctx
	}).
	Feature()

var verifyTanzuJavaWebAppResponseBeforeChange = features.New("verify-app-response").
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

		webpageContainsString, _ := misc.VerifyWebpageContainsString(suiteConfig.Innerloop.Workload.URL, url, suiteConfig.Innerloop.Workload.OriginalString, 10, 10, 30)
		if !webpageContainsString {
			t.Error("webpage does not contains string")
			t.Fail() // DON'T DO t.FailNow() AS WE WANT TO CLEAN UP REGARDLESS OF THE STATE OF THE TEST
		} else {
			t.Log("webpage contains string")
		}

		return ctx
	}).
	Feature()

var makeChangesInFile = features.New("replace-string-in-file").
	Assess("replace-tanzu-to-tap ", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		oldString := suiteConfig.Innerloop.Workload.OriginalString
		newString := suiteConfig.Innerloop.Workload.NewString
		filePath := suiteConfig.Innerloop.Workload.ApplicationFilePath
		t.Logf("Replace from string %s to string %s in file %s", oldString, newString, filePath)
		err := utils.ReplaceStringInFile(filePath, oldString, newString)
		t.Logf("Compiling and building app %s", tiltApp)
		compile()
		if err != nil {
			t.Error(fmt.Errorf("error while replacing string in file %s : %w", filePath, err))
			t.FailNow()
		}
		return ctx
	}).
	Feature()

func compile() {
	app := "./mvnw"
	arg0 := "compile"
	cmd := exec2.Command(app, arg0)
	cmd.Dir = tiltApp
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(stdout))
}

var verifyTanzuJavaWebAppResponseAfterChange = features.New("verify-app-response-after-replace-string").
	Assess("get-externalip-and-check-webpage-for-new-string", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("getting external ip and checking for new string")

		// get external IP
		url, err := client.GetServiceExternalIP("envoy", "tanzu-system-ingress", cfg.Client().RESTConfig(), 2, 30)
		if err != nil {
			t.Error("error while getting external IP")
			t.Fail() // DON'T DO t.FailNow() AS WE WANT TO CLEAN UP REGARDLESS OF THE STATE OF THE TEST
		} else {
			t.Log("external IP retrieved")
		}

		if !strings.HasPrefix(url, "http://") {
			url = "http://" + url
		}

		webpageContainsString, _ := misc.VerifyWebpageContainsString(suiteConfig.Innerloop.Workload.URL, url, suiteConfig.Innerloop.Workload.NewString, 10, 10, 30)
		if !webpageContainsString {
			t.Error("webpage does not contains string")
			t.Fail() // DON'T DO t.FailNow() AS WE WANT TO CLEAN UP REGARDLESS OF THE STATE OF THE TEST
		} else {
			t.Log("webpage contains string")
		}
		return ctx
	}).
	Feature()

var verifyTanzuJavaWebAppBuildStatus = features.New("verify-builds").
	Assess("verify-build-status", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Logf("verify build status")
		buildName = fmt.Sprintf("%s%s", suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.BuildNameSuffix)
		status := kubectl_helper.VerifyBuildStatus(buildName, suiteConfig.Innerloop.Workload.Namespace, 10, 30)
		t.Logf("Build status is : %t", status)
		if !status {
			t.Error(fmt.Errorf("Build is not ready."))
			t.Fail()
		}
		return ctx
	}).
	Feature()

var verifyTanzuJavaWebAppImagesKpacStatus = features.New("verify-images.kpac").
	Assess("verify-images.kpac-status", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Logf("verify latest image status")
		status := kubectl_helper.GetLatestImageStatus(suiteConfig.Innerloop.Workload.Namespace)
		t.Logf("Image status is: %s", status)
		if status != "True" {
			t.Error(fmt.Errorf("Image is not built/ready."))
			t.Fail()
		}
		return ctx
	}).
	Feature()

var verifyTanzuJavaWebAppPodIntentStatus = features.New("verify-pod-intents-annotations-labels").
	Assess("verify-pod-intent-status", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Logf("verify podintent status")
		if !kubectl_helper.VerifyPodIntentStatus(suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace, 5, 30) {
			t.Error(fmt.Errorf("podintent %s is not ready.", suiteConfig.Innerloop.Workload.Name))
			t.Fail()
		}
		return ctx
	}).
	Assess("verify-pod-intent-app-live-view-lables", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Logf("verify if app live view lables are added to podintent")
		status := kubectl_helper.ValidateAppLiveViewLabels(suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace)
		t.Logf("app live view lables status is : %t", status)
		if !status {
			t.Error(fmt.Errorf("App live view lables are not added to the podintent"))
			t.FailNow()
		}
		return ctx
	}).
	Assess("verify-pod-intent-spring-boot-conventions-lables", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Logf("verify if spring-boot-conventions lables are added to podintent")
		status := kubectl_helper.ValidateSpringBootLabels(suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace)
		t.Logf("spring-boot-conventions lables status is : %t", status)
		if !status {
			t.Error(fmt.Errorf("Spring boot conventions lables are not added to the podintent"))
			t.FailNow()
		}
		return ctx
	}).
	Assess("verify-pod-intent-app-live-view-annotations", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Logf("verify if app live view annotations are added to podintent")
		status := kubectl_helper.ValidateAppLiveViewConventions(suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace)
		t.Logf("app live view annotations status is : %t", status)
		if !status {
			t.Error(fmt.Errorf("App live view annotations are not added to the podintent"))
			t.FailNow()
		}
		return ctx
	}).
	Assess("verify-pod-intent-devloper-conventions-annotations", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Logf("verify if devloper-conventions annotations are added to podintent")
		status := kubectl_helper.ValidateDeveloperConventions(suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace)
		t.Logf("devloper-conventions annotations status is : %t", status)
		if !status {
			t.Error(fmt.Errorf("devloper-conventions annotations are not added to the podintent"))
			t.FailNow()
		}
		return ctx
	}).
	Assess("verify-pod-intent-spring-boot-conventions-annotations", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Logf("verify if spring-boot-conventions annotations are added to  the podintent")
		status := kubectl_helper.ValidateSpringBootConventions(suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace)
		t.Logf("spring-boot-conventions annotations status is : %t", status)
		if !status {
			t.Error(fmt.Errorf("spring-boot-conventions annotations are not added to the podintent"))
			t.FailNow()
		}
		return ctx
	}).
	Feature()

var verifyTanzuJavaWebAppImageRepositoryDelivery = features.New("verify-image-repository-delivery").
	Assess("verify-image-repositories", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Logf("verify image-repositories-delivery status")
		imageRepo := suiteConfig.Innerloop.Workload.Name + suiteConfig.Innerloop.Workload.ImageDeliverySuffix
		status := kubectl_helper.VerifyImageRepositoryStatus(imageRepo, suiteConfig.Innerloop.Workload.Namespace, 10, 30)
		t.Logf("ImageRepository %s status is : %t", imageRepo, status)
		if !status {
			t.Error(fmt.Errorf("ImageRepository %s is not ready.", imageRepo))
			t.Fail()
		}
		return ctx
	}).
	Feature()

var verifyTanzuJavaWebAppDeliverable = features.New("verify-deliverables").
	Assess("verify-deliverables-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying deliverables ready status")
		if !kubectl_helper.ValidateDeliverables(suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace, 5, 30) {
			t.Error("deliverables not ready")
			t.FailNow()
		} else {
			t.Log("deliverables ready")
		}
		return ctx
	}).
	Feature()

var verifyTanzuJavaWebAppRevisionStatus = features.New("verify-revision-status").
	Assess("verify-revision-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying revision ready status")

		revisionName = kubectl_helpers.GetLatestRevision(suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace, 1, 30)
		t.Logf("latestRevision set to %s", revisionName)
		revisionReady := kubectl_helpers.ValidateRevisionStatus(revisionName, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace, 5, 30)

		if !revisionReady {
			t.Error("revision not ready")
			t.FailNow()
		} else {
			t.Log("revision ready")
		}
		return ctx
	}).
	Feature()

var verifyTanzuJavaWebAppKsvcStatus = features.New("verify-ksvc-status").
	Assess("verify-ksvc-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Logf("verifying ksvc ready status %s", ksvcLatestReady)

		ksvcReady := kubectl_helpers.VerifyKsvcStatus(suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace, revisionName, 5, 30)
		if !ksvcReady {
			t.Error("ksvc not ready")
			t.FailNow()
		} else {
			t.Log("ksvc ready")
		}

		return ctx
	}).
	Feature()

var verifyTanzuJavaWebAppWorkloadStatus = features.New("verify-workload").
	Assess("verify-workload-status", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Logf("verify workload status")
		status := kubectl_helper.ValidateWorkloadStatus(suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace, 5, 30)
		t.Logf("workload %s validation status : %v", suiteConfig.Innerloop.Workload.Name, status)
		if !status {
			t.Error(fmt.Errorf("workload %s is not ready.", suiteConfig.Innerloop.Workload.Name))
			t.Fail()
		}
		return ctx
	}).
	Feature()

var cleanup = features.New("cleanup").
	Assess("kill-tilt", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Logf("kill tilt process")
		err := (ctx.Value(tiltprocCmdKey).(*os.Process)).Kill()
		if err != nil {
			t.Error(fmt.Errorf("Fail to kill the tilt process"))
			t.FailNow()
		}
		return ctx
	}).
	Assess("delete-workload", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Logf("Deleting workload")
		tanzu_lib.DeleteWorkload(suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace)
		return ctx
	}).
	Assess("remove-dir", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		dir := filepath.Join(utils.GetFileDir(), tiltApp)

		t.Logf("removing directory %s", dir)
		err := os.RemoveAll(dir)
		if err != nil {
			t.Error(fmt.Errorf("error while removing directory %s: %w", dir, err))
			t.FailNow()
		}
		t.Logf("directory %s removed", dir)
		return ctx
	}).
	Assess("update-schema-back-to-default", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		tapValuesSchema, err := getTapValuesSchema()
		if err != nil {
			t.Error(fmt.Errorf("error while getting tap values schema: %w", err))
		}
		tapValuesSchema.Profile = "full"
		tapValuesSchema.SupplyChain = "basic"
		tapValuesSchema.Accelerator.Server.ServiceType = "ClusterIP"
		t.Logf("updating tap values schema %s", suiteConfig.Tap.ValuesSchemaFile)
		err = utils.WriteYAMLFile(suiteConfig.Tap.ValuesSchemaFile, tapValuesSchema)
		if err != nil {
			t.Error(fmt.Errorf("error while updating tap values schema %s: %w", suiteConfig.Tap.ValuesSchemaFile, err))
			t.FailNow()
		}
		t.Logf("tap values schema %s updated", suiteConfig.Tap.ValuesSchemaFile)
		return ctx
	}).
	Feature()

var verifyTanzuJavaWebAppBuildStatusAfterUpdate = features.New("verify-build-status").
	Assess("verify-build-succeeded", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying build succeeded status")

		buildSucceeded := kubectl_helpers.VerifyNewerBuildStatus(buildName, suiteConfig.Innerloop.Workload.Namespace, 15, 60)
		if !buildSucceeded {
			t.Error("build not succeeded")
			t.FailNow()
		} else {
			t.Log("build succeeded")
		}
		return ctx
	}).
	Feature()

var verifyTanzuJavaWebAppKsvcStatusAfterUpdate = features.New("verify-ksvc-status").
	Assess("verify-ksvc-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying ksvc ready status")

		ksvcReady := kubectl_helpers.VerifyNewerKsvcStatus(suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace, revisionName, 5, 30)
		if !ksvcReady {
			t.Error("ksvc not ready")
			t.FailNow()
		} else {
			t.Log("ksvc ready")
		}

		return ctx
	}).
	Feature()

var verifyTanzuJavaWebAppRevisionStatusAfterUpdate = features.New("verify-revision-status").
	Assess("verify-revision-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		t.Log("verifying revision ready status")

		revisionName = kubectl_helpers.GetNewerRevision(revisionName, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace, 5, 30)
		revisionReady := kubectl_helpers.ValidateRevisionStatus(revisionName, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace, 5, 30)
		if !revisionReady {
			t.Error("revision not ready")
			t.FailNow()
		} else {
			t.Log("revision ready")
		}
		return ctx
	}).
	Feature()

var generateAcceleratorProject = features.New("generate-acc-project-and-unzip").
	Assess("generate-project", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		service, accNamespace := "acc-server", "accelerator-system"
		t.Logf("getting external ip for %s (namespace %s)", service, accNamespace)
		serviceExternalIp, err := client.GetServiceExternalIP(service, accNamespace, cfg.Client().RESTConfig(), 2, 30)
		if err != nil {
			t.Error(fmt.Errorf("error while getting external ip for %s (namespace %s): %w", service, accNamespace, err))
			t.FailNow()
		}
		t.Logf("external ip for %s (namespace %s): %s", "server", accNamespace, serviceExternalIp)
		t.Logf("sleeping for 1 minute before generating project")
		t.Logf("generating tanzu java web app accelerator project")
		tapValuesSchema, err := getTapValuesSchema()
		if err != nil {
			t.Error(fmt.Errorf("error while getting tap values schema: %w", err))
		}
		// generate project
		repositoryPrefix := tapValuesSchema.OotbSupplyChainBasic.Registry.Server + "/" + tapValuesSchema.OotbSupplyChainBasic.Registry.Repository
		err = tanzuCmds.TanzuGenerateAccelerator("tanzu-java-web-app", "tanzu-java-web-app", repositoryPrefix, serviceExternalIp, suiteConfig.Tap.Namespace, 4, 30)
		if err != nil {
			t.Error("error while generating accelerator project")
			t.FailNow()
		} else {
			t.Log("accelerator project generated")
		}

		return ctx
	}).
	Assess("unzip-project", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		project := "tanzu-java-web-app"
		zipFile := project + ".zip"

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
		output, err = linux_util.ExecuteCmd(fmt.Sprintf("unzip %s", zipFile))
		t.Logf("command executed: unzip %s. output %s", zipFile, output)
		if err != nil {
			t.Error(fmt.Errorf("error while unzip accelerator project zip file %s: %w: %s", zipFile, err, output))
			t.FailNow()
		}
		t.Logf("Accelerator project zip files %s unzipped successfully", zipFile)

		return ctx
	}).
	Feature()

var updateTiltFile = features.New("update-allow-context-tilt").
	Assess("update-tilt-file", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		newLine := "allow_k8s_contexts(k8s_context())"
		t.Logf("Appending Line %s in tilt file %s", newLine, tiltFile)
		file, err := os.OpenFile(tiltFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			t.Error(fmt.Errorf("error while opening tilt file: %w", err))
			t.FailNow()
		}
		defer file.Close()
		_, err = file.WriteString(newLine)
		if err != nil {
			t.Error(fmt.Errorf("error while updating tilt file: %w", err))
			t.FailNow()
		}
		return ctx
	}).
	Feature()

