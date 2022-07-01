package common_features

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_helpers"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_libs"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubernetes/client"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/misc"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzu_helpers"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzu_libs"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/models"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TanzuDeployWorkloads(t *testing.T, outerloopConfig models.OuterloopConfig) features.Feature {
	return features.New("deploy-buildpack-workloads").
		Assess("deploying-buildpack-workloads-test", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("deploying workloads")

			for _, workload := range outerloopConfig.BuildPacks.Workloads {
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
			}
			return ctx
		}).
		Feature()
}

func VerifyBuildPackWorkloadsGitrepoStatus(t *testing.T, outerloopConfig models.OuterloopConfig) features.Feature {
	return features.New("verify-buildpacks-gitrepo-status").
		Assess("verify-gitrepo-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying gitrepo ready status")

			// check
			for _, workload := range outerloopConfig.BuildPacks.Workloads {
				gitrepoReady := kubectl_helpers.VerifyGitRepoStatus(workload.Name, outerloopConfig.Namespace, 5, 30)
				if !gitrepoReady {
					t.Errorf("%s gitrepo not ready", workload.Name)
					t.Fail()
				} else {
					t.Logf("deployed workload %s", workload.Name)
				}

			}
			return ctx
		}).
		Feature()
}

func DeleteBuildPackWorkloads(t *testing.T, outerloopConfig models.OuterloopConfig) features.Feature {
	return features.New("delete-buildpacks-workloads").
		Assess("delete-buildpack-workloads", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("deleting workloads")

			// delete workload
			for _, workload := range outerloopConfig.BuildPacks.Workloads {

				err := tanzu_libs.DeleteWorkload(workload.Name, outerloopConfig.Namespace)
				if err != nil {
					t.Errorf("error while deleting workload %s", workload.Name)
					t.Fail() // DON'T DO t.FailNow() AS WE WANT TO CLEAN UP REGARDLESS OF THE STATE OF THE TEST
				} else {
					t.Logf("deleted workload %s", workload.Name)
				}
				workloadDeleted := tanzu_helpers.ValidateWorkloadDeleted(workload.Name, outerloopConfig.Namespace, 5, 30)
				if !workloadDeleted {
					t.Errorf("error while validating workload %s deletion", workload.Name)
					t.Fail()
				} else {
					t.Logf("validated workload %s deletion", workload.Name)
				}

			}
			// workaround for kapp-controller issue: https://github.com/vmware-tanzu/carvel-kapp-controller/issues/416
			t.Logf("Waiting for 2 mins after workload deletion to avoid ns getting stuck at deletion")
			time.Sleep(time.Duration(120) * time.Second)
			return ctx
		}).
		Feature()
}

func VerifyBuildPackWorkloadsSourceScanStatus(t *testing.T, outerloopConfig models.OuterloopConfig) features.Feature {
	return features.New("verify-buildpacks-source-scan-status").
		Assess("verify-source-scan-completed", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying source scan status")

			// check
			for _, workload := range outerloopConfig.BuildPacks.Workloads {
				sourceScanCompleted := kubectl_helpers.ValidateSourceScans(workload.Name, outerloopConfig.Namespace, 5, 30)
				if !sourceScanCompleted {
					t.Errorf("source scan %s completed", workload.Name)
					t.Fail()
				} else {
					t.Logf("source scan %s completed successfully", workload.Name)
				}

			}
			return ctx
		}).
		Feature()
}

func VerifyBuildPackWorkloadsBuildStatus(t *testing.T, outerloopConfig models.OuterloopConfig) features.Feature {
	return features.New("verify-buildpacks-build-status").
		Assess("verify-build-succeeded", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying build succeeded status")
			for _, workload := range outerloopConfig.BuildPacks.Workloads {
				buildName = fmt.Sprintf("%s%s", workload.Name, outerloopConfig.Workload.BuildNameSuffix)
				buildSucceeded := kubectl_helpers.VerifyBuildStatus(buildName, outerloopConfig.Namespace, 15, 60)
				if !buildSucceeded {
					t.Errorf("build %s not succeeded", buildName)
					t.Fail()
				} else {
					t.Logf("build %s succeeded", buildName)
				}

			}
			return ctx
		}).
		Feature()
}

