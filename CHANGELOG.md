# Changelog

## [v1.5.0](https://github.com/accurics/terrascan/tree/v1.5.0) (2021-04-23)

[Full Changelog](https://github.com/accurics/terrascan/compare/v1.4.0...v1.5.0)

**Fixed bugs:**

- Recursive loop expanding variables in included module [\#675](https://github.com/accurics/terrascan/issues/675)
- Terrascan doesn't resolve terraform complex variables [\#656](https://github.com/accurics/terrascan/issues/656)
- Panic while resolving floating point variable [\#652](https://github.com/accurics/terrascan/issues/652)
- Terrascan using absolute path for "source" value of resource [\#642](https://github.com/accurics/terrascan/issues/642)
- Failed to initialize terrascan. error : failed to install policies [\#614](https://github.com/accurics/terrascan/issues/614)
- Terrascan not able to read modules within a subdirectory [\#600](https://github.com/accurics/terrascan/issues/600)
- Terrascan init command doesn't work with -c flag [\#550](https://github.com/accurics/terrascan/issues/550)

**Closed issues:**

- Not able to scan repo when google terraform module defined [\#681](https://github.com/accurics/terrascan/issues/681)
- The link referencing the documentation to integrate Terrascan into CI/CD is broken [\#669](https://github.com/accurics/terrascan/issues/669)
- Make saving of "admission request" configurable via an option in the config file for the validating admission webhook [\#664](https://github.com/accurics/terrascan/issues/664)
- Add API\_KEY to the /logs endpoint for the validating admission webhook [\#662](https://github.com/accurics/terrascan/issues/662)
- Panic: not a string [\#647](https://github.com/accurics/terrascan/issues/647)
- unit tests and e2e tests failing on windows [\#639](https://github.com/accurics/terrascan/issues/639)
- Add support for private terraform repos [\#631](https://github.com/accurics/terrascan/issues/631)
- policy not evaluating [\#629](https://github.com/accurics/terrascan/issues/629)
- Terrascan does not support to download modules via SSH [\#621](https://github.com/accurics/terrascan/issues/621)
- terrascan scan fails if path and rego\_subdir are not provided together in the toml configfile [\#619](https://github.com/accurics/terrascan/issues/619)
- Getting error while running scan on our terraform repo [\#607](https://github.com/accurics/terrascan/issues/607)
- Terrascan not found policy id [\#601](https://github.com/accurics/terrascan/issues/601)
- `Policies Violated` and `Violated Policies` are confusing. [\#598](https://github.com/accurics/terrascan/issues/598)
- Invalid categories not being validated from config file [\#594](https://github.com/accurics/terrascan/issues/594)
- Terrascan API server's file scan doesn't work for k8s yaml files [\#584](https://github.com/accurics/terrascan/issues/584)
- Add `/go/bin` to the PATH variable in Docker image [\#577](https://github.com/accurics/terrascan/issues/577)
- terrascan scan command doesn't work with TERRASCAN\_CONFIG env variable [\#570](https://github.com/accurics/terrascan/issues/570)
- Format junit-xml need to have passed test results, not only failed test [\#563](https://github.com/accurics/terrascan/issues/563)
- optimize policy download process in `terrascan init` [\#535](https://github.com/accurics/terrascan/issues/535)

**Merged pull requests:**

- Release v1.5.0 [\#689](https://github.com/accurics/terrascan/pull/689) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- Adds support to configure dashboard mode in k8s validating webhook [\#683](https://github.com/accurics/terrascan/pull/683) ([patilpankaj212](https://github.com/patilpankaj212))
- Updating documentation for k8s admission control [\#679](https://github.com/accurics/terrascan/pull/679) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- Fix recursive variable reference resolution [\#677](https://github.com/accurics/terrascan/pull/677) ([patilpankaj212](https://github.com/patilpankaj212))
- Update mkdocs-material to 7.1.2 [\#676](https://github.com/accurics/terrascan/pull/676) ([pyup-bot](https://github.com/pyup-bot))
- Fixes broken link in README [\#671](https://github.com/accurics/terrascan/pull/671) ([cesar-rodriguez](https://github.com/cesar-rodriguez))
- Docs- fix argo image path.md [\#667](https://github.com/accurics/terrascan/pull/667) ([amirbenv](https://github.com/amirbenv))
- Makes saving of admission requests configurable via a config file option [\#665](https://github.com/accurics/terrascan/pull/665) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- Add authentication with API key for the /logs endpoint [\#663](https://github.com/accurics/terrascan/pull/663) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- Fixes docs format [\#661](https://github.com/accurics/terrascan/pull/661) ([cesar-rodriguez](https://github.com/cesar-rodriguez))
- Update mkdocs.yml [\#660](https://github.com/accurics/terrascan/pull/660) ([amirbenv](https://github.com/amirbenv))
- Support for authenticated tf module download [\#658](https://github.com/accurics/terrascan/pull/658) ([jlk](https://github.com/jlk))
- Fix - terraform complex variables are not getting resolved [\#657](https://github.com/accurics/terrascan/pull/657) ([patilpankaj212](https://github.com/patilpankaj212))
- Reorganized and Updated docs [\#655](https://github.com/accurics/terrascan/pull/655) ([amirbenv](https://github.com/amirbenv))
- Fix- panic when terraform list variable doesn't have a type [\#654](https://github.com/accurics/terrascan/pull/654) ([patilpankaj212](https://github.com/patilpankaj212))
- Fix panic for floating point variables [\#653](https://github.com/accurics/terrascan/pull/653) ([patilpankaj212](https://github.com/patilpankaj212))
- Adding support to scan IAC from atlantis workflow [\#648](https://github.com/accurics/terrascan/pull/648) ([jlk](https://github.com/jlk))
- Fix - k8s resources config data has absolute source paths for resources [\#644](https://github.com/accurics/terrascan/pull/644) ([patilpankaj212](https://github.com/patilpankaj212))
- Fix - terrascan not able to read modules within a subdirectory [\#641](https://github.com/accurics/terrascan/pull/641) ([patilpankaj212](https://github.com/patilpankaj212))
- Add /go/bin to PATH. [\#637](https://github.com/accurics/terrascan/pull/637) ([seancallaway](https://github.com/seancallaway))
- Update mkdocs-material to 7.1.0 [\#636](https://github.com/accurics/terrascan/pull/636) ([pyup-bot](https://github.com/pyup-bot))
- Fix windows tests [\#635](https://github.com/accurics/terrascan/pull/635) ([patilpankaj212](https://github.com/patilpankaj212))
- Fix kustomize scan breakage on windows [\#630](https://github.com/accurics/terrascan/pull/630) ([dev-gaur](https://github.com/dev-gaur))
- Update route53LoggingDisabled.rego to ignore private zones [\#626](https://github.com/accurics/terrascan/pull/626) ([matt-slalom](https://github.com/matt-slalom))
- Adding openssh for downloading modules via ssh [\#625](https://github.com/accurics/terrascan/pull/625) ([sachinar](https://github.com/sachinar))
- Fix - init behavior change [\#624](https://github.com/accurics/terrascan/pull/624) ([patilpankaj212](https://github.com/patilpankaj212))
- Add support for validating admission webhook in terrascan [\#620](https://github.com/accurics/terrascan/pull/620) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- Policy download refactor [\#618](https://github.com/accurics/terrascan/pull/618) ([dev-gaur](https://github.com/dev-gaur))
- Update mkdocs-material to 7.0.6 [\#615](https://github.com/accurics/terrascan/pull/615) ([pyup-bot](https://github.com/pyup-bot))
- Log error in LoadIacDir before continuing [\#613](https://github.com/accurics/terrascan/pull/613) ([jlk](https://github.com/jlk))
- K8S Risk Category Changes [\#608](https://github.com/accurics/terrascan/pull/608) ([Avanti19](https://github.com/Avanti19))
- GCP Risk Category Changes [\#606](https://github.com/accurics/terrascan/pull/606) ([shreyas-phansalkar-189](https://github.com/shreyas-phansalkar-189))
- Category flag e2e tests [\#605](https://github.com/accurics/terrascan/pull/605) ([gaurav-gogia](https://github.com/gaurav-gogia))
- Azure Risk Category Changes [\#604](https://github.com/accurics/terrascan/pull/604) ([gaurav-gogia](https://github.com/gaurav-gogia))
- AWS Risk Category Changes [\#603](https://github.com/accurics/terrascan/pull/603) ([harkirat22](https://github.com/harkirat22))
- Bugfix/revert policies [\#602](https://github.com/accurics/terrascan/pull/602) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- Server mode: take file extension from uploaded file [\#593](https://github.com/accurics/terrascan/pull/593) ([jlk](https://github.com/jlk))
- filepath fixes in e2e tests [\#591](https://github.com/accurics/terrascan/pull/591) ([patilpankaj212](https://github.com/patilpankaj212))
- Update mkdocs-material to 7.0.5 [\#590](https://github.com/accurics/terrascan/pull/590) ([pyup-bot](https://github.com/pyup-bot))
- update helm default chart name and namespace values [\#589](https://github.com/accurics/terrascan/pull/589) ([williepaul](https://github.com/williepaul))
- v1.4.0 doc updates [\#588](https://github.com/accurics/terrascan/pull/588) ([cesar-rodriguez](https://github.com/cesar-rodriguez))
- Terrascan K8s New categories and ruleRef ID changes [\#583](https://github.com/accurics/terrascan/pull/583) ([Avanti19](https://github.com/Avanti19))
- GCP Category Changes  [\#582](https://github.com/accurics/terrascan/pull/582) ([shreyas-phansalkar-189](https://github.com/shreyas-phansalkar-189))
- AWS new Categories [\#581](https://github.com/accurics/terrascan/pull/581) ([harkirat22](https://github.com/harkirat22))
- New Policies for Azure & Category Updates. [\#580](https://github.com/accurics/terrascan/pull/580) ([gaurav-gogia](https://github.com/gaurav-gogia))
- Terrascan init and config handling refactor [\#576](https://github.com/accurics/terrascan/pull/576) ([dev-gaur](https://github.com/dev-gaur))
- Feature: add options to specify desired categories of violations to be reported [\#547](https://github.com/accurics/terrascan/pull/547) ([gaurav-gogia](https://github.com/gaurav-gogia))

## [v1.4.0](https://github.com/accurics/terrascan/tree/v1.4.0) (2021-03-05)

[Full Changelog](https://github.com/accurics/terrascan/compare/v1.3.3...v1.4.0)

**Implemented enhancements:**

- Scanning terraform plan files [\#407](https://github.com/accurics/terrascan/issues/407)
- Adds support for junit xml output [\#527](https://github.com/accurics/terrascan/pull/527)
- Adds e2e test scenarios for help and scan command [\#564](https://github.com/accurics/terrascan/pull/564)
- Adds e2e tests for api server [\#585](https://github.com/accurics/terrascan/pull/585)
- Please checkout our new [Github Action!](https://github.com/marketplace/actions/terrascan-iac-scanner)

**Fixed bugs:**

- Fixed a few bugs in the init command and downloading of fresh policies, including [\#561](https://github.com/accurics/terrascan/issues/561)
- Difference in violated policies for the same terraform file [\#519](https://github.com/accurics/terrascan/issues/519)
- false positive for AWS.Instance.NetworkSecurity.Medium.0506 [\#404](https://github.com/accurics/terrascan/issues/404)
- accurics.gcp.IAM.122 needs to take into account the new name for Uniform bucket-level access flag [\#329](https://github.com/accurics/terrascan/issues/329)
- fix the 'repo already exist' bug and improve error logging for terrascan init [\#552](https://github.com/accurics/terrascan/pull/552) ([dev-gaur](https://github.com/dev-gaur))

**Closed issues:**

- terrascan API server's file scan always returns the resource config [\#578](https://github.com/accurics/terrascan/issues/578)
- Issue on Azure DevOps Agents since 1.3.2 : failed to initialize terrascan [\#561](https://github.com/accurics/terrascan/issues/561)
- Could not get terrascan init to work - would not download policy documents [\#551](https://github.com/accurics/terrascan/issues/551)

**Merged pull requests:**

- release 1.4.0 [\#586](https://github.com/accurics/terrascan/pull/586) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- adds e2e tests for api server [\#585](https://github.com/accurics/terrascan/pull/585) ([patilpankaj212](https://github.com/patilpankaj212))
- adds support to use 'config\_only' attribute in api server's file scan [\#579](https://github.com/accurics/terrascan/pull/579) ([patilpankaj212](https://github.com/patilpankaj212))
- adds support to display passed rules [\#572](https://github.com/accurics/terrascan/pull/572) ([patilpankaj212](https://github.com/patilpankaj212))
- Update mkdocs-material to 7.0.1 [\#567](https://github.com/accurics/terrascan/pull/567) ([pyup-bot](https://github.com/pyup-bot))
- fix filepaths and home directory lookup [\#566](https://github.com/accurics/terrascan/pull/566) ([dev-gaur](https://github.com/dev-gaur))
- adds e2e test scenarios for help and scan command [\#564](https://github.com/accurics/terrascan/pull/564) ([patilpankaj212](https://github.com/patilpankaj212))
- Adds support for scanning tfplan json file [\#562](https://github.com/accurics/terrascan/pull/562) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- fix: renamed the json file to remove spaces [\#560](https://github.com/accurics/terrascan/pull/560) ([harkirat22](https://github.com/harkirat22))
- fix: Changed the description message to handle the violation correctly [\#559](https://github.com/accurics/terrascan/pull/559) ([harkirat22](https://github.com/harkirat22))
- bump versions to v1.3.3 [\#558](https://github.com/accurics/terrascan/pull/558) ([dev-gaur](https://github.com/dev-gaur))
- updated go module files [\#557](https://github.com/accurics/terrascan/pull/557) ([dev-gaur](https://github.com/dev-gaur))
- Initial changes for e2e testing framework [\#553](https://github.com/accurics/terrascan/pull/553) ([patilpankaj212](https://github.com/patilpankaj212))
- Add code of conduct [\#545](https://github.com/accurics/terrascan/pull/545) ([jlk](https://github.com/jlk))
- Fixes incorrect description of RDS encryption policy [\#542](https://github.com/accurics/terrascan/pull/542) ([alex-petrov-vt](https://github.com/alex-petrov-vt))
- changes in log level and messages for load iac functions [\#541](https://github.com/accurics/terrascan/pull/541) ([patilpankaj212](https://github.com/patilpankaj212))
- Update mkdocs-material to 6.2.8 [\#539](https://github.com/accurics/terrascan/pull/539) ([pyup-bot](https://github.com/pyup-bot))
- Updates docs for v1.3.2 [\#537](https://github.com/accurics/terrascan/pull/537) ([cesar-rodriguez](https://github.com/cesar-rodriguez))
- update readme for v1.3.2 [\#534](https://github.com/accurics/terrascan/pull/534) ([dev-gaur](https://github.com/dev-gaur))
- fix - improved description for init command in help [\#532](https://github.com/accurics/terrascan/pull/532) ([nathannaveen](https://github.com/nathannaveen))
- Adds support for junit xml output [\#527](https://github.com/accurics/terrascan/pull/527) ([patilpankaj212](https://github.com/patilpankaj212))
- enhancement: scan terraform registry modules as remote type [\#513](https://github.com/accurics/terrascan/pull/513) ([patilpankaj212](https://github.com/patilpankaj212))
- support for terraform registry remote modules [\#505](https://github.com/accurics/terrascan/pull/505) ([patilpankaj212](https://github.com/patilpankaj212))
- feature: add options to specify desired severity level of violations to be reported [\#501](https://github.com/accurics/terrascan/pull/501) ([dev-gaur](https://github.com/dev-gaur))
- Bump github.com/spf13/cobra from 1.0.0 to 1.1.1 [\#493](https://github.com/accurics/terrascan/pull/493) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.3.2](https://github.com/accurics/terrascan/tree/v1.3.2) (2021-02-03)

[Full Changelog](https://github.com/accurics/terrascan/compare/v1.3.1...v1.3.2)

**Fixed bugs:**

- terrascan init should download new policies [\#521](https://github.com/accurics/terrascan/issues/521)

**Closed issues:**

- How to get rid of "Anonymous, public read access to a container and its blobs can be enabled in Azure Blob storage. This is only recommended if absolutely necessary." [\#405](https://github.com/accurics/terrascan/issues/405)
- False Positive for accurics.azure.NS.161 when Security Groups Association and Subnets are defined indepently from VNet [\#391](https://github.com/accurics/terrascan/issues/391)
- Calico is not supported as a valid Network Security for azurerm\_kubernetes\_cluster [\#376](https://github.com/accurics/terrascan/issues/376)

**Merged pull requests:**

- Update readme for v1.3.2 [\#534](https://github.com/accurics/terrascan/pull/534) ([dev-gaur](https://github.com/dev-gaur))
- bump terrascan version to v1.3.2 [\#533](https://github.com/accurics/terrascan/pull/533) ([dev-gaur](https://github.com/dev-gaur))
- refactor init command for robust policy download checks [\#531](https://github.com/accurics/terrascan/pull/531) ([dev-gaur](https://github.com/dev-gaur))
- terrascan init will download new policies. [\#529](https://github.com/accurics/terrascan/pull/529) ([dev-gaur](https://github.com/dev-gaur))
- bugfix: Checks for security group association defined independently from vnet  [\#526](https://github.com/accurics/terrascan/pull/526) ([harkirat22](https://github.com/harkirat22))
- Update mkdocs-material to 6.2.7 [\#524](https://github.com/accurics/terrascan/pull/524) ([pyup-bot](https://github.com/pyup-bot))
- Fixed typos in docs [\#523](https://github.com/accurics/terrascan/pull/523) ([gauravgahlot](https://github.com/gauravgahlot))
- Enhancement: new set of policies for AWS EC2 instance. [\#522](https://github.com/accurics/terrascan/pull/522) ([harkirat22](https://github.com/harkirat22))
- Harkirat22/bug fix [\#520](https://github.com/accurics/terrascan/pull/520) ([harkirat22](https://github.com/harkirat22))
- fixes \#376 [\#518](https://github.com/accurics/terrascan/pull/518) ([gaurav-gogia](https://github.com/gaurav-gogia))
- fixes \#405 [\#517](https://github.com/accurics/terrascan/pull/517) ([gaurav-gogia](https://github.com/gaurav-gogia))
- Policy/aws launch config [\#516](https://github.com/accurics/terrascan/pull/516) ([harkirat22](https://github.com/harkirat22))
- add support for pod container [\#515](https://github.com/accurics/terrascan/pull/515) ([harkirat22](https://github.com/harkirat22))
- Update mkdocs-material to 6.2.6 [\#514](https://github.com/accurics/terrascan/pull/514) ([pyup-bot](https://github.com/pyup-bot))
- Update README.md and changelog for 1.3.1 [\#509](https://github.com/accurics/terrascan/pull/509) ([amirbenv](https://github.com/amirbenv))

## [v1.3.1](https://github.com/accurics/terrascan/tree/v1.3.1) (2021-01-22)

[Full Changelog](https://github.com/accurics/terrascan/compare/v1.3.0...v1.3.1)

**Implemented enhancements:**

- Support for remote modules
- Tag container image with release version [\#504](https://github.com/accurics/terrascan/issues/504)

**Fixed bugs:**

- Build error on ARM MacOS
- terrascan consider  source = "terraform-aws-modules/vpc/aws"  as local path [\#418](https://github.com/accurics/terrascan/issues/418)
- Failed to read module directory [\#332](https://github.com/accurics/terrascan/issues/332)

**Closed issues:**

- Custom Variable Validation no longer experiemental in 0.13 [\#500](https://github.com/accurics/terrascan/issues/500)

**Merged pull requests:**

- release v1.3.1 [\#508](https://github.com/accurics/terrascan/pull/508) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- fix dependencies that were breaking the darwin/arm64 build [\#507](https://github.com/accurics/terrascan/pull/507) ([williepaul](https://github.com/williepaul))
- support for terraform registry remote modules [\#505](https://github.com/accurics/terrascan/pull/505) ([patilpankaj212](https://github.com/patilpankaj212))
- Readme rule supression [\#503](https://github.com/accurics/terrascan/pull/503) ([amirbenv](https://github.com/amirbenv))
- Bump github.com/hashicorp/go-retryablehttp from 0.6.6 to 0.6.8 [\#496](https://github.com/accurics/terrascan/pull/496) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/go-getter from 1.5.1 to 1.5.2 [\#495](https://github.com/accurics/terrascan/pull/495) ([dependabot[bot]](https://github.com/apps/dependabot))

## [v1.3.0](https://github.com/accurics/terrascan/tree/v1.3.0) (2021-01-19)

[Full Changelog](https://github.com/accurics/terrascan/compare/v1.2.0...v1.3.0)

**Implemented enhancements:**
- Prints output in human friendly format [\#168](https://github.com/accurics/terrascan/issues/168)
- Support for rule suppression using terraform comments,kubernetes annotations, cli arguments, and config file.
- New Policies for Kubernetes [\#480](https://github.com/accurics/terrascan/pull/480)
- Tag released Docker images [\#398](https://github.com/accurics/terrascan/issues/398)
- Add policy for checking insecure\_ssl configuration for github\_repository\_webhook in GitHub provider [\#355](https://github.com/accurics/terrascan/issues/355)
- Introduced support for terraform .14 and .13. Note: This will introduce some breaking changes for terraform v.12 files, even if using --iac-version v.12 flag. Notably we will no longer support multiple providers blocks, and certain references inside provisioner blocks (objects other than self, count or each, where when = destroy) . For more details see: https://github.com/hashicorp/terraform/releases/tag/v0.13.0

**Fixed bugs:**

- terrascan doesn't allow registering multiple versions for an iac-type [\#471](https://github.com/accurics/terrascan/issues/471)
- Debug resource lock [\#432](https://github.com/accurics/terrascan/issues/432)
- terrascan panic: not a string [\#412](https://github.com/accurics/terrascan/issues/412)
- False positive for aws rule vpcFlowLogsNotEnabled [\#408](https://github.com/accurics/terrascan/issues/408)
- accurics.GCP.EKM.132 and accurics.GCP.EKM.131 wrong violation using disk\_encryption\_key [\#382](https://github.com/accurics/terrascan/issues/382)
- s3EnforceUserACL - False Positive [\#359](https://github.com/accurics/terrascan/issues/359)
- How to fix accurics.azure.EKM.20 [\#331](https://github.com/accurics/terrascan/issues/331)
- Why accurics.gcp.IAM.104 suggests enabling a client certificate? [\#330](https://github.com/accurics/terrascan/issues/330)

**Closed issues:**

- terraform can't detect violations in terraform modules [\#468](https://github.com/accurics/terrascan/issues/468)
- uniformBucketEnabled.rego referencing deprecated config [\#453](https://github.com/accurics/terrascan/issues/453)
- Unable to run terrascan scan [\#446](https://github.com/accurics/terrascan/issues/446)
- Terrascan doesn't exit with error on CLI or Parsing errors. [\#442](https://github.com/accurics/terrascan/issues/442)
- Terrascan Failure When Using Terraform 13 + Variable Validation [\#426](https://github.com/accurics/terrascan/issues/426)
- Update policy example in documentation to use latest GitHub implementation [\#422](https://github.com/accurics/terrascan/issues/422)
- Fix link to repo playground in policies documentation [\#421](https://github.com/accurics/terrascan/issues/421)
- terrascan scan  crashes with runtime: goroutine stack exceeds 1000000000-byte limit [\#406](https://github.com/accurics/terrascan/issues/406)
- Typo error in the terrascan Architecture page [\#403](https://github.com/accurics/terrascan/issues/403)
- accurics.gcp.OPS.114 should also check for cos\_containerd image [\#395](https://github.com/accurics/terrascan/issues/395)
- accurics.gcp.NS.112 suggest basic auth is enabled when is not [\#394](https://github.com/accurics/terrascan/issues/394)
- Test coverage missing for kustomize iac-provider [\#379](https://github.com/accurics/terrascan/issues/379)
- Why is vpcFlowLogsNotEnabled determined to be a violation? [\#352](https://github.com/accurics/terrascan/issues/352)

**Merged pull requests:**

- update version to  v1.3.0 [\#502](https://github.com/accurics/terrascan/pull/502) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- Add v13 flag option for terraform iac [\#499](https://github.com/accurics/terrascan/pull/499) ([dev-gaur](https://github.com/dev-gaur))
- Fix: potential bug added in PR \#470 [\#497](https://github.com/accurics/terrascan/pull/497) ([dev-gaur](https://github.com/dev-gaur))
- Bump sigs.k8s.io/kustomize/api from 0.7.1 to 0.7.2 [\#494](https://github.com/accurics/terrascan/pull/494) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/mattn/go-isatty from 0.0.8 to 0.0.12 [\#492](https://github.com/accurics/terrascan/pull/492) ([dependabot[bot]](https://github.com/apps/dependabot))
- solves issue \#382, and improved policy to relate disk with the instance [\#490](https://github.com/accurics/terrascan/pull/490) ([harkirat22](https://github.com/harkirat22))
- solves issue \#331 [\#489](https://github.com/accurics/terrascan/pull/489) ([harkirat22](https://github.com/harkirat22))
- Update mkdocs-material to 6.2.5 [\#488](https://github.com/accurics/terrascan/pull/488) ([pyup-bot](https://github.com/pyup-bot))
- Bump go.uber.org/zap from 1.13.0 to 1.16.0 [\#486](https://github.com/accurics/terrascan/pull/486) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/spf13/afero from 1.3.4 to 1.5.1 [\#485](https://github.com/accurics/terrascan/pull/485) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/iancoleman/strcase from 0.1.1 to 0.1.3 [\#484](https://github.com/accurics/terrascan/pull/484) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/hashicorp/go-version from 1.2.0 to 1.2.1 [\#482](https://github.com/accurics/terrascan/pull/482) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/pelletier/go-toml from 1.8.0 to 1.8.1 [\#481](https://github.com/accurics/terrascan/pull/481) ([dependabot[bot]](https://github.com/apps/dependabot))
- Policy update 2021 01 14 [\#480](https://github.com/accurics/terrascan/pull/480) ([williepaul](https://github.com/williepaul))
- fix panic for list variables [\#479](https://github.com/accurics/terrascan/pull/479) ([patilpankaj212](https://github.com/patilpankaj212))
- adding an else condition to relate management lock with resource group [\#476](https://github.com/accurics/terrascan/pull/476) ([harkirat22](https://github.com/harkirat22))
- adding an else condition to relate the flow log with vpc [\#475](https://github.com/accurics/terrascan/pull/475) ([harkirat22](https://github.com/harkirat22))
- including a check for verifying in-line policy is included  [\#474](https://github.com/accurics/terrascan/pull/474) ([harkirat22](https://github.com/harkirat22))
- adding rule to check if waf is enabled at cloud front distribution [\#473](https://github.com/accurics/terrascan/pull/473) ([harkirat22](https://github.com/harkirat22))
- Added terraform v14 support besides v12. [\#470](https://github.com/accurics/terrascan/pull/470) ([dev-gaur](https://github.com/dev-gaur))
- support comment with rule skipping for resource and scan summary modifications [\#466](https://github.com/accurics/terrascan/pull/466) ([patilpankaj212](https://github.com/patilpankaj212))
- recognize metadata.generateName [\#465](https://github.com/accurics/terrascan/pull/465) ([acc-jon](https://github.com/acc-jon))
- Update mkdocs-material to 6.2.4 [\#464](https://github.com/accurics/terrascan/pull/464) ([pyup-bot](https://github.com/pyup-bot))
- Update README.md [\#463](https://github.com/accurics/terrascan/pull/463) ([amirbenv](https://github.com/amirbenv))
- Deprecated gcs bucket [\#462](https://github.com/accurics/terrascan/pull/462) ([jdyke](https://github.com/jdyke))
- changed the description to include the vulnerable versions [\#460](https://github.com/accurics/terrascan/pull/460) ([harkirat22](https://github.com/harkirat22))
- Fix exit code on error [\#458](https://github.com/accurics/terrascan/pull/458) ([patilpankaj212](https://github.com/patilpankaj212))
- policy for CVE-2020-8555 [\#457](https://github.com/accurics/terrascan/pull/457) ([harkirat22](https://github.com/harkirat22))
- Update README.md [\#456](https://github.com/accurics/terrascan/pull/456) ([amirbenv](https://github.com/amirbenv))
- rule skipping for resources in k8s [\#455](https://github.com/accurics/terrascan/pull/455) ([patilpankaj212](https://github.com/patilpankaj212))
- terrascan argo-cd instructions [\#454](https://github.com/accurics/terrascan/pull/454) ([storebot](https://github.com/storebot))
- Adds CI/CD integration docs [\#452](https://github.com/accurics/terrascan/pull/452) ([cesar-rodriguez](https://github.com/cesar-rodriguez))
- Bump github.com/zclconf/go-cty from 1.2.1 to 1.7.1 [\#449](https://github.com/accurics/terrascan/pull/449) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump sigs.k8s.io/kustomize/api from 0.6.5 to 0.7.1 [\#448](https://github.com/accurics/terrascan/pull/448) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/gorilla/mux from 1.7.4 to 1.8.0 [\#447](https://github.com/accurics/terrascan/pull/447) ([dependabot[bot]](https://github.com/apps/dependabot))
- Update mkdocs-material to 6.2.3 [\#445](https://github.com/accurics/terrascan/pull/445) ([pyup-bot](https://github.com/pyup-bot))
- deps: add dependabot support [\#444](https://github.com/accurics/terrascan/pull/444) ([chenrui333](https://github.com/chenrui333))
- bump go to 1.15 [\#443](https://github.com/accurics/terrascan/pull/443) ([chenrui333](https://github.com/chenrui333))
- implement scan and skip rules [\#441](https://github.com/accurics/terrascan/pull/441) ([patilpankaj212](https://github.com/patilpankaj212))
- scan command refactor [\#436](https://github.com/accurics/terrascan/pull/436) ([patilpankaj212](https://github.com/patilpankaj212))
- Fixes dead link to old getting started page [\#435](https://github.com/accurics/terrascan/pull/435) ([cesar-rodriguez](https://github.com/cesar-rodriguez))
- Add support to extract rules to skip from terraform comments [\#434](https://github.com/accurics/terrascan/pull/434) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- bash output improvements [\#431](https://github.com/accurics/terrascan/pull/431) ([patilpankaj212](https://github.com/patilpankaj212))
- APE-1319: Revamped Getting Started Section [\#430](https://github.com/accurics/terrascan/pull/430) ([acc-jon](https://github.com/acc-jon))
- Add policy AC-K8-NS-SE-M-0188 for CVE-2020-8554 [\#428](https://github.com/accurics/terrascan/pull/428) ([gauravgogia-accurics](https://github.com/gauravgogia-accurics))
- set console mode on windows so colors render [\#427](https://github.com/accurics/terrascan/pull/427) ([acc-jon](https://github.com/acc-jon))
- Update mkdocs-material to 6.1.7 [\#425](https://github.com/accurics/terrascan/pull/425) ([pyup-bot](https://github.com/pyup-bot))
- Update policy example in the documentation [\#424](https://github.com/accurics/terrascan/pull/424) ([HorizonNet](https://github.com/HorizonNet))
- Fix link to rego playground in policies documentation [\#423](https://github.com/accurics/terrascan/pull/423) ([HorizonNet](https://github.com/HorizonNet))
- hopefully remove test failures due to non-deterministic comparisons [\#420](https://github.com/accurics/terrascan/pull/420) ([acc-jon](https://github.com/acc-jon))
- IMDSv1 policy: update category, description [\#419](https://github.com/accurics/terrascan/pull/419) ([acc-jon](https://github.com/acc-jon))
- IMDSv1 check policy [\#417](https://github.com/accurics/terrascan/pull/417) ([harkirat22](https://github.com/harkirat22))
- Add Docker image release tagging on release [\#410](https://github.com/accurics/terrascan/pull/410) ([HorizonNet](https://github.com/HorizonNet))
- Fix typo in architecture documentation [\#409](https://github.com/accurics/terrascan/pull/409) ([HorizonNet](https://github.com/HorizonNet))
- accurics.gcp.IAM.104 Fire rule when client certificate is enabled [\#402](https://github.com/accurics/terrascan/pull/402) ([lucas-giaco](https://github.com/lucas-giaco))
- Update mkdocs-material to 6.1.6 [\#401](https://github.com/accurics/terrascan/pull/401) ([pyup-bot](https://github.com/pyup-bot))
- Added Unit test coverage for Kustomize V3 Iac-provider [\#399](https://github.com/accurics/terrascan/pull/399) ([dev-gaur](https://github.com/dev-gaur))
- Fixes GCP cos node image policy [\#397](https://github.com/accurics/terrascan/pull/397) ([cesar-rodriguez](https://github.com/cesar-rodriguez))
- \#394: recognize that empty values for username and password in master… [\#396](https://github.com/accurics/terrascan/pull/396) ([acc-jon](https://github.com/acc-jon))
- Fix infinite loop on variable resolution [\#393](https://github.com/accurics/terrascan/pull/393) ([dinedal](https://github.com/dinedal))
- Remove demo badge [\#389](https://github.com/accurics/terrascan/pull/389) ([kklin](https://github.com/kklin))
- Update mkdocs-material to 6.1.5 [\#387](https://github.com/accurics/terrascan/pull/387) ([pyup-bot](https://github.com/pyup-bot))

## [v1.2.0](https://github.com/accurics/terrascan/tree/v1.2.0) (2020-11-16)

[Full Changelog](https://github.com/accurics/terrascan/compare/v1.1.0...v1.2.0)

**Implemented enhancements:**

- Add support for Helm [\#353](https://github.com/accurics/terrascan/issues/353)
- Add 'git' to container image, or run container as 'root' user by default [\#349](https://github.com/accurics/terrascan/issues/349)
- Add policy for checking insecure\_ssl configuration for github\_organization\_webhook in GitHub provider [\#339](https://github.com/accurics/terrascan/issues/339)
- Rule for github\_repository seems to be wrongly placed under gcp [\#325](https://github.com/accurics/terrascan/issues/325)

**Fixed bugs:**

- Fail to validate when there are multiple properties with the same name in a resource [\#1](https://github.com/accurics/terrascan/issues/1)

**Closed issues:**

- Deep modules location mis-proccessed.  [\#365](https://github.com/accurics/terrascan/issues/365)
- 20MB binary file included in repo now [\#364](https://github.com/accurics/terrascan/issues/364)
- Private GitHub repositories are not recognized with version 3.0.0+ of GitHub provider [\#326](https://github.com/accurics/terrascan/issues/326)
- Terrascan -var-file=../another dir [\#144](https://github.com/accurics/terrascan/issues/144)
- Error in test\_aws\_security\_group\_inline\_rule\_open and test\_aws\_security\_group\_rule\_open [\#138](https://github.com/accurics/terrascan/issues/138)
- Intial setup after installation [\#136](https://github.com/accurics/terrascan/issues/136)
- Add support for data sources [\#3](https://github.com/accurics/terrascan/issues/3)
- Support from modules [\#2](https://github.com/accurics/terrascan/issues/2)

**Merged pull requests:**

- Bring Go to 1.15 in Github Actions [\#384](https://github.com/accurics/terrascan/pull/384) ([gliptak](https://github.com/gliptak))
- Bring Go to 1.15 in Github Actions [\#383](https://github.com/accurics/terrascan/pull/383) ([gliptak](https://github.com/gliptak))
- fix a bug when rendering subcharts [\#381](https://github.com/accurics/terrascan/pull/381) ([williepaul](https://github.com/williepaul))
- Added kustomize support [\#378](https://github.com/accurics/terrascan/pull/378) ([dev-gaur](https://github.com/dev-gaur))
- Adds support for Helm v3 [\#377](https://github.com/accurics/terrascan/pull/377) ([williepaul](https://github.com/williepaul))
- Update mkdocs-material to 6.1.4 [\#374](https://github.com/accurics/terrascan/pull/374) ([pyup-bot](https://github.com/pyup-bot))
- properly handle nested submodules \(\#365\) [\#373](https://github.com/accurics/terrascan/pull/373) ([acc-jon](https://github.com/acc-jon))
- Address \#365 by properly handling submodule path [\#372](https://github.com/accurics/terrascan/pull/372) ([acc-jon](https://github.com/acc-jon))
- Update mkdocs-material to 6.1.3 [\#371](https://github.com/accurics/terrascan/pull/371) ([pyup-bot](https://github.com/pyup-bot))
- Update mkdocs-material to 6.1.2 [\#370](https://github.com/accurics/terrascan/pull/370) ([pyup-bot](https://github.com/pyup-bot))
- Allow use of multiple policy types \(scan -t x,y or scan -t x -t y\) [\#368](https://github.com/accurics/terrascan/pull/368) ([acc-jon](https://github.com/acc-jon))
- Remove large binary that was included in the repo [\#366](https://github.com/accurics/terrascan/pull/366) ([cesar-rodriguez](https://github.com/cesar-rodriguez))
- fix send request method, previously hardcoded [\#361](https://github.com/accurics/terrascan/pull/361) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- Add git binary to terrascan docker image, required by downloader [\#360](https://github.com/accurics/terrascan/pull/360) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- Adds new policies/regos for AWS serverless services  [\#357](https://github.com/accurics/terrascan/pull/357) ([cesar-rodriguez](https://github.com/cesar-rodriguez))
- Update mkdocs-material to 6.1.0 [\#356](https://github.com/accurics/terrascan/pull/356) ([pyup-bot](https://github.com/pyup-bot))
- Allow configuration of global policy config, fix some typos [\#354](https://github.com/accurics/terrascan/pull/354) ([acc-jon](https://github.com/acc-jon))
- Feature/support resolve variable references [\#351](https://github.com/accurics/terrascan/pull/351) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- Add new policy for checking insecure\_ssl on github\_organization\_webhook [\#350](https://github.com/accurics/terrascan/pull/350) ([HorizonNet](https://github.com/HorizonNet))
- Update mkdocs-material to 6.0.2 [\#348](https://github.com/accurics/terrascan/pull/348) ([pyup-bot](https://github.com/pyup-bot))
- Add support for colorized output [\#347](https://github.com/accurics/terrascan/pull/347) ([acc-jon](https://github.com/acc-jon))
- Update mkdocs-material to 6.0.1 [\#346](https://github.com/accurics/terrascan/pull/346) ([pyup-bot](https://github.com/pyup-bot))
- Adds support for remote Terraform modules and scanning remotely for other IaC tools [\#345](https://github.com/accurics/terrascan/pull/345) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- fix supported providers unit test, sort the wanted result [\#344](https://github.com/accurics/terrascan/pull/344) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- Fix typo on AWS IAM account password policy rego name [\#343](https://github.com/accurics/terrascan/pull/343) ([kmonticolo](https://github.com/kmonticolo))
- Update mkdocs-material to 5.5.14 [\#340](https://github.com/accurics/terrascan/pull/340) ([pyup-bot](https://github.com/pyup-bot))
- Adds docs section for GitHub policies [\#337](https://github.com/accurics/terrascan/pull/337) ([cesar-rodriguez](https://github.com/cesar-rodriguez))
- Automatically populate usage with supported IaC providers, versions, and policies [\#336](https://github.com/accurics/terrascan/pull/336) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- Add line about kubernetes YAML/JSON support [\#335](https://github.com/accurics/terrascan/pull/335) ([williepaul](https://github.com/williepaul))
- Add policy set for GitHub provider [\#334](https://github.com/accurics/terrascan/pull/334) ([HorizonNet](https://github.com/HorizonNet))
- Add check for visibility for github\_repository [\#333](https://github.com/accurics/terrascan/pull/333) ([HorizonNet](https://github.com/HorizonNet))
- Add instructions for booting terrascan demo [\#319](https://github.com/accurics/terrascan/pull/319) ([kklin](https://github.com/kklin))

## [v1.1.0](https://github.com/accurics/terrascan/tree/v1.1.0) (2020-09-16)

[Full Changelog](https://github.com/accurics/terrascan/compare/v1.0.0...v1.1.0)

**Implemented enhancements:**

- Initial kubernetes support [\#313](https://github.com/accurics/terrascan/pull/313) ([williepaul](https://github.com/williepaul))
- Adds different exit code when issues are found [\#299](https://github.com/accurics/terrascan/pull/299) ([cesar-rodriguez](https://github.com/cesar-rodriguez))
- Adding terrascan to Homebrew [\#293](https://github.com/accurics/terrascan/issues/293)

**Fixed bugs:**

- Oudated Docker image [\#294](https://github.com/accurics/terrascan/issues/294)
- Error with XML output [\#290](https://github.com/accurics/terrascan/issues/290)
- Fixed checkIpForward rule \(gcp\) [\#323](https://github.com/accurics/terrascan/pull/323) ([williepaul](https://github.com/williepaul))

**Closed issues:**

- Terrascan wrongly reports a accurics.gcp.NS.130 \(checkIpForward\) violation [\#320](https://github.com/accurics/terrascan/issues/320)
- Allow structure output \(Json\) [\#252](https://github.com/accurics/terrascan/issues/252)
- Throwing Errors when parsing nested brackets in HCL [\#233](https://github.com/accurics/terrascan/issues/233)
- Be able to generate xml/html reports [\#119](https://github.com/accurics/terrascan/issues/119)

**Merged pull requests:**

- Revert "fixed a bug in checkIpForward" [\#322](https://github.com/accurics/terrascan/pull/322) ([cesar-rodriguez](https://github.com/cesar-rodriguez))
- Fixed a bug in checkIpForward [\#321](https://github.com/accurics/terrascan/pull/321) ([williepaul](https://github.com/williepaul))
- Move server command out of ENTRYPOINT and into CMD [\#318](https://github.com/accurics/terrascan/pull/318) ([williepaul](https://github.com/williepaul))
- Send logs to stderr instead of stdout [\#317](https://github.com/accurics/terrascan/pull/317) ([williepaul](https://github.com/williepaul))
- Fix template rendering bug [\#316](https://github.com/accurics/terrascan/pull/316) ([williepaul](https://github.com/williepaul))
- chore\(docs\): add homebrew installation [\#315](https://github.com/accurics/terrascan/pull/315) ([chenrui333](https://github.com/chenrui333))
- Update badges in readme [\#314](https://github.com/accurics/terrascan/pull/314) ([acc-jon](https://github.com/acc-jon))
- Update mkdocs-diagrams to 1.0.0 [\#312](https://github.com/accurics/terrascan/pull/312) ([pyup-bot](https://github.com/pyup-bot))
- Add support to print resource config as an output [\#309](https://github.com/accurics/terrascan/pull/309) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- Manage relative module path [\#308](https://github.com/accurics/terrascan/pull/308) ([guilhem](https://github.com/guilhem))
- Update mkdocs-material to 5.5.12 [\#307](https://github.com/accurics/terrascan/pull/307) ([pyup-bot](https://github.com/pyup-bot))
- chore\(docs\): fix indent of tar extraction [\#306](https://github.com/accurics/terrascan/pull/306) ([zmarouf](https://github.com/zmarouf))
- Fixes issue template and rego capitalization [\#301](https://github.com/accurics/terrascan/pull/301) ([cesar-rodriguez](https://github.com/cesar-rodriguez))
- Update mkdocs-material to 5.5.8 [\#300](https://github.com/accurics/terrascan/pull/300) ([pyup-bot](https://github.com/pyup-bot))
- Update about.md [\#298](https://github.com/accurics/terrascan/pull/298) ([Upa-acc](https://github.com/Upa-acc))
- Updated policies to the latest set [\#297](https://github.com/accurics/terrascan/pull/297) ([williepaul](https://github.com/williepaul))
- Fixes docker latest tag [\#296](https://github.com/accurics/terrascan/pull/296) ([cesar-rodriguez](https://github.com/cesar-rodriguez))
- Typo fixes [\#295](https://github.com/accurics/terrascan/pull/295) ([erichs](https://github.com/erichs))
- Update mkdocs-material to 5.5.7 [\#292](https://github.com/accurics/terrascan/pull/292) ([pyup-bot](https://github.com/pyup-bot))
- Fix xml output [\#291](https://github.com/accurics/terrascan/pull/291) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))

## 1.0.0 (2020-08-16)
Major updates to Terrascan and the underlying architecture including:

- Pluggable architecture written in Golang. We updated the architecture to be easier to extend Terrascan with additional IaC languages and support policies for different cloud providers and cloud native tooling.
- Server mode. This allows Terrascan to be executed as a server and use it's API to perform static code analysis
- Notifications hooks. Will be able to integrate for notifications to external systems (e.g. email, slack, etc.)
- Uses OPA policy engine and policies written in Rego.

## 0.2.3 (2020-07-23)
- Introduces the '-f' flag for passing a list of ".tf" files for linting and the '--version' flag.

## 0.2.2 (2020-07-21)
- Adds Docker image and pipeline to push to DockerHub

## 0.2.1 (2020-06-19)
- Bugfix: The pyhcl hard dependency in the requirements.txt file caused issues if a higher version was installed. This was fixed by using the ">=" operator.

## 0.2.0 (2020-01-11)
- Adds support for terraform 0.12+

## 0.1.2 (2020-01-05)
- Adds ability to setup terrascan as a pre-commit hook

## 0.1.1 (2020-01-01)
- Updates dependent packages to latest versions
- Migrates CI to GitHub Actions from travis

## 0.1.0 (2017-11-26)
- First release on PyPI.

\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
