//go:build all || dummy

package install_tests

import (
	"context"
	"os"
	"testing"

	"path/filepath"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	"time"
)

var (
	clusterName string
	namespace   string
	testEnv     env.Environment
)

func TestMain(m *testing.M) {
	home, _ := os.UserHomeDir()
	testEnv = env.NewWithKubeConfig(filepath.Join(home, ".kube", "config"))
	os.Exit(testEnv.Run(m))
}

func TestPodBringUp(t *testing.T) {
	featureOne := features.New("Feature One").
		Assess("Create Nginx Deployment 1", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			time.Sleep(10 * time.Second)
			t.Log("feature1-assess1")

			return ctx
		}).
		Assess("Wait for Nginx Deployment 1 to be scaled up", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			t.Log("feature1-assess2")
			return ctx

		}).Feature()

	featureTwo := features.New("Feature Two").
		Assess("Create Nginx Deployment 2", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			t.Log("feature2-assess1")
			return ctx
		}).
		Assess("Wait for Nginx Deployment 2 to be scaled up", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			t.Log("feature2-assess2")
			return ctx
		}).Feature()

	testEnv.TestInParallel(t, featureOne, featureTwo)
}
