package tanzu_libs

// Usage:
//   tanzu secret registry [command]

// Available Commands:
//   add         Creates a v1/Secret resource of type kubernetes.io/dockerconfigjson. In case of specifying the --export-to-all-namespaces flag, a SecretExport resource will also get created
//   delete      Deletes v1/Secret resource of type kubernetes.io/dockerconfigjson and the associated SecretExport from the cluster
//   list        Lists all v1/Secret of type kubernetes.io/dockerconfigjson and checks for the associated SecretExport with the same name
//   update      Updates the v1/Secret resource of type kubernetes.io/dockerconfigjson. In case of specifying the --export-to-all-namespaces flag, the SecretExport resource will also get updated. Otherwise, there will be no changes in the SecretExport resource
