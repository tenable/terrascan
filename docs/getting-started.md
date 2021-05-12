# Getting started

Terrascan is a static code analyzer for Infrastructure as Code. It can be installed and run in a number of different ways, and is most commonly used in automated pipelines to identify policy violations before insecure infrastructure is provisioned.

## Quickstart

Quickly get started with common tasks on Terrascan:

- [Running Terrascan for the First Time](#running-terrascan-for-the-first-time): Describes different methods for running Terrascan, one that installs Terrascan locally and another that leverages a Docker container. See  below for details.

- [Scanning with Terrascan](#scanning-with-terrascan) - Explains how to scan vulnerable IaC using Terrascan with the an example including the scan output. See  below for details.

## Running Terrascan for the First Time

Terrascan is a portable executable that does not strictly require installation, and is also available as a container image in Docker Hub. You can use Terrascan in two different methods based on your preference:

1. [Installing Terrascan locally](#Native_executable)
2. [Using a Docker container](#Using_Docker)

The following sections explain how to use it as a [native executable](#native-executable) and how to use the [Docker image](#docker-image). You can choose from one of the following:

### Native executable
Terrascan's [release page](https://github.com/accurics/terrascan/releases) includes latest version of builds for common platforms.  Download and extract the package for your platform. Follow instructions that apply to your platform:

#### macOS and Linux
Download the latest version of builds for macOS and enter the following command.
**Note:** for linux, replace `Darwin` with `Linux`


``` Bash
$ curl -L "$(curl -s https://api.github.com/repos/accurics/terrascan/releases/latest | grep -o -E "https://.+?_Darwin_x86_64.tar.gz")" > terrascan.tar.gz
$ tar -xf terrascan.tar.gz terrascan && rm terrascan.tar.gz
$ install terrascan /usr/local/bin && rm terrascan
$ terrascan
```

If you want to use this executable for the rest of this quickstart, it will help to create an alias or install the executable onto your path. For example with bash you could do something like this:

``` Bash
$ sudo install terrascan /usr/local/bin
```

or:

``` Bash
$ alias terrascan="`pwd`/terrascan"
```
#### Windows

Download the latest version of builds for Windows and enter the following command:

```
tar -zxf terrascan_<version number>_Windows_x86_64.tar.gz
```

### Docker image
Terrascan is also available as a Docker image in Docker Hub and can be used as follows assuming you have Docker installed:

``` Bash
$ docker run --rm accurics/terrascan version
```

If you want to use the Docker image for the rest of this quickstart, it will help to create an alias, script or batch file that reduces the typing necessary.  For example with bash you could do something like this:

``` Bash
$ alias terrascan="docker run --rm -it -v "$(pwd):/iac" -w /iac accurics/terrascan"
```

**Note**: This command includes a few extra options to enable Terrascan has access to the current directory when it is run.

## Scanning with Terrascan

### Example of interactive scan or using CLI


In this example, the [KaiMonkey project](https://github.com/accurics/KaiMonkey) contains some vulnerable Terraform files to scan. To run a scan, follow these steps:

``` Bash
$ git clone https://github.com/accurics/KaiMonkey
...
$ cd KaiMonkey/terraform/aws
$ terrascan scan
```

By default Terrascan will output its findings in human friendly format:

``` sh
Violation Details -

	Description    :	S3 bucket Access is allowed to all AWS Account Users.
	File           :	modules/storage/main.tf
	Line           :	104
	Severity       :	HIGH
	-----------------------------------------------------------------------

	Description    :	S3 bucket Access is allowed to all AWS Account Users.
	File           :	modules/storage/main.tf
	Line           :	112
	Severity       :	HIGH
	-----------------------------------------------------------------------

	Description    :	Ensure that your RDS database has IAM Authentication enabled.
	File           :	modules/storage/main.tf
	Line           :	45
	Severity       :	HIGH
	-----------------------------------------------------------------------

	Description    :	Ensure VPC flow logging is enabled in all VPCs
	File           :	modules/network/main.tf
	Line           :	4
	Severity       :	MEDIUM
	-----------------------------------------------------------------------

	Description    :	EC2 instances should disable IMDS or require IMDSv2
	File           :	modules/compute/main.tf
	Line           :	124
	Severity       :	MEDIUM
	-----------------------------------------------------------------------

	Description    :	http port open to internet
	File           :	modules/network/main.tf
	Line           :	102
	Severity       :	HIGH
	-----------------------------------------------------------------------

	Description    :	Enabling S3 versioning will enable easy recovery from both unintended user actions, like deletes and overwrites
	File           :	modules/storage/main.tf
	Line           :	104
	Severity       :	HIGH
	-----------------------------------------------------------------------

	Description    :	Enabling S3 versioning will enable easy recovery from both unintended user actions, like deletes and overwrites
	File           :	modules/storage/main.tf
	Line           :	112
	Severity       :	HIGH
	-----------------------------------------------------------------------

	Description    :	AWS CloudWatch log group is not encrypted with a KMS CMK
	File           :	modules/compute/main.tf
	Line           :	115
	Severity       :	HIGH
	-----------------------------------------------------------------------

Scan Summary -

	File/Folder         :	/var/folders/2g/9lkfm6ld2lv350svwr15fdgc0000gn/T/x9wqg4/terraform/aws
	IaC Type            :	terraform
	Scanned At          :	2021-01-15 03:11:31.869816 +0000 UTC
	Policies Validated  :	571
	Violated Policies   :	9
	Low                 :	0
	Medium              :	2
	High                :	7
```

You should see a total of 9 violations, which are detailed in the output.

Now that you understand how to run Terrascan, you can explore various options available. The [usage page](./usage.md) covers the options in detail. For more information, see [Related resources](#related_resources).

# Related resources

* The [usage guide](./usage.md) explains general usage, how to scan other types of IaC (such as: Kubernetes, Helm, and Kustomize), List of other IaC providers (e.g. Kubernetes, Helm, etc.), instructions to limit the scan to specific directories or files, and generating the output in different formats.
* [Terrascan Policy Reference](../policies.md)
* The [CI/CD](../cicd.md) page explains how to integrate Terrascan on CI/CD pipelines.
