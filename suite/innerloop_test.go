package suite

import (
	"context"
	"fmt"
	"os"
	//"time"
	"testing"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/client"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/exec"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)
func TestInnerloopBasic(t *testing.T) {
	// f1 := features.New("update-tap-light-supplychainbasic").
	// 	Assess("update-schema", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
	// 		tapValuesSchema.Profile = "light"
	// 		tapValuesSchema.SupplyChain = "basic"
	//         tapValuesSchema.Accelerator.Server.ServiceType = "LoadBalancer"
	// 		t.Logf("updating tap values schema %s", config.Tap.ValuesSchemaFile)
	// 		err := WriteYAMLFile(config.Tap.ValuesSchemaFile, tapValuesSchema)
	// 		if err != nil {
	// 			t.Error(fmt.Errorf("error while updating tap values schema %s: %w", config.Tap.ValuesSchemaFile, err))
	// 			t.FailNow()
	// 		}
	// 		t.Logf("tap values schema %s updated", config.Tap.ValuesSchemaFile)
	// 		return ctx
	// 	}).
	// 	Assess("update-tap", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
	// 		t.Logf("updating package %s", config.Tap.Name)
	// 		cmd, output, err := exec.TanzuUpdatePackage(config.Tap.Name, config.Tap.PackageName, config.Tap.Version, config.Tap.Namespace, config.Tap.ValuesSchemaFile)
	// 		t.Logf("command executed: %s", cmd)
	// 		if err != nil {
	// 			t.Error(fmt.Errorf("error while updating package %s: %w: %s", config.Tap.Name, err, output))
	// 			t.FailNow()
	// 		}
	// 		t.Logf("package %s updated: %s", config.Tap.Name, output)
	// 		t.Logf("sleeping for 1 minute")
 	// 		time.Sleep(time.Minute)
	// 		return ctx
	// 	}).
	// 	Feature()

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
			return context.WithValue(ctx, accServerExternalIpKey, serviceExternalIp)
		}).
		Feature()

	acceleratorNameKey := "acceleratorName"
	f3:= features.New("generate-acc-project").
	Assess("generate-project", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		acceleratorProject := "tanzu-java-web-app"
		acceleratorName := "tanzu-java-web-app"
		repositoryPrefix := tapValuesSchema.OotbSupplyChainBasic.Registry.Server + "/" + tapValuesSchema.OotbSupplyChainBasic.Registry.Repository
		t.Logf("generating accelerator project %s (namespace %s)", acceleratorProject, config.Tap.Namespace)
		cmd, output, err := exec.TanzuGenerateAccelerator(acceleratorName, acceleratorProject, repositoryPrefix, ctx.Value(accServerExternalIpKey).(string) , config.Tap.Namespace)
		t.Logf("command executed: %s", cmd)
		if err != nil {
			t.Error(fmt.Errorf("error while generating accelerator project %s in namespace %s: %w: %s", acceleratorProject, config.Tap.Namespace, err, output))
			t.FailNow()
		}
		t.Logf("Accelerator project %s generated in namespace %s: %s", acceleratorProject, config.Tap.Namespace, output)
		return context.WithValue(ctx, acceleratorNameKey, acceleratorName)
	}).
	Feature()
	
	// f4:= features.New("unzip-acc-project-zip").
	// Assess("unzip-project", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
	// 	zipFile := ctx.Value(acceleratorNameKey).(string) + ".zip"
	// 	t.Logf("Listing accelerator project zip %s)", zipFile)			
	// 	output, err := exec.RunCommand(fmt.Sprintf("ls -lt %s", zipFile))
	// 	t.Logf("command executed: ls -lt %s. output %s", zipFile, output)
	// 	if err != nil {
	// 		t.Error(fmt.Errorf("error while listing accelerator project zip file %s: %w: %s", zipFile, err, output))
	// 		t.FailNow()
	// 	}
	// 	t.Logf("Listing existing project files if exists")
	// 	output, err = exec.RunCommand(fmt.Sprintf("ls -lt %s", ctx.Value(acceleratorNameKey).(string)))
	// 	t.Logf("command executed: ls -lt %s. output %s", ctx.Value(acceleratorNameKey).(string), output)
	// 	if err == nil {
	// 		t.Logf("Deleting %s folder", ctx.Value(acceleratorNameKey))
	// 		output, err := exec.RunCommand(fmt.Sprintf("rm -rf %s", ctx.Value(acceleratorNameKey).(string)))
	// 		t.Logf("command executed: rm -rf %s. output %s", ctx.Value(acceleratorNameKey).(string), output)
	// 		if err != nil {
	// 			t.Error(fmt.Errorf("error while Deleting project files %s: %w: %s",  ctx.Value(acceleratorNameKey).(string), err, output))
	// 			t.FailNow()
	// 		}
	// 	t.Logf("Unzip %s", zipFile)
	// 	output, err = exec.RunCommand(fmt.Sprintf("unzip %s", zipFile))
	// 	t.Logf("command executed: unzip %s. output %s", zipFile, output)
	// 	if err != nil {
	// 		t.Error(fmt.Errorf("error while unzip accelerator project zip file %s: %w: %s", zipFile, err, output))
	// 		t.FailNow()
	// 	}
	// 	}
	// 	t.Logf("Accelerator project zip files %s unzipped successfully", zipFile)
	// 	return ctx
	// }).
	// Feature()

	f5:= features.New("create-workload-tilt-up").
	Assess("tilting-up", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		os.Setenv("NAMESPACE", "tap-install")
		tiltFile := ctx.Value(acceleratorNameKey).(string) + "/Tiltfile"
		tiltCmd := fmt.Sprintf("tilt ci --file %s --port 11222", tiltFile)
		t.Logf("Running tilt command %s", tiltCmd)
		output, err := exec.RunCommand(tiltCmd)
		t.Logf("command executed: %s", tiltCmd)
		if err != nil {
			t.Error(fmt.Errorf("error while tilting-up : %w: %s", err, output))
			t.FailNow()
		}
		return ctx
	}).
	Feature()
	// f6:= features.New("create-workload-tilt-up").
	// Assess("tilting-up", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
	// 	os.Setenv("NAMESPACE", "tap-install")
	// 	tiltFile := ctx.Value(acceleratorNameKey).(string) + "/Tiltfile"
	// 	tiltCmd := fmt.Sprintf("tilt ci --file %s --port 11222", tiltFile)
	// 	t.Logf("Running tilt command %s", tiltCmd)
	// 	output, err := exec.RunCommand(tiltCmd)
	// 	t.Logf("command executed: %s", tiltCmd)
	// 	if err != nil {
	// 		t.Error(fmt.Errorf("error while tilting-up : %w: %s", err, output))
	// 		t.FailNow()
	// 	}
	// 	return ctx
	// }).
	// Feature()
	testenv.Test(t, f2, f3, f5)
}