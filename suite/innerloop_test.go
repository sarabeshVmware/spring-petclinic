package suite

import (
	"context"
	"fmt"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/client"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/exec"
	kubectl_helper "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_helpers"
	tanzu_lib "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzu_libs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	exec2 "os/exec"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	"strings"
	"testing"
	"time"
)

var tiltCmd string

const tiltApp = "tanzu-java-web-app"
const tiltFile = tiltApp + "/Tiltfile"

func TestInnerloopBasic(t *testing.T) {
	f1 := features.New("update-tap-light-supplychainbasic").
		Assess("update-schema", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			tapValuesSchema.Profile = "light"
			tapValuesSchema.SupplyChain = "basic"
	        tapValuesSchema.Accelerator.Server.ServiceType = "LoadBalancer"
			t.Logf("updating tap values schema %s", config.Tap.ValuesSchemaFile)
			err := WriteYAMLFile(config.Tap.ValuesSchemaFile, tapValuesSchema)
			if err != nil {
				t.Error(fmt.Errorf("error while updating tap values schema %s: %w", config.Tap.ValuesSchemaFile, err))
				t.FailNow()
			}
			t.Logf("tap values schema %s updated", config.Tap.ValuesSchemaFile)
			return ctx
		}).
		Assess("update-tap", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("updating package %s", config.Tap.Name)
			cmd, output, err := exec.TanzuUpdatePackage(config.Tap.Name, config.Tap.PackageName, config.Tap.Version, config.Tap.Namespace, config.Tap.ValuesSchemaFile)
			t.Logf("command executed: %s", cmd)
			if err != nil {
				t.Error(fmt.Errorf("error while updating package %s: %w: %s", config.Tap.Name, err, output))
				t.FailNow()
			}
			t.Logf("package %s updated: %s", config.Tap.Name, output)
			t.Logf("sleeping for 1 minute")
			time.Sleep(time.Minute)
			return ctx
		}).
		Feature()

	accServerExternalIpKey := "accServerExternalIp"

	f2 := features.New("get-acc-server-externalip").
		Assess("get-externalip", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			service, accNamespace := "acc-server", "accelerator-system"
			t.Logf("getting external ip for %s (namespace %s)", service, accNamespace)
			serviceExternalIp, err := client.GetServiceExternalIP(service, accNamespace, cfg.Client().RESTConfig())
			if err != nil {
				t.Error(fmt.Errorf("error while getting external ip for %s (namespace %s): %w", service, accNamespace, err))
				t.FailNow()
			}
			t.Logf("external ip for %s (namespace %s): %s", "server", accNamespace, serviceExternalIp)
			t.Logf("sleeping for 1 minute before generating project")
			//time.Sleep(time.Minute)
			return context.WithValue(ctx, accServerExternalIpKey, serviceExternalIp)
		}).
		Feature()

	acceleratorNameKey := "acceleratorName"
	f3 := features.New("generate-acc-project").
		Assess("generate-project", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			acceleratorProject := "tanzu-java-web-app"
			acceleratorName := "tanzu-java-web-app"
			repositoryPrefix := tapValuesSchema.OotbSupplyChainBasic.Registry.Server + "/" + tapValuesSchema.OotbSupplyChainBasic.Registry.Repository
			t.Logf("generating accelerator project %s (namespace %s)", acceleratorProject, config.Tap.Namespace)
			cmd, output, err := exec.TanzuGenerateAccelerator(acceleratorName, acceleratorProject, repositoryPrefix, ctx.Value(accServerExternalIpKey).(string), config.Tap.Namespace)
			t.Logf("command executed: %s", cmd)
			if err != nil {
				t.Error(fmt.Errorf("error while generating accelerator project %s in namespace %s: %w: %s", acceleratorProject, config.Tap.Namespace, err, output))
				t.FailNow()
			}
			t.Logf("Accelerator project %s generated in namespace %s: %s", acceleratorProject, config.Tap.Namespace, output)
			return context.WithValue(ctx, acceleratorNameKey, acceleratorName)
		}).
		Feature()

	f4 := features.New("unzip-acc-project-zip").
		Assess("unzip-project", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			zipFile := ctx.Value(acceleratorNameKey).(string) + ".zip"
			t.Logf("Listing accelerator project zip %s)", zipFile)
			output, err := exec.RunCommand(fmt.Sprintf("ls -lt %s", zipFile))
			t.Logf("command executed: ls -lt %s. output %s", zipFile, output)
			if err != nil {
				t.Error(fmt.Errorf("error while listing accelerator project zip file %s: %w: %s", zipFile, err, output))
				t.FailNow()
			}
			t.Logf("Listing existing project files if exists")
			output, err = exec.RunCommand(fmt.Sprintf("ls -lt %s", ctx.Value(acceleratorNameKey).(string)))
			t.Logf("command executed: ls -lt %s. output %s", ctx.Value(acceleratorNameKey).(string), output)
			if err == nil {
				t.Logf("Deleting %s folder", ctx.Value(acceleratorNameKey))
				output, err := exec.RunCommand(fmt.Sprintf("rm -rf %s", ctx.Value(acceleratorNameKey).(string)))
				t.Logf("command executed: rm -rf %s. output %s", ctx.Value(acceleratorNameKey).(string), output)
				if err != nil {
					t.Error(fmt.Errorf("error while Deleting project files %s: %w: %s", ctx.Value(acceleratorNameKey).(string), err, output))
					t.FailNow()
				}
			}
			t.Logf("Unzip %s", zipFile)
			output, err = exec.RunCommand(fmt.Sprintf("unzip %s", zipFile))
			t.Logf("command executed: unzip %s. output %s", zipFile, output)
			if err != nil {
				t.Error(fmt.Errorf("error while unzip accelerator project zip file %s: %w: %s", zipFile, err, output))
				t.FailNow()
			}
			t.Logf("Accelerator project zip files %s unzipped successfully", zipFile)
			return ctx
		}).
		Feature()

	f5 := features.New("update-allow-context-tilt").
		Assess("update-tilt-file", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			//tiltFile := ctx.Value(acceleratorNameKey).(string) + "/Tiltfile"
			newLine := "allow_k8s_contexts(k8s_context())"
			t.Logf("Appending Line %s in tilt file at %s", newLine, tiltFile)
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
	tiltprocCmdKey := "tiltprocCmd"
	f6 := features.New("create-workload-tilt-up").
		Assess("tilting-up", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			os.Setenv("NAMESPACE", config.Innerloop.Workload.Namespace)
			//tiltFile := ctx.Value(acceleratorNameKey).(string) + "/Tiltfile"
			tiltCmd = fmt.Sprintf("tilt up --file %s --port 11223", tiltFile)
			t.Logf("Running tilt command %s", tiltCmd)
			proc, err := exec.RunCommandWithOutWait(tiltCmd)
			t.Logf("err: %s", err)
			t.Logf("command executed: %s", tiltCmd)
			if err != nil {
				t.Error(fmt.Errorf("error while tilting-up : %w", err))
				t.FailNow()
			}
			t.Logf("sleeping for 2 minute")
			time.Sleep(2 * time.Minute)

			return context.WithValue(ctx, tiltprocCmdKey, proc)
		}).
		Feature()
		
	f7 := features.New("verify-image-repositories").
		Assess("verify-image-repositories", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verify image-repositories status")
			imagerepository := [2]string{config.Innerloop.Workload.Name + "-delivery", config.Innerloop.Workload.Name}
			for _, imageRepo := range imagerepository {
				status := kubectl_helper.GetLatestImageRepositoryStatus(imageRepo, config.Innerloop.Workload.Namespace)
				t.Logf("ImageRepository %s status is : %s", imageRepo, status)
				if status != "True" {
					t.Error(fmt.Errorf("ImageRepository is not ready."))
					t.Fail()
				}
			}

			return ctx
		}).
		Feature()

	f8 := features.New("verify-builds").
		Assess("verify-build-status", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verify build status")
			status := kubectl_helper.VerifyBuildStatus(config.Innerloop.Workload.Namespace)
			t.Logf("Build status is : %t", status)
			if !status {
				t.Error(fmt.Errorf("Build is not ready."))
				t.Fail()
			}
			return ctx
		}).
		Feature()

	f9 := features.New("verify-images.kpac").
		Assess("verify-images.kpac-status", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verify latest image status")
			status := kubectl_helper.GetLatestImageStatus(config.Innerloop.Workload.Namespace)
			t.Logf("Image status is: %s", status)
			if status != "True" {
				t.Error(fmt.Errorf("Image is not built/ready."))
				t.Fail()
			}
			return ctx
		}).
		Feature()

	f10 := features.New("verify-pod-intents-annotations-labels").
		Assess("verify-pod-intent-status", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verify podintent status")
			status := kubectl_helper.GetPodIntentStatus(config.Innerloop.Workload.Name, config.Innerloop.Workload.Namespace)
			t.Logf("podintent status is : %s", status)
			if status != "True" {
				t.Error(fmt.Errorf("podintent is not ready."))
				t.Fail()
			}
			return ctx
		}).
		Assess("verify-pod-intent-app-live-view-lables", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verify if app live view lables are added to podintent")
			status := kubectl_helper.ValidateAppLiveViewLabels(config.Innerloop.Workload.Name, config.Innerloop.Workload.Namespace)
			t.Logf("app live view lables status is : %t", status)
			if !status {
				t.Error(fmt.Errorf("App live view lables are not added to podintent"))
				t.FailNow()
			}
			return ctx
		}).
		Assess("verify-pod-intent-spring-boot-conventions-lables", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verify if spring-boot-conventions lables are added to podintent")
			status := kubectl_helper.ValidateSpringBootLabels(config.Innerloop.Workload.Name, config.Innerloop.Workload.Namespace)
			t.Logf("spring-boot-conventions lables status is : %t", status)
			if !status {
				t.Error(fmt.Errorf("Spring boot conventions lables are not added to podintent"))
				t.FailNow()
			}
			return ctx
		}).
		Assess("verify-pod-intent-app-live-view-annotations", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verify if app live view annotations are added to podintent")
			status := kubectl_helper.ValidateAppLiveViewConventions(config.Innerloop.Workload.Name, config.Innerloop.Workload.Namespace)
			t.Logf("app live view annotations status is : %t", status)
			if !status {
				t.Error(fmt.Errorf("App live view annotations are not added to podintent"))
				t.FailNow()
			}
			return ctx
		}).
		Assess("verify-pod-intent-devloper-conventions-annotations", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verify if devloper-conventions annotations are added to podintent")
			status := kubectl_helper.ValidateDeveloperConventions(config.Innerloop.Workload.Name, config.Innerloop.Workload.Namespace)
			t.Logf("devloper-conventions annotations status is : %t", status)
			if !status {
				t.Error(fmt.Errorf("devloper-conventions annotations are not added to podintent"))
				t.FailNow()
			}
			return ctx
		}).
		Assess("verify-pod-intent-spring-boot-conventions-annotations", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verify if spring-boot-conventions annotations are added to podintent")
			status := kubectl_helper.ValidateSpringBootConventions(config.Innerloop.Workload.Name, config.Innerloop.Workload.Namespace)
			t.Logf("spring-boot-conventions annotations status is : %t", status)
			if !status {
				t.Error(fmt.Errorf("spring-boot-conventions annotations are not added to podintent"))
				t.FailNow()
			}
			return ctx
		}).
		Feature()

	f11 := features.New("verify-ksvc").
		Assess("verify-ksvc-status", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verify ksvc status")
			status := kubectl_helper.GetKsvcStatus(config.Innerloop.Workload.Name, config.Innerloop.Workload.Namespace)
			t.Logf("ksvc status is : %s", status)
			if status != "True" {
				t.Error(fmt.Errorf("ksvc is not ready."))
				t.Fail()
			}
			return ctx
		}).
		Feature()

	f12 := features.New("verify-workload").
		Assess("verify-workload-status", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verify workload status")
			status := kubectl_helper.GetWorkloadStatus(config.Innerloop.Workload.Name, config.Innerloop.Workload.Namespace)
			t.Logf("workload status is : %s", status)
			if status != "True" {
				t.Error(fmt.Errorf("workload is not ready."))
				t.Fail()
			}
			return ctx
		}).
		Feature()
	envoyServerExternalIpKey := "envoyServerExternalIp"

	f13 := features.New("get-envoy-server-externalip").
		Assess("get-externalip", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			service, envoyNamespace := "envoy", "tanzu-system-ingress"
			t.Logf("getting external ip for %s (namespace %s)", service, envoyNamespace)
			serviceExternalIp, err := client.GetServiceExternalIP(service, envoyNamespace, cfg.Client().RESTConfig())
			if err != nil {
				t.Error(fmt.Errorf("error while getting external ip for %s (namespace %s): %w", service, envoyNamespace, err))
				t.FailNow()
			}
			t.Logf("external ip for %s (namespace %s): %s", "server", envoyNamespace, serviceExternalIp)
			return context.WithValue(ctx, envoyServerExternalIpKey, serviceExternalIp)
		}).
		Feature()

	f14 := features.New("verify-app-response").
		Assess("verify-app-response", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verify app response")
			result := GetAppResponse(ctx.Value(envoyServerExternalIpKey).(string), "tanzu-java-web-app.tap-install.example.com")
			t.Logf("App response is : %s", result)
			if result != "Greetings from Spring Boot + Tanzu!" {
				t.Error(fmt.Errorf("App response not valid"))
				t.FailNow()
			}
			return ctx
		}).
		Feature()

	f15 := features.New("replace-string-in-file").
		Assess("replace-tanzu-to-tap ", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			oldString := "Greetings from Spring Boot + Tanzu!"
			newString := "Greetings from Spring Boot + TAP!"
			filePath := "tanzu-java-web-app/src/main/java/com/example/springboot/HelloController.java"
			t.Logf("Replace from string %s to string %s in file %s", oldString, newString, filePath)
			err := exec.ReplaceStringInFile(filePath, oldString, newString)
			compile()
			if err != nil {
				t.Error(fmt.Errorf("error while replacing string in file %s : %w", filePath, err))
				t.FailNow()
			}
			return ctx
		}).
		Feature()

	f16 := features.New("verify-app-response-after-replace-string").
		Assess("verify-app-response", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("verify app response")
			result := GetAppResponse(ctx.Value(envoyServerExternalIpKey).(string), "tanzu-java-web-app.tap-install.example.com")
			t.Logf("App response is : %s", result)
			if result != "Greetings from Spring Boot + TAP!" {
				t.Error(fmt.Errorf("App response not valid"))
				t.FailNow()
			}
			return ctx
		}).
		Feature()

	
	cleanup := features.New("cleanup").
		Assess("kill-tilt", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("kill tilt process")
			// if ctx.Value(tiltCmdKey).(string) != "" {
			// 	err := exec.KillCommandProcess(ctx.Value(tiltCmdKey).(string))
			// 	if err != nil {
			// 		t.Error(fmt.Errorf("Fail to kill the tilt process"))
			// 	}
			// }
			err := (ctx.Value(tiltprocCmdKey).(*os.Process)).Kill()
			if err != nil {
				t.Error(fmt.Errorf("Fail to kill the tilt process"))
				t.FailNow()
			}
			return ctx
		}).
		Assess("delete-workload", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("Deleting workload")
			//tanzu_lib.DeleteWorkload(config.Innerloop.Workload.Name, config.Innerloop.Workload.Namespace)
			tanzu_lib.DeleteWorkload("tanzu-java-web-app", "tap-install")
			return ctx
		}).
		Feature()
	//testenv.Test(t, test1, cleanup)
	testenv.Test(t, f2, f3, f4, f5, f6, f7, f8, f9, f10, f11, f12, f13, f14, f15, f16, cleanup)
}

