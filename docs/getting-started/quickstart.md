# Quickstart
Terrascan is a static code analyzer for Infrastructure as Code. It can be installed and run in a number of different ways, and is most often used in automated pipelines to identify policy violations before insecure infrastructure is provisioned.

This document will show you two different methods for running Terrascan, one that installs Terrascan locally and another that leverages a Docker container.  We will scan some vulnerable IaC to learn how that works and understand the output.  The [usage page](./usage.md) includes more information about how Terrascan can be used.

## Running Terrascan

Terrascan is a portable executable that does not strictly require installation, and is also available as a container image in Docker Hub.  The following sections explain how to use it as a [native executable](#native-executable) and how to use the [Docker image](#using-docker).  Choose which one you'd like to start with; you don't need to do both.

### Native executable
Terrascan's [release page](https://github.com/accurics/terrascan/releases) includes builds for common platforms.  Just download and extract the package for your platform.  For example, if you use a Mac you might do this:

``` Bash
$ curl --location https://github.com/accurics/terrascan/releases/download/v1.2.0/terrascan_1.2.0_Darwin_x86_64.tar.gz --output terrascan.tar.gz
$ tar xzf terrascan.tar.gz
$ ./terrascan version
version: v1.2.0
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
version: v1.2.0
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

By default Terrascan will output its findings in YAML format:

``` YAML
results:
    violations:
        - rule_name: s3Versioning
          description: Enabling S3 versioning will enable easy recovery from both unintended user actions, like deletes and overwrites
          rule_id: AWS.S3Bucket.IAM.High.0370
          severity: HIGH
          category: IAM
          resource_name: km_public_blob
          resource_type: aws_s3_bucket
          file: modules/storage/main.tf
          line: 112
#... lines elided ...
        - rule_name: ec2UsingIMDSv1
          description: EC2 instances should disable IMDS or require IMDSv2
          rule_id: AC-AWS-NS-IN-M-1172
          severity: MEDIUM
          category: Network Security
          resource_name: km_vm
          resource_type: aws_instance
          file: modules/compute/main.tf
          line: 124
    count:
        low: 0
        medium: 2
        high: 7
        total: 9
```

You should see a total of 9 violations, which are detailed in the output.

Now that you understand how to run Terrascan, explore the other options available.  The [usage page](./usage.md) covers the options in detail, including other IaC providers (e.g. Kubernetes, Helm, etc.), limiting the scan to specific directories or files, and outputting in different formats.

# Related resources

* [Terrascan Policy Reference](/policies/)
* The [usage guide](/getting-started/usage/) explains general usage and how to scan other types of IaC, such as Kubernetes, Helm, and Kustomize.

[//]: # (TODO: add info about CI/CD integrations * CI/CD integration )

