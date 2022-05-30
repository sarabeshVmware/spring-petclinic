package kubectl_libs

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
	"gopkg.in/yaml.v3"
)

type GetDeliverablesOutput struct {
	NAME, SOURCE, DELIVERY, READY, REASON, AGE string
}

func GetDeliverables(deliverableName string, namespace string) []GetDeliverablesOutput {
	deliverables := []GetDeliverablesOutput{}
	cmd := "kubectl get deliverable"
	if deliverableName != "" {
		cmd += fmt.Sprintf(" %s", deliverableName)
	}
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += " -A"
	}
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return deliverables
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return deliverables
	}

	ss := linux_util.FieldIndices(temp[0])
	headers := linux_util.GetFields(temp[0], ss)
	for _, element := range temp[1:] {
		words := linux_util.GetFields(element, ss)
		var wl GetDeliverablesOutput
		for index, value := range words {
			reflect.ValueOf(&wl).Elem().FieldByName(headers[index]).SetString(value)
		}
		deliverables = append(deliverables, wl)
	}

	fmt.Printf("deliverables: %+v\n", deliverables)
	return deliverables
}

type GetDeliverablesYamlOutput struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		CreationTimestamp time.Time `yaml:"creationTimestamp"`
		Generation        int       `yaml:"generation"`
		Labels            struct {
			AppKubernetesIoComponent         string `yaml:"app.kubernetes.io/component"`
			AppKubernetesIoPartOf            string `yaml:"app.kubernetes.io/part-of"`
			AppTanzuVmwareComDeliverableType string `yaml:"app.tanzu.vmware.com/deliverable-type"`
			AppsKubernetesIoName             string `yaml:"apps.kubernetes.io/name"`
			AppsTanzuVmwareComWorkloadType   string `yaml:"apps.tanzu.vmware.com/workload-type"`
			CartoRunClusterTemplateName      string `yaml:"carto.run/cluster-template-name"`
			CartoRunResourceName             string `yaml:"carto.run/resource-name"`
			CartoRunSupplyChainName          string `yaml:"carto.run/supply-chain-name"`
			CartoRunTemplateKind             string `yaml:"carto.run/template-kind"`
			CartoRunWorkloadName             string `yaml:"carto.run/workload-name"`
			CartoRunWorkloadNamespace        string `yaml:"carto.run/workload-namespace"`
		} `yaml:"labels"`
		Name            string `yaml:"name"`
		Namespace       string `yaml:"namespace"`
		OwnerReferences []struct {
			APIVersion         string `yaml:"apiVersion"`
			BlockOwnerDeletion bool   `yaml:"blockOwnerDeletion"`
			Controller         bool   `yaml:"controller"`
			Kind               string `yaml:"kind"`
			Name               string `yaml:"name"`
			UID                string `yaml:"uid"`
		} `yaml:"ownerReferences"`
		ResourceVersion string `yaml:"resourceVersion"`
		UID             string `yaml:"uid"`
	} `yaml:"metadata"`
	Spec struct {
		Params []struct {
			Name  string `yaml:"name"`
			Value string `yaml:"value"`
		} `yaml:"params"`
		Source struct {
			Image string `yaml:"image,omitempty"`
			Git   struct {
				Ref struct {
					Branch string `yaml:"branch"`
				} `yaml:"ref"`
				URL string `yaml:"url"`
			} `yaml:"git,omitempty"`
		} `yaml:"source"`
	} `yaml:"spec"`
	Status struct {
		Conditions []struct {
			LastTransitionTime time.Time `yaml:"lastTransitionTime"`
			Message            string    `yaml:"message"`
			Reason             string    `yaml:"reason"`
			Status             string    `yaml:"status"`
			Type               string    `yaml:"type"`
		} `yaml:"conditions"`
		DeliveryRef struct {
		} `yaml:"deliveryRef"`
		ObservedGeneration int `yaml:"observedGeneration"`
	} `yaml:"status"`
}

type Status struct {
	Conditions []struct {
		LastTransitionTime time.Time `yaml:"lastTransitionTime"`
		Message            string    `yaml:"message"`
		Reason             string    `yaml:"reason"`
		Status             string    `yaml:"status"`
		Type               string    `yaml:"type"`
	} `yaml:"conditions"`
	DeliveryRef struct {
	} `yaml:"deliveryRef"`
	ObservedGeneration int `yaml:"observedGeneration"`
}

type OwnerReferences []struct {
	APIVersion         string `yaml:"apiVersion"`
	BlockOwnerDeletion bool   `yaml:"blockOwnerDeletion"`
	Controller         bool   `yaml:"controller"`
	Kind               string `yaml:"kind"`
	Name               string `yaml:"name"`
	UID                string `yaml:"uid"`
}

func GetDeliverablesYaml(name string, namespace string) *GetDeliverablesYamlOutput {
	var raw *GetDeliverablesYamlOutput
	cmd := fmt.Sprintf("kubectl get deliverable %s -n %s -o yaml", name, namespace)
	res1, err1 := linux_util.ExecuteCmd(cmd)
	if err1 != nil {
		return raw
	}
	in := []byte(res1)
	if err := yaml.Unmarshal(in, &raw); err != nil {
		panic(err)
	}
	return raw
}

func DeleteDeliverable(name string, namespace string) (string, error) {
	cmd := fmt.Sprintf("kubectl delete deliverable %s -n %s", name, namespace)
	res, err := linux_util.ExecuteCmd(cmd)
	return res, err
}
