// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"strings"
)

func GitClone(path string, repo string) (string, string, error) {
	cmd := fmt.Sprintf("cd %s; git clone %s", path, repo)
	output, err := RunCommandInBashMode(cmd)
	return cmd, output, err
}

func GitAdd(path string, files []string) (string, string, error) {
	cmd := fmt.Sprintf("cd %s; git add", path)
	for _, file := range files {
		cmd += fmt.Sprintf(" %s", file)
	}
	output, err := RunCommandInBashMode(cmd)
	return cmd, output, err
}

func GitCommit(path string, message string) (string, string, error) {
	cmd := fmt.Sprintf(`cd %s; git commit -m "%s"`, path, message)
	output, err := RunCommandInBashMode(cmd)
	return cmd, output, err
}

func GitPush(path string, force bool) (string, string, error) {
	cmd := fmt.Sprintf("cd %s; git push", path)
	if force {
		cmd += " --force"
	}
	output, err := RunCommandInBashMode(cmd)
	return cmd, output, err
}

func GitResetFromHead(path string, count int) (string, string, error) {
	cmd := fmt.Sprintf("cd %s; git reset --hard HEAD~%d", path, count)
	output, err := RunCommandInBashMode(cmd)
	return cmd, output, err
}

func GitConfig(path string, username string, email string) (string, string, error) {
	cmd := fmt.Sprintf("cd %s; git config --global user.name %s; git config --global user.email %s", path, username, email)
	output, err := RunCommandInBashMode(cmd)
	return cmd, output, err
}

func GitSetUrl(path string, access_token string, repository string) (string, string, error) {
	url := strings.Split(repository, "//")[1]
	cmd := fmt.Sprintf("cd %s; git remote set-url origin https://%s@%s", path, access_token, url)
	output, err := RunCommandInBashMode(cmd)
	return cmd, output, err
}
