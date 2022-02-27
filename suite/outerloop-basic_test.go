//go:build all || outerloop || outerloop_basic

package suite

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzuCmds"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestOuterloopBasic(t *testing.T) {
	t.Log("************** TestCase START: TestOuterloopBasic **************")

	updateTap := features.New("update-tap-full-supplychainbasic").
		Assess("update-package", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Log("updating tap package")

			// get schema and update values
			tapValuesSchema, err := getTapValuesSchema()
			if err != nil {
				t.Error("error while getting tap values schema")
				t.FailNow()
			}
			tapValuesSchema.Profile = "full"
			tapValuesSchema.SupplyChain = "basic"

			// create temporary file
			t.Log("creating tempfile for tap values schema")
			tempFile, err := ioutil.TempFile("", "tap-values*.yaml")
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

			// update tap
			err = tanzuCmds.TanzuUpdatePackage(suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, suiteConfig.Tap.Namespace, tempFile.Name())
			if err != nil {
				t.Error("error while updating tap")
				t.FailNow()
			} else {
				t.Log("updated tap")
				t.Logf("sleeping for 1 minute")
				time.Sleep(time.Minute)
			}

			return ctx
		}).
		Feature()

	deployApps := features.New("deploy-apps-via-yaml-configurations").
		// TODO: remove
		// Assess("deploy-springpetclinic-pipeline", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		// 	return stepfuncs.DeployAppInNamespace(ctx, t, cfg, true, outerloopConfig.SpringPetclinic.Name, []string{outerloopConfig.SpringPetclinic.YamlFile}, outerloopConfig.Namespace)
		// }).
		Assess("deploy-mysqldb", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			name, files, namespace := outerloopConfig.Mysql.Name, []string{outerloopConfig.Mysql.YamlFile}, outerloopConfig.Namespace

			t.Logf("deploying app %s in namespace %s", name, namespace)
			cmd, output, err := exec.KappDeployAppInNamespace(name, files, namespace)
			t.Logf("command executed: %s", cmd)
			if err != nil {
				t.Error(fmt.Errorf("error while deploying app %s in namespace %s: %w: %s", name, namespace, err, output))
				t.FailNow()
			}
			t.Logf("app %s deployed in namespace %s: %s", name, namespace, output)
			return ctx

			// return stepfuncs.DeployAppInNamespace(ctx, t, cfg, true, outerloopConfig.Mysql.Name, []string{outerloopConfig.Mysql.YamlFile}, outerloopConfig.Namespace)
		}).
		Feature()

	// // TODO: servicebinding check

	// TODO: new build check, ksvc revision updation check
	testenv.Test(t,
		updateTap,
		createGithubRepo,
		deployApps,
		deployWorkload,
		verifyGitrepoStatus,
		verifyBuildStatus,
		verifyPodintents,
		verifyKsvcStatus,
		verifyTaskrunStatus,
		verifyWorkloadStatus,
		verifyWebpageOriginal,
		gitUpdate,
		verifyWebpageNew,
		//gitReset,
		removeProjectDir,
		deleteWorkload,
		deleteGithubRepo,
	)
	t.Log("************** TestCase END: TestOuterloopBasic **************")
}
