// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"log"
	"path/filepath"
	"runtime"
)

func CheckError(err error) bool {
	if err != nil {
		log.Fatal(err)
		return true
	}
	return false
}

func GetCurrentDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}

func GetProjectDirectory() string {
	currentDir := GetCurrentDir()
	return filepath.Dir(currentDir)
}

func GetValuesDirectory() string {
	projectDir := GetProjectDirectory()
	return filepath.Join(projectDir, "values")
}
