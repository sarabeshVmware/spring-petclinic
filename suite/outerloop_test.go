package suite

import (
	"context"
	"fmt"
	"time"

	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/client"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/exec"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestOuterloopBasic(t *testing.T) {
	f1 := features.New("update-tap-full-supplychainbasic").
		Assess("update-schema", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			tapValuesSchema.Profile = "full"
			tapValuesSchema.SupplyChain = "basic"

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
			return ctx
		}).
		Feature()

	tapGuiServerExternalIpKey, tapGuiServerPortKey := "tapGuiServerExternalIp", "tapGuiServerPort"

	f2 := features.New("get-tapgui-server-externalip-port").
		Assess("get-externalip", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			service, tapGuiNamespace := "server", "tap-gui"

			t.Logf("getting external ip for %s (namespace %s)", service, tapGuiNamespace)
			serviceExternalIp, err := client.GetServiceExternalIP(service, tapGuiNamespace, cfg.Client().RESTConfig())
			if err != nil {
				t.Error(fmt.Errorf("error while getting external ip for %s (namespace %s): %w", service, tapGuiNamespace, err))
				t.FailNow()
			}
			t.Logf("external ip for %s (namespace %s): %s", "server", tapGuiNamespace, serviceExternalIp)
			return context.WithValue(ctx, tapGuiServerExternalIpKey, serviceExternalIp)
		}).
		Assess("get-port", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			service, tapGuiNamespace := "server", "tap-gui"

			t.Logf("getting port for %s (namespace %s)", service, tapGuiNamespace)
			servicePort, err := client.GetServicePort(service, tapGuiNamespace, cfg.Client().RESTConfig())
			if err != nil {
				t.Error(fmt.Errorf("error while getting port for %s (namespace %s): %w", service, tapGuiNamespace, err))
				t.FailNow()
			}
			t.Logf("port for %s (namespace %s): %d", service, tapGuiNamespace, servicePort)
			return context.WithValue(ctx, tapGuiServerPortKey, servicePort)
		}).
		Feature()

	f3 := features.New("update-tap-tapgui-schema").
		Assess("update-schema", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			baseUrl := ctx.Value(tapGuiServerExternalIpKey).(string) + ":" + fmt.Sprint(ctx.Value(tapGuiServerPortKey).(int))
			appTitle, catalogType := "TAP", "url"

			tapValuesSchema.TapGui.AppConfig.App.Title = appTitle
			tapValuesSchema.TapGui.AppConfig.App.BaseURL = baseUrl
			tapValuesSchema.TapGui.AppConfig.Backend.BaseURL = baseUrl
			tapValuesSchema.TapGui.AppConfig.Backend.Cors.Origin = baseUrl
			tapValuesSchema.TapGui.AppConfig.Catalog.Locations = make([]struct {
				Target string "yaml:\"target,omitempty\""
				Type   string "yaml:\"type,omitempty\""
			}, 1)
			tapValuesSchema.TapGui.AppConfig.Catalog.Locations[0].Type = catalogType
			tapValuesSchema.TapGui.AppConfig.Catalog.Locations[0].Target = config.Outerloop.CatalogInfoYaml

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
			return ctx
		}).
		Feature()

	f4 := features.New("deploy-apps-via-yaml-configurations").
		Assess("deploy-springpetclinic-pipeline", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("deploying app %s in namespace %s", config.Outerloop.SpringPetclinic.Name, config.Outerloop.SpringPetclinic.Namespace)
			cmd, output, err := exec.KappDeployAppInNamespace(config.Outerloop.SpringPetclinic.Name, []string{config.Outerloop.SpringPetclinic.YamlFile}, config.Outerloop.SpringPetclinic.Namespace)
			t.Logf("command executed: %s", cmd)
			if err != nil {
				t.Error(fmt.Errorf("error while deploying app %s in namespace %s: %w: %s", config.Outerloop.SpringPetclinic.Name, config.Outerloop.SpringPetclinic.Namespace, err, output))
				t.FailNow()
			}
			t.Logf("app %s deployed in namespace %s: %s", config.Outerloop.SpringPetclinic.Name, config.Outerloop.SpringPetclinic.Namespace, output)
			return ctx
		}).
		Assess("deploy-mysqldb", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("deploying app %s in namespace %s", config.Outerloop.Mysql.Name, config.Outerloop.Mysql.Namespace)
			cmd, output, err := exec.KappDeployAppInNamespace(config.Outerloop.Mysql.Name, []string{config.Outerloop.Mysql.YamlFile}, config.Outerloop.Mysql.Namespace)
			t.Logf("command executed: %s", cmd)
			if err != nil {
				t.Error(fmt.Errorf("error while deploying app %s in namespace %s: %w: %s", config.Outerloop.Mysql.Name, config.Outerloop.Mysql.Namespace, err, output))
				t.FailNow()
			}
			t.Logf("app %s deployed in namespace %s: %s", config.Outerloop.Mysql.Name, config.Outerloop.Mysql.Namespace, output)
			return ctx
		}).
		Feature()

	f5 := features.New("patch-default-serviceaccount").
		Assess("patch-imagepullsecrets", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			serviceAccount := "default"

			t.Logf("patching imagepullsecrets %s, %s to service account %s", config.TanzunetCredsSecret.Name, config.ImageSecret.Name, serviceAccount)
			cmd, output, err := exec.KubectlPatchServiceAccount(serviceAccount, config.Outerloop.Namespace, fmt.Sprintf(`'{"imagePullSecrets": [{"name": "%s"}, {"name": "%s"}]}'`, config.TanzunetCredsSecret.Name, config.ImageSecret.Name))
			t.Logf("command executed: %s", cmd)
			if err != nil {
				t.Error(fmt.Errorf("error while patching imagepullsecrets %s, %s to service account %s: %w: %s", config.TanzunetCredsSecret.Name, config.ImageSecret.Name, serviceAccount, err, output))
				t.FailNow()
			}
			t.Logf("imagepullsecrets %s, %s patched to service account %s: %s", config.TanzunetCredsSecret.Name, config.ImageSecret.Name, serviceAccount, output)
			return ctx
		}).
		Assess("patch-secrets", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			serviceAccount := "default"

			t.Logf("patching secret %s to service account %s", config.ImageSecret.Name, serviceAccount)
			cmd, output, err := exec.KubectlPatchServiceAccount(serviceAccount, config.Outerloop.Namespace, fmt.Sprintf(`'{"secrets": [{"name": "%s"}]}'`, config.ImageSecret.Name))
			t.Logf("command executed: %s", cmd)
			if err != nil {
				t.Error(fmt.Errorf("error while patching secret %s to service account %s: %w: %s", config.ImageSecret.Name, serviceAccount, err, output))
				t.FailNow()
			}
			t.Logf("secret %s patched to service account %s: %s", config.ImageSecret.Name, serviceAccount, output)
			return ctx
		}).
		Feature()

	f6 := features.New("create-clusterrolebinding").
		Assess("create-clusterrolebinding", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			name, clusterRole, serviceAccount := "apps-admin", "cluster-admin", fmt.Sprintf("%s:default", config.Outerloop.Namespace)

			t.Logf("creating clusterrolebinding %s for clusterrole %s and serviceaccount %s", name, clusterRole, serviceAccount)
			cmd, output, err := exec.KubectlCreateClusterRoleBinding(name, clusterRole, serviceAccount)
			t.Logf("command executed: %s", cmd)
			if err != nil {
				t.Error(fmt.Errorf("error while creating cluster role binding %s for clusterrole %s and serviceaccount %s: %w: %s", name, clusterRole, serviceAccount, err, output))
				t.FailNow()
			}
			t.Logf("clusterrolebinding %s created for clusterrole %s and serviceaccount %s: %s", name, clusterRole, serviceAccount, output)
			return ctx
		}).
		Feature()

	f7 := features.New("deploy-workload").
		Assess("deploy-workload", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("deploying workload %s in namespace %s", config.Outerloop.Workload.YamlFile, config.Outerloop.Workload.Namespace)
			cmd, output, err := exec.TanzuDeployWorkload(config.Outerloop.Workload.YamlFile, config.Outerloop.Workload.Namespace)
			t.Logf("command executed: %s", cmd)
			if err != nil {
				t.Error(fmt.Errorf("error while deploying workload %s in namespace %s: %w: %s", config.Outerloop.Workload.YamlFile, config.Outerloop.Workload.Namespace, err, output))
				t.FailNow()
			}
			t.Logf("workload %s deployed in namespace %s: %s", config.Outerloop.Workload.YamlFile, config.Outerloop.Workload.Namespace, output)
			return ctx
		}).
		Feature()

	f8 := features.New("verify-imagerepository-status").
		Assess("verify-imagerepository", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			ready := false
			iter := 30
			for i := 1; i <= iter; i++ {
				t.Logf("getting image repository status for %s in namespace %s (iteration %d)", config.Outerloop.SpringPetclinic.ImageRepositoryName, config.Outerloop.Workload.Namespace, i)
				cmd, output, err := exec.KubectlGetImageRepositoryStatus(config.Outerloop.SpringPetclinic.ImageRepositoryName, config.Outerloop.Workload.Namespace)
				t.Logf("command executed: %s", cmd)
				if err != nil {
					t.Error(fmt.Errorf("error while getting image repository status for %s in namespace %s: %w: %s", config.Outerloop.SpringPetclinic.ImageRepositoryName, config.Outerloop.Workload.Namespace, err, output))
					t.FailNow()
				}
				t.Logf("workload status for %s in namespace %s: %s", config.Outerloop.SpringPetclinic.ImageRepositoryName, config.Outerloop.Workload.Namespace, output)
				if output == "Ready" {
					ready = true
					break
				}
				t.Logf("sleeping for 1 minute")
				time.Sleep(time.Minute)
			}
			if !ready {
				t.Errorf(`image repository failed to get into ready state after %d iterations`, iter)
				t.Fail()
			}
			return ctx
		}).
		Feature()

	testenv.Test(t, f1, f2, f3, f4, f5, f6, f7, f8)
}
