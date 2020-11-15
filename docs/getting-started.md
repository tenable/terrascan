# Getting Started
Terrascan is a static code analyzer for Infrastructure as Code tooling. It can executed with the native binary/executable or by using the [`docker`](#using-docker) container.

## Installation
Terrascan's binary can be found on the package for each [release](https://github.com/accurics/terrascan/releases). Here's an example of how to install it:

``` Bash
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

### Using Docker
Terrascan is available as a Docker image and can used as follows:

``` Bash
$ docker run accurics/terrascan
```

### Building Terrascan
Terrascan can be built locally. This is helpful if you want to be on the latest version or when developing Terrascan.

``` Bash
$ git clone git@github.com:accurics/terrascan.git
$ cd terrascan
$ make build
$ ./bin/terrascan
```

## Terrascan Commands
Terrascan's interface is divided into subcommands as follows:

*   init = Will initialize Terrascan by downloading the latest Rego policies into ~/.terrascan. Note that the scan command will implicitly call this if it detects that there are no policies found.
*   scan = Will scan IaC files based on the policies contained within the .terrascan directory
*   server = Will start Terrascan's API server
*   help = You can obtain the usage menu by typing `help` or using the `-h` flag on any subcommand (e.g. `terrascan init -h`)

By typing `terrascan` without flags or other arguments, you can display the usage information.

``` Bash
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

Flags:
  -c, --config-path string   config file path
  -h, --help                 help for terrascan
  -l, --log-level string     log level (debug, info, warn, error, panic, fatal) (default "info")
  -x, --log-type string      log output type (console, json) (default "console")
  -o, --output string        output type (json, yaml, xml) (default "yaml")
  -v, --version              version for terrascan

Use "terrascan [command] --help" for more information about a command.
```

### Initializing
The initialization process downloads the latest policies from the [repository](https://github.com/accurics/terrascan) into `~/.terrascan`. The policies are located at `~/.terrascan/pkg/policies/opa/rego` and are fetched when scanning the IaC. This command is implicitly executed if the `scan` command doesn't find policies while executing.

### Scanning
The CLI will default to scanning all supported cloud providers on Terraform HCL files if the `scan` command is used with no arguments. For example, the below two commands will scan the current directory containing Terraform HCL2 files for supported providers (AWS, GCP, Azure, and GitHub) resources:

``` Bash
$ terrascan scan
```

Individual cloud providers can be specified using the -t flag as follows:

``` Bash
$ terrascan scan -t aws
```

By default Terrascan defaults to scanning Terraform HCL files, you can change the IaC provider using the -i flag. Here's an example of scanning kubernetes yaml files:

``` Bash
$ terrascan scan -i k8s
```

The `scan` command supports flags to configure: the directory being scanned, scanning of a specific file, IaC provider type, path to policies, and policy type. The full list of flags can be found by typing `terrascan scan -h`

``` Bash
$ terrascan scan -h
Terrascan

Detect compliance and security violations across Infrastructure as Code to mitigate risk before provisioning cloud native infrastructure.

Usage:
  terrascan scan [flags]

Flags:
      --config-only               will output resource config (should only be used for debugging purposes)
  -h, --help                      help for scan
  -d, --iac-dir string            path to a directory containing one or more IaC files (default ".")
  -f, --iac-file string           path to a single IaC file
  -i, --iac-type string           iac type (helm, k8s, kustomize, terraform)
      --iac-version string        iac version (helm: v3, k8s: v1, kustomize: v3, terraform: v12)
  -p, --policy-path stringArray   policy path directory
  -t, --policy-type strings       policy type (all, aws, azure, gcp, github, k8s) (default [all])
  -r, --remote-type string        type of remote backend (git, s3, gcs, http)
  -u, --remote-url string         url pointing to remote IaC repository
      --use-colors string         color output (auto, t, f) (default "auto")

Global Flags:
  -c, --config-path string   config file path
  -l, --log-level string     log level (debug, info, warn, error, panic, fatal) (default "info")
  -x, --log-type string      log output type (console, json) (default "console")
  -o, --output string        output type (json, yaml, xml) (default "yaml")
```

By default Terrascan will output YAML. This can be changed to JSON or XML by using the `-o` flag.

Terrascan will exit 3 if any issues are found.

#### Scanning code remotely

Terrascan can download and scan remote repositories/code sources by using the `-r` and `-u` flags. Here's and example:

``` Bash
$ terrascan scan -t aws -r git -u git@github.com:accurics/KaiMonkey.git//terraform/aws
```

The URLs for the remote should follow similar naming as the source argument for modules in Terraform. More details [here](https://www.terraform.io/docs/modules/sources.html).

#### Helm

Helm chart can be scanned by specifying "helm" on the -i flag as follows:

```
$ terrascan scan -i helm
```

This command will recursively look for Chart.yaml files in the current directory and scans rendered .yaml, .yml, .tpl template files found under the corresponding /templates directory.

A specific directory to scan can be specified using the `-d` flag. The Helm IaC provider does not support scanning of individual files using the `-f` flag.


#### Kustomize

Kustomize chart can be scanned by specifying "kustomize" on the -i flag as follows:

```
$ terrascan scan -i kustomize
```

This command will look for a kustomization.yaml file in the current directory and scans rendered .yaml or .yml template files.

A specific directory to scan can be specified using the `-d` flag. The Kustomize IaC provider does not support scanning of individual files using the `-f` flag.


### CLI Output types
#### Violations
Terrascan's default output is a list of violations present in the scanned IaC.
``` Bash
$ terrascan scan -t aws
results:
  violations:
  - rule_name: scanOnPushDisabled
    description: Unscanned images may contain vulnerabilities
    rule_id: AWS.ECR.DataSecurity.High.0578
    severity: MEDIUM
    category: Data Security
    resource_name: scanOnPushDisabled
    resource_type: aws_ecr_repository
    file: ecr.tf
    line: 1
  count:
    low: 0
    medium: 1
    high: 0
    total: 1
```
##### Resource Config
Terrascan while scanning the IaC, loads all the IaC files, creates a list of resource configs and then processes this list to report violations. For debugging purposes, it is possible to print this resource configs list as an output by providing the `--config-only` flag to the `terrascan scan` command.
``` Bash
$  terrascan scan -t aws --config-only
aws_ecr_repository:
- id: aws_ecr_repository.scanOnPushDisabled
  name: scanOnPushDisabled
  source: ecr.tf
  line: 1
  type: aws_ecr_repository
  config:
    image_scanning_configuration:
    - scan_on_push:
        value: {}
    image_tag_mutability: MUTABLE
    name: test
- id: aws_ecr_repository.scanOnPushNoSet
  name: scanOnPushNoSet
  source: ecr.tf
  line: 10
  type: aws_ecr_repository
  config:
    image_tag_mutability: MUTABLE
    name: test
```

### Server mode
Server mode will execute Terrascan's API server. This is useful when using Terrascan to enforce policies in a centralized way. By default the server will be started listening in port 9010 and supports the following routes:

* GET /health = Returns the health status of the server
* POST /v1/{iacType}/{iacVersion}/{policyType}/local/file/scan} = The payload for this request should include a `file` parameter with the value being the contents of the file.

You can launch server mode by executing the Terrascan CLI or with the Docker container:

``` Bash
$ terrascan server
```
You can also launch Terrascan using Docker:

``` Bash
$ docker run --rm --name terrascan -p 9010:9010 accurics/terrascan
```

Here's an example of how to send a request to the Terrascan server using curl:

``` Bash
$ curl -i -F "file=@aws_cloudfront_distribution.tf" localhost:9010/v1/terraform/v12/aws/local/file/scan
HTTP/1.1 100 Continue

HTTP/1.1 200 OK
Date: Sun, 16 Aug 2020 02:45:35 GMT
Content-Type: text/plain; charset=utf-8
Transfer-Encoding: chunked

{
  "results": {
    "violations": [
      {
        "rule_name": "cloudfrontNoGeoRestriction",
        "description": "Ensure that geo restriction is enabled for your Amazon CloudFront CDN distribution to whitelist or blacklist a country in order to allow or restrict users in specific locations from accessing web application content.",
        "rule_id": "AWS.CloudFront.Network Security.Low.0568",
        "severity": "LOW",
        "category": "Network Security",
        "resource_name": "s3-distribution-TLS-v1",
        "resource_type": "aws_cloudfront_distribution",
        "file": "terrascan-492583054.tf",
        "line": 7
      },
      {
        "rule_name": "cloudfrontNoHTTPSTraffic",
        "description": "Use encrypted connection between CloudFront and origin server",
        "rule_id": "AWS.CloudFront.EncryptionandKeyManagement.High.0407",
        "severity": "HIGH",
        "category": "Encryption and Key Management",
        "resource_name": "s3-distribution-TLS-v1",
        "resource_type": "aws_cloudfront_distribution",
        "file": "terrascan-492583054.tf",
        "line": 7
      },
      {
        "rule_name": "cloudfrontNoHTTPSTraffic",
        "description": "Use encrypted connection between CloudFront and origin server",
        "rule_id": "AWS.CloudFront.EncryptionandKeyManagement.High.0407",
        "severity": "HIGH",
        "category": "Encryption and Key Management",
        "resource_name": "s3-distribution-TLS-v1",
        "resource_type": "aws_cloudfront_distribution",
        "file": "terrascan-492583054.tf",
        "line": 7
      },
      {
        "rule_name": "cloudfrontNoLogging",
        "description": "Ensure that your AWS Cloudfront distributions have the Logging feature enabled in order to track all viewer requests for the content delivered through the Content Delivery Network (CDN).",
        "rule_id": "AWS.CloudFront.Logging.Medium.0567",
        "severity": "MEDIUM",
        "category": "Logging",
        "resource_name": "s3-distribution-TLS-v1",
        "resource_type": "aws_cloudfront_distribution",
        "file": "terrascan-492583054.tf",
        "line": 7
      },
      {
        "rule_name": "cloudfrontNoSecureCiphers",
        "description": "Secure ciphers are not used in CloudFront distribution",
        "rule_id": "AWS.CloudFront.EncryptionandKeyManagement.High.0408",
        "severity": "HIGH",
        "category": "Encryption and Key Management",
        "resource_name": "s3-distribution-TLS-v1",
        "resource_type": "aws_cloudfront_distribution",
        "file": "terrascan-492583054.tf",
        "line": 7
      }
    ],
    "count": {
      "low": 1,
      "medium": 1,
      "high": 3,
      "total": 5
    }
  }
}
```

### Config File
The `-c` or `--config-path` global variable allows you to provide a TOML configuration file for Terrascan. This file can be used to configure the webhook notifications. Here's an example configuration:

``` TOML
[notifications]
    [notifications.webhook]
    url = "https://httpbin.org/post"
    token = "my_auth_token"
```

### Logging
Logging can be configured by using the `-l` or `--log-level` flags with possible values being: debug, info, warn, error, panic, or fatal. This defaults to "info".

In addition to the default "console" logs, the logs can be configured to be output in JSON by using the `-x` or `--log-type` flag with the value of `json`.
