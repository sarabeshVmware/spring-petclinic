// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func GetFileDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}

func WriteYAMLFile(file string, input interface{}) error {
	log.Printf("creating YAML file %s from input", file)

	// marshal input
	bytes, err := yaml.Marshal(input)
	if err != nil {
		log.Printf("error while marshalling input %s", input)
		log.Printf("error: %s", err)
		return err
	} else {
		log.Printf("marshalled input %s", input)
	}

	// write to file
	err = os.WriteFile(file, bytes, 0677)
	if err != nil {
		log.Printf("error while writing to file %s", file)
		log.Printf("error: %s", err)
	} else {
		log.Printf("file %s written", file)
	}

	return err
}

func ReplaceStringInFile(file string, originalString string, newString string) error {
	log.Printf(`replacing string "%s"->"%s" in file %s`, originalString, newString, file)

	// read file
	inputBytes, err := os.ReadFile(file)
	if err != nil {
		log.Printf("error while reading file %s", file)
		log.Printf("error: %s", err)
		return err
	} else {
		log.Printf("read file %s", file)
	}

	input := string(inputBytes)

	// replace
	input = strings.ReplaceAll(input, originalString, newString)

	// write file
	err = os.WriteFile(file, []byte(input), 0677)
	if err != nil {
		log.Printf("error while writing file %s", file)
		log.Printf("error: %s", err)
	} else {
		log.Printf("wrote file %s", file)
	}

	return err
}

func RemoveDirectory(directory string) error {
	log.Printf("removing directory %s", directory)

	// remove
	err := os.RemoveAll(directory)
	if err != nil {
		log.Printf("error while removing directory %s", directory)
		log.Printf("error: %s", err)
	} else {
		log.Printf("directory %s removed", directory)
	}

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