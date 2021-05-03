# Integration with Atlantis Pull Request Automation
[Atlantis](https://www.runatlantis.io/) is a pull request automation system designed to allow control of terraform execution from github commits.

We have designed two ways to integrate terrascan into an Atlantis setup:
* Atlantis [Workflow-based](https://www.runatlantis.io/docs/custom-workflows.html) integration which sends scan requests to a independently running terraform server
* Custom Atlantis container image which has terrascan built in

In either scenario, the configuration of Atlantis is a diverse topic which will vary from installation to installation. For details around installing, configuring, and using Atlantis, please refer to the [Atlantis documentation](https://www.runatlantis.io/docs/).

## Workflow-based integration
Through this method, you will modify or create a custom workflow for atlantis so your repositories will be scanned by terrascan as part of the pull request automation.

**Requirements**
* The atlantis server must have TCP connectivity to where the terrascan server is running.
* The `curl` command needs to be installed on the system so the `terrascan-remote-scan.sh` script can make the scan request. Atlantis's [docker image](https://hub.docker.com/r/runatlantis/atlantis/) has curl preinstalled.

### Workflow
Next, you will need to modify your workflow to call `terrascan-remote-scan.sh` during the plan stage. In the plan below, the first three `run: terraform` commands are the default for an atlantis workflow; the fourth `run terrascan-remote-scan.sh` is where the terrascan scan is requested. The `terrascan-remote-scan.sh` script can be found under the `scripts` directory in this project; you will need to copy it to a location where it can be executed by the atlantis server. If the `terrascan-remote-scan.sh` script is not in the directory where the atlantis server command is being run to, you will have to specify the path to the script.

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
(the variables `$WORKSPACE` and `$PLANFILE` referenced in the above yaml are populated by atlantis)

### Script configuration
Next, the `terrascan-remote-scan.sh` script will need to be modified for your environment. The script is [located here](https://github.com/accurics/terrascan/tree/master/scripts). Open the script with your favorite editor and review the following six settings near the top of the file:

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
* `TERRASCAN_PORT` is the TCP port which terrascan server is listening on. By default, this is `9010`.
* `IGNORE_LOW_SERVERITY` specifies if the `terrascan-remote-scan.sh` should fail a build if a low-severity finding is found. Some users will want to set this to `true` so they may ignore low-severity findings.
* `IAC`, `IAC_VERSION`, and `CLOUD_PROVIDER` are terrascan options. Descriptions and valid values can be found by running `terrascan scan -h`.

### Running atlantis
Run atlantis with your terrascan-workflow.yaml as a [server-side repo configuration](https://www.runatlantis.io/docs/server-side-repo-config.html). This can depend on how you choose to [deploy atlantis](https://www.runatlantis.io/docs/deployment.html#deployment-2).
If running the atlantis binary directly, note the following command:

```bash
$ atlantis server \
--atlantis-url="$URL" \
--gh-user="$USERNAME" \
--gh-token="$TOKEN" \
--gh-webhook-secret="$SECRET" \
--repo-allowlist="$REPO_ALLOWLIST" \
--repo-config=terrascan-workflow.yaml
```
(the variables in the example above must be set separately using `export` or similar shell methods)

Additionally, before the first pull request is processed, terrascan must be running in `server` mode:

```
terrascan server
```

Once the systems are running, when atlantis is called via pull request, or a comment of `atlantis plan`, terrascan will be called as part of the atlantis plan flow. Scan results will be placed in a comment on the pull request, and if issues are found the test will be marked as failed.

## Custom Atlantis Container

We have come up with a custom container built on top of the official atlantis container image, that will allow users to
run IaC scans with terrascan, besides the usual atlantis usage. There's a default atlantis workflow setup inside in the
container to be used. Users may also override that default workflow using the --repo-config flag.

### Usage

On the code repository's end, usage is exactly the same as atlantis, you comment `atlantis plan` and `atlantis plan` on
your Pull Requests to trigger the custom atlantis-terrascan workflow.

##### To use our container image:
```
docker pull accurics/terrascan_atlantis
```

##### To build your own container image:
```
docker build ./atlantis -t <image_name>
```

###### Running the container:

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

The syntax of the atlantis server command here is same as in [atlantis docs](https://www.runatlantis.io/docs/),
except for an optional `-c` flag which can be used to pass in the toml config filepath for terrascan.
Another way to provide the toml config filepath would be the TERRASCAN_CONFIG environment variable.

You need to provide all the environment variables that terraform requires to operate with your respective cloud providers.
It's a good practice to use a [specific tag](https://hub.docker.com/r/accurics/terrascan_atlantis/tags) of the container
image rather than the latest tag.

The default workflow.yaml file used is the `atlantis/workflow.yaml` in this repo. You can override the default workflow
using the `--repo-config` flag. It will be up to you how you want to trigger `terrascan` for your usage. You can do
something along the lines of `atlantis/workflow.yaml` and `atlantis/terrascan.sh` itself.
