package kubectl_libs

import (
	"fmt"
	"log"
	linux_util "pkg/utils/linux_util"
	"reflect"
	"strings"
)

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
