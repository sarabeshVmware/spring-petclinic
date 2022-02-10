//go:build outerloop

package suite

import (
	"context"
	"fmt"
	"testing"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/exec"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_helpers"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzu_helpers"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestOuterloopTestScan(t *testing.T) {
	t.Log("****************TestOuterloopTestScan execution started****************")
	updateTap := features.New("update-tap-full-supplychain-testing_scanning").
		Assess("update-schema", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			tapValuesSchema.Profile = "full"
			tapValuesSchema.SupplyChain = "testing_scanning"
			if err := utils.WriteYAMLFile(suiteConfig.Tap.ValuesSchemaFile, tapValuesSchema); err != nil {
				t.Error(fmt.Errorf("error while writing file %s", suiteConfig.Tap.ValuesSchemaFile))
				t.FailNow()
			}
			return ctx

			// return stepfuncs.WriteFile(ctx, t, cfg, true, suiteConfig.Tap.ValuesSchemaFile, tapValuesSchema)
		}).
		Assess("update-tap", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			name, packageName, version, namespace, valuesSchemaFile := suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, suiteConfig.Tap.Namespace, suiteConfig.Tap.ValuesSchemaFile

			t.Logf("updating package %s", name)
			cmd, output, err := exec.TanzuUpdatePackage(name, packageName, version, namespace, valuesSchemaFile)
			t.Logf("command executed: %s", cmd)
			if err != nil {
				t.Error(fmt.Errorf("error while updating package %s: %w: %s", name, err, output))
				t.FailNow()
			}
			t.Logf("package %s updated: %s", name, output)
			return ctx

			// return stepfuncs.UpdatePackage(ctx, t, cfg, true, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, suiteConfig.Tap.Namespace, suiteConfig.Tap.ValuesSchemaFile)
		}).
		Feature()

	verifyPackageInstalled := features.New("check-grype-and-scanning-package-installed").
		Assess("check-grype", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("Checking grype is installed in the %s namespace", suiteConfig.Tap.Namespace)
			if !tanzu_helpers.IsGrypeInstalled(suiteConfig.Tap.Namespace) {
				t.Error(fmt.Errorf("grype is not installed in the %s namespace", suiteConfig.Tap.Namespace))
				t.Fail()
			}
			return ctx
		}).
		Assess("check-scanning", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("Checking scanning is installed in the %s namespace", suiteConfig.Tap.Namespace)
			if !tanzu_helpers.IsScanningInstalled(suiteConfig.Tap.Namespace) {
				t.Error(fmt.Errorf("scanning is not installed in the %s namespace", suiteConfig.Tap.Namespace))
				t.Fail()
			}
			return ctx
		}).
		Feature()

	deployScanPolicy := features.New("deploy-scan-policy-app-via-yaml-configurations").
	Assess("deploy-scan-policy", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		file, namespace := outerloopConfig.ScanPolicy.YamlFile, outerloopConfig.Namespace

		t.Logf("deploying scan-policy %s in namespace %s", file, namespace)
		cmd, output, err := exec.KubectlApplyConfiguration(file, namespace)
		t.Logf("command executed: %s", cmd)
		if err != nil {
			t.Error(fmt.Errorf("error while deploying scan-policy %s in namespace %s: %w: %s", file, namespace, err, output))
			t.FailNow()
		}
		t.Logf("scan policy %s deployed in namespace %s: %s", file, namespace, output)
		return ctx
	}).
	Feature()
	
	verifyPipelineStatus := features.New("verify-pipeline-status").
	Assess("verify-pipeline-installed", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		if !kubectl_helpers.ValidatePipelineExists("", outerloopConfig.Namespace) {
			t.Error(fmt.Errorf("Pipeline is not installed in namespace %s", outerloopConfig.Namespace))
			t.FailNow()
		}
		return ctx

	}).
	Feature()

	verifySourceScanStatus := features.New("verify-source-scan-status").
	Assess("verify-source-scan-completed", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		if !kubectl_helpers.ValidateSourceScans("", outerloopConfig.Namespace) {
			t.Error(fmt.Errorf("Source scan is not Completed in namespace %s", outerloopConfig.Namespace))
			t.FailNow()
		}
		return ctx

	}).
	Feature()

	verifyImageScanStatus := features.New("verify-image-scan-status").
	Assess("verify-image-scan-completed", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		if !kubectl_helpers.ValidateImageScans("", outerloopConfig.Namespace) {
			t.Error(fmt.Errorf("Image scan is not Completed in namespace %s", outerloopConfig.Namespace))
			t.FailNow()
		}
		return ctx

	}).
	Feature()

	verifyPipelineRunStatus := features.New("verify-pipeline-runs-status").
	Assess("verify-pipeline-runs-succeded", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		if !kubectl_helpers.ValidatePipelineRuns("", outerloopConfig.Namespace) {
			t.Error(fmt.Errorf("Pipeline Runs succeded in namespace %s", outerloopConfig.Namespace))
			t.FailNow()
		}
		return ctx

	}).
	Feature()

	/* updateWorkloadFileuncomment := features.New("uncomment-test-label").
	Assess("uncomment-test-label-from workload-file", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		oldString := "# apps.tanzu.vmware.com/has-tests: true"
		newString := "apps.tanzu.vmware.com/has-tests: true"
		filePath := outerloopConfig.Workload.YamlFile
		t.Logf("Replace from string %s to string %s in file %s", oldString, newString, filePath)
		err := exec.ReplaceStringInFile(filePath, oldString, newString)
		if err != nil {
			t.Error(fmt.Errorf("error while replacing string in file %s : %w", filePath, err))
			t.FailNow()
		}
		return ctx

	}).
	Feature()

	updateWorkloadFilecomment := features.New("comment-test-label").
	Assess("comment-test-label-from workload-file", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		oldString := "apps.tanzu.vmware.com/has-tests: true"
		newString := "# apps.tanzu.vmware.com/has-tests: true"
		filePath := outerloopConfig.Workload.YamlFile
		t.Logf("Replace from string %s to string %s in file %s", oldString, newString, filePath)
		err := exec.ReplaceStringInFile(filePath, oldString, newString)
		if err != nil {
			t.Error(fmt.Errorf("error while replacing string in file %s : %w", filePath, err))
			t.FailNow()
		}
		return ctx

	}).
	Feature() */
	testenv.Test(t,
		updateTap,
		verifyPackageInstalled,
		deployScanPolicy,
		deployMysqlService,
		deployPipeline,
		verifyPipelineStatus,
		deployWorkloadWithTest,
		verifyGitrepoStatus,
		verifyPipelineRunStatus,
		verifySourceScanStatus,
		verifyImageskpac,
		verifyBuildStatus,
		verifyImageScanStatus,
		verifyPodintents,
		verifyTaskrunStatus,
		verifyKsvcStatus,
		verifyWorkloadStatus,
		getEnvoyExternalIP,
		verifyApplicationRunningOriginal,
		gitUpdate,
		verifyApplicationRunningNew,
		gitReset,
		cleanRemoveProjectDir,
		deleteWorkload,
	)
}
