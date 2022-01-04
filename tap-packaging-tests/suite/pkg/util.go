// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func RunCommand(command string) (string, error) {
	commandName := strings.Split(command, " ")[0]
	arguments := strings.Split(command, " ")[1:]
	cmd := exec.Command(commandName, arguments...)
	stdoutStderr, err := cmd.CombinedOutput()

	return string(stdoutStderr), err
}

func GetFileDir() string {
	_, filename, _, _ := runtime.Caller(1)

	return filepath.Dir(filename)
}
