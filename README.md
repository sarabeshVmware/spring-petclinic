**TAP Packages**

Carvel packages that are included in a TAP repo bundle.

**Repository Layout**
- **.imgpkg/**: images.yml contains package bundle references used by package CRs.

- **packages**/: packages directory contains package metadata and package CRs for each package. 
packages directory contains package metadata and package CRs.

**Dev workflow**
1. Create a directory for your package.
   for ex. tbs.tanzu.vmware.com
2. Add Package metadata CR into above created package directory with metadata.yml name.
   for ex. tbs.tanzu.vmware.com/metadata.yml
3. Add Package CR into above created package directory with file name as <version>.yml.
   for ex. tbs.tanzu.vmware.com/1.2.0.yml

Generate TAP repo bundle
kbld -f ./packages --imgpkg-lock-output ./.imgpkg/images.yml

Push it to target registry
imgpkg push <>
