// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import "fmt"

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
