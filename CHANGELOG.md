# Changelog

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
