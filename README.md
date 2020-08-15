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

### Docker
Terrascan is also available as a Docker image and can be used as follows

```
$ docker run accurics/terrascan
```

### Homebrew
Terrascan can be installed using Homebrew on macOS:

```
brew install terrascan
```

## Getting started

To scan your code for security issues you can run the following

```
$ terrascan scan -t aws
```

The following commands are available:

```
$ terrascan
Terrascan

An advanced IaC (Infrastructure-as-Code) file scanner written in Go.
Secure your cloud deployments at design time.
For more information, please visit https://www.accurics.com

Usage:
  terrascan [command]

Available Commands:
  help        Help about any command
  init        Initialize Terrascan
  scan        Scan IaC (Infrastructure-as-Code) files for vulnerabilities.
  server      Run Terrascan as an API server

Flags:
  -c, --config-path string   config file path
  -h, --help                 help for terrascan
  -l, --log-level string     log level (debug, info, warn, error, panic, fatal) (default "info")
  -x, --log-type string      log output type (console, json) (default "console")
  -o, --output-type string   output type (json, yaml, xml) (default "yaml")
  -v, --version              version for terrascan

Use "terrascan [command] --help" for more information about a command.
```

## Documentation

To learn more about Terrascan check out the documentation https://docs.accurics.com where we include a getting started guide, Terrascan's architecture, a break down of it's commands, and a deep dive into policies.

## Developing Terrascan
To learn more about developing and contributing to Terrascan refer to our [contributing guide](CONTRIBUTING.md).


To learn more about compiling Terraform and contributing suggested changes, please refer to the contributing guide.

## License

Terrascan is licensed under the [Apache 2.0 License](LICENSE).
