//go:build all || install

package install_tests

import (
	"context"
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

var PackagesResourcesDir = filepath.Join(utils.GetFileDir(), "../../resources/components")

func getPackagesList() (Packages, error) {
	log.Printf("getting package list")

	pkgList := Packages{}
	file := filepath.Join(PackagesResourcesDir, "install-metadata.yaml")

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

func installUnistallPackage(t *testing.T, packageName string) {
	pkgList, _ := getPackagesList()
	for _, pkg := range pkgList {
		if pkg.Name != packageName {
			continue
		}
		t.Logf("Installing and Uninstalling package: %s", pkg.Name)
		availablePkgs := tanzu_libs.ListAvailablePackages(pkg.Package, suiteConfig.PackageRepository.Namespace)
		for index, pkgVersion := range availablePkgs {
			installIndividualPackages := features.New(fmt.Sprintf("install-uninstall-%s", packageName)).
				Assess(fmt.Sprintf("version-%s", pkgVersion.VERSION), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
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
					return ctx
				}).
				Feature()
			testenv.Test(t, installIndividualPackages)
		}
		break
	}
}

func TestInstallPackages(t *testing.T) {
	t.Log("************** TestCase START: TestInstallUninstallAllComponentAllVersionInPackageRepo **************")

	pkgList, _ := getPackagesList()

	uninstallAllIndividualPackages := features.New("uninstall-individual-packages").
		Assess("uninstallation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			for i := len(pkgList) - 1; i >= 0; i-- {
				uninstallPkg := features.New(pkgList[i].Name).
					Assess("uninstall", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
						t.Logf("uninstalling package: %s", pkgList[i].Name)
						err := tanzu_libs.DeleteInstalledPackage(pkgList[i].Name, suiteConfig.PackageRepository.Namespace)
						if err != nil {
							t.Error(fmt.Errorf("Uninstallation FAILED for package : %s", pkgList[i].Name))
							t.Fail()
						}
						return ctx
					}).
					Feature()
				testenv.Test(t, uninstallPkg)
			}
			return ctx
		}).
		Feature()

	installUninstallTapPackages := features.New("install-uninstall-tap-packages").
		Assess("install-uninstall", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			availablePkgs := tanzu_libs.ListAvailablePackages(suiteConfig.Tap.PackageName, suiteConfig.PackageRepository.Namespace)
			for _, pkgVersion := range availablePkgs {
				installUninstallTapPkg := features.New(suiteConfig.Tap.Name).
					Assess(fmt.Sprintf("version-%s", pkgVersion.VERSION), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {

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
						return ctx
					}).
					Feature()
				testenv.Test(t, installUninstallTapPkg)
			}
			return ctx
		}).
		Feature()

	installCertManagerPackages := features.New("install-cert-manager-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
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
						installCertMgr := features.New(pkg.Name).
							Assess(fmt.Sprintf("version-%s", pkgVersion.VERSION), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
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
						testenv.Test(t, installCertMgr)
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
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
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
						installContour := features.New(pkg.Name).
							Assess(fmt.Sprintf("version-%s", pkgVersion.VERSION), func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
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
						testenv.Test(t, installContour)
					}
					break
				} else {
					continue
				}
			}
			return ctx
		}).
		Feature()

	installServiceBindingsPackages := features.New("install-service-bindings-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "service-bindings")
			return ctx
		}).
		Feature()

	installSourceControllerPackages := features.New("install-source-controller-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "source-controller")
			return ctx
		}).
		Feature()

	installServicesToolkitPackages := features.New("install-services-toolkit-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "services-toolkit")
			return ctx
		}).
		Feature()

	installScanControllerPackages := features.New("install-scan-controller-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "scan-controller")
			return ctx
		}).
		Feature()

	installApiPortalPackages := features.New("install-api-portal-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "api-portal")
			return ctx
		}).
		Feature()

	installBuildServicePackages := features.New("install-buildservice-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "buildservice")
			return ctx
		}).
		Feature()

	installFluxcdSourceControllerPackages := features.New("install-fluxcd-source-controller-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "fluxcd-source-controller")
			return ctx
		}).
		Feature()

	installTektonPipelinesPackages := features.New("install-tekton-pipelines-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "tekton-pipelines")
			return ctx
		}).
		Feature()

	installTapTelemetryPackages := features.New("install-tap-telemetry-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "tap-telemetry")
			return ctx
		}).
		Feature()

	installConventionsControllerPackages := features.New("install-conventions-controller-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "conventions-controller")
			return ctx
		}).
		Feature()

	installImageWebhookPolicyPackages := features.New("install-image-policy-webhook-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "image-policy-webhook")
			return ctx
		}).
		Feature()

	installMetadataStorePackages := features.New("install-metadata-store-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "metadata-store")
			return ctx
		}).
		Feature()

	installCartographerPackages := features.New("install-cartographer-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "cartographer")
			return ctx
		}).
		Feature()

	installGrypeScannerPackages := features.New("install-grype-scanner-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "grype-scanner")
			return ctx
		}).
		Feature()

	installAppLiveViewPackages := features.New("install-appliveview-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "appliveview")
			return ctx
		}).
		Feature()

	installAppLiveViewConventionsPackages := features.New("install-appliveview-conventions-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "appliveview-conventions")
			return ctx
		}).
		Feature()

	installAcceleratorPackages := features.New("install-accelerator-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "accelerator")
			return ctx
		}).
		Feature()

	installDeveloperConventionsPackages := features.New("install-developer-conventions-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "developer-conventions")
			return ctx
		}).
		Feature()

	installCloudNativeRuntimesPackages := features.New("install-cnrs-packages").
		Assess("iinstallation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "cnrs")
			return ctx
		}).
		Feature()

	installTapGuiPackages := features.New("install-tap-gui-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "tap-gui")
			return ctx
		}).
		Feature()

	installLearningCenterPackages := features.New("install-learningcenter-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "learningcenter")
			return ctx
		}).
		Feature()

	installOotbTemplatesPackages := features.New("install-ootb-templates-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "ootb-templates")
			return ctx
		}).
		Feature()

	installSpringBootConventionsPackages := features.New("install-spring-boot-conventions-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "spring-boot-conventions")
			return ctx
		}).
		Feature()

	installLearningCenterWorkshopsPackages := features.New("install-learningcenter-workshops-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "learningcenter-workshops")
			return ctx
		}).
		Feature()

	installOotbSupplyChainBasicPackages := features.New("install-ootb-supply-chain-basic-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "ootb-supply-chain-basic")
			return ctx
		}).
		Feature()

	installOotbSupplyChainTestingPackages := features.New("install-ootb-supply-chain-testing-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "ootb-supply-chain-testing")
			return ctx
		}).
		Feature()

	installOotbSupplyChainTestingScanningPackages := features.New("install-ootb-supply-chain-testing-scanning-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "ootb-supply-chain-testing-scanning")
			return ctx
		}).
		Feature()

	installOotbDeliveryBasicPackages := features.New("install-ootb-delivery-basic-packages").
		Assess("installation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			installUnistallPackage(t, "ootb-delivery-basic")
			return ctx
		}).
		Feature()

	testenv.TestInParallel(t,
		installCertManagerPackages,
		installServiceBindingsPackages,
		installSourceControllerPackages,
		installServicesToolkitPackages,
		installScanControllerPackages,
		installApiPortalPackages,
		installBuildServicePackages,
		installFluxcdSourceControllerPackages,
		installTektonPipelinesPackages,
		installTapTelemetryPackages)
	testenv.TestInParallel(t,
		installContourPackages,
		installConventionsControllerPackages,
		installImageWebhookPolicyPackages,
		installMetadataStorePackages,
		installCartographerPackages,
		installGrypeScannerPackages)
	testenv.TestInParallel(t,
		installAppLiveViewPackages,
		installAppLiveViewConventionsPackages,
		installAcceleratorPackages,
		installDeveloperConventionsPackages,
		installCloudNativeRuntimesPackages,
		installTapGuiPackages,
		installLearningCenterPackages,
		installOotbTemplatesPackages,
		installSpringBootConventionsPackages)
	testenv.TestInParallel(t,
		installLearningCenterWorkshopsPackages,
		installOotbSupplyChainBasicPackages,
		installOotbSupplyChainTestingPackages,
		installOotbSupplyChainTestingScanningPackages,
		installOotbDeliveryBasicPackages)
	testenv.Test(t,
		uninstallAllIndividualPackages,
		installUninstallTapPackages,
	)
	t.Log("************** TestCase END: TestInstallUninstallAllComponentAllVersionInPAckageRepo **************")
}
