// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package stepfuncs

// import (
// 	"context"
// 	"fmt"
// 	"os"

// 	// "strings"

// 	"testing"

// 	"gitlab.eng.vmware.com/tap/tap-packages/suite/client"
// 	"gitlab.eng.vmware.com/tap/tap-packages/suite/exec"
// 	"gopkg.in/yaml.v3"
// 	"sigs.k8s.io/e2e-framework/pkg/envconf"
// )

// func WriteFile(ctx context.Context, t *testing.T, cfg *envconf.Config, failNow bool, file string, contents interface{}) context.Context {
// 	t.Logf("writing file %s", file)
// 	bytes, err := yaml.Marshal(contents)
// 	if err != nil {
// 		t.Error(fmt.Errorf("error while writing file %s: %w", file, err))
// 		if failNow {
// 			t.FailNow()
// 		} else {
// 			t.Fail()
// 		}
// 	}
// 	err = os.WriteFile(file, bytes, 0677)
// 	if err != nil {
// 		t.Error(fmt.Errorf("error while writing file %s: %w", file, err))
// 		if failNow {
// 			t.FailNow()
// 		} else {
// 			t.Fail()
// 		}
// 	}
// 	t.Logf("file %s written", file)
// 	return ctx
// }

// func UpdatePackage(ctx context.Context, t *testing.T, cfg *envconf.Config, failNow bool, name string, packageName string, version string, namespace string, valuesSchemaFile string) context.Context {
// 	t.Logf("updating package %s", name)
// 	cmd, output, err := exec.TanzuUpdatePackage(name, packageName, version, namespace, valuesSchemaFile)
// 	t.Logf("command executed: %s", cmd)
// 	if err != nil {
// 		t.Error(fmt.Errorf("error while updating package %s: %w: %s", name, err, output))
// 		if failNow {
// 			t.FailNow()
// 		} else {
// 			t.Fail()
// 		}
// 	}
// 	t.Logf("package %s updated: %s", name, output)
// 	return ctx
// }

// func GetServiceExternalIp(ctx context.Context, t *testing.T, cfg *envconf.Config, service string, namespace string) (context.Context, string) {
// 	t.Logf("getting external ip for %s (namespace %s)", service, namespace)
// 	serviceExternalIp, err := client.GetServiceExternalIP(service, namespace, cfg.Client().RESTConfig())
// 	if err != nil {
// 		t.Error(fmt.Errorf("error while getting external ip for %s (namespace %s): %w", service, namespace, err))
// 		t.FailNow()
// 	}
// 	t.Logf("external ip for %s (namespace %s): %s", "server", namespace, serviceExternalIp)
// 	return ctx, serviceExternalIp
// }

// func GetServiceExternalPort(ctx context.Context, t *testing.T, cfg *envconf.Config, service string, namespace string) (context.Context, int) {
// 	t.Logf("getting port for %s (namespace %s)", service, namespace)
// 	servicePort, err := client.GetServicePort(service, namespace, cfg.Client().RESTConfig())
// 	if err != nil {
// 		t.Error(fmt.Errorf("error while getting port for %s (namespace %s): %w", service, namespace, err))
// 		t.FailNow()
// 	}
// 	t.Logf("port for %s (namespace %s): %d", service, namespace, servicePort)
// 	return ctx, servicePort
// }

// func PatchServiceAccount(ctx context.Context, t *testing.T, cfg *envconf.Config, patch string, serviceAccount string, namespace string) context.Context {
// 	t.Logf("patching %s to service account %s in namespace %s", patch, serviceAccount, namespace)
// 	cmd, output, err := exec.KubectlPatchServiceAccount(serviceAccount, namespace, patch)
// 	t.Logf("command executed: %s", cmd)
// 	if err != nil {
// 		t.Error(fmt.Errorf("error while patching %s to service account %s in namespace %s: %w: %s", patch, serviceAccount, namespace, err, output))
// 		t.FailNow()
// 	}
// 	t.Logf("%s patched to service account %s in namespace %s: %s", patch, serviceAccount, namespace, output)
// 	return ctx
// }

// func CreateClusterRoleBinding(ctx context.Context, t *testing.T, cfg *envconf.Config, name string, clusterRole string, serviceAccount string) context.Context {
// 	t.Logf("creating clusterrolebinding %s for clusterrole %s and serviceaccount %s", name, clusterRole, serviceAccount)
// 	cmd, output, err := exec.KubectlCreateClusterRoleBinding(name, clusterRole, serviceAccount)
// 	t.Logf("command executed: %s", cmd)
// 	if err != nil {
// 		t.Error(fmt.Errorf("error while creating cluster role binding %s for clusterrole %s and serviceaccount %s: %w: %s", name, clusterRole, serviceAccount, err, output))
// 		t.FailNow()
// 	}
// 	t.Logf("clusterrolebinding %s created for clusterrole %s and serviceaccount %s: %s", name, clusterRole, serviceAccount, output)
// 	return ctx
// }

