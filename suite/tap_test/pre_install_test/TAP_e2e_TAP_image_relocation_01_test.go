//go:build all || relocation

package pre_install_test

import (
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
	"testing"
)

func TestTapImageRelocation(t *testing.T) {

	t.Log("************** TestCase START: TestTapImageRelocation **************")
	testenv.Test(t,
		common_features.DockerLoginGCP(t, suiteConfig.NonTanzuRepository.Server, suiteConfig.NonTanzuRepository.Username, suiteConfig.NonTanzuRepository.Password),
		common_features.ImgPkgCopyToRepo(t, suiteConfig.PackageRepository.Image, suiteConfig.NonTanzuRepository.Repository),
		common_features.CreateSecret(t, suiteConfig.TapRegistrySecret.Name, suiteConfig.NonTanzuRepository.Server, suiteConfig.NonTanzuRepository.Username, suiteConfig.NonTanzuRepository.Password, suiteConfig.TapRegistrySecret.Namespace, suiteConfig.TapRegistrySecret.Export, true),
		common_features.AddPackageRepository(t, suiteConfig.PackageRepository.Name, suiteConfig.NonTanzuRepository.Repository, suiteConfig.PackageRepository.Version, suiteConfig.PackageRepository.Namespace),
		common_features.InstallPackage(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.UpgradeVersions.TapVersion, suiteConfig.Tap.Namespace, suiteConfig.Tap.ValuesSchemaFile, suiteConfig.Tap.PollTimeout),

		//outerloop
		// common_features.CreateGithubRepo(t, outerloopConfig.Project.Name, outerloopConfig.Project.RepoTemplate, outerloopConfig.Project.AccessToken),
		// common_features.ApplyKubectlConfigurationFile(t, outerloopConfig.Mysql.YamlFile, outerloopConfig.Namespace),
		// common_features.TanzuDeployWorkload(t, outerloopConfig.Workload.YamlFile, outerloopConfig.Namespace),
		// common_features.VerifyTanzuWorkloadStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		// common_features.VerifyWorkloadResponse(t, outerloopConfig.Project.Host, outerloopConfig.Project.OriginalString, outerloopConfig.Project.WebpageRelativePath),
		// common_features.UpdateGitRepository(t, outerloopConfig.Project.Username, outerloopConfig.Project.Email, outerloopConfig.Project.Repository, outerloopConfig.Project.Name, outerloopConfig.Project.AccessToken, outerloopConfig.Project.File, outerloopConfig.Project.OriginalString, outerloopConfig.Project.NewString, outerloopConfig.Project.CommitMessage),
		// common_features.VerifyTanzuWorkloadStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		// common_features.VerifyWorkloadResponse(t, outerloopConfig.Project.Host, outerloopConfig.Project.NewString, outerloopConfig.Project.WebpageRelativePath),
		// common_features.DeleteGithubRepo(t, outerloopConfig.Project.Name, outerloopConfig.Project.AccessToken),
		// common_features.OuterloopCleanUp(t, outerloopConfig.Workload.Name, outerloopConfig.Project.Name, outerloopConfig.Namespace),

		common_features.DeletePackage(t, suiteConfig.Tap.Name, suiteConfig.Tap.Namespace),
		common_features.DeletePackageRepository(t, suiteConfig.PackageRepository.Name, suiteConfig.PackageRepository.Namespace),
		common_features.DeleteGCRImageRepository(t, suiteConfig.NonTanzuRepository.Repository),
	)
	t.Log("************** TestCase END: TestTapImageRelocation **************")
}
