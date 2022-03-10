//go:build all || install

package suite

import (
	"context"
	// "io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"

	// "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_helpers"
	// "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzuCmds"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzu_libs"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

type Packages []struct {
	Name       string `yaml:"name"`
	Package    string `yaml:"package"`
	ValuesFile string `yaml:"values_file,omitempty"`
	RbacFile   string `yaml:"rbac_file,omitempty"`
	PollTimout string `yaml:"poll_timout,omitempty"`
}

var PackagesResourcesDir = filepath.Join(utils.GetFileDir(), "resources", "components")

func getPackagesList() (Packages, error) {
	log.Printf("getting package list")

	pkgList := Packages{}
	file := filepath.Join(PackagesResourcesDir, "packages.yaml")

	// read file
	pkgListBytes, err := os.ReadFile(file)
	if err != nil {
		log.Printf("error while reading packages file %s", file)
		log.Printf("error: %s", err)
		return pkgList, err
	} else {
		log.Printf("read packages file %s", file)
	}

	// unmarshall
	err = yaml.Unmarshal(pkgListBytes, &pkgList)
	if err != nil {
		log.Printf("error while unmarshalling packages file %s", file)
		log.Printf("error: %s", err)
		return pkgList, err
	} else {
		log.Printf("unmarshalled file %s", file)
	}

	return pkgList, nil
}

func TestInstallUninstallAllComponentAllVersionInPackageRepo(t *testing.T) {
	t.Log("************** TestCase START: TestInstallUninstallAllComponentAllVersionInPackageRepo **************")

	pkgList, _ := getPackagesList()
	f1 := features.New("install-individual-packages").
		Assess("install-individual-packages", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			for _, pkg := range pkgList {
				t.Logf("pkgname: %s", pkg.Name)
				availablePkgs := tanzu_libs.ListAvailablePackages(pkg.Package, "tap-install")
				for _, pkgVersion := range availablePkgs {
					t.Logf("version: %s", pkg.VERSION)
				}
			}
			return ctx
		}).
		Feature()

	f2 := features.New("uninstall-individual-packages").
		Assess("uninstall-individual-packages", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			for _, pkg := range pkgList {
				t.Logf("pkgname: %s", pkg.Name)
			}
			return ctx
		}).
		Feature()

	f3 := features.New("install-uninstall-tap-packages").
		Assess("install-uninstall-tap-packages", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			for _, pkg := range pkgList {
				t.Logf("pkgname: %s", pkg.Name)
			}
			return ctx
		}).
		Feature()

	/*
	   for _, pkg := pkgList   {
	           log.Printf("package : %s", pkg.Name)
	           //fetch versions n list them first
	           f1 := features.New("install-uninstall-cert-manager").
	                   Assess("install-component", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
	                           kubectl_helpers.GetWorkloadStatus("")
	                           return ctx
	                   }).
	                   Assess("update-tap", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
	                           t.Logf("updating package %s", suiteConfig.Tap.Name)
	                           cmd, output, err := exec.TanzuUpdatePackage(suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, suiteConfig.Tap.Namespace, suiteConfig.Tap.ValuesSchemaFile)
	                           t.Logf("command executed: %s", cmd)
	                           if err != nil {
	                                   t.Error(fmt.Errorf("error while updating package %s: %w: %s", suiteConfig.Tap.Name, err, output))
	                                   t.FailNow()
	                           }
	                           t.Logf("package %s updated: %s", suiteConfig.Tap.Name, output)
	                           t.Logf("sleeping for 1 minute")
	                           time.Sleep(time.Minute)
	                           return ctx
	           }).
	           Feature()
	   }

	*/
	testenv.Test(t,
		f1,
		f2,
		f3,
	)
	t.Log("************** TestCase END: TestInstallUninstallAllComponentAllVersionInPAckageRepo **************")
}
