// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"log"
	"os/exec"
	"strings"
)

func RunCommand(command string, allowError bool, dontLogCommand bool) ([]byte, error) {
	commandName := strings.Split(command, " ")[0]
	arguments := strings.Split(command, " ")[1:]
	cmd := exec.Command(commandName, arguments...)
	stdoutStderr, err := cmd.CombinedOutput()
	if !dontLogCommand {
		log.Printf("Command executed: %s %s", commandName, strings.Join(arguments, " "))
	}
	log.Printf("Output: \n%s", string(stdoutStderr))
	if !allowError {
		CheckError(err)
	}
	return stdoutStderr, err
}

func Run(command string) ([]byte, error) {
	return RunCommand(command, false, false)
}

func Run_AllowError(command string) ([]byte, error) {
	return RunCommand(command, true, false)
}

func Run_DontLogCommand(command string) ([]byte, error) {
	return RunCommand(command, false, true)
}

func Run_DontLogCommand_AllowError(command string) ([]byte, error) {
	return RunCommand(command, true, true)
}