func VerifyBuildPackWorkloadsPodintents(t *testing.T, outerloopConfig models.OuterloopConfig) features.Feature {
	return features.New("verify-buildpacks-podintents-labels-conventions").
		Assess("verify-podintent-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying podintent ready status")

			// check
			for _, workload := range outerloopConfig.BuildPacks.Workloads {
				if !kubectl_helpers.VerifyPodIntentStatus(workload.Name, outerloopConfig.Namespace, 5, 30) {
					t.Errorf("podintent %s not ready", workload.Name)
					t.Fail()
				} else {
					t.Logf("podintent %s ready", workload.Name)
				}
			}

			return ctx
		}).
		Assess("verify-podintent-alv-lables", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying appliveview labels present in podintent")

			// check
			for _, workload := range outerloopConfig.BuildPacks.Workloads {
				alvLabelsPresent := kubectl_helpers.ValidateAppLiveViewLabels(workload.Name, outerloopConfig.Namespace)
				if alvLabelsPresent && workload.ContainsConventions {
					t.Logf("appliveview labels present in podintent %s", workload.Name)
				} else if !workload.ContainsConventions {
					t.Logf("appliveview lables absent in podintent %s", workload.Name)
				} else {
					t.Errorf("appliveview lables absent in podintent %s", workload.Name)
					t.Fail()
				}

			}

			return ctx
		}).
		Assess("verify-podintent-springbootconventions-lables", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying springbootconventions labels present in podintent")

			// check
			for _, workload := range outerloopConfig.BuildPacks.Workloads {
				springbootconventionsLabelsPresent := kubectl_helpers.ValidateSpringBootLabels(workload.Name, outerloopConfig.Namespace)
				if springbootconventionsLabelsPresent && workload.ContainsConventions {
					t.Logf("springbootconventions labels present in podintent %s", workload.Name)
				} else if !workload.ContainsConventions {
					t.Logf("springbootconventions labels absent in podintent %s", workload.Name)
				} else {
					t.Errorf("springbootconventions lables absent in podintent %s", workload.Name)
					t.Fail()
				}

			}

			return ctx
		}).
		Assess("verify-podintent-alv-conventions", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying appliveview conventions present in podintent")

			// check
			for _, workload := range outerloopConfig.BuildPacks.Workloads {
				appliveviewConventionsPresent := kubectl_helpers.ValidateAppLiveViewConventions(workload.Name, outerloopConfig.Namespace)
				if appliveviewConventionsPresent && workload.ContainsConventions {
					t.Logf("appliveview conventions present in podintent %s", workload.Name)
				} else if !workload.ContainsConventions {
					t.Logf("appliveview conventions absent in podintent %s", workload.Name)
				} else {
					t.Errorf("appliveview conventions absent in podintent %s", workload.Name)
					t.Fail()
				}

			}

			return ctx
		}).
		Assess("verify-podintent-springbootconventions-conventions", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying springbootconventions conventions present in podintent")

			// check
			for _, workload := range outerloopConfig.BuildPacks.Workloads {
				springbootconventionsConventionsPresent := kubectl_helpers.ValidateSpringBootConventions(workload.Name, outerloopConfig.Namespace)
				if springbootconventionsConventionsPresent && workload.ContainsConventions {
					t.Logf("springbootconventions conventions present in podintent %s", workload.Name)
				} else if !workload.ContainsConventions {
					t.Logf("springbootconventions conventions absent in podintent %s", workload.Name)
				} else {
					t.Errorf("springbootconventions conventions absent in podintent %s", workload.Name)
					t.Fail()
				}

			}

			return ctx
		}).
		Feature()
}

func VerifyBuildPackWorkloadsImageScanStatus(t *testing.T, outerloopConfig models.OuterloopConfig) features.Feature {
	return features.New("verify-buildpacks-imagescan-status").
		Assess("verify-imagescan-completed", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying image scan status")

			// checking image scan status, but not failing test as the main feature is to test buildpacks
			for _, workload := range outerloopConfig.BuildPacks.Workloads {
				imageScanCompleted := kubectl_helpers.ValidateImageScans(workload.Name, outerloopConfig.Namespace, 5, 30)
				if !imageScanCompleted {
					t.Logf("image scan %s failed", workload.Name)
				} else {
					t.Logf("image scan %s completed successfully", workload.Name)
				}
			}
			return ctx
		}).
		Feature()
}

