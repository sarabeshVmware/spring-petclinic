package kubectl_libs

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

type GetPodintentOutput struct {
	NAME, READY, REASON, AGE string
}

func GetPodintent(name string, namespace string) []GetPodintentOutput {
	podIntents := []GetPodintentOutput{}
	cmd := "kubectl get podintent"
	if name != "" {
		cmd += fmt.Sprintf(" %s", name)
	}
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += " -A"
	}
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return podIntents
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return podIntents
	}

	ss := linux_util.FieldIndices(temp[0])
	headers := linux_util.GetFields(temp[0], ss)
	for _, element := range temp[1:] {
		words := linux_util.GetFields(element, ss)
		var podIntent GetPodintentOutput
		for index, value := range words {
			reflect.ValueOf(&podIntent).Elem().FieldByName(headers[index]).SetString(value)
		}
		podIntents = append(podIntents, podIntent)
	}
	fmt.Printf("podIntents: %+v\n", podIntents)
	return podIntents
}

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