// func DeployAppInNamespace(ctx context.Context, t *testing.T, cfg *envconf.Config, failNow bool, name string, files []string, namespace string) context.Context {
// 	t.Logf("deploying app %s in namespace %s", name, namespace)
// 	cmd, output, err := exec.KappDeployAppInNamespace(name, files, namespace)
// 	t.Logf("command executed: %s", cmd)
// 	if err != nil {
// 		t.Error(fmt.Errorf("error while deploying app %s in namespace %s: %w: %s", name, namespace, err, output))
// 		if failNow {
// 			t.FailNow()
// 		} else {
// 			t.Fail()
// 		}
// 	}
// 	t.Logf("app %s deployed in namespace %s: %s", name, namespace, output)
// 	return ctx
// }

// func DeployWorkload(ctx context.Context, t *testing.T, cfg *envconf.Config, file string, namespace string) context.Context {
// 	t.Logf("deploying workload %s in namespace %s", file, namespace)
// 	cmd, output, err := exec.TanzuDeployWorkload(file, namespace)
// 	t.Logf("command executed: %s", cmd)
// 	if err != nil {
// 		t.Error(fmt.Errorf("error while deploying workload %s in namespace %s: %w: %s", file, namespace, err, output))
// 		t.FailNow()
// 	}
// 	t.Logf("workload %s deployed in namespace %s: %s", file, namespace, output)
// 	return ctx
// }

// func GitClone(ctx context.Context, t *testing.T, cfg *envconf.Config, path string, repo string) context.Context {
// 	t.Logf("cloning repository %s at %s", repo, path)
// 	cmd, output, err := exec.GitClone(path, repo)
// 	t.Logf("command executed: %s", cmd)
// 	if err != nil {
// 		t.Error(fmt.Errorf("error while cloning repository %s at %s: %w: %s", repo, path, err, output))
// 		t.FailNow()
// 	}
// 	t.Logf("repository %s cloned at %s: %s", repo, path, output)
// 	return ctx
// }

// func GitAdd(ctx context.Context, t *testing.T, cfg *envconf.Config, path string, files []string) context.Context {
// 	t.Logf("adding files %s for repository at %s", files, path)
// 	cmd, output, err := exec.GitAdd(path, files)
// 	t.Logf("command executed: %s", cmd)
// 	if err != nil {
// 		t.Error(fmt.Errorf("error while adding files %s for repository at %s: %w: %s", files, path, err, output))
// 		t.FailNow()
// 	}
// 	t.Logf("files %s added for repository at %s: %s", files, path, output)
// 	return ctx
// }

// func GitCommit(ctx context.Context, t *testing.T, cfg *envconf.Config, path string, message string) context.Context {
// 	t.Logf("committing files for repository at %s (message %s)", path, message)
// 	cmd, output, err := exec.GitCommit(path, message)
// 	t.Logf("command executed: %s", cmd)
// 	if err != nil {
// 		t.Error(fmt.Errorf("error while committing files for repository at %s: %w: %s", path, err, output))
// 		t.FailNow()
// 	}
// 	t.Logf("committed files for repository at %s (message %s): %s", path, message, output)
// 	return ctx
// }

// func GitPush(ctx context.Context, t *testing.T, cfg *envconf.Config, path string, force bool) context.Context {
// 	t.Logf("pushing commits for repository at %s", path)
// 	cmd, output, err := exec.GitPush(path, force)
// 	t.Logf("command executed: %s", cmd)
// 	if err != nil {
// 		t.Error(fmt.Errorf("error while pushing commits for repository at %s: %w: %s", path, err, output))
// 		t.FailNow()
// 	}
// 	t.Logf("pushed commits for repository at %s: %s", path, output)
// 	return ctx
// }

// func GitResetFromHead(ctx context.Context, t *testing.T, cfg *envconf.Config, path string, count int) context.Context {
// 	t.Logf("resetting commits at HEAD~%d for repository at %s", count, path)
// 	cmd, output, err := exec.GitResetFromHead(path, count)
// 	t.Logf("command executed: %s", cmd)
// 	if err != nil {
// 		t.Error(fmt.Errorf("error while resetting commits at HEAD~%d for repository at %s: %w: %s", count, path, err, output))
// 		t.FailNow()
// 	}
// 	t.Logf("resetted commits at HEAD~%d for repository at %s: %s", count, path, output)
// 	return ctx
// }

// func RemoveDirectory(ctx context.Context, t *testing.T, cfg *envconf.Config, dir string) context.Context {
// 	t.Logf("removing directory %s", dir)
// 	err := os.RemoveAll(dir)
// 	if err != nil {
// 		t.Error(fmt.Errorf("error while removing directory %s: %w", dir, err))
// 		t.FailNow()
// 	}
// 	t.Logf("directory %s removed", dir)
// 	return ctx
// }

// // func UpdateFileReplaceString(ctx context.Context, t *testing.T, cfg *envconf.Config, file string, originalString string, newString string) context.Context {
// // 	t.Logf("updating file %s", file)
// // 	inputBytes, err := os.ReadFile(file)
// // 	if err != nil {
// // 		t.Error(fmt.Errorf("error while updating file %s: %w", file, err))
// // 		t.FailNow()
// // 	}
// // 	input := strings.ReplaceAll(string(inputBytes), originalString, newString)
// // 	return WriteFile(ctx, t, cfg, file, input)
// // }
