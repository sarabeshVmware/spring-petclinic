# e2e-suite

## Directory structure
<pre>
.<font color="#303030">
├── TAP_e2e_innerloop_01_test.go                           -> Innerloop test: basic supply chain, git source live update
├── TAP_e2e_innerloop_02_test.go                           -> Innerloop test: basic supply chain, local source
├── TAP_e2e_outerloop_01_test.go                           -> Outerloop test: basic supply chain, git source
├── TAP_e2e_outerloop_02_test.go                           -> Outerloop test: testing supply chain, git source
├── TAP_e2e_outerloop_03_test.go                           -> Outerloop test: testing & scanning supply chain, git source
├── TAP_e2e_outerloop_04_test.go                           -> Outerloop test: basic supply chain, gitops delivery
├── TAP_e2e_outerloop_05_test.go                           -> Outerloop test: testing & scanning supply chain, multiple apps</font><font color="#505050">
├── common_innerloop_features_test.go                      -> Common innerloop features
├── common_outerloop_features_test.go                      -> Common outerloop features</font><font color="#707070">
├── envfuncs                                               -> Package for suite level env functions
│   ├── cluster-essentials-funcs.go                        -> Suite-level functions for cluster essentials
│   ├── clusterrole-funcs.go                               -> Suite-level functions for cluster role
│   ├── namespace-funcs.go                                 -> Suite-level functions for namespace
│   ├── package-funcs.go                                   -> Suite-level functions for package
│   ├── package-repository-funcs.go                        -> Suite-level functions for package repository
│   └── secret-funcs.go                                    -> Suite-level functions for secret</font><font color="#909090">
├── go.mod
├── go.sum
├── main_test.go                                           -> Main test function for suite setup/cleanup</font><font color="#B0B0B0">
├── pkg
│   ├── docker                                             -> Package for docker functions
│   │   └── docker.go
│   ├── git                                                -> Package for git functions
│   │   └── git.go
│   ├── github                                             -> Package for github functions
│   │   └── github.go
│   ├── imgpkg                                             -> Package for imgpkg functions
│   │   └── imgpkg.go
│   ├── kubectl                                            -> Package for kubectl functions
│   │   ├── kubectlCmds
│   │   │   └── kubectl-commands.go
│   │   ├── kubectl_helpers
│   │   │   └── kubectl-helper.go
│   │   ├── kubectl_libs
│   │   │   ├── carto-deliverables.go
│   │   │   ├── carto-runnables.go
│   │   │   ├── carto-workload.go
│   │   │   ├── carvel-pkgi.go
│   │   │   ├── conventions-podintents.go
│   │   │   ├── flux-gitrepo.go
│   │   │   ├── kubectl.go
│   │   │   ├── learningcenter-trainingportals.go
│   │   │   ├── ootb-supplychain-scanning.go
│   │   │   ├── servicebindings.go
│   │   │   ├── sourcecontroller-imagerepositories.go
│   │   │   ├── tekton-pipelines.go
│   │   │   ├── tekton-prs.go
│   │   │   └── tekton-taskruns.go
│   │   └── unit_tests
│   │       └── main.go
│   ├── kubernetes                                         -> Package for kubernetes api functions
│   │   └── client
│   │       └── client.go
│   ├── misc                                               -> Package for miscellaneous functions
│   │   └── misc.go
│   ├── pivnet                                             -> Package for pivnet functions
│   │   ├── pivnet_helpers
│   │   │   └── pivnet_helpers.go
│   │   ├── pivnet_libs
│   │   │   ├── artifacts.go
│   │   │   ├── filegroups.go
│   │   │   ├── login.go
│   │   │   ├── productfiles.go
│   │   │   ├── release.go
│   │   │   └── user-group.go
│   │   └── scripts
│   │       ├── config.yaml
│   │       ├── create-release-version-tag.go
│   │       └── create-tanzunet-release.go
│   ├── tanzu                                              -> Package for tanzu cli functions
│   │   ├── tanzuCmds
│   │   │   └── tanzu-commands.go
│   │   ├── tanzu_helpers
│   │   │   └── tanzu-helpers.go
│   │   ├── tanzu_libs
│   │   │   ├── accelerator.go
│   │   │   ├── apps-clustersupplychain.go
│   │   │   ├── apps-workload.go
│   │   │   ├── insight-config.go
│   │   │   ├── insight-images.go
│   │   │   ├── package-available.go
│   │   │   ├── package-install.go
│   │   │   ├── package-installed.go
│   │   │   ├── package-repository.go
│   │   │   └── secret-registry.go
│   │   └── unit_tests
│   │       └── main.go
│   └── utils                                              -> Package for utilities' functions
│       ├── linux_util
│       │   └── linux_util.go
│       └── utils.go</font><font color="#D0D0D0">
├── resources                                              -> Resources directory
│   ├── components                                         -> Resources for components
│   │   ├── cert-manager-install.yaml                      -> PackageInstall CR for cert-manager
│   │   ├── cert-manager-rbac.yaml                         -> RBAC for cert-manager
│   │   ├── cnrs.yaml                                      -> Values schema for cnrs
│   │   ├── contour-install.yaml                           -> PackageInstall CR for contour
│   │   ├── contour-rbac.yaml                              -> RBAC for contour
│   │   ├── install-metadata.yaml                          -> List of packages sorted by dependency
│   │   ├── learning-center-config.yaml                    -> Values schema for learning center
│   │   ├── ootb-supply-chain-basic-values.yaml            -> Values schema for OOTB supply chain basic
│   │   ├── ootb-supply-chain-testing-scanning-values.yaml -> Values schema for OOTB supply chain testing & scanning
│   │   ├── ootb-supply-chain-testing-values.yaml          -> Values schema for OOTB supply chain testing
│   │   ├── tap-values.yaml                                -> Values schema for tap
│   │   └── tbs-values.yaml                                -> Values schema for build service
│   ├── innerloop                                          -> Resources for innerloop tests
│   │   └── tanzu-web-app-workload.yaml                    -> Workload YAML for tanzu-java-web-app
│   ├── outerloop                                          -> Resources for outerloop
│   │   ├── git-ssh-secrets.yaml                           -> Secret YAML for git ssh
│   │   ├── lenient-scan-policy.yaml                       -> Lenient scan policy YAML (violatingSeverities := ["UnknownSeverity"])
│   │   ├── mysql-service.yaml                             -> YAML for mysql deployment
│   │   ├── outerloop-config.yaml                          -> Config file for outerloop tests
│   │   ├── pipeline-buildpacks-test.yaml                  -> Tekton pipeline for buildpacks test
│   │   ├── scan-policy.yaml                               -> Scan policy YAML
│   │   ├── spring-petclinic-tests-pipeline.yaml           -> Tekton pipeline for spring-petclinic
│   │   ├── workload-gitops.yaml                           -> Workload YAML for spring-petclinic with gitops
│   │   ├── workload-test.yaml                             -> Workload YAML for spring-petclinic with testing (has-tests: true)
│   │   └── workload.yaml                                  -> Workload YAML for spring-petclinic
│   └── suite                                              -> Resources for suite-level configuration
│       ├── developer-namespace.yaml                       -> Developer namespace YAML
│       ├── suite-config.yaml                              -> Config file for suite
│       └── tap-values.yaml                                -> Values schema for tap</font><font color="#FFFFFF">
└── tap_test                                               -> New directory for tests
    ├── common_features                                    -> Package for common features
    │   ├── common_features.go                             -> Common features
    │   ├── common_innerloop_features.go                   -> Common innerloop features
    │   └── common_outerloop_features.go                   -> Common outerloop features
    ├── install_test                                       -> Package for install tests
    │   ├── TAP_e2e_TAP_upgrade_downgrade_01_test.go       -> Upgrade/downgrade test
    │   ├── TAP_e2e_install_uninstall_01_test.go           -> Install/uninstall test
    │   └── main_test.go                                   -> Main test function for suite setup/cleanup
    ├── models                                             -> YAML files' read and data structure handling
    │   ├── outerloopConfig.go                             -> Functions to handle resources/outerloop/outerloop-config.yaml
    │   ├── suiteConfig.go                                 -> Functions to handle resources/suite/suite-config.yaml
    │   └── tapValuesSchema.go                             -> Functions to handle resources/suite/tap-values.yaml
    └── pre_install_test                                   -> Package for pre-install tests
        ├── TAP_e2e_TAP_image_relocation_01_test.go        -> Image relocation test
        └── main_test.go                                   -> Main test function for suite setup/cleanup</font>
</pre>

## Usage
`go test . <flags> [-tags=TAG1[,TAG2..]]`

### Tags:
<pre>
TestInnerloopBasicSupplychainGitSourceLiveUpdate -> all || innerloop || innerloop_basic_git_source
TestInnerloopBasicSupplychainLocalSource         -> all || innerloop || innerloop_basic
TestOuterloopBasicSupplychainGitSource           -> all || outerloop || outerloop_basic
TestOuterloopTestSupplychainGitSource            -> all || outerloop || outerloop_testing
TestOuterloopScanSupplychainGitSource            -> all || outerloop || outerloop_testing_scanning
TestOuterloopBasicSupplychainGitopsDelivery      -> all || outerloop || outerloop_basic_delivery
TestOuterloopScanSupplychainMultipleApps         -> all || outerloop || outerloop_scan_multiple_apps
</pre>

Example: `go test . -v -timeout=60m -tags=innerloop_basic`

## Miscellaneous
- [CI pipeline](https://gitlab.eng.vmware.com/dap-engineering-operations/tap-pipeline/-/blob/main/ci/tasks/e2e-suite.sh)
