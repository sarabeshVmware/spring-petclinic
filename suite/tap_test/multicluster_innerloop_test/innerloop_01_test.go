//go:build all || multicluster_innerloop || multicluster_innerloop_basic

package multicluster_test

import (
	"context"
	"io/ioutil"
	"os"
	 "testing"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzuCmds"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_libs"
	//"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_helpers"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
	"sigs.k8s.io/e2e-framework/pkg/features"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/models"
)

func TestMulticlusterInnerloopBasicSupplychainLocalSource(t *testing.T) {
	t.Log("************** TestCase START: TestMulticlusterInnerloopBasicSupplychainLocalSource **************")
	updateTap := features.New("update-tap-full-supplychainbasic").
		Assess("update-package", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("updating tap package")

			// get schema and update values
			tapValuesSchema, err := models.GetProfileTapValuesSchema("view")
			if err != nil {
				t.Error("error while getting tap values schema")
				t.FailNow()
			}
			// switch to view cluster
			t.Log("changing to view context")
			kubectl_libs.UseContext(suiteConfig.Multicluster.ViewClusterContext)
			t.Log("switched cluster")
			tapValuesSchema.Profile = "view"
			tapValuesSchema.Accelerator.Server.ServiceType = "LoadBalancer"
	// 		tapGuiExternalIP := kubectl_helpers.GetServiceExternalIP("server", "tap-gui", 2, 30)
	// 		tapGuiUrl := fmt.Sprintf("http://%s:7000", tapGuiExternalIP)
	// 		tapValuesSchema.TapGui.AppConfig.Backend.BaseURL = tapGuiUrl
	// 		tapValuesSchema.TapGui.AppConfig.Backend.Cors.Origin = tapGuiUrl
	// 		tapValuesSchema.TapGui.AppConfig.App.BaseURL = tapGuiUrls
	// 		// switch to iterate cluster
	// 		kubectl_libs.UseContext(suiteConfig.Multicluster.IterateClusterContext)
	// 		tapValuesSchema.TapGui.AppConfig.Kubernetes.ClusterLocatorMethods[0].Clusters[0].URL = kubectl_helper.GetCurrentClusterURL()
	// 		// tapValuesSchema.TapGui.AppConfig.Kubernetes.ClusterLocatorMethods[0].Clusters[0].Name = "iterate-cluster"
	// 		// tapValuesSchema.TapGui.AppConfig.Kubernetes.ClusterLocatorMethods[0].Clusters[0].AuthProvider = "serviceAccount"
	// 		tapValuesSchema.TapGui.AppConfig.Kubernetes.ClusterLocatorMethods[0].Clusters[0].ServiceAccountToken = kubectl_helper.GetClusterToken("tap-gui-viewer", "tap-gui")
	// 		tapValuesSchema.TapGui.AppConfig.Kubernetes.ClusterLocatorMethods[0].Clusters[0].SkipTLSVerify = true
			
	// 		// create temporary file
			t.Log("creating tempfile for tap values schema")
			tempFile, err := ioutil.TempFile("", "view-tap-values*.yaml")
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
			// switch to view cluster
			kubectl_libs.UseContext(suiteConfig.Multicluster.ViewClusterContext)
			//update tap
			err = tanzuCmds.TanzuUpdatePackage(suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, suiteConfig.Tap.Namespace, tempFile.Name())
			if err != nil {
				t.Error("error while updating tap")
				t.FailNow()
			} else {
				t.Log("updated tap")
			}

			return ctx
		}).
		Feature()

	testenv.Test(t,
		updateTap,
		// Switch to View cluster
		common_features.ChangeContext(t, suiteConfig.Multicluster.ViewClusterContext),
		common_features.GenerateAcceleratorProject(t, suiteConfig.Tap.Namespace),
		// switch to Iterate cluster
		common_features.ChangeContext(t, suiteConfig.Multicluster.IterateClusterContext),
		common_features.UpdateTanzuJavaWebAppTiltFile(t, suiteConfig.Innerloop.Workload.Name),
		common_features.TiltUp(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyTanzuJavaWebAppImageRepository(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyTanzuJavaWebAppBuildStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.BuildNameSuffix, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyTanzuJavaWebAppImagesKpacStatus(t, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyTanzuJavaWebAppPodIntentStatus(t, suiteConfig.Innerloop.Workload.Name,  suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyTanzuJavaWebAppImageRepositoryDelivery(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.ImageDeliverySuffix, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyTanzuJavaWebAppRevisionStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyTanzuJavaWebAppKsvcStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyTanzuWorkloadStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyWorkloadResponse(t, suiteConfig.Innerloop.Workload.URL, suiteConfig.Innerloop.Workload.OriginalString, ""),
		common_features.ReplaceStringInFile(t, suiteConfig.Innerloop.Workload.OriginalString, suiteConfig.Innerloop.Workload.NewString, suiteConfig.Innerloop.Workload.ApplicationFilePath, suiteConfig.Innerloop.Workload.Name),
		common_features.VerifyWorkloadResponse(t, suiteConfig.Innerloop.Workload.URL, suiteConfig.Innerloop.Workload.NewString, ""),
		common_features.VerifyTanzuJavaWebAppDeliverable(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.InnerloopCleanUp(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
	)
	t.Log("************** TestCase END: TestMulticlusterInnerloopBasicSupplychainLocalSource **************")
}
