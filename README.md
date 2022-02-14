**TAP Packages**

Carvel packages that are included in a TAP repo bundle.

**Repository Layout**
- **packages/**: packages directory contains package metadata and package CRs for each package. 

**Dev workflow**
1. Create a directory for your package.
   for ex. app-accelerator
2. Add Package metadata CR into above created package directory with metadata.yml name.
   for ex. app-accelerator/metadata.yml
3. Create a directory for package version inside package directory.
   for ex. app-accelerator/0.2.0
3. Add Package CR into above created version directory with file name as package.yml.
   for ex. app-accelerator/0.2.0/package.yaml

**TAP Packaging**

Generate TAP repo bundle and push to dev.registry.tanzu.vmware.com
go run create-package-repo.go main <tag for tap repo>


Note: While adding/modifying package CRs, If the product is released to User groups only in Tanzunet, Please add **DAP Internal Users** and **VMware Internal Early Access group** to the product early group on TanzuNet so that TAP team can validate it.

**For Component teams - Creating TAP Repo Bundle with local changes, kindly check [TAP Requirements Doc](https://docs.google.com/document/d/1Af66UbjEABF8GGYGbN-bGHRuakGsf-oAdQ4JvuR5HUY/edit#heading=h.jgnkc0jhnafj)**

1. Create a package CR for the newer version. 

2. Check out the latest tap-packages. Add a new directory with the version under `./packages/<component>/<version>`. Add the `package.yml` to the directory
   
   eg: /packages/api-portal/1.0.9/package.yaml

3. Update the min version constraint in tap-pkg under `./tap-pkg/config/<component>.yaml`
   
   eg: /tap-pkg/config/api-portal.yaml

4. Update the package version to be included in the tap repo bundle under `./repos/<release>.yaml`
   
   eg: /repos/1.1.0.yaml

5. Do a docker login with the user that has access to all the tap components. For example, `tap-dev-internal-user`

6. Run the `create-package-repo.go` script under ./scripts with paramenters to generate the tap-repo bundle
   
   ```
   go run scripts/create-package-repo.go <release> <tag> <repository> <registry>
   ```
   
   Where

      `release` is the release for which the bundle to be created. This should match the release file name under `./repos`. For example, if you want to create a  repo bundle from `1.1.0.yaml` release file, then the release should be `1.1.0`

      `tag` is the aftifact tag for the bundle. For example, `1.1.0-testbuild.1`

      `repository` is the path where to push the bundle in the OCI registry. For example, `tanzu-application-platform/tap-packages-trial`

      `registry` is the OCI registry to which the bundle to be pushed. For example, `dev.registry.tanzu.vmware.com`

   ```
   go run scripts/create-package-repo.go 1.0.1 1.0.1-testbuild.1 tanzu-application-platform/tap-packages-trial dev.registry.tanzu.vmware.com

   Executing command :  kbld [--file repos/generated/1.0.1/packages --imgpkg-lock-output repos/generated/1.0.1/.imgpkg/images.yml]
   Executing command :  imgpkg [push --tty --bundle dev.registry.tanzu.vmware.com/tanzu-application-platform/tap-packages-trial:1.0.1-testbuild.1 --file repos/generated/1.0.1 --lock-output output.yaml]
   Package Repository pushed to dev.registry.tanzu.vmware.com/tanzu-application-platform/tap-packages-trial@sha256:c595181f8079d4d420e9484f91a3457117d2e067de81dec1be8574e98b596062
   ```