func VerifyBuildPackWorkloadsTaskrunStatus(t *testing.T, outerloopConfig models.OuterloopConfig) features.Feature {

	return features.New("verify-buildpacks-taskrun-status").
		Assess("verify-taskrun-succeeded", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying taskrun succeeded status")

			for _, workload := range outerloopConfig.BuildPacks.Workloads {
				taskRunPrefix := fmt.Sprintf("%s%s", workload.Name, outerloopConfig.Workload.TaskRunInfix)
				taskrunSucceeded := kubectl_helpers.VerifyTaskrunStatus(taskRunPrefix, outerloopConfig.Namespace, 5, 30)
				if !taskrunSucceeded {
					t.Errorf("taskrun %s not succeeded", workload.Name)
					t.Fail()
				} else {
					t.Logf("taskrun %s succeeded", workload.Name)
				}

			}
			return ctx
		}).
		Feature()
}

func VerifyBuildPackWorkloadsTestTaskrunStatus(t *testing.T, outerloopConfig models.OuterloopConfig) features.Feature {
	return features.New("verify-buildpacks-test-taskrun-status").
		Assess("verify-test-taskrun-succeeded", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying test taskrun succeeded status")

			for _, workload := range outerloopConfig.BuildPacks.Workloads {
				taskrunSucceeded := kubectl_helpers.VerifyTestTaskrunStatus(workload.Name, outerloopConfig.Workload.TaskRunTestSuffix, outerloopConfig.Namespace, 5, 30)
				if !taskrunSucceeded {
					t.Errorf("taskrun %s not succeeded", workload.Name)
					t.Fail()
				} else {
					t.Logf("taskrun %s succeeded", workload.Name)
				}

			}

			return ctx
		}).
		Feature()
}

func VerifyBuildPackWorkloadsWorkloadStatus(t *testing.T, outerloopConfig models.OuterloopConfig) features.Feature {
	return features.New("verify-buildpacks-workload-status").
		Assess("verify-workload-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying workload ready status")

			// check
			for _, workload := range outerloopConfig.BuildPacks.Workloads {
				workloadStatus := kubectl_helpers.GetWorkloadStatus(workload.Name, outerloopConfig.Namespace)
				if workloadStatus != "True" {
					t.Errorf("workload %s not ready", workload.Name)
					t.Fail()
				} else {
					t.Logf("workload %s ready", workload.Name)
				}

			}
			return ctx
		}).
		Feature()
}

func VerifyBuildPackWorkloadsRevisionStatus(t *testing.T, outerloopConfig models.OuterloopConfig) features.Feature {
	return features.New("verify-buildpacks-revision-status").
		Assess("verify-revision-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying revision ready status")

			for _, workload := range outerloopConfig.BuildPacks.Workloads {
				revisionName = kubectl_helpers.GetLatestRevision(workload.Name, outerloopConfig.Namespace, 5, 30)
				revisionReady := kubectl_helpers.ValidateRevisionStatus(revisionName, workload.Name, outerloopConfig.Namespace, 10, 30)
				if !revisionReady {
					t.Errorf("revision %s not ready", revisionName)
					t.Fail()
				} else {
					t.Logf("revision %s ready", revisionName)
				}

			}
			return ctx
		}).
		Feature()
}

func VerifyBuildPackWorkloadsKsvcStatus(t *testing.T, outerloopConfig models.OuterloopConfig) features.Feature {
	return features.New("verify-buildpacks-ksvc-status").
		Assess("verify-ksvc-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verifying ksvc ready status")

			for _, workload := range outerloopConfig.BuildPacks.Workloads {
				revisionName = kubectl_helpers.GetLatestRevision(workload.Name, outerloopConfig.Namespace, 5, 30)
				ksvcReady := kubectl_helpers.VerifyKsvcStatus(workload.Name, outerloopConfig.Namespace, revisionName, 5, 30)
				if !ksvcReady {
					t.Errorf("ksvc %s not ready", revisionName)
					t.Fail()
				} else {
					t.Logf("ksvc %s ready", revisionName)
				}

			}
			return ctx
		}).
		Feature()
}

