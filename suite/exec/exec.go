// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"os"
	"os/exec"
	"strings"
)

func RunCommand(command string) (string, error) {
	commandName := strings.Split(command, " ")[0]
	arguments := strings.Split(command, " ")[1:]
	cmd := exec.Command(commandName, arguments...)
	stdoutStderr, err := cmd.CombinedOutput()
	return string(stdoutStderr), err
}

func RunCommandInBashMode(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	stdoutStderr, err := cmd.CombinedOutput()
	return string(stdoutStderr), err
}

func RunCommandWithOutWait(command string) (*os.Process,error) {
	var proc *os.Process
	commandName := strings.Split(command, " ")[0]
	arguments := strings.Split(command, " ")[1:]
	cmd := exec.Command(commandName, arguments...)
	err := cmd.Start()
	proc = cmd.Process
	return proc, err
}

