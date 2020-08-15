# Changelog

## 1.0.0 (UNRELEASED)
Major updates to Terrascan and the underlying architecture including:
* Pluggable architecture written in Golang. We updated the architecture to be easier to extend Terrascan with additional IaC languages and support policies for different cloud providers and cloud native tooling.
* Server mode. This allows Terrascan to be executed as a server and use it's API to perform static code analysis
* Notifications hooks. Will be able to integrate for notifications to external systems (e.g. email, slack, etc.)
* Uses OPA policy engine and policies written in Rego.

## 0.2.3 (2020-07-23)
* Introduces the '-f' flag for passing a list of ".tf" files for linting and the '--version' flag.

## 0.2.2 (2020-07-21)
* Adds Docker image and pipeline to push to DockerHub

## 0.2.1 (2020-06-19)
* Bugfix: The pyhcl hard dependency in the requirements.txt file caused issues if a higher version was installed. This was fixed by using the ">=" operator.

## 0.2.0 (2020-01-11)
* Adds support for terraform 0.12+

## 0.1.2 (2020-01-05)
* Adds ability to setup terrascan as a pre-commit hook

## 0.1.1 (2020-01-01)
* Updates dependent packages to latest versions
* Migrates CI to GitHub Actions from travis

## 0.1.0 (2017-11-26)
* First release on PyPI.
