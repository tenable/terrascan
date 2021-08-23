# Using Terrascan in command line mode

This section contains the following information:

* [Terrascan commands](#terrascan-commands)
* [Scanning](#scanning) with examples
* [Configuring the output format for a scan](#configuring-the-output-format-for-a-scan)


The following is a description of all the commands available. Terrascan's interface is divided into subcommands as follows:

*   `init` = Initializes Terrascan by downloading the latest Rego policies into ~/.terrascan. The scan command will implicitly run this before a scan if it detects that there are no policies found.
*   `scan` = scans Infrastructure as code files based on the policies contained within the ".terrascan" directory
*   `server` = Starts the Terrascan's API server
*   `help` = You can view the usage menu by typing `help` or using the `-h` flag on any subcommand (e.g. `terrascan init -h`). You can also view this by typing `terrascan` without flags or other arguments.

## Terrascan Commands

``` Bash
$ terrascan
Terrascan

Detect compliance and security violations across Infrastructure as Code to mitigate risk before provisioning cloud native infrastructure.
For more information, please visit https://docs.accurics.com

Usage:
  terrascan [command]

Available Commands:
  help        Provides usage info about any command
  init        Initialize Terrascan
  scan        Start scan to detect compliance and security violations across Infrastructure as Code.
  server      Run Terrascan as an API server
  version     Shows the Terrascan version you are currently using.

Flags:
  -c, --config-path string   config file path
  -h, --help                 help for terrascan
  -l, --log-level string     log level (debug, info, warn, error, panic, fatal) (default "info")
  -x, --log-type string      log output type (console, json) (default "console")
  -o, --output string        output type (human, json, yaml, xml) (default "human")

Use "terrascan [command] --help" for more information about a command.
```

## Initializing (optional)

The initialization process downloads the latest policies from the [repository](https://github.com/accurics/terrascan) into `~/.terrascan`.
By default the policies are installed here: `~/.terrascan/pkg/policies/opa/rego` and are fetched while scanning an IaC.
Use the following command to start the initialization process if you are updating the policies:

``` Bash
$ terrascan init
```
>**Note**: The `init` command is implicitly executed if the `scan` command does not find policies while executing.

## Scanning

If the `scan` command is used with no arguments (as shown below), the scan will include all supported cloud providers on Terraform HCL files:

``` Bash
$ terrascan scan
```

The `scan` command supports flags to configure the following:
- Specify a directory to be scanned
- Specify a particular IaC file to be scanned
- Configure IaC provider type
- Directory path to policies
- Specify policy type
- Retrieve vulnerability scanning results from docker images referenced in IaC

The full list of flags for the scan command can be found by typing
`terrascan scan -h`

### Scanning current directory containing terraform files for AWS Resources

The following will scan the current directory containing Terraform HCL2 files for AWS resources:

``` Bash
$ terrascan scan -t aws
```
### Scanning for a specific IaC provider
By default, Terrascan defaults to scanning Terraform HCL files. Use the `-i` flag to change the IaC provider. Here's an example of scanning kubernetes yaml files:

```Bash
$ terrascan scan -i k8s
```
### Scanning code remotely
Terrascan can be installed remotely to scan remote repositories or code resources using the `-r` and `-u` flags. Here's an example:

``` Bash
$ terrascan scan -t aws -r git -u git@github.com:accurics/KaiMonkey.git//terraform/aws
```

> **Important**: The URLs for the remote repositories should follow similar naming conventions as the source argument for modules in Terraform. For more details, see [this article](https://www.terraform.io/docs/modules/sources.html).

#### Scanning private Terraform module repositories
When scanning Terraform code, Terrascan checks for the availability of the file `~/.terraformrc`. This file contains credential information to authenticate a private terraform module registry. If this file is present, Terrascan will attempt to use the credentials when authenticating the private repository. For more details on the format of this file, please see Terraform's [config file documentation](https://www.terraform.io/docs/cli/config/config-file.html).

## Configuring the output format for a scan

By default, Terrascan output is displayed in a human friendly format. Use the `-o` flag to change this to **YAML**, **JSON**, **XML**, **JUNIT-XML** and **SARIF** formats.

> **Note**: Terrascan will exit with an error code if any errors or violations are found during a scan.

> #### List of possible Exit Codes
> | Scenario      | Exit Code |
> | ----------- | ----------- |
> | scan summary has errors and violations | 5 |
> | scan summary has errors but no violations | 4 |
> | scan summary has violations but no errors | 3 |
> | scan summary has no violations or errors | 0 |
> | scan command errors out due to invalid inputs | 1 |


Terrascan's output is a list of security violations present in the scanned IaC files. The example below is terrascan's output in YAML.
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


### Scanning a Helm Chart

Helm chart can be scanned by specifying "helm" on the -i flag as follows:

```
$ terrascan scan -i helm
```

This command will recursively look for `Chart.yaml` files in the current directory and scan rendered `.yaml`, `.yml`, `.tpl` template files found under the corresponding `/templates` directory.

A specific directory to scan can be specified using the `-d` flag. The Helm IaC provider does not support scanning of individual files using the `-f` flag.

### Scanning a Kustomize Chart

A Kustomize chart can be scanned by specifying "kustomize" on the -i flag as follows:

```
$ terrascan scan -i kustomize
```

This command looks for a `kustomization.yaml` file in the current directory and scans rendered .yaml or .yml template files.

Terrascan considers Kustomize v4 as the default version. Other supported versions (v2 and v3) of Kustomize could be scanned by specifying --iac-version flag as follows:

```
$ terrascan scan -i kustomize --iac-version v2
```

Scanning v2 and v3 requires the corresponding Kustomize binary and the path to the binary must be specified in the `KUSTOMIZE_<VERSION>` ENV variable.

e.g: For --iac-version v2, we need to have:

    KUSTOMIZE_V2=path/to/kustomize/v2/binary

To install Kustomize one can use [this script](https://github.com/accurics/terrascan/tree/master/scripts/install_kustomize.sh)

A specific directory to scan can be specified using the `-d` flag. The Kustomize IaC provider does not support scanning of individual files using the `-f` flag.

### Scanning a Dockerfile

A Dockerfile can be scanned by specifying "docker" on the -i flag as follows:

```
$ terrascan scan -i docker
```

This command looks for a `Dockerfile` in the current directory and scans that file.

A specific directory to scan can be specified using the `-d` flag. With the `-d` flag, it will check for all the docker files (named as `Dockerfile`) in the provided directory recursively. A specific dockerfile can be scanned using `-f` flag by providing a path to the file.

### Retrieve Docker Image Vulnerabilities

Terrascan can display vulnerabilities for Docker images present in the IaC files being scanned by specifying the `--find-vuln` flag along with the scan command as follows:

```
$ terrascan scan -i <IaC Provider> --find-vuln
```

This command looks for all the Docker images present in the IaC files being scanned and retrieves any vulnerabilities as reported by it's container registry. Supported container registries include: AWS Elastic Container Registry (ECR), Azure Container Registry, Google Container Registry, and Google Artifact Registry.

The following environment variables are required when connecting to the container registries:

#### AWS Elastic Container Registry (ECR)

ECR requires your environment to be configured similar to the requirements of AWS's SDK. For example, the `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION` environment variables can be set when connecting to AWS using API keys for an AWS user. More information [here](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html).

#### Google Container Registry and Artifact Registry

Terrascan requires a service account with access to the Container Analysis and Container Registry permissions. The `GOOGLE_APPLICATION_CREDENTIALS` environment variable can be set to the path of the service account's key when scanning. More information about GCP authentication available [here](https://cloud.google.com/docs/authentication/getting-started).

#### Azure Container Registry

When integrating vulnerability results from Azure, Terrascan requires the `AZURE_AUTH_LOCATION`, and `AZURE_ACR_PASSWORD` environment variables.

The `AZURE_AUTH_LOCATION` should contain the path to your azure authentication json. You can generate this as follows:

 ``` Bash
 az ad sp create-for-rbac --sdk-auth > azure.auth
 ```

After generating the file, set the `azure.auth` file path as the `AZURE_AUTH_LOCATION` environment variable. More information about using file based authentication for the Azure SDK available [here](https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization#use-file-based-authentication).

Terrascan also requires the password to the registry set into the `AZURE_ACR_PASSWORD` environment variable. This can be fetched using the az cli as follows:

``` Bash
az acr credential show --name RegistryName
```

### Resource Config
While scanning a IaC, Terrascan loads all the IaC files, creates a list of resource configs and then processes this list to report violations. For debugging purposes, you can print this resource configs list as an output by using the `--config-only` flag to the `terrascan scan` command.

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
## More details on scan command

### List of options for scan command:

| Flag      | Description | Options (default highlighted )
| ----------- | ----------- |------------|
| -h | Help for scan command | See list of all flags supported with descriptions, default options in all commands are highlighted in bold|
| -d | Use this to scan a specific directory. Use "." for current directory | AWS, GCP, Azure, and GitHub|
| -f | Use this command to scan a specific file | <tbd any formats/limitations for example file size> |
| -i type  | Use this to change the IaC provider | arm, cft, docker, helm, k8s, kustomize, **terraform**|
| -i version  | Use this in conjunction with `- i type` to specify the version of IaC provider | Supported versions of each IaC are: `arm: v1, cft: v1, docker: v1, helm: v3, k8s: v1, kustomize: v2, v3, v4, terraform: v12, v13, v14, v15`|
| -p | Use this to specify directory path for policies | By default policies are installed here: <tbd specify a default path> |
| -t  | Use this to specify individual cloud providers | **all**, aws, azure, gcp, github, k8s|
| -r | Use this to specify directory path for remote backend | git, s3, gcs, http |
| -u | Use this to specify directory URL for remote IaC repositories | see options below |
| |scan-rules|Specify rules to scan, example: --scan-rules="ruleID1,ruleID2"|
| |skip-rules|Specify one or more rules to skip while scanning. Example: --skip-rules="ruleID1,ruleID2"|
| |use-colours |Configure the color for output (**auto**, t, f) |
|--non-recursive |Use this for non recursive directories and modules scan | By default directory is scanned recursively, if this flag is used then only provided root directory will be scanned|
|--use-terraform-cache |Use this to refer terraform remote modules from terraform init cache rather than downloading | By default remote module will be downloaded in temporary directory. If this flag is set then modules will be refered from terraform init cache if module is not present in terraform init cache it will be downloaded. Directory will be scanned non recurively if this flag is used.(applicable only with terraform IaC provider)|
| --find-vuln | find vulnerbilities | Use this to fetch vulnerabilities identified on the registry for docker images present in IaC the files scanned |
| -v | verbose | Displays violations with all details |

| Global flags | Description | Options |
| ----------- | ----------- |------------|
| -c | Use this to specify config file settings | Format supported is `*.TOML` |
| -l | Use this to specify what log settings | debug, **info**, warn, error, panic, fatal  |
| -x | Use this to specify the log file format | **console**, json |
| -o | Use this to specify the scan output type | **human**, json, yaml, xml, junit-xml, sarif, github-sarif |



### Full help for scan command:


``` Bash
$ terrascan scan -h
Terrascan

Detect compliance and security violations across Infrastructure as Code to mitigate risk before provisioning cloud native infrastructure.

Usage:
  terrascan scan [flags]

Flags:
     --categories strings        list of categories of violations to be reported by terrascan (example: --categories="category1,category2")
      --config-only               will output resource config (should only be used for debugging purposes)
      --find-vuln                 fetches vulnerabilities identified in Docker images
  -h, --help                      help for scan
  -d, --iac-dir string            path to a directory containing one or more IaC files (default ".")
  -f, --iac-file string           path to a single IaC file
  -i, --iac-type string           iac type (arm, cft, docker, helm, k8s, kustomize, terraform, tfplan)
      --iac-version string        iac version (arm: v1, cft: v1, docker: v1, helm: v3, k8s: v1, kustomize: v2, v3, v4, terraform: v12, v13, v14, v15, tfplan: v1)
      --non-recursive             do not scan directories and modules recursively
  -p, --policy-path stringArray   policy path directory
  -t, --policy-type strings       policy type (all, aws, azure, gcp, github, k8s) (default [all])
  -r, --remote-type string        type of remote backend (git, s3, gcs, http, terraform-registry)
  -u, --remote-url string         url pointing to remote IaC repository
      --scan-rules strings        one or more rules to scan (example: --scan-rules="ruleID1,ruleID2")
      --severity string           minimum severity level of the policy violations to be reported by terrascan
      --show-passed               display passed rules, along with violations
      --skip-rules strings        one or more rules to skip while scanning (example: --skip-rules="ruleID1,ruleID2")
      --use-colors string         color output (auto, t, f) (default "auto")
      --use-terraform-cache       use terraform init cache for remote modules (when used directory scan will be non recursive,flag applicable only with terraform IaC provider)
  -v, --verbose                   will show violations with details (applicable for default output)

Global Flags:
  -c, --config-path string   config file path
  -l, --log-level string     log level (debug, info, warn, error, panic, fatal) (default "info")
  -x, --log-type string      log output type (console, json) (default "console")
  -o, --output string        output type (human, json, yaml, xml, junit-xml, sarif, github-sarif) (default "human")
```
