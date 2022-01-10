package kubectl_lib

import (
	"fmt"
	"log"
	linux_util "pkg/utils/linux_util"
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
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return workloads
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return workloads
	}

	ss := linux_util.FieldIndices(temp[0])
	headers := linux_util.GetFields(temp[0], ss)
	for _, element := range temp[1:] {
		words := linux_util.GetFields(element, ss)
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

func GetKsvc(name string, namespace string) GetKsvcOutput {
	var ksvc GetKsvcOutput
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
		return ksvc
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return ksvc
	}

	ss := linux_util.FieldIndices(temp[0])
	headers, words := linux_util.GetFields(temp[0], ss), linux_util.GetFields(temp[1], ss)

	for index, value := range words {
		reflect.ValueOf(&ksvc).Elem().FieldByName(headers[index]).SetString(value)
	}
	fmt.Printf("ksvc: %+v\n", ksvc)
	return ksvc
}

type GetSourceScanOutput struct {
	NAME, PHASE, SCANNEDREVISION, SCANNEDREPOSITORY, AGE, CRITICAL, HIGH, MEDIUM, LOW, UNKNOWN, CVETOTAL string
}

func GetSourceScan(name string, namespace string) GetSourceScanOutput {
	var sourceScan GetSourceScanOutput
	cmd := "kubectl get sourcescan"
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
		return sourceScan
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return sourceScan
	}

	ss := linux_util.FieldIndices(temp[0])
	headers, words := linux_util.GetFields(temp[0], ss), linux_util.GetFields(temp[1], ss)

	for index, value := range words {
		reflect.ValueOf(&sourceScan).Elem().FieldByName(headers[index]).SetString(value)
	}
	fmt.Printf("sourceScan: %+v\n", sourceScan)
	return sourceScan
}

type GetPipelineOutput struct {
	NAME, AGE string
}

func GetPipeline(name string, namespace string) GetPipelineOutput {
	var pipeline GetPipelineOutput
	cmd := "kubectl get pipeline"
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
		return pipeline
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return pipeline
	}

	ss := linux_util.FieldIndices(temp[0])
	headers, words := linux_util.GetFields(temp[0], ss), linux_util.GetFields(temp[1], ss)

	for index, value := range words {
		reflect.ValueOf(&pipeline).Elem().FieldByName(headers[index]).SetString(value)
	}
	fmt.Printf("pipeline: %+v\n", pipeline)
	return pipeline
}

type GetPipelineRunsOutput struct {
	NAME, SUCCEEDED, REASON, STARTTIME, COMPLETIONTIME string
}

func GetPipelineRuns(name string, namespace string) GetPipelineRunsOutput {
	var pipeline GetPipelineRunsOutput
	cmd := "kubectl get prs"
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
		return pipeline
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return pipeline
	}

	ss := linux_util.FieldIndices(temp[0])
	headers, words := linux_util.GetFields(temp[0], ss), linux_util.GetFields(temp[1], ss)

	for index, value := range words {
		reflect.ValueOf(&pipeline).Elem().FieldByName(headers[index]).SetString(value)
	}
	fmt.Printf("pipeline: %+v\n", pipeline)
	return pipeline
}

type GetImageScanOutput struct {
	NAMESPACE, NAME, PHASE, SCANNEDIMAGE, AGE, CRITICAL, HIGH, MEDIUM, LOW, UNKNOWN, CVETOTAL string
}

func GetImageScan(name string, namespace string) GetImageScanOutput {
	var imgScan GetImageScanOutput
	cmd := "kubectl get imagescan"
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
		return imgScan
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return imgScan
	}

	ss := linux_util.FieldIndices(temp[0])
	headers, words := linux_util.GetFields(temp[0], ss), linux_util.GetFields(temp[1], ss)

	for index, value := range words {
		reflect.ValueOf(&imgScan).Elem().FieldByName(headers[index]).SetString(value)
	}
	fmt.Printf("imgScan: %+v\n", imgScan)
	return imgScan
}

type GetServiceBindingsOutput struct {
	NAME, READY, REASON, AGE string
}

func GetServiceBindings(name string, namespace string) GetServiceBindingsOutput {
	var svcBindings GetServiceBindingsOutput
	cmd := "kubectl get servicebindings"
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
		return svcBindings
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return svcBindings
	}

	ss := linux_util.FieldIndices(temp[0])
	headers, words := linux_util.GetFields(temp[0], ss), linux_util.GetFields(temp[1], ss)

	for index, value := range words {
		reflect.ValueOf(&svcBindings).Elem().FieldByName(headers[index]).SetString(value)
	}
	fmt.Printf("svcBindings: %+v\n", svcBindings)
	return svcBindings
}

type GetTrainingPortalsOutput struct {
	NAME, URL, ADMINUSERNAME, ADMINPASSWORD, STATUS string
}

func GetTrainingPortals(name string, namespace string) GetTrainingPortalsOutput {
	var tps GetTrainingPortalsOutput
	cmd := "kubectl get trainingportals"
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
		return tps
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return tps
	}

	ss := linux_util.FieldIndices(temp[0])
	headers, words := linux_util.GetFields(temp[0], ss), linux_util.GetFields(temp[1], ss)

	for index, value := range words {
		reflect.ValueOf(&tps).Elem().FieldByName(headers[index]).SetString(value)
	}
	fmt.Printf("trainingPortals: %+v\n", tps)
	return tps
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
