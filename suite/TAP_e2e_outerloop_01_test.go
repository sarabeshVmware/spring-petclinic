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

func TestOuterloopBasicSupplychainGitSource(t *testing.T) {
	t.Log("************** TestCase START: TestOuterloopBasicSupplychainGitSource **************")

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

	testenv.Test(t,
		updateTap,
		createGithubRepo,
		deployMysqldbService,
		deployWorkload,
		verifyGitrepoStatus,
		verifyBuildStatus,
		verifyPodintents,
		verifyDeliverables,
		verifyServiceBindings,
		verifyRevisionStatus,
		verifyKsvcStatus,
		verifyTaskrunStatus,
		verifyWorkloadStatus,
		verifyWebpageOriginal,
		gitUpdate,
		verifyBuildStatusAfterUpdate,
		verifyRevisionStatusAfterUpdate,
		verifyKsvcStatusAfterUpdate,
		verifyWebpageNew,
		removeProjectDir,
		deleteWorkload,
		deleteGithubRepo,
	)
	t.Log("************** TestCase END: TestOuterloopBasicSupplychainGitSource **************")
}
