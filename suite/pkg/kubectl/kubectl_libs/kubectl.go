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
	cmd := "kubectl get builds.kpack.io"
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
	APIVersion string `json:"apiVersion"`
	Data       struct {
		CaCrt     string `json:"ca.crt"`
		Namespace string `json:"namespace"`
		Token     string `json:"token"`
	} `json:"data"`
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

func RestartScanLinkController() (string, error) {
	cmd := "kubectl rollout restart deployment.apps/scan-link-controller-manager -n scan-link-system"
	res, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while restarting scan controller")
	}
	return res, err
}

func GetCurrentContext() string {
	clusterName := ""
	cmd := "kubectl config current-context"
	clusterName, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while fetching config current context")
	}
	return clusterName
}

func UseContext(context string) (string, error) {
	cmd := fmt.Sprintf("kubectl config use-context %s", context)
	res, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while switching context")
	}
	return res, err
}

type CurrentConfigViewOutput struct {
	APIVersion string `yaml:"apiVersion"`
	Clusters   []struct {
		Cluster struct {
			CertificateAuthorityData string `yaml:"certificate-authority-data"`
			Server                   string `yaml:"server"`
		} `yaml:"cluster"`
		Name string `yaml:"name"`
	} `yaml:"clusters"`
	Contexts []struct {
		Context struct {
			Cluster string `yaml:"cluster"`
			User    string `yaml:"user"`
		} `yaml:"context"`
		Name string `yaml:"name"`
	} `yaml:"contexts"`
	CurrentContext string `yaml:"current-context"`
	Kind           string `yaml:"kind"`
	Preferences    struct {
	} `yaml:"preferences"`
	Users []struct {
		Name string `yaml:"name"`
		User struct {
			ClientCertificateData string `yaml:"client-certificate-data"`
			ClientKeyData         string `yaml:"client-key-data"`
			Token                 string `yaml:"token"`
		} `yaml:"user"`
	} `yaml:"users"`
}

func GetCurrentConfigView() *CurrentConfigViewOutput {
	var raw *CurrentConfigViewOutput
	cmd := fmt.Sprintf("kubectl config view --minify")
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
func KubectlApplyConfiguration(file string, namespace string) error {
	log.Printf("applying configuration %s in namespace %s", file, namespace)

	// execute cmd
	cmd := fmt.Sprintf("kubectl apply -f %s", file)
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	}
	output, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while applying configuration %s in namespace %s", file, namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("configuration %s applied in namespace %s", file, namespace)
		log.Printf("output: %s", output)
	}

	return err
}

type GetDeploymentsOutput struct {
	NAMESPACE, NAME, READY, UP_TO_DATE, AVAILABLE, AGE string
}

func GetDeployments(deploymentName string, namespace string) []GetDeploymentsOutput {
	deployments := []GetDeploymentsOutput{}
	cmd := "kubectl get deployments"
	if deploymentName != "" {
		cmd += fmt.Sprintf(" %s", deploymentName)
	}
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += " -A"
	}
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return deployments
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")

	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return deployments
	}

	ss := linux_util.FieldIndices(temp[0])
	headers := linux_util.GetFields(temp[0], ss)
	for index, ele := range headers {
		headers[index] = strings.ReplaceAll(ele, "-", "_")
	}
	for _, element := range temp[1:] {
		words := linux_util.GetFields(element, ss)
		var deployment GetDeploymentsOutput
		for index, value := range words {
			reflect.ValueOf(&deployment).Elem().FieldByName(headers[index]).SetString(value)
		}
		deployments = append(deployments, deployment)
	}

	fmt.Printf("deployments: %+v\n", deployments)
	return deployments
}
