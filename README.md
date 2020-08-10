# Terrascan
![CI](https://github.com/accurics/terrascan/workflows/build/badge.svg)
[![codecov](https://codecov.io/gh/accurics/terrascan/branch/master/graph/badge.svg)](https://codecov.io/gh/accurics/terrascan)
[![community](https://img.shields.io/discourse/status?server=https%3A%2F%2Fcommunity.accurics.com)](https://community.accurics.com)
[![Documentation](https://readthedocs.org/projects/terrascan/badge/?version=latest)](https://terrascan.readthedocs.io/en/latest/?badge=latest)
[![downloads](https://img.shields.io/github/downloads/accurics/terrascan/total)](https://github.com/accurics/terrascan/releases)


Terrascan is a static code analyzer and linter for security weanesses in Infrastructure as Code (IaC).

* GitHub Repo: https://github.com/accurics/terrascan
* Documentation: https://docs.accurics.com
* Discuss: https://community.accurics.com

## Features
* 500+ Policies for security best practices
* Scanning of Terraform 12+ (HCL2)
* Support for AWS, Azure, and GCP

## Installing
Terrascan's binary for your architecture can be found on the releases page. Here's an example of how to install it:

```
$ curl --location https://github.com/accurics/terrascan/releases/download/v1.0.0/terrascan_darwin_amd64.zip --output terrascan_darwin_amd64.zip
$ unzip terrascan_darwin_amd64.zip
Archive:  terrascan_darwin_amd64.zip
  inflating: terrascan
$ install terrascan /usr/local/bin
$ terrascan --help
```

### Homebrew
Terrascan can be installed using Homebrew on macOS:

```
brew install terrascan
```

### Chocolatey
Terrascan can be installed on Windows using Chocolatey:

```
choco install terrascan
```

### Docker
Terrascan is also available as a Docker image and can be used as follows

	$ docker run accurics/terrascan

## Getting started

To scan your code for security weaknesses you can run the following

```
$ terrascan --iac terraform --iac-version v12 --cloud aws -d pkg/iac-providers/terraform/v12/testdata/moduleconfigs
```

The following flags are available:

```
$ terrascan --help
Usage of ./bin/terrascan:
  -cloud string
        cloud provider (supported values: aws)
  -d string
        IaC directory path
  -f string
        IaC file path
  -iac string
        IaC provider (supported values: terraform)
  -iac-version string
        IaC version (supported values: 'v12' for terraform) (default "default")
  -log-level string
        logging level (debug, info, warn, error, panic, fatal) (default "info")
  -log-type string
        log type (json, console) (default "console")
  -server
        run terrascan in server mode
```

## Documentation

To learn more about Terrascan check out the documentation https://docs.accurics.com where we include a getting started guide, Terrascan's architecture, a break down of it's commands, and how to write your own policies.

## Developing Terrascan
To learn more about developing and contributing to Terrascan refer to our (contributing guide)[CONTRIBUTING.md].


To learn more about compiling Terraform and contributing suggested changes, please refer to the contributing guide.

## License

Terrascan is licensed under the [Apache 2.0 License](LICENSE).
