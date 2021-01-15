![Terrascan](https://raw.githubusercontent.com/accurics/terrascan/master/docs/img/Terrascan_By_Accurics_Logo_38B34A-333F48.svg)

[![GitHub release](https://img.shields.io/github/release/accurics/terrascan)](https://github.com/accurics/terrascan/releases/latest)
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202-blue)](https://github.com/accurics/terrascan/blob/master/LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/accurics/terrascan/pulls)
![CI](https://github.com/accurics/terrascan/workflows/build/badge.svg)
[![codecov](https://codecov.io/gh/accurics/terrascan/branch/master/graph/badge.svg)](https://codecov.io/gh/accurics/terrascan)
[![community](https://img.shields.io/discourse/status?server=https%3A%2F%2Fcommunity.accurics.com)](https://community.accurics.com)
[![Documentation Status](https://readthedocs.com/projects/accurics-terrascan/badge/?version=latest)](https://docs.accurics.com/projects/accurics-terrascan/en/latest/?badge=latest)

Terrascan detects security vulnerabilities and compliance violations across your Infrastructure as Code. Mitigate risks before provisioning cloud native infrastructure. Run locally or integrate with your CI\CD.


* Documentation: https://docs.accurics.com/projects/accurics-terrascan
* Discuss: https://community.accurics.com

## Features
* 500+ Policies for security best practices
* Scanning of Terraform (HCL2)
* Scanning of Kubernetes (JSON/YAML), Helm v3, and Kustomize v3
* Support for AWS, Azure, GCP, Kubernetes and GitHub


## Quick Start
### Step 1: Install
Terrascan's supports multiple ways to install, including [brew](https://github.com/accurics/terrascan#install-via-brew).
Here, we will download the terrascan binary directly from the [releases](https://github.com/accurics/terrascan/releases) page. Make sure to select the right binary for your machine. Here's an example of how to install it:

```sh
$ curl --location https://github.com/accurics/terrascan/releases/download/v1.3.0/terrascan_1.3.0_Darwin_x86_64.tar.gz --output terrascan.tar.gz
$ tar -xvf terrascan.tar.gz
  x CHANGELOG.md
  x LICENSE
  x README.md
  x terrascan
$ install terrascan /usr/local/bin
$ terrascan
```
### Step 2: Run
To scan your code for security issues you can run the following (defaults to scanning Terraform).

```sh
$ terrascan scan
```
Terrascan will exit 3 if any issues are found.

The following commands are available:

```sh
$ terrascan
Terrascan

Detect compliance and security violations across Infrastructure as Code to mitigate risk before provisioning cloud native infrastructure.
For more information, please visit https://docs.accurics.com

Usage:
  terrascan [command]

Available Commands:
  help        Help about any command
  init        Initialize Terrascan
  scan        Detect compliance and security violations across Infrastructure as Code.
  server      Run Terrascan as an API server
  version     Terrascan version

Flags:
  -c, --config-path string   config file path
  -h, --help                 help for terrascan
  -l, --log-level string     log level (debug, info, warn, error, panic, fatal) (default "info")
  -x, --log-type string      log output type (console, json) (default "console")
  -o, --output string        output type (human, json, yaml, xml) (default "human")

Use "terrascan [command] --help" for more information about a command.
```

### Step 3: Integrate with CI\CD
Please refer to our [documentation to integrate with your pipeline](https://docs.accurics.com/projects/accurics-terrascan/en/latest/cicd/).


## Other Installation Options


### Install via `brew`

[Homebrew](https://brew.sh/) users can install by:

```sh
$ brew install terrascan
```

### Docker
Terrascan is also available as a Docker image and can be used as follows

```sh
$ docker run accurics/terrascan
```

### Install via go get, if you have Go installed
```
$ export GO111MODULE=on
$ go get -u github.com/accurics/terrascan/cmd/terrascan
  go: downloading github.com/accurics/terrascan v1.3.0
  go: found github.com/accurics/terrascan/cmd/terrascan in github.com/accurics/terrascan v1.3.0
  ...
$ terrascan
```

### Building Terrascan
Terrascan can be built locally. This is helpful if you want to be on the latest version or when developing Terrascan.

```sh
$ git clone git@github.com:accurics/terrascan.git
$ cd terrascan
$ make build
$ ./bin/terrascan
```


## Developing Terrascan
To learn more about developing and contributing to Terrascan refer to the [contributing guide](CONTRIBUTING.md).

## License

Terrascan is licensed under the [Apache 2.0 License](LICENSE).
