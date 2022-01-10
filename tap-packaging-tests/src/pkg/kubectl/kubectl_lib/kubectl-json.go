package kubectl_lib

import (
	"encoding/json"
	"fmt"
	"log"
	linux_util "pkg/utils/linux_util"
	"time"
)

type GetPodintentJsonOutput struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		CreationTimestamp time.Time `json:"creationTimestamp"`
		Generation        int       `json:"generation"`
		Labels            struct {
			AppKubernetesIoComponent       string `json:"app.kubernetes.io/component"`
			CartoRunClusterSupplyChainName string `json:"carto.run/cluster-supply-chain-name"`
			CartoRunClusterTemplateName    string `json:"carto.run/cluster-template-name"`
			CartoRunResourceName           string `json:"carto.run/resource-name"`
			CartoRunTemplateKind           string `json:"carto.run/template-kind"`
			CartoRunWorkloadName           string `json:"carto.run/workload-name"`
			CartoRunWorkloadNamespace      string `json:"carto.run/workload-namespace"`
		} `json:"labels"`
		Name            string `json:"name"`
		Namespace       string `json:"namespace"`
		OwnerReferences []struct {
			APIVersion         string `json:"apiVersion"`
			BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
			Controller         bool   `json:"controller"`
			Kind               string `json:"kind"`
			Name               string `json:"name"`
			UID                string `json:"uid"`
		} `json:"ownerReferences"`
		ResourceVersion string `json:"resourceVersion"`
		UID             string `json:"uid"`
	} `json:"metadata"`
	Spec struct {
		ServiceAccountName string `json:"serviceAccountName"`
		Template           struct {
			Metadata struct {
				Annotations struct {
					AppsTanzuVmwareComLiveUpdate         string `json:"apps.tanzu.vmware.com/live-update"`
					DeveloperConventionsTargetContainers string `json:"developer.conventions/target-containers"`
				} `json:"annotations"`
				Labels struct {
					AppKubernetesIoComponent       string `json:"app.kubernetes.io/component"`
					AppsTanzuVmwareComWorkloadType string `json:"apps.tanzu.vmware.com/workload-type"`
					CartoRunWorkloadName           string `json:"carto.run/workload-name"`
				} `json:"labels"`
			} `json:"metadata"`
			Spec struct {
				Containers []struct {
					Image     string `json:"image"`
					Name      string `json:"name"`
					Resources struct {
					} `json:"resources"`
					SecurityContext struct {
						RunAsUser int `json:"runAsUser"`
					} `json:"securityContext"`
				} `json:"containers"`
				ServiceAccountName string `json:"serviceAccountName"`
			} `json:"spec"`
		} `json:"template"`
	} `json:"spec"`
	Status struct {
		Conditions []struct {
			LastTransitionTime time.Time `json:"lastTransitionTime"`
			Status             string    `json:"status"`
			Type               string    `json:"type"`
		} `json:"conditions"`
		ObservedGeneration int `json:"observedGeneration"`
		Template           struct {
			Metadata struct {
				Annotations struct {
					AppsTanzuVmwareComLiveUpdate                    string `json:"apps.tanzu.vmware.com/live-update"`
					AutoscalingKnativeDevMaxScale                   string `json:"autoscaling.knative.dev/maxScale"`
					AutoscalingKnativeDevMinScale                   string `json:"autoscaling.knative.dev/minScale"`
					BootSpringIoActuator                            string `json:"boot.spring.io/actuator"`
					BootSpringIoVersion                             string `json:"boot.spring.io/version"`
					ConventionsAppsTanzuVmwareComAppliedConventions string `json:"conventions.apps.tanzu.vmware.com/applied-conventions"`
					DeveloperAppsTanzuVmwareComImageSourceDigest    string `json:"developer.apps.tanzu.vmware.com/image-source-digest"`
					DeveloperConventionsTargetContainers            string `json:"developer.conventions/target-containers"`
				} `json:"annotations"`
				Labels struct {
					AppKubernetesIoComponent               string `json:"app.kubernetes.io/component"`
					AppsTanzuVmwareComWorkloadType         string `json:"apps.tanzu.vmware.com/workload-type"`
					CartoRunWorkloadName                   string `json:"carto.run/workload-name"`
					ConventionsAppsTanzuVmwareComFramework string `json:"conventions.apps.tanzu.vmware.com/framework"`
					TanzuAppLiveView                       string `json:"tanzu.app.live.view"`
					TanzuAppLiveViewApplicationFlavours    string `json:"tanzu.app.live.view.application.flavours"`
					TanzuAppLiveViewApplicationName        string `json:"tanzu.app.live.view.application.name"`
				} `json:"labels"`
			} `json:"metadata"`
			Spec struct {
				Containers []struct {
					Env []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"env"`
					Image string `json:"image"`
					Name  string `json:"name"`
					Ports []struct {
						ContainerPort int    `json:"containerPort"`
						Protocol      string `json:"protocol"`
					} `json:"ports"`
					Resources struct {
					} `json:"resources"`
					SecurityContext struct {
						RunAsUser int `json:"runAsUser"`
					} `json:"securityContext"`
				} `json:"containers"`
				ServiceAccountName string `json:"serviceAccountName"`
			} `json:"spec"`
		} `json:"template"`
	} `json:"status"`
}

func GetPodintentJson(name string, namespace string) *GetPodintentJsonOutput {
	cmd := fmt.Sprintf("kubectl get podintent %s -n %s -o json", name, namespace)
	res1, err1 := linux_util.ExecuteCmd(cmd)
	if err1 != nil {
		log.Println("something bad happened")
	}
	in := []byte(res1)
	var raw *GetPodintentJsonOutput
	if err := json.Unmarshal(in, &raw); err != nil {
		panic(err)
	}
	return raw
}

type GetImageRepositoriesJsonOutput struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		CreationTimestamp time.Time `json:"creationTimestamp"`
		Generation        int       `json:"generation"`
		Labels            struct {
			AppKubernetesIoComponent       string `json:"app.kubernetes.io/component"`
			CartoRunClusterSupplyChainName string `json:"carto.run/cluster-supply-chain-name"`
			CartoRunClusterTemplateName    string `json:"carto.run/cluster-template-name"`
			CartoRunResourceName           string `json:"carto.run/resource-name"`
			CartoRunTemplateKind           string `json:"carto.run/template-kind"`
			CartoRunWorkloadName           string `json:"carto.run/workload-name"`
			CartoRunWorkloadNamespace      string `json:"carto.run/workload-namespace"`
		} `json:"labels"`
		Name            string `json:"name"`
		Namespace       string `json:"namespace"`
		OwnerReferences []struct {
			APIVersion         string `json:"apiVersion"`
			BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
			Controller         bool   `json:"controller"`
			Kind               string `json:"kind"`
			Name               string `json:"name"`
			UID                string `json:"uid"`
		} `json:"ownerReferences"`
		ResourceVersion string `json:"resourceVersion"`
		UID             string `json:"uid"`
	} `json:"metadata"`
	Spec struct {
		Image    string `json:"image"`
		Interval string `json:"interval"`
	} `json:"spec"`
	Status struct {
		Artifact struct {
			Checksum       string    `json:"checksum"`
			LastUpdateTime time.Time `json:"lastUpdateTime"`
			Path           string    `json:"path"`
			Revision       string    `json:"revision"`
			URL            string    `json:"url"`
		} `json:"artifact"`
		Conditions []struct {
			LastTransitionTime time.Time `json:"lastTransitionTime"`
			Status             string    `json:"status"`
			Type               string    `json:"type"`
		} `json:"conditions"`
		ObservedGeneration int    `json:"observedGeneration"`
		URL                string `json:"url"`
	} `json:"status"`
}

func GetImageRepositoriesJson(name string, namespace string) *GetImageRepositoriesJsonOutput {
	cmd := fmt.Sprintf("kubectl get imagerepositories %s -n %s -o json", name, namespace)
	res1, err1 := linux_util.ExecuteCmd(cmd)
	if err1 != nil {
		log.Println("something bad happened")
	}
	in := []byte(res1)
	var raw *GetImageRepositoriesJsonOutput
	if err := json.Unmarshal(in, &raw); err != nil {
		panic(err)
	}
	log.Printf("apiVersion: %s", raw.APIVersion)
	log.Printf("Kind: %s", raw.Kind)
	fmt.Print("Conditions are:")
	for _, element := range raw.Status.Conditions {
		fmt.Printf("%s ==> %s", element.Type, element.Status)
	}
	return raw
}

type GetRunnablesJsonOutput struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		CreationTimestamp time.Time `json:"creationTimestamp"`
		Generation        int       `json:"generation"`
		Labels            struct {
			CartoRunClusterSupplyChainName string `json:"carto.run/cluster-supply-chain-name"`
			CartoRunClusterTemplateName    string `json:"carto.run/cluster-template-name"`
			CartoRunResourceName           string `json:"carto.run/resource-name"`
			CartoRunTemplateKind           string `json:"carto.run/template-kind"`
			CartoRunWorkloadName           string `json:"carto.run/workload-name"`
			CartoRunWorkloadNamespace      string `json:"carto.run/workload-namespace"`
		} `json:"labels"`
		Name            string `json:"name"`
		Namespace       string `json:"namespace"`
		OwnerReferences []struct {
			APIVersion         string `json:"apiVersion"`
			BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
			Controller         bool   `json:"controller"`
			Kind               string `json:"kind"`
			Name               string `json:"name"`
			UID                string `json:"uid"`
		} `json:"ownerReferences"`
		ResourceVersion string `json:"resourceVersion"`
		UID             string `json:"uid"`
	} `json:"metadata"`
	Spec struct {
		Inputs struct {
			SourceRevision string `json:"source-revision"`
			SourceURL      string `json:"source-url"`
		} `json:"inputs"`
		RunTemplateRef struct {
			Kind string `json:"kind"`
			Name string `json:"name"`
		} `json:"runTemplateRef"`
		Selector struct {
			MatchingLabels struct {
				AppsTanzuVmwareComPipeline string `json:"apps.tanzu.vmware.com/pipeline"`
			} `json:"matchingLabels"`
			Resource struct {
				APIVersion string `json:"apiVersion"`
				Kind       string `json:"kind"`
			} `json:"resource"`
		} `json:"selector"`
	} `json:"spec"`
	Status struct {
		Conditions []struct {
			LastTransitionTime time.Time `json:"lastTransitionTime"`
			Message            string    `json:"message"`
			Reason             string    `json:"reason"`
			Status             string    `json:"status"`
			Type               string    `json:"type"`
		} `json:"conditions"`
		ObservedGeneration int `json:"observedGeneration"`
		Outputs            struct {
			Revision string `json:"revision"`
			URL      string `json:"url"`
		} `json:"outputs"`
	} `json:"status"`
}

func GetRunnablesJson(name string, namespace string) *GetRunnablesJsonOutput {
	cmd := fmt.Sprintf("kubectl get runnables %s -n %s -o json", name, namespace)
	res1, err1 := linux_util.ExecuteCmd(cmd)
	if err1 != nil {
		log.Println("something bad happened")
	}
	in := []byte(res1)
	var raw *GetRunnablesJsonOutput
	if err := json.Unmarshal(in, &raw); err != nil {
		panic(err)
	}
	log.Printf("apiVersion: %s", raw.APIVersion)
	log.Printf("Kind: %s", raw.Kind)
	return raw
}
