// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
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
