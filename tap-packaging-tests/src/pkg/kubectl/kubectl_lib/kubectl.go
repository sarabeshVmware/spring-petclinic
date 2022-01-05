package kubectl_lib

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"unicode"
)

type span struct {
	start int
	end   int
}

func FieldIndices(s string) []span {
	f := unicode.IsSpace
	spans := make([]span, 0, 32)
	start := -1 // valid span start if >= 0
	for end, rune := range s {
		if f(rune) {
			if start >= 0 {
				spans = append(spans, span{start, end})
				start = ^start
			}
		} else {
			if start < 0 {
				start = end
			}
		}
	}
	// Last field might end at EOF.
	if start >= 0 {
		spans = append(spans, span{start, len(s)})
	}
	for index := range spans {
		if index == 0 {
			continue
		}
		spans[index-1].end = spans[index].start
	}
	return spans
}

func GetFields(s string, spans []span) []string {
	// Create strings from field indices.
	if len(s) < spans[len(spans)-1].end { // if last few column values are empty - padding string with spaces to the right
		b := fmt.Sprintf("%s%d%s", "%-", spans[len(spans)-1].end, "v")
		s = fmt.Sprintf(b, s)
	}
	if len(s) > spans[len(spans)-1].end { // if column values exceed column header length
		spans[len(spans)-1].end = len(s)
	}
	a := make([]string, len(spans))
	for i, span := range spans {
		a[i] = strings.TrimSpace(s[span.start:span.end])
	}
	return a
}

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

	ss := FieldIndices(temp[0])
	headers := GetFields(temp[0], ss)
	for _, element := range temp[1:] {
		words := GetFields(element, ss)
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
	response, err := executeCmd(cmd)
	if err != nil {
		return workloads
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return workloads
	}

	ss := FieldIndices(temp[0])
	headers := GetFields(temp[0], ss)
	for _, element := range temp[1:] {
		words := GetFields(element, ss)
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

	ss := FieldIndices(temp[0])
	headers := GetFields(temp[0], ss)
	for _, element := range temp[1:] {
		words := GetFields(element, ss)
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
	response, err := executeCmd(cmd)
	if err != nil {
		return builds
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return builds
	}

	ss := FieldIndices(temp[0])
	headers := GetFields(temp[0], ss)
	for _, element := range temp[1:] {
		words := GetFields(element, ss)
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

	ss := FieldIndices(temp[0])
	headers, words := GetFields(temp[0], ss), GetFields(temp[1], ss)

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
	response, err := executeCmd(cmd)
	if err != nil {
		return ksvc
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return ksvc
	}

	ss := FieldIndices(temp[0])
	headers, words := GetFields(temp[0], ss), GetFields(temp[1], ss)

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
	response, err := executeCmd(cmd)
	if err != nil {
		return sourceScan
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return sourceScan
	}

	ss := FieldIndices(temp[0])
	headers, words := GetFields(temp[0], ss), GetFields(temp[1], ss)

	for index, value := range words {
		reflect.ValueOf(&sourceScan).Elem().FieldByName(headers[index]).SetString(value)
	}
	fmt.Printf("sourceScan: %+v\n", sourceScan)
	return sourceScan
}
