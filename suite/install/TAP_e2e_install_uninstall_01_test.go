//go:build all || install

package install_tests

import (
	"context"
	// "io/ioutil"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectlCmds"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/kubectl/kubectl_helpers"
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
	//latestPkgList := tanzu_libs.ListAllAvailablePackages(suiteConfig.PackageRepository.Namespace)
	installAllIndividualPackages := features.New("install-individual-packages").
		Assess("install-individual-packages", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			for _, pkg := range pkgList {
				if pkg.Name == "cert-manager" || pkg.Name == "contour" {
					continue
				}
				t.Logf("Installing package: %s", pkg.Name)
				availablePkgs := tanzu_libs.ListAvailablePackages(pkg.Package, suiteConfig.PackageRepository.Namespace)
				for index, pkgVersion := range availablePkgs {
					t.Logf("package name: %s, version: %s", pkg.Name, pkgVersion.VERSION)
					tanzu_libs.InstallPackage(pkg.Name, pkg.Package, pkgVersion.VERSION, suiteConfig.PackageRepository.Namespace, pkg.ValuesFile, pkg.PollTimout)
					installed := tanzu_helpers.ValidateInstalledPackageStatus(pkg.Name, suiteConfig.PackageRepository.Namespace, 5, 30)
					if installed {
						t.Logf("Installed package : %s, version: %s successfully", pkg.Name, pkgVersion.VERSION)
					} else {
						t.Error(fmt.Errorf("Installation FAILED for package : %s, version: %s", pkg.Name, pkgVersion.VERSION))
						t.Fail()
					}
					if index != len(availablePkgs)-1 {
						err := tanzu_libs.DeleteInstalledPackage(pkg.Name, suiteConfig.PackageRepository.Namespace)
						if err != nil {
							t.Error(fmt.Errorf("Uninstallation FAILED for package : %s, version: %s", pkg.Name, pkgVersion.VERSION))
							t.Fail()
						}
					}
				}
			}
			return ctx
		}).
		Feature()

	uninstallAllIndividualPackages := features.New("uninstall-individual-packages").
		Assess("uninstall-individual-packages", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			for i := len(pkgList) - 1; i >= 0; i-- {
				t.Logf("pkgname: %s", pkgList[i].Name)
				err := tanzu_libs.DeleteInstalledPackage(pkgList[i].Name, suiteConfig.PackageRepository.Namespace)
				if err != nil {
					t.Error(fmt.Errorf("Uninstallation FAILED for package : %s", pkgList[i].Name))
					t.Fail()
				}
			}
			return ctx
		}).
		Feature()

	installUninstallTapPackages := features.New("install-uninstall-tap-packages").
		Assess("install-uninstall-tap-packages", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			availablePkgs := tanzu_libs.ListAvailablePackages(suiteConfig.Tap.PackageName, suiteConfig.PackageRepository.Namespace)
			for _, pkgVersion := range availablePkgs {
				tanzu_libs.InstallPackage(suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, pkgVersion.VERSION, suiteConfig.Tap.Namespace, suiteConfig.Tap.ValuesSchemaFile, suiteConfig.Tap.PollTimeout)
				installed := kubectl_helpers.ValidateTAPInstallation(suiteConfig.Tap.Name, suiteConfig.Tap.Namespace, 10, 60)
				if !installed {
					kubectl_helpers.LogFailedResourcesDetails(suiteConfig.Tap.Namespace)
					t.Error(fmt.Errorf("error while installing package %s (%s)", suiteConfig.Tap.Name, suiteConfig.Tap.Namespace))
					t.Fail()
				} else {
					t.Logf("Installed package : %s, version: %s successfully", suiteConfig.Tap.Name, pkgVersion.VERSION)
				}
				err := tanzu_libs.DeleteInstalledPackage(suiteConfig.Tap.Name, suiteConfig.PackageRepository.Namespace)
				if err != nil {
					t.Error(fmt.Errorf("Uninstallation FAILED for package : %s, version: %s", suiteConfig.Tap.Name, pkgVersion.VERSION))
					t.Fail()
				}
			}
			return ctx
		}).
		Feature()

	installCertManagerPackages := features.New("install-cert-manager-packages").
		Assess("install-cert-manager-packages", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			for _, pkg := range pkgList {
				if pkg.Name == "cert-manager" {
					availablePkgs := tanzu_libs.ListAvailablePackages(pkg.Package, suiteConfig.PackageRepository.Namespace)
					err := kubectlCmds.KubectlApplyConfiguration(pkg.RbacFile, suiteConfig.PackageRepository.Namespace)
					if err != nil {
						t.Error("error while applying cert-manager rbac file")
						t.Fail()
					} else {
						t.Log("applied cert-manager rbac")
					}
					oldText := "<VERSION>"
					for index, pkgVersion := range availablePkgs {
						// write new feature every time
						installCertMgr := features.New("install-cert-mgr").
							Assess(fmt.Sprintf("install-cert-mgr-version-%s", pkgVersion.VERSION), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
								t.Logf("pkg: %s, version: %s", pkg.Name, pkgVersion.VERSION)
								utils.ReplaceStringInFile(pkg.ValuesFile, oldText, pkgVersion.VERSION)
								oldText = pkgVersion.VERSION
								kubectlCmds.KubectlApplyConfiguration(pkg.ValuesFile, suiteConfig.PackageRepository.Namespace)
								installed := tanzu_helpers.ValidateInstalledPackageStatus(pkg.Name, suiteConfig.PackageRepository.Namespace, 5, 30)
								if installed {
									t.Logf("Installed package : %s, version: %s successfully", pkg.Name, pkgVersion.VERSION)
								} else {
									t.Error(fmt.Errorf("Installation FAILED for package : %s, version: %s", pkg.Name, pkgVersion.VERSION))
									t.Fail()
								}
								if index != len(availablePkgs)-1 {
									err := tanzu_libs.DeleteInstalledPackage(pkg.Name, suiteConfig.PackageRepository.Namespace)
									if err != nil {
										t.Error(fmt.Errorf("Uninstallation FAILED for package : %s, version: %s", pkg.Name, pkgVersion.VERSION))
										t.Fail()
									}
								}
								return ctx
							}).
							Feature()
					}
					break
				} else {
					continue
				}
			}
			return ctx
		}).
		Feature()

	installContourPackages := features.New("install-contour-packages").
		Assess("install-contour-packages", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			for _, pkg := range pkgList {
				if pkg.Name == "contour" {
					availablePkgs := tanzu_libs.ListAvailablePackages(pkg.Package, suiteConfig.PackageRepository.Namespace)
					err := kubectlCmds.KubectlApplyConfiguration(pkg.RbacFile, suiteConfig.PackageRepository.Namespace)
					if err != nil {
						t.Error("error while applying contour rbac file")
						t.Fail()
					} else {
						t.Log("applied contour rbac")
					}
					oldText := "<VERSION>"
					for index, pkgVersion := range availablePkgs {
						t.Logf("pkg: %s, version: %s", pkg.Name, pkgVersion.VERSION)
						utils.ReplaceStringInFile(pkg.ValuesFile, oldText, pkgVersion.VERSION)
						oldText = pkgVersion.VERSION
						kubectlCmds.KubectlApplyConfiguration(pkg.ValuesFile, suiteConfig.PackageRepository.Namespace)
						installed := tanzu_helpers.ValidateInstalledPackageStatus(pkg.Name, suiteConfig.PackageRepository.Namespace, 5, 30)
						if installed {
							t.Logf("Installed package : %s, version: %s successfully", pkg.Name, pkgVersion.VERSION)
						} else {
							t.Error(fmt.Errorf("Installation FAILED for package : %s, version: %s", pkg.Name, pkgVersion.VERSION))
							t.Fail()
						}
						if index != len(availablePkgs)-1 {
							err := tanzu_libs.DeleteInstalledPackage(pkg.Name, suiteConfig.PackageRepository.Namespace)
							if err != nil {
								t.Error(fmt.Errorf("Uninstallation FAILED for package : %s, version: %s", pkg.Name, pkgVersion.VERSION))
								t.Fail()
							}
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
		installCertManagerPackages,
		installContourPackages,
		installAllIndividualPackages,
		uninstallAllIndividualPackages,
		installUninstallTapPackages,
	)
	t.Log("************** TestCase END: TestInstallUninstallAllComponentAllVersionInPAckageRepo **************")
}
