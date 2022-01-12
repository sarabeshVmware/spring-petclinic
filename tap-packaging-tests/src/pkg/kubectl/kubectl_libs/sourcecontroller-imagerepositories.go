package kubectl_libs

import (
	"encoding/json"
	"fmt"
	"log"
	linux_util "pkg/utils/linux_util"
	"reflect"
	"strings"
	"time"
)

type GetImageRepositoriesOutput struct {
	NAME, IMAGE, URL, READY, REASON, AGE string
}

func GetImageRepositories(name string, namespace string) []GetImageRepositoriesOutput {
	imagerepos := []GetImageRepositoriesOutput{}
	cmd := "kubectl get imagerepositories"
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
		return imagerepos
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return imagerepos
	}

	ss := linux_util.FieldIndices(temp[0])
	headers := linux_util.GetFields(temp[0], ss)
	for _, element := range temp[1:] {
		words := linux_util.GetFields(element, ss)
		var imagerepo GetImageRepositoriesOutput
		for index, value := range words {
			reflect.ValueOf(&imagerepo).Elem().FieldByName(headers[index]).SetString(value)
		}
		imagerepos = append(imagerepos, imagerepo)
	}

	fmt.Printf("imagerepos: %+v\n", imagerepos)
	return imagerepos
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
	return raw
}
