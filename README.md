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


Note: While adding/modifying package CRs, If the product is released to User groups only in Tanzunet, Please add **VMware Internal Early Access group** to the product early group on TanzuNet so that TAP team can validate it.
