package kubectl_lib

import (
	"fmt"
	"log"
	"reflect"
	"strings"
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
	response, err := executeCmd(cmd)
	if err != nil {
		return podIntents
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return podIntents
	}

	headers := strings.Fields(temp[0])
	for _, element := range temp[1:] {
		words := strings.Fields(element)
		var podIntent GetPodintentOutput
		index_inc := false
		for index, value := range words {
			if len(words) < len(headers) && headers[index] == "REASON" {
				index_inc = true
			}
			if index_inc {
				reflect.ValueOf(&podIntent).Elem().FieldByName(headers[index+1]).SetString(value)
			} else {
				reflect.ValueOf(&podIntent).Elem().FieldByName(headers[index]).SetString(value)
			}

		}
		podIntents = append(podIntents, podIntent)
	}
	fmt.Printf("podIntents: %+v\n", podIntents)
	return podIntents
}

type GetWorkloadOutput struct {
	NAME, SOURCE, SUPPLYCHAIN, READY, REASON string
}

func GetWorkload(workloadName string, namespace string) []GetWorkloadOutput {
	workloads := []GetWorkloadOutput{}
	cmd := "kubectl get workload"
	if workloadName != "" {
		cmd += fmt.Sprintf(" %s", workloadName)
	}
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += " -A"
	}
	response, err := executeCmd(cmd)
	if err != nil {
		return workloads
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return workloads
	}

	headers := strings.Fields(temp[0])
	for _, element := range temp[1:] {
		words := strings.Fields(element)
		var wl GetWorkloadOutput
		for index, value := range words {
			reflect.ValueOf(&wl).Elem().FieldByName(headers[index]).SetString(value)
		}
		workloads = append(workloads, wl)
	}
	fmt.Printf("workloads: %+v\n", workloads)
	return workloads
}

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
	response, err := executeCmd(cmd)
	if err != nil {
		return imagerepos
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return imagerepos
	}

	headers := strings.Fields(temp[0])
	for _, element := range temp[1:] {
		words := strings.Fields(element)
		var imagerepo GetImageRepositoriesOutput
		index_inc := false
		for index, value := range words {
			if len(words) < len(headers) && headers[index] == "REASON" {
				index_inc = true
			}
			if index_inc {
				reflect.ValueOf(&imagerepo).Elem().FieldByName(headers[index+1]).SetString(value)
			} else {
				reflect.ValueOf(&imagerepo).Elem().FieldByName(headers[index]).SetString(value)
			}
		}
		imagerepos = append(imagerepos, imagerepo)
	}
	fmt.Printf("imagerepos: %+v\n", imagerepos)
	return imagerepos
}

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
	response, err := executeCmd(cmd)
	if err != nil {
		return builds
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return builds
	}

	headers := strings.Fields(temp[0])
	for _, element := range temp[1:] {
		words := strings.Fields(element)
		var build GetBuildsOutput
		for index, value := range words {
			reflect.ValueOf(&build).Elem().FieldByName(headers[index]).SetString(value)
		}
		builds = append(builds, build)
	}
	fmt.Printf("builds: %+v\n", builds)
	return builds
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
	response, err := executeCmd(cmd)
	if err != nil {
		return latestImage
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return latestImage
	}

	headers, words := strings.Fields(temp[0]), strings.Fields(temp[1])

	for index, value := range words {
		reflect.ValueOf(&latestImage).Elem().FieldByName(headers[index]).SetString(value)
	}
	fmt.Printf("latestImage: %+v\n", latestImage)
	return latestImage
}
