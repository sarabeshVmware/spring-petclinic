package kubectl_libs

import (
	"encoding/json"
	"fmt"
	"log"
	linux_util "pkg/utils/linux_util"
	"time"
)

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
	return raw
}
