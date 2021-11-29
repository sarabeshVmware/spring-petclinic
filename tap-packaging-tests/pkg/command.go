// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"fmt"
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

func RunCommandWithBash(command string, allowError bool, dontLogCommand bool) ([]byte, error) {
	cmd := exec.Command("bash", "-c", command)
	stdoutStderr, err := cmd.CombinedOutput()
	if !dontLogCommand {
		log.Printf("Command executed: %s", command)
	}
	log.Printf("Output: \n%s", string(stdoutStderr))
	if !allowError {
		CheckError(err)
	}
	return stdoutStderr, err
}

func RunWithBash(command string) ([]byte, error) {
	return RunCommandWithBash(command, false, false)
}

func RunWithBash_AllowError(command string) ([]byte, error) {
	return RunCommandWithBash(command, true, false)
}

func RunWithBash_DontLogCommand(command string) ([]byte, error) {
	return RunCommandWithBash(command, false, true)
}

func RunWithBash_DontLogCommand_AllowError(command string) ([]byte, error) {
	return RunCommandWithBash(command, true, true)
}

func RunAndDisown(command string) (int, *exec.Cmd) {
	commandName := strings.Split(command, " ")[0]
	arguments := strings.Split(command, " ")[1:]
	cmd := exec.Command(commandName, arguments...)
	err := cmd.Start()
	pid := cmd.Process.Pid
	log.Printf("Command executed: %s %s", commandName, strings.Join(arguments, " "))
	CheckError(err)
	log.Printf("Command started, PID: %d", pid)
	return pid, cmd
}

func KillPID(pid int) {
	log.Printf("Killing PID: %d", pid)
	Run(fmt.Sprintf("kill %d", pid))
}
