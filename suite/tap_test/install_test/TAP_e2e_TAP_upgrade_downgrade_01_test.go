//go:build upgrade

package install_tests

import (
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
	"path/filepath"
	"testing"
)

func TestTapUgradeDowngrade(t *testing.T) {
	t.Log("************** TestCase START: TestTapUpgradeDowngrade **************")

	tap_1_0_2_values_file := filepath.Join(filepath.Join(utils.GetFileDir(), "../../resources/components"), "tap-values.yaml")

	testenv.Test(t,
		common_features.UpdatePackageRepository(t, suiteConfig.PackageRepository.Name, suiteConfig.UpgradeVersions.Image, suiteConfig.Tap.Namespace),
		common_features.InstallPackage(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.UpgradeVersions.TapVersion, suiteConfig.Tap.Namespace, tap_1_0_2_values_file, suiteConfig.Tap.PollTimeout),

		//innerloop before tap update
		common_features.TanzuDeployWorkload(t, suiteConfig.Innerloop.Workload.YamlFile, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyTanzuWorkloadStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyWorkloadResponse(t, suiteConfig.Innerloop.Workload.URL, suiteConfig.Innerloop.Workload.OriginalString, ""),
		common_features.GitClone(t, suiteConfig.GitCredentials.Username, suiteConfig.GitCredentials.Email, suiteConfig.Innerloop.Workload.Gitrepository),
		common_features.UpdateTiltFile(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace,""),
		common_features.TiltUp(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyTanzuWorkloadStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.ReplaceStringInFile(t, suiteConfig.Innerloop.Workload.OriginalString, suiteConfig.Innerloop.Workload.NewString, suiteConfig.Innerloop.Workload.ApplicationFilePath, suiteConfig.Innerloop.Workload.Name),
		common_features.VerifyTanzuWorkloadStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyWorkloadResponse(t, suiteConfig.Innerloop.Workload.URL, suiteConfig.Innerloop.Workload.NewString, ""),

		//outerloop before tap update
		common_features.CreateGithubRepo(t, outerloopConfig.Project.Name, outerloopConfig.Project.RepoTemplate, outerloopConfig.Project.AccessToken),
		common_features.ApplyKubectlConfigurationFile(t, outerloopConfig.Mysql.YamlFile, outerloopConfig.Namespace),
		common_features.TanzuDeployWorkload(t, outerloopConfig.Workload.YamlFile, outerloopConfig.Namespace),
		common_features.VerifyTanzuWorkloadStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyWorkloadResponse(t, outerloopConfig.Project.Host, outerloopConfig.Project.OriginalString, outerloopConfig.Project.WebpageRelativePath),
		common_features.UpdateGitRepository(t, outerloopConfig.Project.Username, outerloopConfig.Project.Email, outerloopConfig.Project.Repository, outerloopConfig.Project.Name, outerloopConfig.Project.AccessToken, outerloopConfig.Project.File, outerloopConfig.Project.OriginalString, outerloopConfig.Project.NewString, outerloopConfig.Project.CommitMessage),
		common_features.VerifyTanzuWorkloadStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyWorkloadResponse(t, outerloopConfig.Project.Host, outerloopConfig.Project.NewString, outerloopConfig.Project.WebpageRelativePath),

		//tap update
		common_features.UpdatePackageRepository(t, suiteConfig.PackageRepository.Name, suiteConfig.UpgradeVersions.Image, suiteConfig.Tap.Namespace),
		common_features.UpdateTapVersion(t, suiteConfig.Tap.Name, suiteConfig.Tap.PackageName, suiteConfig.Tap.Namespace, suiteConfig.Tap.ValuesSchemaFile, suiteConfig.UpgradeVersions.UpgradeTapVersion, suiteConfig.Tap.PollTimeout),

		//checking existing innerloop and deleting it
		common_features.VerifyTanzuWorkloadStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyWorkloadResponse(t, suiteConfig.Innerloop.Workload.URL, suiteConfig.Innerloop.Workload.NewString, ""),
		common_features.ReplaceStringInFile(t, suiteConfig.Innerloop.Workload.NewString, suiteConfig.Innerloop.Workload.OriginalString, suiteConfig.Innerloop.Workload.ApplicationFilePath, suiteConfig.Innerloop.Workload.Name),
		common_features.VerifyTanzuWorkloadStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyWorkloadResponse(t, suiteConfig.Innerloop.Workload.URL, suiteConfig.Innerloop.Workload.NewString, ""),
		common_features.InnerloopCleanUp(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),

		//checking existing outerloop and deleting it
		common_features.VerifyTanzuWorkloadStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyWorkloadResponse(t, outerloopConfig.Project.Host, outerloopConfig.Project.NewString, outerloopConfig.Project.WebpageRelativePath),
		common_features.UpdateGitRepository(t, outerloopConfig.Project.Username, outerloopConfig.Project.Email, outerloopConfig.Project.Repository, outerloopConfig.Project.Name, outerloopConfig.Project.AccessToken, outerloopConfig.Project.File, outerloopConfig.Project.NewString, outerloopConfig.Project.OriginalString, outerloopConfig.Project.CommitMessage),
		common_features.VerifyTanzuWorkloadStatus(t, outerloopConfig.Workload.Name, outerloopConfig.Namespace),
		common_features.VerifyWorkloadResponse(t, outerloopConfig.Project.Host, outerloopConfig.Project.OriginalString, outerloopConfig.Project.WebpageRelativePath),
		common_features.DeleteGithubRepo(t, outerloopConfig.Project.Name, outerloopConfig.Project.AccessToken),
		common_features.OuterloopCleanUp(t, outerloopConfig.Workload.Name, outerloopConfig.Project.Name, outerloopConfig.Namespace),

		// innerloop after tap update
		common_features.TanzuDeployWorkload(t, suiteConfig.Innerloop.Workload.YamlFile, suiteConfig.Innerloop.Workload.Namespace),
		common_features.GitClone(t, suiteConfig.GitCredentials.Username, suiteConfig.GitCredentials.Email, suiteConfig.Innerloop.Workload.Gitrepository),
		common_features.VerifyTanzuWorkloadStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyWorkloadResponse(t, suiteConfig.Innerloop.Workload.URL, suiteConfig.Innerloop.Workload.OriginalString, ""),
		common_features.UpdateTiltFile(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.TiltUp(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.ReplaceStringInFile(t, suiteConfig.Innerloop.Workload.OriginalString, suiteConfig.Innerloop.Workload.NewString, suiteConfig.Innerloop.Workload.ApplicationFilePath, suiteConfig.Innerloop.Workload.Name),
		common_features.VerifyTanzuWorkloadStatus(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),
		common_features.VerifyWorkloadResponse(t, suiteConfig.Innerloop.Workload.URL, suiteConfig.Innerloop.Workload.NewString, ""),
		common_features.InnerloopCleanUp(t, suiteConfig.Innerloop.Workload.Name, suiteConfig.Innerloop.Workload.Namespace),

		// outerloop after tap update
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

		// final cleanup,
		common_features.DeletePackage(t, suiteConfig.Tap.Name, suiteConfig.Tap.Namespace),
	)
	t.Log("************** TestCase END: TestTapUpgradeDowngrade **************")
}
