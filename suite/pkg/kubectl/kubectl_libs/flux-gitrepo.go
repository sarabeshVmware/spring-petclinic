package kubectl_libs

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

type GetGitrepoOutput struct {
	NAME, URL, READY, STATUS, AGE string
}

func GetGitrepo(gitrepoName string, namespace string) []GetGitrepoOutput {
	gitRepos := []GetGitrepoOutput{}
	cmd := "kubectl get gitrepo"
	if gitrepoName != "" {
		cmd += fmt.Sprintf(" %s", gitrepoName)
	}
	if namespace != "" {
		cmd += fmt.Sprintf(" -n %s", namespace)
	} else {
		cmd += " -A"
	}
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		return gitRepos
	}

	temp := strings.Split(strings.TrimSuffix(response, "\n"), "\n")
	if len(temp) <= 1 {
		log.Printf("Output : %s", temp[0])
		return gitRepos
	}

	ss := linux_util.FieldIndices(temp[0])
	headers := linux_util.GetFields(temp[0], ss)
	for _, element := range temp[1:] {
		words := linux_util.GetFields(element, ss)
		var gitRepo GetGitrepoOutput
		for index, value := range words {
			reflect.ValueOf(&gitRepo).Elem().FieldByName(headers[index]).SetString(value)
		}
		gitRepos = append(gitRepos, gitRepo)
	}
	fmt.Printf("gitRepos: %+v\n", gitRepos)
	return gitRepos
}
