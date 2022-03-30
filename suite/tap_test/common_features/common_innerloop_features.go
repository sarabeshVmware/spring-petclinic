package common_features

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzu_libs"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/models"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func UpdateTiltFile(t *testing.T, workloadName string, namespace string) features.Feature {
	return features.New("update-allow-context-tilt").
		Assess("update-tilt-file", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			tiltFile := filepath.Join(rootDir, workloadName, "/Tiltfile")
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
			tiltFile := filepath.Join(rootDir, workloadName, "/Tiltfile")
			tapValuesSchema, err := models.GetTapValuesSchema()
			if err != nil {
				t.Error("error while updating tilt file")
				t.FailNow()
			}
			source_image := fmt.Sprintf("%s/%s/%s-source", tapValuesSchema.OotbSupplyChainBasic.Registry.Server, tapValuesSchema.OotbSupplyChainBasic.Registry.Repository, workloadName)
			err = utils.ReplaceStringInFile(tiltFile, "<SOURCE_IMAGE>", source_image)
			if err != nil {
				t.Error(fmt.Errorf("Error while editing tiltfile: %w", err))
				t.FailNow()
			}
			return ctx
		}).
		Assess("update-namespace", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			t.Logf("updating namespace image in tiltfile")
			tiltFile := filepath.Join(rootDir, workloadName, "/Tiltfile")
			err := utils.ReplaceStringInFile(tiltFile, "<DEVELOPMENT_NAMESPACE>", namespace)
			if err != nil {
				t.Error(fmt.Errorf("Error while editing tiltfile: %w", err))
				t.FailNow()
			}
			return ctx
		}).
		Feature()
}

func TiltUp(t *testing.T, workloadName string, namespace string) features.Feature {
	return features.New("create-workload-tilt-up").
		Assess("tilting-up", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			//under the expectation that the clone repo has workloadName as folder name and Tiltfile in root
			tiltFile := filepath.Join(rootDir, workloadName, "Tiltfile")
			t.Logf("Setting NAMESPACE environment variable to %s", namespace)
			os.Setenv("NAMESPACE", namespace)
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
}

func InnerloopCleanUp(t *testing.T, workloadName string, namespace string) features.Feature {
	return features.New("innerloop cleanup").
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
			tanzu_libs.DeleteWorkload(workloadName, namespace)
			return ctx
		}).
		Assess("remove-dir", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			dir := filepath.Join(rootDir, workloadName)

			t.Logf("removing directory %s", dir)
			err := os.RemoveAll(dir)
			if err != nil {
				t.Error(fmt.Errorf("error while removing directory %s: %w", dir, err))
				t.FailNow()
			}
			t.Logf("directory %s removed", dir)
			return ctx
		}).
		Feature()
}
