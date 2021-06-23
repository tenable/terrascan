# Integration with Atlantis Pull Request Automation

[Atlantis](https://www.runatlantis.io/) is a pull request automation system designed to control Terraform execution from Github commits.

You can integrate Terrascan into an Atlantis setup using one of the two ways:

* Method 1: Atlantis [Workflow-based](https://www.runatlantis.io/docs/custom-workflows.html) integration which sends scan requests to an independently running terraform server
* Method 2: Custom Atlantis container image which has an integrated Terrascan

In either scenario, the configuration of Atlantis varies from installation to installation. For instructions to install, configure, and use Atlantis, see the [Atlantis documentation](https://www.runatlantis.io/docs/).

## Method 1: Workflow-based integration
In this method, you can modify or create a custom workflow for Atlantis so your repositories will be scanned by Terrascan as part of the pull request automation.

**Requirements**

The following requirements must be met before starting the integration workflow:

* The atlantis server must have TCP connectivity to where the terrascan server is running.
* The `curl` command must be installed on the system so the `terrascan-remote-scan.sh` script can make the scan request. Atlantis's [docker image](https://hub.docker.com/r/runatlantis/atlantis/) has curl preinstalled.

## Integration steps for Workflow based integration

- Modify Workflow
- Configure the Script
- Run Atlantis

### Modify Workflow

1. Modify your workflow to call `terrascan-remote-scan.sh` during the plan stage.
2. See the 'plan' detailed below:
  - the first three `run: terraform` commands are the default for an atlantis workflow.
  >**Note**: The values for the variables `$WORKSPACE` and `$PLANFILE` referenced in the second and third run commands in the yaml below are automatically provided by atlantis
  - The fourth `run terrascan-remote-scan.sh` initiates the Terrascan scan request.

>**Note**: By default, the `terrascan-remote-scan.sh` script can be found under the `scripts` directory in this project; copy this to a location where it can be executed by the Atlantis server.
If the `terrascan-remote-scan.sh` script is not in the directory where the Atlantis server command is being run to, you will have to specify the path to the script in the fourth run command.

```yaml
repos:
- id: /.*/
  workflow: terrascan

workflows:
  terrascan:
    plan:
      steps:
        - run: terraform init -input=false -no-color
        - run: terraform workspace select -no-color $WORKSPACE
        - run: terraform plan -input=false -refresh -no-color --out $PLANFILE
        - run: terrascan-remote-scan.sh
```
### Script configuration

Modify the `terrascan-remote-scan.sh` script according your environment. The script is [located here](https://github.com/accurics/terrascan/tree/master/scripts). Open the script with your any editor of your choice and review the following six settings which is found at the top of the file:

```
TERRASCAN_SERVER=192.168.1.55
TERRASCAN_PORT=9010
IGNORE_LOW_SEVERITY=false
IAC=terraform
IAC_VERSION=v14
CLOUD_PROVIDER=aws
```
Descriptions of these settings are as follows:
* `TERRASCAN_SERVER` is the hostname or IP address of the host running the terrascan server. This will be used by the script to submit the scan request.
* `TERRASCAN_PORT` is the TCP port which Terrascan server is listening on. By default, this is `9010`.
* `IGNORE_LOW_SERVERITY` allows you to specify the scan response for low-severity findings in the code. During a scan if the `terrascan-remote-scan.sh` should fail a build if a low-severity finding is found. Some users will want to set this to `true` so they may ignore low-severity findings.
* `IAC`, `IAC_VERSION`, and `CLOUD_PROVIDER` are terrascan options. Descriptions and valid values can be found by running `terrascan scan -h`.

### Running atlantis
Run Atlantis with the `terrascan-workflow.yaml` as a [server-side repo configuration](https://www.runatlantis.io/docs/server-side-repo-config.html). The command for this depends on how you choose to [deploy Atlantis](https://www.runatlantis.io/docs/deployment.html#deployment-2).
If running the Atlantis binary directly, use the following command:

```bash
$ atlantis server \
--atlantis-url="$URL" \
--gh-user="$USERNAME" \
--gh-token="$TOKEN" \
--gh-webhook-secret="$SECRET" \
--repo-allowlist="$REPO_ALLOWLIST" \
--repo-config=terrascan-workflow.yaml
```
> **Note**: The variables in the example above must be configured separately using `export` or similar shell methods.

[comment]: <> (Instructions/link to configure would be useful here)

**Important**: Before the first pull request is processed, run Terrascan in `server` mode using the following command:

```
terrascan server
```
### Automated scanning and results

When the systems are running, if Atlantis is initiated either via a pull request, or via a comment of `atlantis plan`, Terrascan will be also be invoked as part of the atlantis plan flow. Scan results are reported as part of the pull request as comments, this notifies the reviewers before approving a requests.  If any issues are found the test will be marked as failed.

## Method 2: Custom Atlantis Container

Terrascan offers a custom container built on top of the official Atlantis container image, which allows users to run IaC scans with Terrascan, in addition to the usual atlantis usage. There's a built-in atlantis workflow configured inside the
container which is ready to be used.
The default `workflow.yaml` file used is the `atlantis/workflow.yaml` in the Terrascan repo.
Alternatively, you can also override that default workflow using the `--repo-config` flag.

### Steps to use the custom Atlantis container

In code repository, usage is exactly the same as atlantis, add a comment `atlantis plan` and `atlantis plan` to your Pull Requests to trigger the custom atlantis-terrascan workflow.

##### To use the default built-in container image:

```
docker pull accurics/terrascan_atlantis
```

##### To build your own container image:
```
docker build ./integrations/atlantis -t <image_name>
```

### Run the container:

```bash
docker run \
--env-file=<.env-file> \
-p 4141:4141 \
-v <pwd>/config_data/:/etc/terrascan/ \
accurics/terrascan_atlantis server \
--gh-user="$USERNAME" --gh-token="$TOKEN" --gh-webhook-secret="$SECRET" \
--repo-allowlist="$REPO_ALLOWLIST" \
-c /etc/terrascan/config.toml
```

The syntax of the Atlantis server command here is same as in [atlantis docs](https://www.runatlantis.io/docs/), except for an optional `-c` flag which can be used to specify the file path for the toml config to be used by Terrascan. Another way to provide the toml config filepath would be the TERRASCAN_CONFIG environment variable. You need to provide all the environment variables that terraform requires to operate with your respective cloud providers.

> **Note**: As a good practice, Terrascan recommends use of a [specific tag](https://hub.docker.com/r/accurics/terrascan_atlantis/tags) of the container image rather than the latest tag.

[comment]: <> (Moved the workflow yaml note to above where its mentioned)

### Running a scan

With everything configured, a local Terrascan scan will be triggered as part of the Atlantis plan workflow.
