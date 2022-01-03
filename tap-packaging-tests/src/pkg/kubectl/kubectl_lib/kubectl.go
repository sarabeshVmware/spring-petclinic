package kubectl_lib

import (
	"fmt"
	"log"
	"strings"
)

type GetPodintentOutput struct {
	NAME   string
	READY  string
	REASON string
	AGE    string
}

func GetPodintent(name string, namespace string) []GetPodintentOutput {
	cmd := "kubectl get podintent"
	if name != "" {
		cmd += fmt.Sprintf(" %s", name)
	}
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += fmt.Sprintf(" -A")
	}
	res1, err1 := executeCmd(cmd)
	if err1 != nil {
		log.Println("something bad happened")
	}
	res1 = strings.TrimSuffix(res1, "\n")
	temp := strings.Split(res1, "\n")
	podIntents := []GetPodintentOutput{}
	for _, element := range temp[1:] {
		words := strings.Fields(element)
		var podIntent GetPodintentOutput
		if len(words) == 4 {
			podIntent = GetPodintentOutput{
				NAME:   words[0],
				READY:  words[1],
				REASON: words[2],
				AGE:    words[3],
			}
		} else {
			podIntent = GetPodintentOutput{
				NAME:  words[0],
				READY: words[1],
				AGE:   words[2],
			}
		}
		podIntents = append(podIntents, podIntent)
	}
	fmt.Printf("%+v\n", podIntents)
	fmt.Printf("Name: %s, Ready: %s, Reason: %s, Age: %s ", podIntents[0].NAME, podIntents[0].READY, podIntents[0].REASON, podIntents[0].AGE)
	return podIntents
}

type GetWorkloadOutput struct {
	NAME        string
	SOURCE      string
	SUPPLYCHAIN string
	READY       string
	REASON      string
}

func GetWorkload(workloadName string, namespace string) []GetWorkloadOutput {
	cmd := "kubectl get workload"
	if workloadName != "" {
		cmd += fmt.Sprintf(" %s", workloadName)
	}
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += fmt.Sprintf(" -A")
	}
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
	return workloads
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
	cmd := "kubectl get imagerepositories"
	if name != "" {
		cmd += fmt.Sprintf(" %s", name)
	}
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += fmt.Sprintf(" -A")
	}
	res1, err1 := executeCmd(cmd)
	if err1 != nil {
		log.Printf("something bad happened")
	}
	res1 = strings.TrimSuffix(res1, "\n")
	temp := strings.Split(res1, "\n")
	temp1 := strings.Fields(temp[0])
	println("fields: ", len(temp1))
	for _, element := range temp1 {
		println(element, "index: ", strings.Index(temp[0], element))
	}
	imagerepos := []GetImageRepositoriesOutput{}
	for _, element := range temp[1:] {
		words := strings.Fields(element)
		var wl GetImageRepositoriesOutput
		if len(words) == 6 {
			wl = GetImageRepositoriesOutput{
				NAME:   words[0],
				IMAGE:  words[1],
				URL:    words[2],
				READY:  words[3],
				REASON: words[4],
				AGE:    words[5],
			}
		} else {
			wl = GetImageRepositoriesOutput{
				NAME:  words[0],
				IMAGE: words[1],
				URL:   words[2],
				READY: words[3],
				AGE:   words[4],
			}

		}
		imagerepos = append(imagerepos, wl)
	}
	fmt.Printf("%+v\n", imagerepos)
	return imagerepos
}
