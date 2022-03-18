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

type GetBuildsOutput struct {
	NAME, IMAGE, SUCCEEDED string
}

func GetBuilds(buildName string, namespace string) []GetBuildsOutput {
	builds := []GetBuildsOutput{}
	cmd := "kubectl get builds"
	if buildName != "" {
		cmd += fmt.Sprintf(" %s", buildName)
	}
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += " -A"
	}
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return builds
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return builds
	}

	ss := linux_util.FieldIndices(temp[0])
	headers := linux_util.GetFields(temp[0], ss)
	for _, element := range temp[1:] {
		words := linux_util.GetFields(element, ss)
		var build GetBuildsOutput
		for index, value := range words {
			reflect.ValueOf(&build).Elem().FieldByName(headers[index]).SetString(value)
		}
		builds = append(builds, build)
	}

	fmt.Printf("builds: %+v\n", builds)
	return builds
}

type GetImagesOutput struct {
	NAME, LATESTIMAGE, READY string
}

func GetImages(imageName string, namespace string) []GetImagesOutput {
	images := []GetImagesOutput{}
	cmd := "kubectl get images.kpac"
	if imageName != "" {
		cmd += fmt.Sprintf(" %s", imageName)
	}
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += " -A"
	}
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return images
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return images
	}

	ss := linux_util.FieldIndices(temp[0])
	headers := linux_util.GetFields(temp[0], ss)

	for _, element := range temp[1:] {
		words := linux_util.GetFields(element, ss)
		var image GetImagesOutput
		for index, value := range words {
			reflect.ValueOf(&image).Elem().FieldByName(headers[index]).SetString(value)
		}
		images = append(images, image)
	}

	fmt.Printf("images: %+v\n", images)
	return images
}

type GetLatestImageOutput struct {
	NAME, LATESTIMAGE, READY string
}

func GetLatestImage(namespace string) GetLatestImageOutput {
	var latestImage GetLatestImageOutput
	cmd := "kubectl get images.kpac"
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += " -A"
	}
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return latestImage
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return latestImage
	}

	ss := linux_util.FieldIndices(temp[0])
	headers, words := linux_util.GetFields(temp[0], ss), linux_util.GetFields(temp[1], ss)

	for index, value := range words {
		reflect.ValueOf(&latestImage).Elem().FieldByName(headers[index]).SetString(value)
	}
	fmt.Printf("latestImage: %+v\n", latestImage)
	return latestImage
}

type GetKsvcOutput struct {
	NAME, URL, LATESTCREATED, LATESTREADY, READY, REASON string
}

func GetKsvc(name string, namespace string) []GetKsvcOutput {
	ksvcs := []GetKsvcOutput{}
	cmd := "kubectl get ksvc"
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
		return ksvcs
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return ksvcs
	}

	ss := linux_util.FieldIndices(temp[0])
	headers := linux_util.GetFields(temp[0], ss)
	for _, element := range temp[1:] {
		words := linux_util.GetFields(element, ss)
		var ksvc GetKsvcOutput
		for index, value := range words {
			reflect.ValueOf(&ksvc).Elem().FieldByName(headers[index]).SetString(value)
		}
		ksvcs = append(ksvcs, ksvc)
	}

	fmt.Printf("ksvc: %+v\n", ksvcs)
	return ksvcs
}

type GetIngressOutput struct {
	NAME, CLASS, HOSTS, ADDRESS, PORTS, AGE string
}

func GetIngress(name string, namespace string) GetIngressOutput {
	var ingress GetIngressOutput
	cmd := "kubectl get ingress"
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
		return ingress
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return ingress
	}

	ss := linux_util.FieldIndices(temp[0])
	headers, words := linux_util.GetFields(temp[0], ss), linux_util.GetFields(temp[1], ss)

	for index, value := range words {
		reflect.ValueOf(&ingress).Elem().FieldByName(headers[index]).SetString(value)
	}
	fmt.Printf("Ingress: %+v\n", ingress)
	return ingress
}

