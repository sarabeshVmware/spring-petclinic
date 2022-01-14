// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import "fmt"

func KappDeployAppInNamespace(name string, files []string, namespace string) (string, string, error) {
	cmd := fmt.Sprintf("kapp deploy -a %s -n %s -y", name, namespace)
	for _, file := range files {
		cmd += fmt.Sprintf(" -f %s", file)
	}
	output, err := RunCommand(cmd)
	return cmd, output, err
}
