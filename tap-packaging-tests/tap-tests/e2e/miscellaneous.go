// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"log"
	"os"
	"strings"

	tap "gitlab.eng.vmware.com/tap/tap-packages/tap-packaging-tests/pkg"
)

func UpdateFile(filePath string, oldString string, newString string) {
	log.Printf("Updating file %s: %s -> %s", filePath, oldString, newString)
	inputBytes, err := os.ReadFile(filePath)
	tap.CheckError(err)
	input := strings.ReplaceAll(string(inputBytes), oldString, newString)
	err = os.WriteFile(filePath, []byte(input), 0666)
	tap.CheckError(err)
}
