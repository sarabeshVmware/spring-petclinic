//go:build all || service_toolkit

package service_toolkit_test

import (
	"testing"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/tap_test/common_features"
	"path/filepath"
	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils"
)

func TestServicetoolkit(t *testing.T) {
	t.Log("************** TestCase START: TestServicetoolkit **************")
	resource_claims_and_rbac_rmq_file:= filepath.Join(filepath.Join(utils.GetFileDir(), "../../resources/service_toolkit"), "resource-claims-rmq.yaml")
	rabbitmq_cluster_service_instance_file:=filepath.Join(filepath.Join(utils.GetFileDir(), "../../resources/service_toolkit"), "rabbitmq-cluster-service-instance.yaml")
	testenv.Test(t,
		common_features.InstallAndDeployRmqOperator(t, suiteConfig.ServiceToolkit.Name, suiteConfig.ServiceToolkit.Gitrepository),
		common_features.ApplyKubectlConfigurationFile(t, resource_claims_and_rbac_rmq_file, suiteConfig.RegistryCredentialsSecret.Namespace),

		//Create a RabbitMQ service instance 
		common_features.ApplyKubectlConfigurationFile(t, rabbitmq_cluster_service_instance_file, suiteConfig.RegistryCredentialsSecret.Namespace),
		common_features.VerifyRabbitmqClustersStatus(t, "", suiteConfig.RegistryCredentialsSecret.Namespace),
		common_features.ServiceInstanceList(t,suiteConfig.RegistryCredentialsSecret.Namespace),

		// // Create workload
		common_features.TanzuCreateWorkloadWithGitRepo(t, suiteConfig.ServiceToolkit.WorkloadName, suiteConfig.ServiceToolkit.WorkloadRepository, suiteConfig.RegistryCredentialsSecret.Namespace),
		common_features.VerifyBuildStatus(t, suiteConfig.ServiceToolkit.WorkloadName, suiteConfig.ServiceToolkit.BuildNameSuffix, suiteConfig.RegistryCredentialsSecret.Namespace),
		common_features.VerifyRevisionStatus(t, suiteConfig.ServiceToolkit.WorkloadName, suiteConfig.RegistryCredentialsSecret.Namespace),
		common_features.VerifyKsvcStatus(t, suiteConfig.ServiceToolkit.WorkloadName, suiteConfig.RegistryCredentialsSecret.Namespace),
		common_features.GetKsvcUrl(t, suiteConfig.ServiceToolkit.WorkloadName, suiteConfig.RegistryCredentialsSecret.Namespace),
		common_features.VerifyWorkloadResponse(t, suiteConfig.ServiceToolkit.WorkloadURL, suiteConfig.ServiceToolkit.Message, ""),
		common_features.TanzuDeleteWorkload(t, suiteConfig.ServiceToolkit.WorkloadName, suiteConfig.CreateNamespaces[0]),
	)
	t.Log("************** TestCase END: TestServicetoolkit **************")
}
