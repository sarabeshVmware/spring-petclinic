// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package kubectlCmds

import (
	"fmt"
	"log"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

func KubectlCreateNamespace(namespace string) error {
	log.Printf("creating namespace %s", namespace)

	// execute cmd
	cmd := fmt.Sprintf("kubectl create ns %s", namespace)
	output, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while creating namespace %s", namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("namespace %s created", namespace)
		log.Printf("output: %s", output)
	}

	return err
}

func KubectlDeleteNamespace(namespace string) error {
	log.Printf("deleting namespace %s", namespace)

	// execute cmd
	cmd := fmt.Sprintf("kubectl get namespace %[1]s -o json | tr -d \"\n\" | sed 's/\"finalizers\": [[^]]+]/\"finalizers\": []/' | kubectl replace --raw /api/v1/namespaces/%[1]s/finalize -f -; kubectl delete ns %[1]s --wait=false", namespace)
	output, err := linux_util.ExecuteCmdInBashMode(cmd)
	if err != nil {
		log.Printf("error while deleting namespace %s", namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("namespace %s deleted", namespace)
		log.Printf("output: %s", output)
	}

	return err
}

func KubectlApplyConfiguration(file string, namespace string) error {
	log.Printf("applying configuration %s in namespace %s", file, namespace)

	// execute cmd
	cmd := fmt.Sprintf("kubectl apply -n %s -f %s", namespace, file)
	output, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while applying configuration %s in namespace %s", file, namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("configuration %s applied in namespace %s", file, namespace)
		log.Printf("output: %s", output)
	}

	return err
}

func KubectlDeleteConfiguration(file string, namespace string) error {
	log.Printf("deleting configuration %s from namespace %s", file, namespace)

	// execute cmd
	cmd := fmt.Sprintf("kubectl delete -n %s -f %s", namespace, file)
	output, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Printf("error while deleting configuration %s from namespace %s", file, namespace)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("configuration %s deleted from namespace %s", file, namespace)
		//log.Printf("output: %s", output)
	}

	return err
}
