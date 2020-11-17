# Changelog

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
