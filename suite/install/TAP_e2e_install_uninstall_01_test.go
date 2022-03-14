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
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectlCmds"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/tanzu/tanzu_helpers"
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

var PackagesResourcesDir = filepath.Join(utils.GetFileDir(), "../resources/components")

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
		//return pkgList, err
	} else {
		log.Printf("unmarshalled file %s", file)
	}

	for index, _ := range pkgList {
		if pkgList[index].ValuesFile != "" {
			pkgList[index].ValuesFile = filepath.Join(PackagesResourcesDir, pkgList[index].ValuesFile)
		}
		if pkgList[index].RbacFile != "" {
			pkgList[index].RbacFile = filepath.Join(PackagesResourcesDir, pkgList[index].RbacFile)
		}
	}

	return pkgList, err
}

func TestInstallUninstallAllComponentAllVersionInPackageRepo(t *testing.T) {
	t.Log("************** TestCase START: TestInstallUninstallAllComponentAllVersionInPackageRepo **************")

	pkgList, _ := getPackagesList()
	//latestPkgList := tanzu_libs.ListAllAvailablePackages("tap-install")
	f1 := features.New("install-individual-packages").
		Assess("install-individual-packages", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			for _, pkg := range pkgList {
				t.Logf("pkgname: %s", pkg.Name)
				availablePkgs := tanzu_libs.ListAvailablePackages(pkg.Package, "tap-install")
				for index, pkgVersion := range availablePkgs {
					t.Logf("version: %s", pkgVersion.VERSION)
					tanzu_libs.InstallePackage(pkg.Name, pkg.Package, pkgVersion.VERSION, "tap-install", pkg.ValuesFile, pkg.PollTimout)
					tanzu_helpers.ValidateInstalledPackageStatus(pkg.Name, "tap-install", 5, 30)
					if index != len(availablePkgs)-1 {
						tanzu_libs.DeleteInstalledPackage(pkg.Package, "tap-install")
					}
				}
			}
			return ctx
		}).
		Feature()

	f2 := features.New("uninstall-individual-packages").
		Assess("uninstall-individual-packages", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			for _, pkg := range pkgList {
				t.Logf("pkgname: %s", pkg.Name)
				tanzu_libs.DeleteInstalledPackage(pkg.Package, "tap-install")
			}
			return ctx
		}).
		Feature()

	f3 := features.New("install-uninstall-tap-packages").
		Assess("install-uninstall-tap-packages", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			for _, pkg := range pkgList {
				if pkg.Name == "tap" {
					availablePkgs := tanzu_libs.ListAvailablePackages(pkg.Package, "tap-install")
					for _, pkgVersion := range availablePkgs {
						tanzu_libs.InstallePackage(pkg.Name, pkg.Package, pkgVersion.VERSION, "tap-install", pkg.ValuesFile, pkg.PollTimout)
						tanzu_helpers.ValidateInstalledPackageStatus(pkg.Name, "tap-install", 15, 60)
						tanzu_libs.DeleteInstalledPackage(pkg.Package, "tap-install")
					}
					break
				} else {
					continue
				}
			}
			return ctx
		}).
		Feature()

	f4 := features.New("install-cert-manager-packages").
		Assess("install-cert-manager-packages", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			for _, pkg := range pkgList {
				if pkg.Name == "cert-manager" {
					availablePkgs := tanzu_libs.ListAvailablePackages(pkg.Package, "tap-install")
					err := kubectlCmds.KubectlApplyConfiguration(pkg.RbacFile, "tap-install")
					if err != nil {
						t.Error("error while applying cert-manager rbac file")
						t.FailNow()
					} else {
						t.Log("applied cert-manager rbac")
					}
					for index, pkgVersion := range availablePkgs {
						t.Logf("version: %s", pkgVersion.VERSION)
						oldText := "<VERSION>"
						utils.ReplaceStringInFile(pkg.ValuesFile, oldText, pkgVersion.VERSION)
						oldText = pkgVersion.VERSION
						kubectlCmds.KubectlApplyConfiguration(pkg.ValuesFile, "tap-install")
						tanzu_helpers.ValidateInstalledPackageStatus(pkg.Name, "tap-install", 5, 30)
						if index != len(availablePkgs)-1 {
							tanzu_libs.DeleteInstalledPackage(pkg.Package, "tap-install")
						}
					}
					break
				} else {
					continue
				}
			}
			return ctx
		}).
		Feature()

	f5 := features.New("install-contour-packages").
		Assess("install-contour-packages", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			for _, pkg := range pkgList {
				if pkg.Name == "contour" {
					availablePkgs := tanzu_libs.ListAvailablePackages(pkg.Package, "tap-install")
					err := kubectlCmds.KubectlApplyConfiguration(pkg.RbacFile, "tap-install")
					if err != nil {
						t.Error("error while applying contour rbac file")
						t.FailNow()
					} else {
						t.Log("applied contour rbac")
					}
					for index, pkgVersion := range availablePkgs {
						t.Logf("version: %s", pkgVersion.VERSION)
						oldText := "<VERSION>"
						utils.ReplaceStringInFile(pkg.ValuesFile, oldText, pkgVersion.VERSION)
						oldText = pkgVersion.VERSION
						kubectlCmds.KubectlApplyConfiguration(pkg.ValuesFile, "tap-install")
						tanzu_helpers.ValidateInstalledPackageStatus(pkg.Name, "tap-install", 5, 30)
						if index != len(availablePkgs)-1 {
							tanzu_libs.DeleteInstalledPackage(pkg.Package, "tap-install")
						}
					}
					break
				} else {
					continue
				}
			}
			return ctx
		}).
		Feature()

	testenv.Test(t,
		f4,
		f5,
		f1,
		f2,
		f3,
	)
	t.Log("************** TestCase END: TestInstallUninstallAllComponentAllVersionInPAckageRepo **************")
}
