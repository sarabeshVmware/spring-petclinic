package suite

import (
	"os"
	"path/filepath"
	"testing"

	tap "gitlab.eng.vmware.com/tap/tap-packages/tap-packaging-tests/pkg"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
)

var testenv env.Environment

func TestMain(m *testing.M) {
	tap.CheckPrerequisites() // TODO: create function to move to Setup and Finish

	home, err := os.UserHomeDir()
	tap.CheckError(err)
	testenv = env.NewWithKubeConfig(filepath.Join(home, ".kube", "config"))
	namespace := "tap-install"

	testenv.Setup(
		envfuncs.CreateNamespace(namespace),
	)

	testenv.Finish(
		envfuncs.DeleteNamespace(namespace),
	)

	os.Exit(testenv.Run(m))
}
