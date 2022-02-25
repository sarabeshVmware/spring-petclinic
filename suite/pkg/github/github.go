// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package github

import (
	"fmt"
	"log"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

func CreateGithubRepo(repo_name string, template_name string, token string) error {
	log.Printf("creating repo %s", repo_name)

	// execute cmd
	cmd := fmt.Sprintf("echo %s > token; gh auth login -h Github.com --with-token < token; rm token; gh repo create %s --public --template %s", token, repo_name, template_name)
	output, err := linux_util.ExecuteCmdInBashMode(cmd)
	if err != nil {
		log.Printf("error while creating repo %s ", repo_name)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("repo %s created ", repo_name)
		log.Printf("output: %s", output)
	}
	return err
}

func DeleteGithubRepo(repo_name string, token string) error {
	log.Printf("deleting repo %s ", repo_name)

	// execute cmd
	cmd := fmt.Sprintf("echo %s > token; gh auth login -h Github.com --with-token < token; rm token; gh repo delete %s --confirm", token, repo_name)
	output, err := linux_util.ExecuteCmdInBashMode(cmd)
	if err != nil {
		log.Printf("error while deleting repo %s ", repo_name)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("repo %s deleted ", repo_name)
		log.Printf("output: %s", output)
	}
	return err
}
