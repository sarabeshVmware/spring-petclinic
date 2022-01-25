// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package stepfuncs

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"testing"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/exec"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

func VerifyImagerepositoryReady(ctx context.Context, t *testing.T, cfg *envconf.Config, name string, namespace string) context.Context {
	ready := false
	iter := 30
	for i := 1; i <= iter; i++ {
		t.Logf("getting image repository ready status for %s in namespace %s (iteration %d)", name, namespace, i)
		cmd, output, err := exec.KubectlIsImageRepositoryReady(name, namespace)
		t.Logf("command executed: %s", cmd)
		if err != nil {
			t.Error(fmt.Errorf("error while getting image repository ready status for %s in namespace %s: %w: %t", name, namespace, err, output))
			t.FailNow()
		}
		t.Logf("image repository ready status for %s in namespace %s: %t", name, namespace, output)
		if output == true {
			ready = true
			break
		}
		t.Logf("sleeping for 1 minute")
		time.Sleep(time.Minute)
	}
	if !ready {
		t.Errorf(`image repository failed to get into ready state after %d iterations`, iter)
		t.Fail()
	}
	return ctx
}

func VerifyGitrepoReady(ctx context.Context, t *testing.T, cfg *envconf.Config, name string, namespace string) context.Context {
	ready := false
	iter := 30
	for i := 1; i <= iter; i++ {
		t.Logf("getting git repository ready status for %s in namespace %s (iteration %d)", name, namespace, i)
		cmd, output, err := exec.KubectlIsGitRepositoryReady(name, namespace)
		t.Logf("command executed: %s", cmd)
		if err != nil {
			t.Error(fmt.Errorf("error while getting git repository ready status for %s in namespace %s: %w: %t", name, namespace, err, output))
			t.FailNow()
		}
		t.Logf("git repository ready status for %s in namespace %s: %t", name, namespace, output)
		if output == true {
			ready = true
			break
		}
		t.Logf("sleeping for 1 minute")
		time.Sleep(time.Minute)
	}
	if !ready {
		t.Errorf(`git repository failed to get into ready state after %d iterations`, iter)
		t.Fail()
	}
	return ctx
}

func VerifyBuildSucceeded(ctx context.Context, t *testing.T, cfg *envconf.Config, prefix string, namespace string) context.Context {
	succeeded := false
	iter := 30
	for i := 1; i <= iter; i++ {
		t.Logf("getting build succeeded status for %s in namespace %s (iteration %d)", prefix, namespace, i)
		cmd, output, err := exec.KubectlIsBuildSucceeded(prefix, namespace)
		t.Logf("command executed: %s", cmd)
		if err != nil {
			t.Error(fmt.Errorf("error while getting build succeeded status for %s in namespace %s: %w: %t", prefix, namespace, err, output))
			t.FailNow()
		}
		t.Logf("build succeeded status for %s in namespace %s: %t", prefix, namespace, output)
		if output == true {
			succeeded = true
			break
		}
		t.Logf("sleeping for 1 minute")
		time.Sleep(time.Minute)
	}
	if !succeeded {
		t.Errorf(`build failed to get into succeeded state after %d iterations`, iter)
		t.Fail()
	}
	return ctx
}

func VerifyPodintentReady(ctx context.Context, t *testing.T, cfg *envconf.Config, name string, namespace string) context.Context {
	ready := false
	iter := 30
	for i := 1; i <= iter; i++ {
		t.Logf("getting podintent ready status for %s in namespace %s (iteration %d)", name, namespace, i)
		cmd, output, err := exec.KubectlIsPodintentReady(name, namespace)
		t.Logf("command executed: %s", cmd)
		if err != nil {
			t.Error(fmt.Errorf("error while getting podintent ready status for %s in namespace %s: %w: %t", name, namespace, err, output))
			t.FailNow()
		}
		t.Logf("podintent ready status for %s in namespace %s: %t", name, namespace, output)
		if output == true {
			ready = true
			break
		}
		t.Logf("sleeping for 1 minute")
		time.Sleep(time.Minute)
	}
	if !ready {
		t.Errorf(`podintent failed to get into ready state after %d iterations`, iter)
		t.Fail()
	}
	return ctx
}

func VerifyKsvcReady(ctx context.Context, t *testing.T, cfg *envconf.Config, name string, namespace string) context.Context {
	ready := false
	iter := 30
	for i := 1; i <= iter; i++ {
		t.Logf("getting ksvc ready status for %s in namespace %s (iteration %d)", name, namespace, i)
		cmd, output, err := exec.KubectlIsKsvcReady(name, namespace)
		t.Logf("command executed: %s", cmd)
		if err != nil {
			t.Error(fmt.Errorf("error while getting ksvc ready status for %s in namespace %s: %w: %t", name, namespace, err, output))
			t.FailNow()
		}
		t.Logf("ksvc ready status for %s in namespace %s: %t", name, namespace, output)
		if output == true {
			ready = true
			break
		}
		t.Logf("sleeping for 1 minute")
		time.Sleep(time.Minute)
	}
	if !ready {
		t.Errorf(`ksvc failed to get into ready state after %d iterations`, iter)
		t.Fail()
	}
	return ctx
}