type GetServiceAccountJsonOutput struct {
	APIVersion       string `json:"apiVersion"`
	ImagePullSecrets []struct {
		Name string `json:"name"`
	} `json:"imagePullSecrets"`
	Kind     string `json:"kind"`
	Metadata struct {
		Annotations struct {
			KubectlKubernetesIoLastAppliedConfiguration string `json:"kubectl.kubernetes.io/last-applied-configuration"`
		} `json:"annotations"`
		CreationTimestamp string `json:"creationTimestamp"`
		Name              string `json:"name"`
		Namespace         string `json:"namespace"`
		ResourceVersion   string `json:"resourceVersion"`
		UID               string `json:"uid"`
	} `json:"metadata"`
	Secrets []struct {
		Name string `json:"name"`
	} `json:"secrets"`
}

func GetServiceAccountJson(name string, namespace string) *GetServiceAccountJsonOutput {
	var raw *GetServiceAccountJsonOutput
	cmd := fmt.Sprintf("kubectl get sa %s -n %s -o json", name, namespace)
	res, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return raw
	}
	in := []byte(res)
	if err := json.Unmarshal(in, &raw); err != nil {
		panic(err)
	}
	return raw
}

func PatchServiceAccount(sa string, namespace string, patch string) bool {
	log.Printf("patching sa")
	cmd := fmt.Sprintf("kubectl patch serviceaccount %s -n %s -p %s", sa, namespace, patch)
	output, err := linux_util.ExecuteCmdInBashMode(cmd)
	if err != nil {
		log.Printf("error while applying patch %s for sa %s in namespace %s", patch, sa, namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
		return false
	} else {
		log.Printf("applied patch %s for sa %s in namespace %s", patch, sa, namespace)
		log.Printf("output: %s", output)
	}
	return true
}

type GetRevisionsOutput struct {
	NAME, CONFIG_NAME, K8S_SERVICE_NAME, GENERATION, READY, REASON, ACTUAL_REPLICAS, DESIRED_REPLICAS string
}

func GetRevisions(name string, namespace string) []GetRevisionsOutput {
	revisions := []GetRevisionsOutput{}
	cmd := "kubectl get revision"
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
		return revisions
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return revisions
	}

	ss := linux_util.FieldIndicesWithSingleSpace(temp[0])
	headers := linux_util.GetFields(temp[0], ss)
	for index, ele := range headers {
		headers[index] = strings.ReplaceAll(ele, " ", "_")
	}

	for _, element := range temp[1:] {
		words := linux_util.GetFields(element, ss)
		var revision GetRevisionsOutput
		for index, value := range words {
			reflect.ValueOf(&revision).Elem().FieldByName(headers[index]).SetString(value)
		}
		revisions = append(revisions, revision)
	}

	fmt.Printf("revisions: %+v\n", revisions)
	return revisions
}

type GetSecretsOutput struct {
	APIVersion string `json:"apiVersion"`
	Data       struct {
		CaCrt     string `json:"ca.crt"`
		Namespace string `json:"namespace"`
		Token     string `json:"token"`
	} `json:"data"`
	Kind     string `json:"kind"`
	Metadata struct {
		Annotations struct {
			KubernetesIoServiceAccountName string `json:"kubernetes.io/service-account.name"`
			KubernetesIoServiceAccountUID  string `json:"kubernetes.io/service-account.uid"`
		} `json:"annotations"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
		Name              string    `json:"name"`
		Namespace         string    `json:"namespace"`
		ResourceVersion   string    `json:"resourceVersion"`
		UID               string    `json:"uid"`
	} `json:"metadata"`
	Type string `json:"type"`
}

func GetSecrets(name string, namespace string) *GetSecretsOutput {
	var raw *GetSecretsOutput
	cmd := fmt.Sprintf("kubectl get secrets %s -n %s -o json", name, namespace)
	res, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return raw
	}
	in := []byte(res)
	if err := json.Unmarshal(in, &raw); err != nil {
		panic(err)
	}
	return raw
}
