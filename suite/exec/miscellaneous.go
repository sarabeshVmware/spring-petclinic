// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"log"
	"os"
	"strings"
)

func ReplaceStringInFile(filePath string, oldString string, newString string) error {
	log.Printf("Updating file %s: %s -> %s", filePath, oldString, newString)
	inputBytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	input := strings.ReplaceAll(string(inputBytes), oldString, newString)
	err = os.WriteFile(filePath, []byte(input), 0666)
	return err
}