func GetAppResponse(envoyExternalIP string, url string) string {
	// resp, err := http.Get(url)
	// if err != nil {
	// 	fmt.Println(fmt.Errorf("Failed to get app response: %d", err))
	// 	return ""
	// }
	// defer resp.Body.Close()
	// body, _ := io.ReadAll(resp.Body)
	// return string(body)
	time.Sleep(time.Minute)
	if !strings.HasPrefix(envoyExternalIP, "http://") {
		envoyExternalIP = "http://" + envoyExternalIP
	}
	req, err := http.NewRequest("GET", envoyExternalIP, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Host = url

	var retries int = 10
	for retries > 0 {
		resp, err := http.DefaultClient.Do(req)
		log.Println(resp.StatusCode)
		if err == nil {
			log.Println("Status code is :", resp.StatusCode)
			break
		} else {
			log.Println("err:%w", err)
			retries -= 1
			log.Printf("Retry after 30 seconds")
			time.Sleep(30 * time.Second)
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Bad HTTP Response: %s", resp.Status)
	}
	defer resp.Body.Close()
	resultStringBytes, _ := ioutil.ReadAll(resp.Body)
	resultString := string(resultStringBytes)
	return resultString
}

func compile() {
	app := "./mvnw"
	arg0 := "compile"
	cmd := exec2.Command(app, arg0)
	dir, _ := os.Getwd()
	fmt.Println(dir, "wd")
	cmd.Dir = "tanzu-java-web-app"
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(stdout))
}