func VerifyTaskrunSucceeded(ctx context.Context, t *testing.T, cfg *envconf.Config, prefix string, namespace string) context.Context {
	succeeded := false
	iter := 30
	for i := 1; i <= iter; i++ {
		t.Logf("getting taskrun succeeded status for %s in namespace %s (iteration %d)", prefix, namespace, i)
		cmd, output, err := exec.KubectlIsTaskrunSucceeded(prefix, namespace)
		t.Logf("command executed: %s", cmd)
		if err != nil {
			t.Error(fmt.Errorf("error while getting taskrun succeeded status for %s in namespace %s: %w: %t", prefix, namespace, err, output))
			t.FailNow()
		}
		t.Logf("taskrun succeeded status for %s in namespace %s: %t", prefix, namespace, output)
		if output == true {
			succeeded = true
			break
		}
		t.Logf("sleeping for 1 minute")
		time.Sleep(time.Minute)
	}
	if !succeeded {
		t.Errorf(`taskrun failed to get into succeeded state after %d iterations`, iter)
		t.Fail()
	}
	return ctx
}

func VerifyPodintentAnnotation(ctx context.Context, t *testing.T, cfg *envconf.Config, annotationKey string, annotationValue string, checkOnlyKey bool, podintent string, namespace string) context.Context {
	exists := false
	iter := 30
	for i := 1; i <= iter; i++ {
		t.Logf("getting podintent annotation %s for %s in namespace %s (iteration %d)", annotationKey, podintent, namespace, i)
		cmd, output, err := exec.KubectlIsPodintentAnnotationExists(annotationKey, annotationValue, checkOnlyKey, podintent, namespace)
		t.Logf("command executed: %s", cmd)
		if err != nil {
			t.Error(fmt.Errorf("error while getting podintent annotation %s for %s in namespace %s: %w: %t", annotationKey, podintent, namespace, err, output))
			t.FailNow()
		}
		t.Logf("podintent annotation %s status for %s in namespace %s: %t", annotationKey, podintent, namespace, output)
		if output == true {
			exists = true
			break
		}
		t.Logf("sleeping for 1 minute")
		time.Sleep(time.Minute)
	}
	if !exists {
		t.Errorf(`podintent annotation doesn't exist after %d iterations`, iter)
		t.Fail()
	}
	return ctx
}

func VerifyPodintentLabel(ctx context.Context, t *testing.T, cfg *envconf.Config, labelKey string, labelValue string, checkOnlyKey bool, podintent string, namespace string) context.Context {
	exists := false
	iter := 30
	for i := 1; i <= iter; i++ {
		t.Logf("getting podintent label %s for %s in namespace %s (iteration %d)", labelKey, podintent, namespace, i)
		cmd, output, err := exec.KubectlIsPodintentLabelExists(labelKey, labelValue, checkOnlyKey, podintent, namespace)
		t.Logf("command executed: %s", cmd)
		if err != nil {
			t.Error(fmt.Errorf("error while getting podintent label %s for %s in namespace %s: %w: %t", labelKey, podintent, namespace, err, output))
			t.FailNow()
		}
		t.Logf("podintent label %s status for %s in namespace %s: %t", labelKey, podintent, namespace, output)
		if output == true {
			exists = true
			break
		}
		t.Logf("sleeping for 1 minute")
		time.Sleep(time.Minute)
	}
	if !exists {
		t.Errorf(`podintent label doesn't exist after %d iterations`, iter)
		t.Fail()
	}
	return ctx
}

func VerifyApplicationRunningWithValidationString(ctx context.Context, t *testing.T, cfg *envconf.Config, envoyExternalIP string, host string, validationString string) context.Context {
	t.Logf("checking application %s for result: %s", host, validationString)

	validated := false
	iter := 30
	for i := 1; i <= iter; i++ {
		if !strings.HasPrefix(envoyExternalIP, "http://") {
			envoyExternalIP = "http://" + envoyExternalIP
		}
		req, err := http.NewRequest("GET", envoyExternalIP, nil)
		if err != nil {
			t.Error(fmt.Errorf("error while giving http request: %w", err))
			t.FailNow()
		}
		req.Host = host

		var retries int = 10
		for retries > 0 {
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				retries -= 1
				t.Logf("didn't get response")
				t.Logf("sleeping for 30 seconds")
				time.Sleep(30 * time.Second)
			} else {
				t.Logf("status code is: %d", resp.StatusCode)
				break
			}
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(fmt.Errorf("error while giving http response: %w", err))
			t.FailNow()
		}
		if resp.StatusCode != http.StatusOK {
			t.Error(fmt.Errorf("bad HTTP Response: %s", resp.Status))
			t.FailNow()
		}
		defer resp.Body.Close()
		resultStringBytes, _ := ioutil.ReadAll(resp.Body)
		resultString := string(resultStringBytes)
		t.Logf(resultString)
		if strings.Contains(resultString, validationString) {
			t.Logf("application %s validated, got result: %s", host, validationString)
			validated = true
			break
		} else {
			t.Logf("getting string %s", resultString)
			t.Logf("sleeping for 30 seconds")
			time.Sleep(30 * time.Second)
		}
	}

	if !validated {
		t.Errorf(`application %s not validated %d iterations`, host, iter)
		t.Fail()
	}
	return ctx
}
