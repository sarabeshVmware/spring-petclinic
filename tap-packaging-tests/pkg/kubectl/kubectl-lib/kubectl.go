package kubectl_lib

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type GetPodintentOutput struct {
	NAME   string
	READY  string
	REASON string
	AGE    string
}

func GetPodintent(appName string, namespace string) {
	cmd := fmt.Sprintf("kubectl get podintent %s -n %s", appName, namespace)
	//  cmd := "kubectl get imagerepositories -A"
	res1, err1 := executeCmd(cmd)
	if err1 != nil {
		log.Println("something bad happened")
	}
	res1 = strings.TrimSuffix(res1, "\n")
	temp := strings.Split(res1, "\n")
	podIntents := []GetPodintentOutput{}
	for _, element := range temp[1:] {
		words := strings.Fields(element)
		podIntent := GetPodintentOutput{
			NAME:   words[0],
			READY:  words[1],
			REASON: "", // if ready is false then reason is populated
			AGE:    words[2],
		}
		podIntents = append(podIntents, podIntent)
	}
	fmt.Printf("%+v\n", podIntents)
	fmt.Printf("Name: %s, Ready: %s, Reason: %s, Age: %s ", podIntents[0].NAME, podIntents[0].READY, podIntents[0].REASON, podIntents[0].AGE)
}

type GetWorkloadOutput struct {
	NAME        string
	SOURCE      string
	SUPPLYCHAIN string
	READY       string
	REASON      string
}

func GetWorkload(workloadName string, namespace string) {
	if namespace == "" {
		namespace = "tap-install"
	}
	cmd := "kubectl get workload"
	if workloadName != "" {
		cmd += fmt.Sprintf(" %s", workloadName)
	}
	cmd += fmt.Sprintf(" -n %s", namespace)
	res1, err1 := executeCmd(cmd)
	if err1 != nil {
		log.Printf("something bad happened")
	}
	res1 = strings.TrimSuffix(res1, "\n")
	temp := strings.Split(res1, "\n")
	workloads := []GetWorkloadOutput{}
	for _, element := range temp[1:] {
		words := strings.Fields(element)
		wl := GetWorkloadOutput{
			NAME:        words[0],
			SOURCE:      words[1],
			SUPPLYCHAIN: words[2],
			READY:       words[3],
			REASON:      words[4],
		}
		workloads = append(workloads, wl)
	}
	fmt.Printf("%+v\n", workloads)
}

type GetImageRepositoriesOutput struct {
	NAME   string
	IMAGE  string
	URL    string
	READY  string
	REASON string
	AGE    string
}

func GetImageRepositories(name string, namespace string) []GetImageRepositoriesOutput {
	if namespace == "" {
		namespace = "tap-install"
	}
	cmd := "kubectl get imagerepositories"
	if name != "" {
		cmd += fmt.Sprintf(" %s", name)
	}
	cmd += fmt.Sprintf(" -n %s", namespace)
	res1, err1 := executeCmd(cmd)
	if err1 != nil {
		log.Printf("something bad happened")
	}
	res1 = strings.TrimSuffix(res1, "\n")
	temp := strings.Split(res1, "\n")
	imagerepos := []GetImageRepositoriesOutput{}
	for _, element := range temp[1:] {
		words := strings.Fields(element)
		wl := GetImageRepositoriesOutput{
			NAME:   words[0],
			IMAGE:  words[1],
			URL:    words[2],
			READY:  words[3],
			REASON: words[4],
			AGE:    words[5],
		}
		imagerepos = append(imagerepos, wl)
	}
	fmt.Printf("%+v\n", imagerepos)
	return imagerepos
}

func executeCmd1(command string) (string, error) {
	commandName := strings.Split(command, " ")[0]
	arguments := strings.Split(command, " ")[1:]
	cmd := exec.Command(commandName, arguments...)
	stdoutStderr, err := cmd.CombinedOutput()
	log.Printf("Command executed: %s", command)
	if err != nil {
		fmt.Printf(fmt.Sprint(err) + ": " + string(stdoutStderr))
	} else {
		log.Printf("Output: \n%s", string(stdoutStderr))
	}
	return string(stdoutStderr), err
}

func main() {
	GetPodintent("tanzu-java-web-app", "tap-install")
	GetWorkload("", "")
}
