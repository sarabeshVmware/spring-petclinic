//go:build all || relocation

package pre_install_test

import (
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
	"strings"
	"testing"
)

func RelocateImagesAndInstallTapPackageFeature(t *testing.T, server string, username string, password string, passwordType string, repository string) {
	tapPackageSourceBundle := fmt.Sprintf("%s:%s", suiteConfig.PackageRepository.Image, suiteConfig.PackageRepository.Version)
	testenv.Test(t,
		common_features.DockerLogin(t, server, username, password),
		common_features.ImgPkgCopyToRepo(t, tapPackageSourceBundle, repository),
		common_features.CreateSecret(t, suiteConfig.TapRegistrySecret.Name, server, username, password, passwordType, suiteConfig.TapRegistrySecret.Namespace, suiteConfig.TapRegistrySecret.Export),
		common_features.CreateSecret(t, suiteConfig.RegistryCredentialsSecret.Name, suiteConfig.RegistryCredentialsSecret.Registry, suiteConfig.RegistryCredentialsSecret.Username, suiteConfig.RegistryCredentialsSecret.Password, "string", suiteConfig.RegistryCredentialsSecret.Namespace, suiteConfig.RegistryCredentialsSecret.Export),
		common_features.AddPackageRepository(t, suiteConfig.PackageRepository.Name, repository, suiteConfig.PackageRepository.Version, suiteConfig.PackageRepository.Namespace),
		common_features.InstallPackage(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Version, suiteConfig.Tap.Namespace, suiteConfig.Tap.ValuesSchemaFile, suiteConfig.Tap.PollTimeout),
		common_features.CheckIfPackageInstalled(suiteConfig.Tap.Name, suiteConfig.Tap.Namespace, 10, 60),
	)
}

func OuterloopTestFeature(t *testing.T) {

	testenv.Test(t,
		common_features.CreateGithubRepo(t, outerloopConfig.Project.Name, outerloopConfig.Project.RepoTemplate, outerloopConfig.Project.AccessToken),
		common_features.ApplyKubectlConfigurationFile(t, outerloopConfig.Mysql.YamlFile, outerloopConfig.Namespace),
		common_features.TanzuDeployWorkload(t, outerloopConfig.Workload.YamlFile, outerloopConfig.Namespace),
		common_features.VerifyTanzuWorkloadStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyWorkloadResponse(t, outerloopConfig.Project.Host, outerloopConfig.Project.OriginalString, outerloopConfig.Project.WebpageRelativePath),
		common_features.UpdateGitRepository(t, outerloopConfig.Project.Username, outerloopConfig.Project.Email, outerloopConfig.Project.Repository, outerloopConfig.Project.Name, outerloopConfig.Project.AccessToken, outerloopConfig.Project.File, outerloopConfig.Project.OriginalString, outerloopConfig.Project.NewString, outerloopConfig.Project.CommitMessage),
		common_features.VerifyTanzuWorkloadStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyWorkloadResponse(t, outerloopConfig.Project.Host, outerloopConfig.Project.NewString, outerloopConfig.Project.WebpageRelativePath),
		common_features.DeleteGithubRepo(t, outerloopConfig.Project.Name, outerloopConfig.Project.AccessToken),
		common_features.OuterloopCleanUp(t, outerloopConfig.Workload.Name, outerloopConfig.Project.Name, outerloopConfig.Namespace),
	)
}

func CleanupResourcesFeature(t *testing.T) {

	testenv.Test(t,
		common_features.DeletePackage(t, suiteConfig.Tap.Name, suiteConfig.Tap.Namespace),
		common_features.DeletePackageRepository(t, suiteConfig.PackageRepository.Name, suiteConfig.PackageRepository.Namespace),
		common_features.DeleteSecret(t, suiteConfig.TapRegistrySecret.Name, suiteConfig.TapRegistrySecret.Namespace),
		common_features.DeleteSecret(t, suiteConfig.RegistryCredentialsSecret.Name, suiteConfig.RegistryCredentialsSecret.Namespace),
		// cleanup of image repositories is tracked in https://jira.eng.vmware.com/browse/DAPEO-132
		//common_features.DeleteImageRepository(t, repository),
	)
}

func TestTapImageRelocation(t *testing.T) {

	t.Log("************** TestCase START: TestTapImageRelocation **************")

	for _, repository := range suiteConfig.NonTanzuRepository {
		t.Logf("testing imgpkg copy for %s", repository.Server)
		RelocateImagesAndInstallTapPackageFeature(t, repository.Server, repository.Username, repository.Password, repository.PasswordType, repository.Repository)
		OuterloopTestFeature(t)
		CleanupResourcesFeature(t)
	}

	t.Log("************** TestCase END: TestTapImageRelocation **************")
}
