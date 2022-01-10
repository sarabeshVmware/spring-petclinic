package tanzu_libs

import (
	"encoding/json"
	"fmt"
	"log"
	linux_util "pkg/utils/linux_util"
)

type GetInstalledPackagesOutput []struct {
	Conditions         string `json:"conditions"`
	Name               string `json:"name"`
	PackageName        string `json:"package-name"`
	PackageVersion     string `json:"package-version"`
	Status             string `json:"status"`
	UsefulErrorMessage string `json:"useful-error-message"`
}

func GetInstalledPackages(name string, namespace string) *GetInstalledPackagesOutput {

	cmd := fmt.Sprintf("tanzu package installed get %s", name)
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += " -A"
	}
	cmd += " -o json"

	res1, err1 := linux_util.ExecuteCmd(cmd)
	if err1 != nil {
		log.Println("something bad happened")
	}
	in := []byte(res1)
	var raw *GetInstalledPackagesOutput
	if err := json.Unmarshal(in, &raw); err != nil {
		panic(err)
	}
	return raw
}
