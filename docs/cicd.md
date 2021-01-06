# Integrating Terrascan into CI/CD

Terrascan can be integrated into CI/CD pipelines to enforce security best practices as codified in the OPA rego policies included as part of Terrascan or any custom policies. This section contains examples on how to configure Terrascan in popular CI/CD tooling.

## GitHub Actions

Terrascan can be configured as a job within GitHub actions workflows. A straightforward way to accomplish this is by using the [super-linter](https://github.com/github/super-linter) GitHub action which includes Terrascan. Note that at the moment super-linter only supports scanning Terraform HCL files.

When using super-linter you can pass the environment variable "VALIDATE_TERRAFORM_TERRASCAN: true" to ensure that Terraform configuration files are evaluated using Terrascan. To configure your GitHub actions workflow a file with the below YAML content can be included within the .github/workflows/ directory of your repository.

``` YAML
---
name: Scan Code Base
on:
  push:
  pull_request:
    branches: [master]
jobs:
  build:
    name: Scan Code Base
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v2
      - name: Scan Code Base
        uses: github/super-linter@v3
        env:
          VALIDATE_ALL_CODEBASE: true
          DEFAULT_BRANCH: master
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          VALIDATE_TERRAFORM_TERRASCAN: true
```

Documentation on the GitHub actions workflow syntax is available [here](https://help.github.com/en/articles/workflow-syntax-for-github-actions).


## GitLab CI

[GitLab CI](https://docs.gitlab.com/ee/ci/README.html) can use [Docker images](https://docs.gitlab.com/ee/ci/docker/using_docker_images.html) as part of a pipeline. We can take advantage of this functionality and use Terrascan's docker image as part of your [pipeline](https://docs.gitlab.com/ee/ci/pipelines/) to scan infrastructure as code.

To do this you can update your .gitlab-ci.yml file to use the "accurics/terrascan:latest" image with the ["bin/sh", "-c"] entrypoint. Terrascan can be found on "/go/bin" in the image and you can use any [Terrascan command line options](http://ubusvr:8000/getting-started/usage/#terrascan-commands) according to your needs. Here's an example .gitlab-ci.yml file:

``` YAML
stages:
  - scan

terrascan:
  image:
    name: accurics/terrascan:latest
    entrypoint: ["/bin/sh", "-c"]
  stage: scan
  script:
    - /go/bin/terrascan scan .
```


## ArgoCD Application PreSync Hooks


Terrascan can be configured as a job during the application sync process using [resource hooks](https://argoproj.github.io/argo-cd/user-guide/resource_hooks). The PreSync resource hook is the best way to evaluate the kubernetes deployment configuration and report any violations. 


![picture](img/terrascan-argo-cd-pipeline.png)

See example hooks yaml where one can simply add it to an existing kubernetes configuration.


``` YAML
apiVersion: batch/v1
kind: Job
metadata:
  generateName: terrascan-hook-
  annotations:
    argocd.argoproj.io/hook: PreSync
spec:
  ttlSecondsAfterFinished: 3600
  template:
    spec:
      volumes:
      - name: secret-volume
        secret:
          secretName: ssh-key-secret    
      containers:
      - name: terrascan-argocd
        image: accurics/terrascan-argocd:latest
        command: ["/bin/ash", "-c"]
        args:
        - >
          cp /etc/secret-volume/ssh-privatekey /home/terrascan/.ssh/id_rsa &&
          chmod 400 /home/terrascan/.ssh/id_rsa &&
          /go/bin/terrascan scan -r git -u git@bitbucket.org:example/argo-cd-nginx-sample.git -i k8s -t k8s | /home/terrascan/bin/notify_slack.sh webhook-tests argo-cd https://hooks.slack.com/services/TXXXXXXXX/XXXXXXXXXXX/0XXXXXXXXXXXXXXXXXX
        volumeMounts:
          - name: secret-volume
            readOnly: true
            mountPath: "/etc/secret-volume"
      restartPolicy: Never
  backoffLimit: 1
```

As shown, the PreSync requires access to the repository and using the same branch (default) as the Argo CD application pipeline. 

For non-public repositories, the private key needs to be added as a kubernetes secret.

``` CONSOLE
  kubectl create secret generic ssh-key-secret \
    --from-file=ssh-privatekey=/path/to/.ssh/id_rsa \
    --from-file=ssh-publickey=/path/to/.ssh/id_rsa.pub
```

Configuring the job to delete only after the specified time see `ttlSecondsAfterFinished` will allow users to check for violations in the User Interface, the alternative is through **notifications**.

![picture](img/terrascan-argo-cd-resource-hook-logs.png)

Below is the full example of building the terrascan-argo-cd integration container. 

`known_hosts`

```
github.com,192.30.255.113 ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==
bitbucket.org,104.192.141.1 ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAubiN81eDcafrgMeLzaFPsw2kNvEcqTKl/VqLat/MaB33pZy0y3rJZtnqwR2qOOvbwKZYKiEO1O6VqNEBxKvJJelCq0dTXWT5pbO2gDXC6h6QDXCaHo6pOHGPUy+YBaGQRGuSusMEASYiWunYN0vCAI8QaXnWMXNMdFP3jHAJH0eDsoiGnLPBlBp4TNm6rYI74nMzgz3B9IikW4WVK+dc8KZJZWYjAuORU3jc1c/NPskD2ASinf8v3xnfXeukU0sJ5N6m5E8VLjObPEO+mN2t/FZTMZLiFqPWc/ALSqnMnnhwrNi2rbfg/rd/IpL8Le3pSBne8+seeFVBoGqzHM9yXw==
gitlab.com,172.65.251.78 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBFSMqzJeV9rUzU4kWitGjeR4PWSa29SPqJ1fVkhtj3Hw9xjLVXVYrU9QlYWrOLXBpQ6KWjbjTDTdDkoohFzgbEY=
```

`notify_slack.sh` the notification example, for this case we used slack.

``` SH
#!/bin/ash

function send_slack_notificaton {
  channel=$1
  username=$2
  slack_hook=$3

  curl -X POST --data-urlencode payload="{\"channel\": \"#${channel}\", \"username\": \"${username}\", \"text\": \" \`\`\` $(cat results.out) \`\`\` \", \"icon_emoji\": \":ghost:\"}" ${slack_hook}
}

if [ -p /dev/stdin ]; then
  echo "processing terrascan results"
  while IFS= read line; do
          echo "${line}" | tr '\\"' ' ' >> results.out
  done

  cat results.out

  send_slack_notificaton $1 $2 $3

  echo "notification exit code: $?"
else
  echo "no response skipping"
fi
```

`Dockerfile`

``` SH
  FROM accurics/terrascan:929e377

  ENTRYPOINT []

  USER root 

  RUN apk add --no-cache openssh curl

  WORKDIR /home/terrascan

  RUN mkdir -p .ssh && mkdir -p bin

  COPY known_hosts .ssh

  COPY notify_slack.sh bin/

  RUN chown -R terrascan:terrascan .ssh && \
      chown -R terrascan:terrascan bin && \
      chmod 400 .ssh/known_hosts && \
      chmod u+x bin/notify_slack.sh

  USER terrascan

  CMD ["ash"]
```