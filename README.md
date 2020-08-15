# Terrascan
![CI](https://github.com/accurics/terrascan/workflows/build/badge.svg)
[![codecov](https://codecov.io/gh/accurics/terrascan/branch/master/graph/badge.svg)](https://codecov.io/gh/accurics/terrascan)
[![community](https://img.shields.io/discourse/status?server=https%3A%2F%2Fcommunity.accurics.com)](https://community.accurics.com)
[![Documentation Status](https://readthedocs.com/projects/accurics-terrascan/badge/?version=latest)](https://docs.accurics.com/projects/accurics-terrascan/en/latest/?badge=latest)
[![downloads](https://img.shields.io/github/downloads/accurics/terrascan/total)](https://github.com/accurics/terrascan/releases)


Detect compliance and security violations across Infrastructure as Code to mitigate risk before provisioning cloud native infrastructure.


* GitHub Repo: https://github.com/accurics/terrascan
* Documentation: https://docs.accurics.com
* Discuss: https://community.accurics.com

## Features
* 500+ Policies for security best practices
* Scanning of Terraform 12+ (HCL2)
* Support for AWS, Azure, and GCP

## Installing
Terrascan's binary for your architecture can be found on the [releases](https://github.com/accurics/terrascan/releases) page. Here's an example of how to install it:

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

### Docker
Terrascan is also available as a Docker image and can be used as follows

```
$ docker run accurics/terrascan
```

## Getting started

To scan your code for security issues you can run the following

```
$ terrascan --iac terraform --iac-version v12 --cloud aws -d pkg/iac-providers/terraform/v12/testdata/moduleconfigs
```

The following flags are available:

```
$ terrascan -h

Terrascan

Scan IaC files for security violations

Usage

    terrascan -cloud [aws|azure|gcp] [options...]

Options

Cloud
    -cloud                Required. Cloud provider (supported values: aws, azure, gcp)

IaC (Infrastructure as Code)
    -d                    IaC directory path (default: current working directory)
    -f                    IaC file path
    -iac                  IaC provider (supported values: terraform, default: terraform)
    -iac-version          IaC version (supported values: 'v12' for Terraform, default: v12)
    -p                    Policy directory path

Mode
    -server               Run Terrascan in server mode

Logging
    -log-level            Logging level (supported values: debug, info, warn, error, panic, fatal)
    -log-type             Logging type (supported values: json, yaml, console, default: console)

Miscellaneous
    -config               Configuration file path
    -version              Print the Terrascan version
```

## Documentation

To learn more about Terrascan check out the documentation https://docs.accurics.com where we include a getting started guide, Terrascan's architecture, a break down of it's commands, and how to write your own policies.

## Developing Terrascan
To learn more about developing and contributing to Terrascan refer to our [contributing guide](CONTRIBUTING.md).


To learn more about compiling Terraform and contributing suggested changes, please refer to the contributing guide.

## License

Terrascan is licensed under the [Apache 2.0 License](LICENSE).