func VerifyBuildPackWorkloadsReachability(t *testing.T, outerloopConfig models.OuterloopConfig) features.Feature {
	return features.New("verify-buildpacks-webpage-reachability").
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
				if workload.Name == "java-maven-app" || workload.Name == "java-native-image-app" {
					t.Logf("Skipping the reachability check for java maven and java native image app as it's not working, but it has no problems with tbs")
					return ctx
				}
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

			}
			return ctx
		}).
		Feature()
}

func ProcessDeliverableForBuildPackWorkloads(t *testing.T, outerloopConfig models.OuterloopConfig, buildContext string, runContext string, targetRepo string) features.Feature {
	return features.New("getting deliverable and changing file").
		Assess("getting deliverable files", func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {

			for _, workload := range outerloopConfig.BuildPacks.Workloads {

				// changing to build cluster
				_, err := kubectl_libs.UseContext(buildContext)
				if err != nil {
					t.Errorf("error while changing context to %s", buildContext)
					t.FailNow()
				} else {
					t.Logf("context changed to %s", buildContext)
				}

				valid := kubectl_helpers.ValidateBuildClusterDeliverableStatus(workload.Name, outerloopConfig.Namespace, 5, 5)
				if !valid {
					t.Errorf("error while getting deliverable %s", workload.Name)
					t.FailNow()
				} else {
					t.Logf("validated deliverable %s success", workload.Name)
					deliverable := kubectl_libs.GetDeliverablesYaml(workload.Name, outerloopConfig.Namespace)
					if targetRepo != "" {
						sourceImage := kubectl_libs.GetDeliverables(workload.Name, outerloopConfig.Namespace)[0].SOURCE
						imageTag := strings.Split(sourceImage, ":")[1]
						newSourceImage := fmt.Sprintf("%s:%s", targetRepo, imageTag)
						deliverable.Spec.Source.Image = newSourceImage
					}
					deliverable.Status = kubectl_libs.Status{}
					deliverable.Metadata.OwnerReferences = kubectl_libs.OwnerReferences{}
					// create temporary deliverable file
					t.Log("creating tempfile for deliverable manifest")
					tempFile, err := ioutil.TempFile("", "deliverable*.yaml")
					if err != nil {
						t.Error("error while creating tempfile for deliverable manifest")
						t.FailNow()
					} else {
						t.Log("created tempfile")
					}
					defer os.Remove(tempFile.Name())

					// write the updated manifest to the temporary file
					err = utils.WriteYAMLFile(tempFile.Name(), deliverable)
					if err != nil {
						t.Error("error while writing updated deliverable manifest to YAML file")
						t.FailNow()
					} else {
						t.Log("wrote deliverable manifest to file")
					}

					// changing to run cluster
					_, err = kubectl_libs.UseContext(runContext)
					if err != nil {
						t.Errorf("error while changing context to %s", runContext)
						t.FailNow()
					} else {
						t.Logf("context changed to %s", runContext)
					}

					t.Log("generated deliverable.yaml to be applied :")
					linux_util.ExecuteCmd(fmt.Sprintf("cat %s", tempFile.Name()))

					//deploying deliverable
					err = kubectl_libs.KubectlApplyConfiguration(tempFile.Name(), outerloopConfig.Namespace)
					if err != nil {
						t.Error("error deploying deliverable")
						t.FailNow()
					} else {
						t.Log("deployed deliverable")
					}

					valid = kubectl_helpers.ValidateDeliverables(workload.Name, outerloopConfig.Namespace, 10, 5)
					if !valid {
						t.Error("error deploying deliverable")
						t.FailNow()
					} else {
						t.Log("deployed deliverable")
					}
				}
			}
			return ctx
		}).Feature()
}

