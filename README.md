[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/dlminvestments/terrascan)

# Terrascan

[![GitHub release](https://img.shields.io/github/release/accurics/terrascan)](https://github.com/accurics/terrascan/releases/latest)
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202-blue)](https://github.com/accurics/terrascan/blob/master/LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/accurics/terrascan/pulls)
![CI](https://github.com/accurics/terrascan/workflows/build/badge.svg)
[![codecov](https://codecov.io/gh/accurics/terrascan/branch/master/graph/badge.svg)](https://codecov.io/gh/accurics/terrascan)
[![community](https://img.shields.io/discourse/status?server=https%3A%2F%2Fcommunity.accurics.com)](https://community.accurics.com)
[![Documentation Status](https://readthedocs.com/projects/accurics-terrascan/badge/?version=latest)](https://docs.accurics.com/projects/accurics-terrascan/en/latest/?badge=latest)
[![Blimp demo badge](https://blimpup.io/demo-badge.svg?repo=https://github.com/accurics/terrascan.git)](https://blimpup.io/preview-env/?repo=https://github.com/accurics/terrascan.git&composeFiles=deploy/docker-compose.yml&port=terrascan:9010)


Detect compliance and security violations across Infrastructure as Code to mitigate risk before provisioning cloud native infrastructure.


* GitHub Repo: https://github.com/accurics/terrascan
* Documentation: https://docs.accurics.com
* Discuss: https://community.accurics.com

## Features
* 500+ Policies for security best practices
* Scanning of Terraform 12+ (HCL2)
* Scanning of Kubernetes YAML/JSON
* Support for AWS, Azure, GCP, Kubernetes and GitHub

## Installing
Terrascan's binary for your architecture can be found on the [releases](https://github.com/accurics/terrascan/releases) page. Here's an example of how to install it:

```sh
$ curl --location https://github.com/accurics/terrascan/releases/download/v1.1.0/terrascan_1.1.0_Darwin_x86_64.tar.gz --output terrascan.tar.gz
$ tar -xvf terrascan.tar.gz
  x CHANGELOG.md
  x LICENSE
  x README.md
  x terrascan
$ install terrascan /usr/local/bin
$ terrascan
```

If you have go installed, Terrascan can be installed with `go get`
```
$ export GO111MODULE=on
$ go get -u github.com/accurics/terrascan/cmd/terrascan
  go: downloading github.com/accurics/terrascan v1.1.0
  go: found github.com/accurics/terrascan/cmd/terrascan in github.com/accurics/terrascan v1.1.0
  ...
$ terrascan
```

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

### Building Terrascan
Terrascan can be built locally. This is helpful if you want to be on the latest version or when developing Terrascan.

```sh
$ git clone git@github.com:accurics/terrascan.git
$ cd terrascan
$ make build
$ ./bin/terrascan
```

### Demoing
If you want to play around with Terrascan without running it locally, you can
[boot a personal demo copy
](https://blimpup.io/preview-env/?repo=https://github.com/accurics/terrascan.git&composeFiles=deploy/docker-compose.yml&port=terrascan:9010)
from your browser without downloading or setting up anything.

1. Click the [demo
   link](https://blimpup.io/preview-env/?repo=https://github.com/accurics/terrascan.git&composeFiles=deploy/docker-compose.yml&port=terrascan:9010)
   to boot this repo in the Blimp cloud.

1. Once the sandbox is booted, get its public URL by clicking "Connect" on the terrascan service on the left.

   The page will 404, but that's OK because we just need the domain name to
   create our URL that we'll hit with `curl`.

1. Run the following command in your terminal to scan a simple Terraform file.
   Make sure to replace `<YOUR PUBLIC URL>` with the URL from the previous
   step.

   ```
   curl -i -F "file=@-" https://<YOUR PUBLIC URL>/v1/terraform/v12/aws/local/file/scan << EOF
   variable "my-variable" {
     default = "default"
     type    = string
   }
   EOF
   ```

   The full URL will look something like `https://a98c0197112b7a4a96b72ea21ac0802b.blimp.dev/v1/terraform/v12/aws/local/file/scan`.

   The command will output something like this:
   ```
   {
     "ResourceConfig": {},
     "Violations": {
       "results": {
         "violations": [],
         "count": {
           "low": 0,
           "medium": 0,
           "high": 0,
           "total": 0
         }
       }
     }
   }
   ```

See the [server mode
docs](https://docs.accurics.com/projects/accurics-terrascan/en/latest/getting-started/#server-mode)
for more information on how to use the server endpoints.

## Getting started

To scan your code for security issues you can run the following

```sh
$ terrascan scan
```
Terrascan will exit 3 if any issues are found.

The following commands are available:

```sh
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
To learn more about developing and contributing to Terrascan refer to the [contributing guide](CONTRIBUTING.md).

## License

Terrascan is licensed under the [Apache 2.0 License](LICENSE).
