// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package suite

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"gopkg.in/yaml.v3"
)

func GetFileDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}

func WriteYAMLFile(file string, input interface{}) error {
	bytes, err := yaml.Marshal(input)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, bytes, 0677)
	return err
}

func SetLogger(directory string) (string, error) {
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		return "", err
	}
	logFilePath := filepath.Join("logs", fmt.Sprintf("log_%s.log", time.Now().Format(time.RFC3339Nano)))
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return logFilePath, err
	}
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetFlags(log.LstdFlags | log.Llongfile)
	return logFilePath, err
}