func MulticlusterOuterloopCleanupforBuildPackWorkloads(t *testing.T, outerloopConfig models.OuterloopConfig, buildContext string, runContext string) features.Feature {
	return features.New("outerloop cleanup").
		Assess("delete-deliverables", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			// changing to run cluster
			_, err := kubectl_libs.UseContext(runContext)
			if err != nil {
				t.Errorf("error while changing context to %s", runContext)
				t.FailNow()
			} else {
				t.Logf("context changed to %s", runContext)
			}

			t.Logf("Deleting deliverable")
			for _, workload := range outerloopConfig.BuildPacks.Workloads {
				kubectl_libs.DeleteDeliverable(workload.Name, outerloopConfig.Namespace)
			}
			return ctx
		}).
		Assess("delete-workloads", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			// changing to build cluster
			_, err := kubectl_libs.UseContext(buildContext)
			if err != nil {
				t.Errorf("error while changing context to %s", buildContext)
				t.FailNow()
			} else {
				t.Logf("context changed to %s", buildContext)
			}

			t.Logf("Deleting workload")
			for _, workload := range outerloopConfig.BuildPacks.Workloads {
				tanzu_libs.DeleteWorkload(workload.Name, outerloopConfig.Namespace)
			}
			return ctx
		}).
		Feature()
}

func ListBuildPackWorkloadsVulnerabilities(t *testing.T, outerloopConfig models.OuterloopConfig, skipVulnerabilityCheck bool, metadataStoreDomain string, viewContext string, buildContext string) features.Feature {
	return features.New("list-buildpacks-vulnerabilities").
		Assess("setup insight plugin configs", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("Setup metadata store and  insight config")
			setupInsightPluginConfigForMulticluster(t, cfg, metadataStoreDomain, viewContext)
			return ctx
		}).
		Assess("list vulnerabilities", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("listing vulnerabilities")
			_, err := kubectl_libs.UseContext(buildContext)
			if err != nil {
				t.Errorf("error while changing context to %s", buildContext)
				t.FailNow()
			} else {
				t.Logf("context changed to %s", buildContext)
			}

			for _, workload := range outerloopConfig.BuildPacks.Workloads {
				imageDigest := kubectl_helpers.GetImageDigest(workload.Name, outerloopConfig.Namespace, 2, 30)
				log.Printf("imageDigest: %s", imageDigest)
				vulnerabilitiesData, err := tanzu_libs.ListInsightImagesVulnerabilities(imageDigest)
				if err != nil {
					t.Errorf("error while getting vulnerabilities for %s", workload.Name)
					t.Fail()
				}
				if !skipVulnerabilityCheck {
					// Exception only for CVE-2016-1000027, since it can not be fixed at this point: https://github.com/spring-projects/spring-framework/issues/24434
					if !strings.Contains(vulnerabilitiesData, "CVE-2016-1000027 (Critical)") {
						t.Error("CVE(s) detected in image scans")
						t.Fail()
					}
				}
			}
			return ctx
		}).
		Feature()
}

func VerifyBuildPackWorkloadsDataExistInMetadata(t *testing.T, outerloopConfig models.OuterloopConfig) features.Feature {
	return features.New("verify-buildpacks-metadata").
		Assess("verify metadata", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("verify metadata")

			for _, workload := range outerloopConfig.BuildPacks.Workloads {

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

			}
			return ctx
		}).
		Feature()
}

func setupInsightPluginConfigForMulticluster(t *testing.T, cfg *envconf.Config, metadataStoreDomain string, viewContext string) {

	// changing to view cluster
	_, err := kubectl_libs.UseContext(viewContext)
	if err != nil {
		t.Errorf("error while changing context to %s", viewContext)
		t.FailNow()
	} else {
		t.Logf("context changed to %s", viewContext)
	}

	//getting metadata store app access token
	serviceAccount := kubectl_libs.GetServiceAccountJson("metadata-store-read-write-client", "metadata-store")
	secretName := serviceAccount.Secrets[0].Name
	secret := kubectl_libs.GetSecret(secretName, "metadata-store")
	encodedToken := string(secret.Data.Token)
	authToken, err := base64.StdEncoding.DecodeString(encodedToken)
	if err != nil {
		t.Error("error while decoding token")
		t.FailNow()
	}

	//getting ingresss ca cert
	caSecret := kubectl_libs.GetSecret("ingress-cert", "metadata-store").Data.CaCrt
	caEncodedToken := string(caSecret)
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

	// set url
	if !strings.HasPrefix(metadataStoreDomain, "https://") {
		metadataStoreDomain = "https://" + metadataStoreDomain
	}

	//configure tanzu insight config set-target command
	err = tanzu_libs.TanzuConfigureInsight(tempFile.Name(), string(authToken), metadataStoreDomain)
	if err != nil {
		t.FailNow()
	}

}
