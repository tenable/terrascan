# Changelog

## [v1.10.0](https://github.com/accurics/terrascan/tree/v1.10.0) (2021-08-24)

[Full Changelog](https://github.com/accurics/terrascan/compare/v1.9.0...v1.10.0)

**Implemented enhancements:**

- Add capability to extract references to container images in K8s ecosystem IaC [\#881](https://github.com/accurics/terrascan/issues/881)

**Fixed bugs:**

- Terrascan does not exit with error code in pipeline or CLI [\#950](https://github.com/accurics/terrascan/issues/950)

**Closed issues:**

- Links are Not formatted Properly in Contributor Doc [\#969](https://github.com/accurics/terrascan/issues/969)
- Enabling dependabot or renovate for automatic dependency update [\#959](https://github.com/accurics/terrascan/issues/959)
- AC\_K8S\_0131 triggers on a Namespace resource [\#957](https://github.com/accurics/terrascan/issues/957)
- Integrity issue with Kustomize v4 support [\#956](https://github.com/accurics/terrascan/issues/956)
- Add Support For ECR [\#927](https://github.com/accurics/terrascan/issues/927)
- Add capability to extract references to container images in terraform [\#898](https://github.com/accurics/terrascan/issues/898)
- Kustomize support says v3 but is actually v4 [\#891](https://github.com/accurics/terrascan/issues/891)

**Merged pull requests:**

- Extract images from Dockerfiles [\#1002](https://github.com/accurics/terrascan/pull/1002) ([nasir-rabbani](https://github.com/nasir-rabbani))
- Revert "update resource type to map\[string\]bool" [\#1001](https://github.com/accurics/terrascan/pull/1001) ([cesar-rodriguez](https://github.com/cesar-rodriguez))
- Upgrade helm.sh/helm/v3 to version 3.6.1 [\#1000](https://github.com/accurics/terrascan/pull/1000) ([patilpankaj212](https://github.com/patilpankaj212))
- Bump github.com/pelletier/go-toml from 1.8.1 to 1.9.3 [\#999](https://github.com/accurics/terrascan/pull/999) ([dependabot[bot]](https://github.com/apps/dependabot))
- Adds additional policies for dockerfile [\#996](https://github.com/accurics/terrascan/pull/996) ([pavniii](https://github.com/pavniii))
- terrascan should exit with non zero exit code when scan error are present [\#994](https://github.com/accurics/terrascan/pull/994) ([patilpankaj212](https://github.com/patilpankaj212))
- Bump github.com/hashicorp/go-getter from 1.5.2 to 1.5.7 [\#993](https://github.com/accurics/terrascan/pull/993) ([dependabot[bot]](https://github.com/apps/dependabot))
- update resource type to map\[string\]bool [\#992](https://github.com/accurics/terrascan/pull/992) ([patilpankaj212](https://github.com/patilpankaj212))
- docs: fixes links in contributing documentation [\#990](https://github.com/accurics/terrascan/pull/990) ([cesar-rodriguez](https://github.com/cesar-rodriguez))
- Initial support for image scanning [\#989](https://github.com/accurics/terrascan/pull/989) ([Rchanger](https://github.com/Rchanger))
- added binary based support for kustomize v2 and v3 [\#988](https://github.com/accurics/terrascan/pull/988) ([nasir-rabbani](https://github.com/nasir-rabbani))
- Docs: adds brew instructions to release checklist [\#987](https://github.com/accurics/terrascan/pull/987) ([cesar-rodriguez](https://github.com/cesar-rodriguez))
- Update mkdocs-material to 7.2.4 [\#985](https://github.com/accurics/terrascan/pull/985) ([pyup-bot](https://github.com/pyup-bot))
- modify wait logic for service account creation in e2e validating webhook test [\#979](https://github.com/accurics/terrascan/pull/979) ([patilpankaj212](https://github.com/patilpankaj212))
- Update mkdocs-material to 7.2.3 [\#978](https://github.com/accurics/terrascan/pull/978) ([pyup-bot](https://github.com/pyup-bot))
- Bump github.com/hashicorp/hcl/v2 from 2.10.0 to 2.10.1 [\#972](https://github.com/accurics/terrascan/pull/972) ([dependabot[bot]](https://github.com/apps/dependabot))
- Bump github.com/spf13/afero from 1.5.1 to 1.6.0 [\#970](https://github.com/accurics/terrascan/pull/970) ([dependabot[bot]](https://github.com/apps/dependabot))
- Adds: e2e test for docker IaC provider [\#968](https://github.com/accurics/terrascan/pull/968) ([Rchanger](https://github.com/Rchanger))
- Fix dependency issue that caused dependabot to fail [\#966](https://github.com/accurics/terrascan/pull/966) ([patilpankaj212](https://github.com/patilpankaj212))
- fix\(policies\): removing false-positive for K8s namespaces [\#961](https://github.com/accurics/terrascan/pull/961) ([danmx](https://github.com/danmx))
- Extract Docker images from Terraform templates [\#937](https://github.com/accurics/terrascan/pull/937) ([dev-gaur](https://github.com/dev-gaur))
- Fixes supported Kustomize version \(should be v4\) [\#932](https://github.com/accurics/terrascan/pull/932) ([dev-gaur](https://github.com/dev-gaur))
- Extract Docker images from k8s YAML files [\#905](https://github.com/accurics/terrascan/pull/905) ([dev-gaur](https://github.com/dev-gaur))

## [v1.9.0](https://github.com/accurics/terrascan/tree/v1.9.0) (2021-08-06)

[Full Changelog](https://github.com/accurics/terrascan/compare/v1.8.1...v1.9.0)

**Implemented enhancements:**

- Dockerfile Support [\#798](https://github.com/accurics/terrascan/issues/798)
- pre-commit hook [\#311](https://github.com/accurics/terrascan/issues/311)
- Add support for CFT nested stacks [\#949](https://github.com/accurics/terrascan/pull/949)
- Adds support for using Terraform modules cached locally [\#940](https://github.com/accurics/terrascan/pull/940)

**Fixed bugs:**

- Helm chart scans use only 4 policies [\#946](https://github.com/accurics/terrascan/issues/946)

**Closed issues:**

- Link to docks in README  [\#944](https://github.com/accurics/terrascan/issues/944)
- Ensure remote modules are downloaded only once [\#936](https://github.com/accurics/terrascan/issues/936)
- Rule supression for specific resources [\#868](https://github.com/accurics/terrascan/issues/868)

**Merged pull requests:**

- Fixes k8s policy filtering [\#963](https://github.com/accurics/terrascan/pull/963) ([patilpankaj212](https://github.com/patilpankaj212))
- Update mkdocs-material to 7.2.2 [\#954](https://github.com/accurics/terrascan/pull/954) ([pyup-bot](https://github.com/pyup-bot))
- Adds Terrascan pre-commit [\#953](https://github.com/accurics/terrascan/pull/953) ([mihirhasan](https://github.com/mihirhasan))
- Add support for CFT nested stacks [\#949](https://github.com/accurics/terrascan/pull/949) ([sigmabaryon](https://github.com/sigmabaryon))
- fix - remote repo scan with config only option generates panic [\#948](https://github.com/accurics/terrascan/pull/948) ([patilpankaj212](https://github.com/patilpankaj212))
- Update mkdocs-material to 7.2.1 [\#947](https://github.com/accurics/terrascan/pull/947) ([pyup-bot](https://github.com/pyup-bot))
- Update README.md [\#945](https://github.com/accurics/terrascan/pull/945) ([sangam14](https://github.com/sangam14))
- Update helm chart progress checklist [\#943](https://github.com/accurics/terrascan/pull/943) ([dev-gaur](https://github.com/dev-gaur))
- Adds support for using Terraform modules cached locally [\#940](https://github.com/accurics/terrascan/pull/940) ([Rchanger](https://github.com/Rchanger))
- Update mkdocs-material to 7.2.0 [\#939](https://github.com/accurics/terrascan/pull/939) ([pyup-bot](https://github.com/pyup-bot))
- Dockerfile support  [\#849](https://github.com/accurics/terrascan/pull/849) ([Rchanger](https://github.com/Rchanger))

## [v1.8.1](https://github.com/accurics/terrascan/tree/v1.8.1) (2021-07-22)

[Full Changelog](https://github.com/accurics/terrascan/compare/v1.8.0...v1.8.1)

**Closed issues:**

- terrascan init should not be triggered if the user only wants to generate normalised json. [\#926](https://github.com/accurics/terrascan/issues/926)
- No rules are processed in GitlabCI [\#925](https://github.com/accurics/terrascan/issues/925)
- Scanning remote modules doesn't have same results as for scanning Terraform plan itself [\#923](https://github.com/accurics/terrascan/issues/923)
- Module AWS.KMS.Logging.High.0400 seems to serve no purpose [\#917](https://github.com/accurics/terrascan/issues/917)
- Secure ciphers are not used in CloudFront distribution [\#875](https://github.com/accurics/terrascan/issues/875)
- Correct point in time recovery for DynamoDB still leads to violation [\#838](https://github.com/accurics/terrascan/issues/838)

**Merged pull requests:**

- fix go mod files [\#941](https://github.com/accurics/terrascan/pull/941) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- Update mkdocs-material to 7.1.11 [\#938](https://github.com/accurics/terrascan/pull/938) ([pyup-bot](https://github.com/pyup-bot))
- Update mkdocs to 1.2.2 [\#935](https://github.com/accurics/terrascan/pull/935) ([pyup-bot](https://github.com/pyup-bot))
- K8s Policy to detect a service type Loadbalancer without a selector [\#931](https://github.com/accurics/terrascan/pull/931) ([harkirat22](https://github.com/harkirat22))
- Fix \#926: Do not initiate policy engine incase of --config-only flag [\#930](https://github.com/accurics/terrascan/pull/930) ([dev-gaur](https://github.com/dev-gaur))
- Update mkdocs-material to 7.1.10 [\#929](https://github.com/accurics/terrascan/pull/929) ([pyup-bot](https://github.com/pyup-bot))
- fix\(sws/cloudfront\): wrong check tls version [\#928](https://github.com/accurics/terrascan/pull/928) ([frediana](https://github.com/frediana))
- fixes: broken doc links [\#921](https://github.com/accurics/terrascan/pull/921) ([Rchanger](https://github.com/Rchanger))
- update getting started and Usage, fix links [\#920](https://github.com/accurics/terrascan/pull/920) ([amirbenv](https://github.com/amirbenv))
- Update overview.md [\#919](https://github.com/accurics/terrascan/pull/919) ([sangam14](https://github.com/sangam14))
- Remove unnecessary KMS deletion window code [\#918](https://github.com/accurics/terrascan/pull/918) ([matt-slalom](https://github.com/matt-slalom))
- minor-doc-fix [\#916](https://github.com/accurics/terrascan/pull/916) ([amirbenv](https://github.com/amirbenv))
- fix confusing error log message [\#914](https://github.com/accurics/terrascan/pull/914) ([dev-gaur](https://github.com/dev-gaur))
- add integrations overview and minor fixes [\#913](https://github.com/accurics/terrascan/pull/913) ([amirbenv](https://github.com/amirbenv))
- Updating the dax cluster policy [\#909](https://github.com/accurics/terrascan/pull/909) ([shreyas-phansalkar-189](https://github.com/shreyas-phansalkar-189))
- add github-sarif writer for github suited sarif output [\#907](https://github.com/accurics/terrascan/pull/907) ([dev-gaur](https://github.com/dev-gaur))
- Add support for arm linked templates [\#903](https://github.com/accurics/terrascan/pull/903) ([sigmabaryon](https://github.com/sigmabaryon))
- terraform 0.15 support [\#860](https://github.com/accurics/terrascan/pull/860) ([dev-gaur](https://github.com/dev-gaur))

## [v1.8.0](https://github.com/accurics/terrascan/tree/v1.8.0) (2021-07-02)

[Full Changelog](https://github.com/accurics/terrascan/compare/v1.7.0...v1.8.0)

**Implemented enhancements:**

- Add Support for new reference id field [\#786](https://github.com/accurics/terrascan/issues/786)

**Fixed bugs:**

- Sarif output has wrong file path value for file scans [\#861](https://github.com/accurics/terrascan/issues/861)
- 'k8s' key updated multiple times in policy package [\#439](https://github.com/accurics/terrascan/issues/439)

**Closed issues:**

- Terrascan is failing in scan [\#887](https://github.com/accurics/terrascan/issues/887)
- Refactor to Disable CGO [\#884](https://github.com/accurics/terrascan/issues/884)
- Issue on Azure Pipelines: failed to initialize terrascan 1.7.0 [\#864](https://github.com/accurics/terrascan/issues/864)
- Can't skip rules with underscore [\#856](https://github.com/accurics/terrascan/issues/856)
- Recursive Loop Scanning Terraform [\#851](https://github.com/accurics/terrascan/issues/851)
- Improve filenames in remote modules [\#841](https://github.com/accurics/terrascan/issues/841)
- Issues running terrascan in azure pipelines [\#835](https://github.com/accurics/terrascan/issues/835)

**Merged pull requests:**

- fix error messages reported from hcl diags [\#911](https://github.com/accurics/terrascan/pull/911) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- add in-file instrumentation segment [\#910](https://github.com/accurics/terrascan/pull/910) ([amirbenv](https://github.com/amirbenv))
- Minor documentation fixes [\#908](https://github.com/accurics/terrascan/pull/908) ([brandysnaps](https://github.com/brandysnaps))
- Use CGO independent package for sqlite [\#906](https://github.com/accurics/terrascan/pull/906) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- GH action doc - fix code block [\#902](https://github.com/accurics/terrascan/pull/902) ([amirbenv](https://github.com/amirbenv))
- Update cicd-fix code block.md [\#901](https://github.com/accurics/terrascan/pull/901) ([amirbenv](https://github.com/amirbenv))
- fixes: recursive loop when parent and child module has same local block [\#900](https://github.com/accurics/terrascan/pull/900) ([Rchanger](https://github.com/Rchanger))
- Update mkdocs-material to 7.1.9 [\#895](https://github.com/accurics/terrascan/pull/895) ([pyup-bot](https://github.com/pyup-bot))
- Updates documentation on Terrascan github action [\#894](https://github.com/accurics/terrascan/pull/894) ([cesar-rodriguez](https://github.com/cesar-rodriguez))
- fix usage overview links.md [\#893](https://github.com/accurics/terrascan/pull/893) ([amirbenv](https://github.com/amirbenv))
- Split usage docs [\#890](https://github.com/accurics/terrascan/pull/890) ([amirbenv](https://github.com/amirbenv))
- add proper values via metadata [\#888](https://github.com/accurics/terrascan/pull/888) ([gaurav-gogia](https://github.com/gaurav-gogia))
- Update mkdocs to 1.2.1 [\#886](https://github.com/accurics/terrascan/pull/886) ([pyup-bot](https://github.com/pyup-bot))
- Update Integration Docs.md [\#885](https://github.com/accurics/terrascan/pull/885) ([amirbenv](https://github.com/amirbenv))
- k8s policies refactor [\#879](https://github.com/accurics/terrascan/pull/879) ([gaurav-gogia](https://github.com/gaurav-gogia))
- mod azure policies to improve parity with siac [\#878](https://github.com/accurics/terrascan/pull/878) ([gaurav-gogia](https://github.com/gaurav-gogia))
- Fix authorization header for http request [\#877](https://github.com/accurics/terrascan/pull/877) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- Adding Id fix for github policies [\#874](https://github.com/accurics/terrascan/pull/874) ([shreyas-phansalkar-189](https://github.com/shreyas-phansalkar-189))
- Bugfix/k8s id field [\#873](https://github.com/accurics/terrascan/pull/873) ([gaurav-gogia](https://github.com/gaurav-gogia))
- Add ID Field in Azure Policies [\#872](https://github.com/accurics/terrascan/pull/872) ([gaurav-gogia](https://github.com/gaurav-gogia))
- adding ID field for aws policies [\#871](https://github.com/accurics/terrascan/pull/871) ([harkirat22](https://github.com/harkirat22))
- Adding missing Id field for GCP policies [\#870](https://github.com/accurics/terrascan/pull/870) ([shreyas-phansalkar-189](https://github.com/shreyas-phansalkar-189))
- Updating network security policies for GCP [\#869](https://github.com/accurics/terrascan/pull/869) ([shreyas-phansalkar-189](https://github.com/shreyas-phansalkar-189))
- improves: filename in remote module [\#867](https://github.com/accurics/terrascan/pull/867) ([Rchanger](https://github.com/Rchanger))
- Adding AWS Network Security Policies [\#866](https://github.com/accurics/terrascan/pull/866) ([shreyas-phansalkar-189](https://github.com/shreyas-phansalkar-189))
- Change api, Add support for s3 bucket resource and better cft loader [\#865](https://github.com/accurics/terrascan/pull/865) ([sigmabaryon](https://github.com/sigmabaryon))
- Fixes incorrect filepath reporting in sarif output & added e2e tests for sarif output [\#863](https://github.com/accurics/terrascan/pull/863) ([dev-gaur](https://github.com/dev-gaur))
- Bugfix/az nw sec policies [\#862](https://github.com/accurics/terrascan/pull/862) ([gaurav-gogia](https://github.com/gaurav-gogia))
- Update mkdocs-material to 7.1.8 [\#859](https://github.com/accurics/terrascan/pull/859) ([pyup-bot](https://github.com/pyup-bot))
- Fix AC\_AZURE\_0185 policy [\#858](https://github.com/accurics/terrascan/pull/858) ([maxgio92](https://github.com/maxgio92))
- fixed sarif unit tests hardcoding code smell [\#857](https://github.com/accurics/terrascan/pull/857) ([dev-gaur](https://github.com/dev-gaur))
- fix broken link to `usage.md` [\#855](https://github.com/accurics/terrascan/pull/855) ([dan-hill2802](https://github.com/dan-hill2802))
- Added "id" field support & policy validation tests [\#843](https://github.com/accurics/terrascan/pull/843) ([nasir-rabbani](https://github.com/nasir-rabbani))
- Add Microsoft Azure ARM as an IaC Provider  [\#736](https://github.com/accurics/terrascan/pull/736) ([gauravgahlot](https://github.com/gauravgahlot))

## [v1.7.0](https://github.com/accurics/terrascan/tree/v1.7.0) (2021-06-09)

[Full Changelog](https://github.com/accurics/terrascan/compare/v1.6.0...v1.7.0)

**Implemented enhancements:**

- Enhancement: Support sarif as output format [\#775](https://github.com/accurics/terrascan/issues/775)
- Admission Controller e2e tests [\#749](https://github.com/accurics/terrascan/issues/749)
- Enhance terrascan docker to support all terrascan run modes [\#748](https://github.com/accurics/terrascan/issues/748)
- Config file changes for server and admission controller [\#747](https://github.com/accurics/terrascan/issues/747)
- Create Helm charts for the terrascan admission webhook setup. [\#685](https://github.com/accurics/terrascan/issues/685)
- Enhancement: Use module instance name for download directory [\#672](https://github.com/accurics/terrascan/issues/672)

**Fixed bugs:**

- Azure AKS failling to check the network policy status. [\#789](https://github.com/accurics/terrascan/issues/789)
- Scan for terraform doesn't error out if a module definition refers to a directory with no tf files [\#782](https://github.com/accurics/terrascan/issues/782)
- Wrong detection of MemoryRequestsCheck,CpuRequestsCheck,noReadinessProbe and nolivenessProbe policy in k8s Job spec  [\#767](https://github.com/accurics/terrascan/issues/767)
- Update Docker build for terrascan to use numeric UID [\#766](https://github.com/accurics/terrascan/issues/766)
- Wrong detection of AllowPrivilegeEscalation \(policy AC-K8-CA-PO-H-0165\) in K8s pod spec  [\#721](https://github.com/accurics/terrascan/issues/721)
- Failed to run prepared query error in opa/engine.go [\#709](https://github.com/accurics/terrascan/issues/709)
- tfplan should use resource address for id field [\#702](https://github.com/accurics/terrascan/issues/702)
- Rule IDs with spaces cannot be skipped [\#610](https://github.com/accurics/terrascan/issues/610)
- AWS.CloudFront.Network Security.Low.0568 Doesn't allow skipping due to space in filename [\#549](https://github.com/accurics/terrascan/issues/549)
- Error parsing syntax if using complex query for dynamic ip\_restriction in azurerm\_function\_app or azurerm\_app\_service ressource [\#433](https://github.com/accurics/terrascan/issues/433)

**Closed issues:**

- Add support for YAML format for terrascan config file [\#807](https://github.com/accurics/terrascan/issues/807)
- Add ID field [\#805](https://github.com/accurics/terrascan/issues/805)
- Add a middleware to log incoming http\(s\) requests on terrascan server [\#784](https://github.com/accurics/terrascan/issues/784)
- terrascan server: validation missing for --cert-path and --key-path [\#769](https://github.com/accurics/terrascan/issues/769)
- show-passed should report passes only for the existing resources [\#757](https://github.com/accurics/terrascan/issues/757)
- Out of the box handling of certificates in helm charts for terrascan in Server mode  [\#756](https://github.com/accurics/terrascan/issues/756)
- In-file Instrumentation [\#755](https://github.com/accurics/terrascan/issues/755)
- Release 1.5.2 or 1.6.0 [\#745](https://github.com/accurics/terrascan/issues/745)
- Issue in GCP Policyfile unrestrictedRdpAccess.rego [\#735](https://github.com/accurics/terrascan/issues/735)
- accurics.azure.AKS.3 is defective [\#711](https://github.com/accurics/terrascan/issues/711)
- Rule `lambdaNotEncryptedWithKms` should not check for KMS when env vars are not being used [\#682](https://github.com/accurics/terrascan/issues/682)
- Terrascan does not resolve env var for aws\_rds\_cluster attribute storage\_encrypted [\#678](https://github.com/accurics/terrascan/issues/678)
- Valid Terraform configuration fails with `s3EnforceUserAcl` [\#659](https://github.com/accurics/terrascan/issues/659)
-   kmsKeyExposedPolicy:22: eval\_builtin\_error: json.unmarshal: invalid character '$' looking for beginning of value} [\#627](https://github.com/accurics/terrascan/issues/627)
- Terrascan not able to find terraform config files in a sub directory, but it works in case of k8s infrastructure type [\#622](https://github.com/accurics/terrascan/issues/622)
- Potential nil-dereference found while fuzzing [\#611](https://github.com/accurics/terrascan/issues/611)
- terrascan should have a `category-list` command [\#597](https://github.com/accurics/terrascan/issues/597)
- Improved Documentation [\#416](https://github.com/accurics/terrascan/issues/416)
- Improve test coverage for k8s [\#400](https://github.com/accurics/terrascan/issues/400)

**Merged pull requests:**

- Fixing the bug for google\_kms\_crypto\_key policies [\#848](https://github.com/accurics/terrascan/pull/848) ([shreyas-phansalkar-189](https://github.com/shreyas-phansalkar-189))
- Fix AWS dynamo Db policy for point in time recovery [\#847](https://github.com/accurics/terrascan/pull/847) ([harkirat22](https://github.com/harkirat22))
- Bugfix/use ref id old format [\#846](https://github.com/accurics/terrascan/pull/846) ([gaurav-gogia](https://github.com/gaurav-gogia))
- reference ids with & and \<space\> fixed [\#845](https://github.com/accurics/terrascan/pull/845) ([gaurav-gogia](https://github.com/gaurav-gogia))
- fixes: Terraform inner block reference resolution [\#844](https://github.com/accurics/terrascan/pull/844) ([Rchanger](https://github.com/Rchanger))
- Bump up to Go/1.16 [\#836](https://github.com/accurics/terrascan/pull/836) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- \[fix\] Add Alternate names for k8s services [\#834](https://github.com/accurics/terrascan/pull/834) ([rahulchheda](https://github.com/rahulchheda))
- Support for spaces in policy reference\_id [\#833](https://github.com/accurics/terrascan/pull/833) ([nasir-rabbani](https://github.com/nasir-rabbani))
- fix - type assertion check for hcl.Body in terraform iac provider [\#832](https://github.com/accurics/terrascan/pull/832) ([patilpankaj212](https://github.com/patilpankaj212))
- Add ID Field for AWS Policies' Metadata [\#831](https://github.com/accurics/terrascan/pull/831) ([gaurav-gogia](https://github.com/gaurav-gogia))
- Policy to check CVE-2021-25737 [\#830](https://github.com/accurics/terrascan/pull/830) ([harkirat22](https://github.com/harkirat22))
- Enhancing AWS policies [\#829](https://github.com/accurics/terrascan/pull/829) ([harkirat22](https://github.com/harkirat22))
- aws s3 policy `s3EnforceUserAcl` update [\#828](https://github.com/accurics/terrascan/pull/828) ([gaurav-gogia](https://github.com/gaurav-gogia))
- add check for env vars and kms [\#827](https://github.com/accurics/terrascan/pull/827) ([gaurav-gogia](https://github.com/gaurav-gogia))
- Add ID Field for K8s Policies' Metadata [\#826](https://github.com/accurics/terrascan/pull/826) ([Avanti19](https://github.com/Avanti19))
- Do not trim resource id from tfplan json [\#825](https://github.com/accurics/terrascan/pull/825) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
- Add ID Field for GCP Policies' Metadata [\#824](https://github.com/accurics/terrascan/pull/824) ([gaurav-gogia](https://github.com/gaurav-gogia))
- fix - source path for k8s file scan is absolute [\#821](https://github.com/accurics/terrascan/pull/821) ([patilpankaj212](https://github.com/patilpankaj212))
- added pending test changes for config reader [\#820](https://github.com/accurics/terrascan/pull/820) ([patilpankaj212](https://github.com/patilpankaj212))
- fix: moves the pending test to running [\#819](https://github.com/accurics/terrascan/pull/819) ([Rchanger](https://github.com/Rchanger))
- fix multierror variable issue [\#818](https://github.com/accurics/terrascan/pull/818) ([patilpankaj212](https://github.com/patilpankaj212))
- \[feat.\] Merge Webhook and Server Helm Chart [\#817](https://github.com/accurics/terrascan/pull/817) ([rahulchheda](https://github.com/rahulchheda))
- add support for YAML format for terrascan config file [\#816](https://github.com/accurics/terrascan/pull/816) ([kanchwala-yusuf](https://github.com/kanchwala-yusuf))
-  Add AWS  CFT as an IaC Provider  [\#815](https://github.com/accurics/terrascan/pull/815) ([mahendrabagul](https://github.com/mahendrabagul))
- fix failing e2e test [\#812](https://github.com/accurics/terrascan/pull/812) ([patilpankaj212](https://github.com/patilpankaj212))
- Adding Aws new policies cloudTrail [\#810](https://github.com/accurics/terrascan/pull/810) ([Avanti19](https://github.com/Avanti19))
- Feature/az id field [\#808](https://github.com/accurics/terrascan/pull/808) ([gaurav-gogia](https://github.com/gaurav-gogia))
- added support for sarif formatted violation reports [\#806](https://github.com/accurics/terrascan/pull/806) ([dev-gaur](https://github.com/dev-gaur))
- Adds support to scan config resources with applicable policies & Refactors filteration [\#803](https://github.com/accurics/terrascan/pull/803) ([patilpankaj212](https://github.com/patilpankaj212))
- Adds: in-file instrumentation for resource prioritizing [\#802](https://github.com/accurics/terrascan/pull/802) ([Rchanger](https://github.com/Rchanger))
- shifted opa engine warning message to debug log level [\#800](https://github.com/accurics/terrascan/pull/800) ([dev-gaur](https://github.com/dev-gaur))
- fix: added validation for module local source dir [\#793](https://github.com/accurics/terrascan/pull/793) ([Rchanger](https://github.com/Rchanger))
- policy metadata changes to include `policy\_type` and `resource\_type` [\#792](https://github.com/accurics/terrascan/pull/792) ([patilpankaj212](https://github.com/patilpankaj212))
- Fix pod level securityContext support [\#790](https://github.com/accurics/terrascan/pull/790) ([harkirat22](https://github.com/harkirat22))
- Fix policy code for securityContext and Probes [\#787](https://github.com/accurics/terrascan/pull/787) ([harkirat22](https://github.com/harkirat22))
- add logging middleware for server [\#785](https://github.com/accurics/terrascan/pull/785) ([dev-gaur](https://github.com/dev-gaur))
- config file changes for terrascan server [\#780](https://github.com/accurics/terrascan/pull/780) ([patilpankaj212](https://github.com/patilpankaj212))
- Automate generation of TLS Certs using Helm [\#779](https://github.com/accurics/terrascan/pull/779) ([rahulchheda](https://github.com/rahulchheda))
- Add webhook setup capability and remote repo scan capability in the helm charts [\#778](https://github.com/accurics/terrascan/pull/778) ([dev-gaur](https://github.com/dev-gaur))
- Changed description of policy file to match port. [\#777](https://github.com/accurics/terrascan/pull/777) ([menzbua](https://github.com/menzbua))
- Added source\_range 0.0.0.0/0 \(any\) to avoid rule violations [\#776](https://github.com/accurics/terrascan/pull/776) ([menzbua](https://github.com/menzbua))
- support for `module name` in violation summary  [\#774](https://github.com/accurics/terrascan/pull/774) ([patilpankaj212](https://github.com/patilpankaj212))
- Modified the Dockerfile to use numeric UID  [\#773](https://github.com/accurics/terrascan/pull/773) ([Rchanger](https://github.com/Rchanger))
- adds e2e tests for validating webhook [\#772](https://github.com/accurics/terrascan/pull/772) ([patilpankaj212](https://github.com/patilpankaj212))
- add validation for tls private key and cert file values [\#771](https://github.com/accurics/terrascan/pull/771) ([dev-gaur](https://github.com/dev-gaur))
- Documentation [\#768](https://github.com/accurics/terrascan/pull/768) ([lalchand12](https://github.com/lalchand12))
- change docs to include docker subcommands.md [\#765](https://github.com/accurics/terrascan/pull/765) ([amirbenv](https://github.com/amirbenv))
- shifted custom atlantis container source under integrations/ directory [\#758](https://github.com/accurics/terrascan/pull/758) ([dev-gaur](https://github.com/dev-gaur))
- Update mkdocs-material to 7.1.4 [\#746](https://github.com/accurics/terrascan/pull/746) ([pyup-bot](https://github.com/pyup-bot))
- Add a kustomize based guide for setting up terrascan server and validating webhook in kubernetes [\#739](https://github.com/accurics/terrascan/pull/739) ([dev-gaur](https://github.com/dev-gaur))
- Fix accurics.azure.AKS.3 [\#712](https://github.com/accurics/terrascan/pull/712) ([xortim](https://github.com/xortim))
- Update mkdocs-redirects to 1.0.3 [\#710](https://github.com/accurics/terrascan/pull/710) ([pyup-bot](https://github.com/pyup-bot))
- Initial addition of terrascan helm chart [\#688](https://github.com/accurics/terrascan/pull/688) ([jlk](https://github.com/jlk))

## [v1.6.0](https://github.com/accurics/terrascan/tree/v1.6.0) (2021-05-10)

[Full Changelog](https://github.com/accurics/terrascan/compare/v1.5.1...v1.6.0)

**Implemented enhancements:**

- Atlantis Integration [\#686](https://github.com/accurics/terrascan/issues/686)
- Enhancement: support for all iac scan for cli [\#673](https://github.com/accurics/terrascan/issues/673)
- Feature request: scan sub-folders too [\#411](https://github.com/accurics/terrascan/issues/411)

**Fixed bugs:**

- Admission Controller Doesn't display feedback for kubectl "create" and "apply" [\#731](https://github.com/accurics/terrascan/issues/731)

**Closed issues:**

- GKE Control Plane is exposed to few public IP addresses [\#743](https://github.com/accurics/terrascan/issues/743)
- Error with finding Enable AWS CloudWatch Logs for APIs [\#730](https://github.com/accurics/terrascan/issues/730)
- Task: Add to github actions ability to build/push terrascan\_atlantis image [\#728](https://github.com/accurics/terrascan/issues/728)
- accurics.azure.NS.161 does not work with tfplan [\#725](https://github.com/accurics/terrascan/issues/725)
- terrascan "latest" docker image broken for tfplan [\#718](https://github.com/accurics/terrascan/issues/718)
- Local expansion recursive infinite loop [\#690](https://github.com/accurics/terrascan/issues/690)

**Merged pull requests:**

- Feature/aws new policies sp [\#751](https://github.com/accurics/terrascan/pull/751) ([shreyas-phansalkar-189](https://github.com/shreyas-phansalkar-189))
- Argocd doc volume field modification [\#742](https://github.com/accurics/terrascan/pull/742) ([Rchanger](https://github.com/Rchanger))
- Update mkdocs.yml [\#741](https://github.com/accurics/terrascan/pull/741) ([amirbenv](https://github.com/amirbenv))
- fix failing test [\#740](https://github.com/accurics/terrascan/pull/740) ([patilpankaj212](https://github.com/patilpankaj212))
- AWS policy pack update [\#737](https://github.com/accurics/terrascan/pull/737) ([harkirat22](https://github.com/harkirat22))
- Adding release checklist [\#734](https://github.com/accurics/terrascan/pull/734) ([jlk](https://github.com/jlk))
- Gh action terrscan\_atlantis release [\#733](https://github.com/accurics/terrascan/pull/733) ([dev-gaur](https://github.com/dev-gaur))
- adds agrocd integration dockerfile, scripts, doc  and examples [\#732](https://github.com/accurics/terrascan/pull/732) ([Rchanger](https://github.com/Rchanger))
- Fix NSG associations [\#727](https://github.com/accurics/terrascan/pull/727) ([xortim](https://github.com/xortim))
- changes for argocd integration [\#724](https://github.com/accurics/terrascan/pull/724) ([patilpankaj212](https://github.com/patilpankaj212))
- Update admission-controller-webhooks-usage.md [\#722](https://github.com/accurics/terrascan/pull/722) ([amirbenv](https://github.com/amirbenv))
- fix - \#718 [\#720](https://github.com/accurics/terrascan/pull/720) ([patilpankaj212](https://github.com/patilpankaj212))
- doc: add homebrew badge [\#714](https://github.com/accurics/terrascan/pull/714) ([chenrui333](https://github.com/chenrui333))
- update version [\#713](https://github.com/accurics/terrascan/pull/713) ([chenrui333](https://github.com/chenrui333))
- adds skipped tests for server file scan when file is k8s yaml [\#706](https://github.com/accurics/terrascan/pull/706) ([Rchanger](https://github.com/Rchanger))
- fixes infinte loop while local variable resolution [\#700](https://github.com/accurics/terrascan/pull/700) ([patilpankaj212](https://github.com/patilpankaj212))
- add terrascan atlantis container files, scripts and doc. [\#684](https://github.com/accurics/terrascan/pull/684) ([dev-gaur](https://github.com/dev-gaur))
- adds support to scan directory with all iac providers in cli mode [\#674](https://github.com/accurics/terrascan/pull/674) ([patilpankaj212](https://github.com/patilpankaj212))
- adds support to scan sub folders for terraform iac provider [\#640](https://github.com/accurics/terrascan/pull/640) ([patilpankaj212](https://github.com/patilpankaj212))

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
