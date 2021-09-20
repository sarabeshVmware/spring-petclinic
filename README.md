# tap-tests

Repo to maintain the test scripts for Tanzu Application Platform Packaging

## CLI Usage
```
➜ go run test.go -h               
TAP packaging tests CLI

Usage:
  tap-tests [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  install     Install packages

Flags:
  -h, --help   help for tap-tests

Use "tap-tests [command] --help" for more information about a command.
```

## Installing packages
### User Input
Configurable values such as credentials, repository image, etc. are provided via `user_input.yaml`. The supported fields are:
```yaml
namespace:
secrets:
  - name:
    registry:
    username:
    password:
package_repository:
  name:
  image:
packages:
  - name:
    installed_name:
    version:
    use_values_file:
```
### Steps
1. Currently, the script can add the package repository and installs packages in a configured environment. To set that up, refer: [Installing Tanzu Application Platform](https://docs-staging.vmware.com/en/VMware-Tanzu-Application-Platform/0.2/tap-0-2/GUID-install-intro.html)
2. Populate `user_input.yaml`. To use the default setup:
    - Provide credentials:
      ```yaml
      secrets:
      - name: tap-registry
        registry: registry.tanzu.vmware.com
        username:
        password:
      ```
    - Select the list of packages you want to install in `user_input.yaml`. For example, for installing `cloud-native-runtimes` and `app-accelerator`, just keep:
      ```yaml
      packages:
        - name: accelerator.apps.tanzu.vmware.com
          installed_name: accelerator-apps
          version: <version>
        - name: cnrs.tanzu.vmware.com
          installed_name: cnrs
          version: <version>
      ```
3. Run: `go run test.go install [--pre-cleanup] [--post-cleanup]`
    ```
    ➜ go run test.go install -h
    Install packages

    Usage:
      tap-tests install [flags]

    Flags:
      -h, --help           help for install
          --post-cleanup   Cleanup namespace, secrets, repository and packages after installation.
          --pre-cleanup    Cleanup namespace, secrets, repository and packages before installation.
    ```
