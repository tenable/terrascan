# Getting Started
Terrascan is a static code analyzer for Infrastructure as Code tooling. It can executed with the native binary/executable or by using the [`docker`](#using-docker) container.

## Installation
Terrascan's binary can be found on the package for each [release](https://github.com/accurics/terrascan/releases). Here's an example of how to install it:

``` Bash linenums="1"
$ curl --location https://github.com/accurics/terrascan/releases/download/v1.0.0/terrascan_darwin_amd64.zip --output terrascan_darwin_amd64.zip
$ unzip terrascan_darwin_amd64.zip
Archive:  terrascan_darwin_amd64.zip
  inflating: terrascan
$ install terrascan /usr/local/bin
$ terrascan --help
```

### Using Docker
Terrascan is available as a Docker image and can used as follows:

``` Bash linenums="1"
$ docker run accurics/terrascan
```

### Installing on macOS
For Mac users, Terrascan can be installed using Homebrew:

``` Bash linenums="1"
brew install terrascan
```

### Building Terrascan
Terrascan can be built locally. This is helpful if you want to be on the latest version or when developing Terrascan.

``` Bash linenums="1"
$ git clone git@github.com:accurics/terrascan.git
$ cd terrascan
$ make build
$ ./bin/terrascan
```

## Scanning

By typing `terrascan` without flags or other arguments, you can display the usage information.

``` Bash linenums="1"
$ terrascan

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


``` Bash linenums="1"
$ terrascan -cloud aws -iac terraform -iac-version v12 -p $REGO_POLICIES -d . --output json
```

### Example scanning Terraform (HCL2)

Here's an example of scanning Terraform HCL2 files containing AWS resources:

``` Bash linenums="1"
terrascan -cloud aws -d ~/iac_folder
```
In the example above, the `-cloud` flag is used to specify AWS as the cloud provider and the `-d` flag is used to specify the directory to scan.

### Launch Terrascan in server mode

To launch Terrascan in server mode you can execute the following:

``` Bash linenums="1"
terrascan -server
```
