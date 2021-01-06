# Integrating Terrascan into CI/CD

Terrascan can be integrated into CI/CD pipelines to enforce security best practices as codified in the OPA rego policies included as part of Terrascan or any custom policies. This section contains examples on how to configure Terrascan in popular CI/CD tooling.

## GitHub Actions

Terrascan can be configured as a job within GitHub actions workflows. A straightforward way to accomplish this is by using the [super-linter](https://github.com/github/super-linter) GitHub action which includes Terrascan. Note that at the moment super-linter only supports scanning Terraform HCL files.

When using super-linter you can pass the environment variable "VALIDATE_TERRAFORM_TERRASCAN: true" to ensure that Terraform configuration files are evaluated using Terrascan. To configure your GitHub actions workflow a file with the below YAML content can be included within the .github/workflows/ directory of your repository.

``` YAML
---
name: Scan Code Base
on:
  push:
  pull_request:
    branches: [master]
jobs:
  build:
    name: Scan Code Base
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v2
      - name: Scan Code Base
        uses: github/super-linter@v3
        env:
          VALIDATE_ALL_CODEBASE: true
          DEFAULT_BRANCH: master
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          VALIDATE_TERRAFORM_TERRASCAN: true
```

Documentation on the GitHub actions workflow syntax is available [here](https://help.github.com/en/articles/workflow-syntax-for-github-actions).


## GitLab CI

[GitLab CI](https://docs.gitlab.com/ee/ci/README.html) can use [Docker images](https://docs.gitlab.com/ee/ci/docker/using_docker_images.html) as part of a pipeline. We can take advantage of this functionality and use Terrascan's docker image as part of your [pipeline](https://docs.gitlab.com/ee/ci/pipelines/) to scan infrastructure as code.

To do this you can update your .gitlab-ci.yml file to use the "accurics/terrascan:latest" image with the ["bin/sh", "-c"] entrypoint. Terrascan can be found on "/go/bin" in the image and you can use any [Terrascan command line options](http://ubusvr:8000/getting-started/usage/#terrascan-commands) according to your needs. Here's an example .gitlab-ci.yml file:

``` YAML
stages:
  - scan

terrascan:
  image:
    name: accurics/terrascan:latest
    entrypoint: ["/bin/sh", "-c"]
  stage: scan
  script:
    - /go/bin/terrascan scan .
```
