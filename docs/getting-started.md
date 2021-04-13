# Quickstart
Terrascan is a static code analyzer for Infrastructure as Code. It can be installed and run in a number of different ways, and is most often used in automated pipelines to identify policy violations before insecure infrastructure is provisioned.

This document will show you two different methods for running Terrascan, one that installs Terrascan locally and another that leverages a Docker container.  We will scan some vulnerable IaC to learn how that works and understand the output.  The [usage page](./usage.md) includes more information about how Terrascan can be used.

## Running Terrascan

Terrascan is a portable executable that does not strictly require installation, and is also available as a container image in Docker Hub.  The following sections explain how to use it as a [native executable](#native-executable) and how to use the [Docker image](#using-docker).  Choose which one you'd like to start with; you don't need to do both.

### Native executable
Terrascan's [release page](https://github.com/accurics/terrascan/releases) includes builds for common platforms.  Just download and extract the package for your platform.  For example, if you use a Mac you might do this:

``` Bash
$ curl --location https://github.com/accurics/terrascan/releases/download/v1.4.0/terrascan_1.4.0_Darwin_x86_64.tar.gz --output terrascan.tar.gz
$ tar xzf terrascan.tar.gz
$ ./terrascan version
version: v1.4.0
```

If you want to use this executable for the rest of this quickstart, it will help to create an alias or install the executable onto your path. For example with bash you could do something like this:

``` Bash
$ sudo install terrascan /usr/local/bin
```

or:

``` Bash
$ alias terrascan="`pwd`/terrascan"
```

### Using Docker
Terrascan is also available as a Docker image in Docker Hub and can be used as follows assuming you have Docker installed:

``` Bash
$ docker run --rm accurics/terrascan version
version: v1.4.0
```

If you want to use the Docker image for the rest of this quickstart, it will help to create an alias, script or batch file that reduces the typing necessary.  For example with bash you could do something like this:

``` Bash
$ alias terrascan="docker run --rm -it -v "$(pwd):/iac" -w /iac accurics/terrascan"
```

This command includes a few extra options that ensure Terrascan has access to the current directory when it is run.

## Scanning with Terrascan

When scanning with Terrascan, it defaults to scanning all supported cloud providers on Terraform HCL files.

Our [KaiMonkey project](https://github.com/accurics/KaiMonkey) contains some vulnerable Terraform files to scan. Try the following:

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

Now that you understand how to run Terrascan, explore the other options available.  The [usage page](./usage.md) covers the options in detail, including other IaC providers (e.g. Kubernetes, Helm, etc.), limiting the scan to specific directories or files, and outputting in different formats.

# Related resources

* [Terrascan Policy Reference](../policies.md)
* The [usage guide](./usage.md) explains general usage and how to scan other types of IaC, such as Kubernetes, Helm, and Kustomize.
* The [CI/CD](../cicd.md) page explains how to integrate Terrascan on CI/CD pipelines.
