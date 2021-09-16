// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"log"
	"os/exec"
	"strings"
)

type Command struct {
	CommandName    string
	Arguments      []string
	DontLogCommand bool // default false
	AllowError     bool // default false
}

func RunCommand(command Command) ([]byte, error) {
	cmd := exec.Command(command.CommandName, command.Arguments...)
	stdoutStderr, err := cmd.CombinedOutput()
	if !command.DontLogCommand {
		log.Printf("Command executed: %s %s", command.CommandName, strings.Join(command.Arguments, " "))
	}
	log.Printf("Output: \n%s", string(stdoutStderr))
	if !command.AllowError {
		CheckError(err)
	}
	return stdoutStderr, err
}
