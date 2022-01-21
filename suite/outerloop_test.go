package suite

import (
	"context"
	"fmt"

	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/stepfuncs"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestOuterloopBasic(t *testing.T) {
	f1 := features.New("update-tap-full-supplychainbasic").
		Assess("update-schema", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			tapValuesSchema.Profile = "full"
			tapValuesSchema.SupplyChain = "basic"
			return stepfuncs.WriteFile(ctx, t, cfg, config.Tap.ValuesSchemaFile, tapValuesSchema)
		}).
		Assess("update-tap", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			return stepfuncs.UpdatePackage(ctx, t, cfg, config.Tap.Name, config.Tap.PackageName, config.Tap.Version, config.Tap.Namespace, config.Tap.ValuesSchemaFile)
		}).
		Feature()

	// tapGuiServerExternalIpKey, tapGuiServerPortKey := "tapGuiServerExternalIp", "tapGuiServerPort"

	// f2 := features.New("get-tapgui-server-externalip-port").
	// 	Assess("get-externalip", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
	// 		ctx, serviceExternalIp := stepfuncs.GetServiceExternalIp(ctx, t, cfg, "server", "tap-gui")
	// 		return context.WithValue(ctx, tapGuiServerExternalIpKey, serviceExternalIp)
	// 	}).
	// 	Assess("get-port", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
	// 		ctx, serviceExternalPort := stepfuncs.GetServiceExternalPort(ctx, t, cfg, "server", "tap-gui")
	// 		return context.WithValue(ctx, tapGuiServerPortKey, serviceExternalPort)
	// 	}).
	// 	Feature()

	// f3 := features.New("update-tap-tapgui-schema").
	// 	Assess("update-schema", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
	// 		baseUrl := ctx.Value(tapGuiServerExternalIpKey).(string) + ":" + fmt.Sprint(ctx.Value(tapGuiServerPortKey).(int))
	// 		appTitle, catalogType := "TAP", "url"
	// 		tapValuesSchema.TapGui.AppConfig.App.Title = appTitle
	// 		tapValuesSchema.TapGui.AppConfig.App.BaseURL = baseUrl
	// 		tapValuesSchema.TapGui.AppConfig.Backend.BaseURL = baseUrl
	// 		tapValuesSchema.TapGui.AppConfig.Backend.Cors.Origin = baseUrl
	// 		tapValuesSchema.TapGui.AppConfig.Catalog.Locations = make([]struct {
	// 			Target string "yaml:\"target,omitempty\""
	// 			Type   string "yaml:\"type,omitempty\""
	// 		}, 1)
	// 		tapValuesSchema.TapGui.AppConfig.Catalog.Locations[0].Type = catalogType
	// 		tapValuesSchema.TapGui.AppConfig.Catalog.Locations[0].Target = config.Outerloop.CatalogInfoYaml
	// 		return stepfuncs.WriteFile(ctx, t, cfg, config.Tap.ValuesSchemaFile, tapValuesSchema)
	// 	}).
	// 	Assess("update-tap", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
	// 		return stepfuncs.UpdatePackage(ctx, t, cfg, config.Tap.Name, config.Tap.PackageName, config.Tap.Version, config.Tap.Namespace, config.Tap.ValuesSchemaFile)
	// 	}).
	// 	Feature()

	f4 := features.New("deploy-apps-via-yaml-configurations").
		Assess("deploy-springpetclinic-pipeline", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			return stepfuncs.DeployAppInNamespace(ctx, t, cfg, config.Outerloop.SpringPetclinic.Name, []string{config.Outerloop.SpringPetclinic.YamlFile}, config.Outerloop.SpringPetclinic.Namespace)
		}).
		Assess("deploy-mysqldb", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			return stepfuncs.DeployAppInNamespace(ctx, t, cfg, config.Outerloop.Mysql.Name, []string{config.Outerloop.Mysql.YamlFile}, config.Outerloop.Mysql.Namespace)
		}).
		Feature()

	f5 := features.New("patch-default-serviceaccount").
		Assess("patch-imagepullsecrets", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			return stepfuncs.PatchServiceAccount(ctx, t, cfg, fmt.Sprintf(`'{"imagePullSecrets": [{"name": "%s"}, {"name": "%s"}]}'`, config.TanzunetCredsSecret.Name, config.ImageSecret.Name), "default", config.Outerloop.Namespace)
		}).
		Assess("patch-secrets", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			return stepfuncs.PatchServiceAccount(ctx, t, cfg, fmt.Sprintf(`'{"secrets": [{"name": "%s"}]}'`, config.ImageSecret.Name), "default", config.Outerloop.Namespace)
		}).
		Feature()

	f6 := features.New("create-clusterrolebinding").
		Assess("create-clusterrolebinding", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			return stepfuncs.CreateClusterRoleBinding(ctx, t, cfg, "apps-admin", "cluster-admin", fmt.Sprintf("%s:default", config.Outerloop.Namespace))
		}).
		Feature()

	f7 := features.New("deploy-workload").
		Assess("deploy-workload", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			return stepfuncs.DeployWorkload(ctx, t, cfg, config.Outerloop.Workload.YamlFile, config.Outerloop.Workload.Namespace)
		}).
		Feature()

	f8 := features.New("verify-imagerepository-status").
		Assess("verify-imagerepository-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			return stepfuncs.VerifyImagerepositoryReady(ctx, t, cfg, config.Outerloop.SpringPetclinic.ImagerepositoryName, config.Outerloop.Workload.Namespace)
		}).
		Feature()

	f9 := features.New("verify-gitrepo-status").
		Assess("verify-gitrepo-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			return stepfuncs.VerifyGitrepoReady(ctx, t, cfg, config.Outerloop.SpringPetclinic.GitrepositoryName, config.Outerloop.Workload.Namespace)
		}).
		Feature()

	f10 := features.New("verify-build-status").
		Assess("verify-build-succeeded", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			return stepfuncs.VerifyBuildSucceeded(ctx, t, cfg, config.Outerloop.SpringPetclinic.BuildNamePrefix, config.Outerloop.Workload.Namespace)
		}).
		Feature()

	f11 := features.New("verify-podintent-status").
		Assess("verify-podintent-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			return stepfuncs.VerifyPodintentReady(ctx, t, cfg, config.Outerloop.SpringPetclinic.PodintentName, config.Outerloop.Workload.Namespace)
		}).
		Feature()

	f12 := features.New("verify-podintent-annotations").
		Assess("verify-applied-conventions", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			return stepfuncs.VerifyPodintentAnnotation(ctx, t, cfg, "conventions.apps.tanzu.vmware.com/applied-conventions", "", true, config.Outerloop.SpringPetclinic.PodintentName, config.Outerloop.Workload.Namespace)
		}).
		Assess("verify-developer-conventions", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			return stepfuncs.VerifyPodintentAnnotation(ctx, t, cfg, "developer.conventions/target-containers", "workload", false, config.Outerloop.SpringPetclinic.PodintentName, config.Outerloop.Workload.Namespace)
		}).
		Feature()

	f13 := features.New("verify-podintent-labels").
		Assess("verify-appliveview", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			return stepfuncs.VerifyPodintentLabel(ctx, t, cfg, "tanzu.app.live.view", "true", false, config.Outerloop.SpringPetclinic.PodintentName, config.Outerloop.Workload.Namespace)
		}).
		Assess("verify-appliveview-applicatoin-flavours", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			return stepfuncs.VerifyPodintentLabel(ctx, t, cfg, "tanzu.app.live.view.application.flavours", "spring-boot", false, config.Outerloop.SpringPetclinic.PodintentName, config.Outerloop.Workload.Namespace)
		}).
		Assess("verify-appliveview-application-name", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			return stepfuncs.VerifyPodintentLabel(ctx, t, cfg, "tanzu.app.live.view.application.name", "petclinic", false, config.Outerloop.SpringPetclinic.PodintentName, config.Outerloop.Workload.Namespace)
		}).
		Assess("verify-springboot-conventions", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			return stepfuncs.VerifyPodintentLabel(ctx, t, cfg, "conventions.apps.tanzu.vmware.com/framework", "spring-boot", false, config.Outerloop.SpringPetclinic.PodintentName, config.Outerloop.Workload.Namespace)
		}).
		Feature()

	f14 := features.New("verify-ksvc-status").
		Assess("verify-ksvc-ready", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			return stepfuncs.VerifyKsvcReady(ctx, t, cfg, config.Outerloop.SpringPetclinic.KsvcName, config.Outerloop.Workload.Namespace)
		}).
		Feature()

	f15 := features.New("verify-taskrun-status").
		Assess("verify-taskrun-succeeded", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			return stepfuncs.VerifyTaskrunSucceeded(ctx, t, cfg, config.Outerloop.SpringPetclinic.TaskrunNamePrefix, config.Outerloop.Workload.Namespace)
		}).
		Feature()

	testenv.Test(t, f1, f4, f5, f6, f7, f8, f9, f10, f11, f12, f13, f14, f15)
}
