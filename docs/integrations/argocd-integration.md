# Integration of Terrascan with Argo CD

Terrascan can be configured as an Argo CD job during the application sync process using ArgoCD’s resource hook. The PreSync resource hook is the best way to evaluate the kubernetes deployment configuration and report any violations.

## Terrascan can be integrated with Argo CD in two ways
___
1. Use terrascan as a pre-sync hook to scan remote repositories
2. Use terrascan’s k8s admission controller along with a pre-sync that scans a configured repository with the admission controller webhook


### Method 1. Configure terrascan as a PreSync hook and scan the remote repository.
___

#### Configure a PreSync hook

The following example of a hook yaml is nearly ready to be added to an existing kubernetes configuration. To complete the configutation, you need to:
- Ensure that the secrets,  `known_hosts`, and `ssh_config` volume are relevant for your specific environment.
- Specify a terrascan image.

You can also map a slack notification script to the container which will send notifications to your Slack webhook endpoint after the embedded script scans the repo.

```yaml
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
     securityContext:
       seccompProfile:
         type: RuntimeDefault
     volumes:
       #add a configmap for the slack notification scripts
       - name: notification-scripts
         configMap:
           name: slack-notifications
       #add all required ssh keys need to clone your repos
       - name: ssh-key-secret
         secret:
           secretName: ssh-key-secret
       #add a secret for git config file   
       - name: ssh-config
         secret:
           secretName: ssh-config-secret
       #add a configmap for the ssh known_hosts file
       - name: ssh-known-hosts
         configMap:
           name: known-hosts-config
     containers:
     - name: terrascan-argocd
       image: <terrscan-image>
       resources:
         requests:
           cpu: "1"
           memory: "256Mi"
         limits:
           cpu: "1"
           memory: "256Mi"
       command: ["/bin/sh", "-c"]
       args:
       - >
         cp /etc/secret-volume/ssh-private-key /home/terrascan/.ssh/id_ed25519_github &&
         cp /etc/ssh-config-secret/ssh-config /home/terrascan/.ssh/config &&
         cp /etc/ssh-known-hosts-secret/ssh-known-hosts /home/terrascan/.ssh/known_hosts &&
         chmod -R 400 /home/terrascan/.ssh/* &&
         /go/bin/terrascan scan -r git -u <YOUR REPOSITORY PATH>-i k8s -t k8s | /data/notify_slack.sh webhook-tests argo-cd https://hooks.slack.com/services/TXXXXXXXX/XXXXXXXXXXX/0XXXXXXXXXXXXXXXXXX
       securityContext:
         seccompProfile:
           type: RuntimeDefault
         allowPrivilegeEscalation: false
         runAsNonRoot: true
         runAsUser: 101
       livenessProbe:
         exec:
           command:
           - /go/bin/terrascan
           - version
         periodSeconds: 10
         initialDelaySeconds: 10
       readinessProbe:
         exec:
           command:
           - /go/bin/terrascan
           - version
         periodSeconds: 10
       #if want to use private repo
       volumeMounts:
         - mountPath: /etc/secret-volume
           name: ssh-key-secret
           readOnly: true
         - mountPath: /etc/ssh-config-secret
           name: ssh-config
           readOnly: true
         - mountPath: /etc/ssh-known-hosts-secret
           name: ssh-known-hosts
           readOnly: true
         - mountPath: /data
           name: notification-scripts
           readOnly: true

     restartPolicy: Never
 backoffLimit: 1
```

> **Note:** As shown above, the PreSync requires access to the repository where IaC is stored, using the same branch (default) as the ArgoCD application pipeline.

To allow users to check for violations in the web interface, configure the job to delete after the specified time, using the parameter `ttlSecondsAfterFinished`. In addition, violation can be reported as webhook notifications, as shown below.

##### Example slack notification script



```sh
#!/bin/sh

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

For private repositories, the private following keys must be added as kubernetes secret:
 - `private key` and ssh `config` as Secret
 - `known_hosts`as ConfigMap
 
```
 kubectl create secret generic ssh-key-secret \
   --from-file=ssh-privatekey= < path to your private key > \
    --from-file=ssh-publickey=< path to your public key >
```

**Config-map**:

```
  kubectl  create configmap ssh-known-hosts --from-file=< path to your known hosts file >
```   

```
  kubectl  create configmap slack-notifications --from-file=< path to your notification script >
```

**ssh config secret**

```
 kubectl create secret generic ssh-config-secret \
   --from-file=< path to your ssh config file >
```   

##### Example ssh config file

```
 Host github.com
  HostName github.com
  IdentityFile ~/.ssh/id_ed25519_github
```

After configuring the presynchook yaml file, add the file to the relevant repository folder to configure Argo CD.


### Method 2. Use PreSyncHook to trigger the Terrascan Server Service
___
You can use a pre-deployed terrascan server service in K8s cluster to scan the remote repository from Argo CD PreSync hook.
To configure, follow these steps:

#### Step 1: Configure Terrascan Server webhook deployment yaml file with required keys and volumes and service to expose the controller pod.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
name: terrascan-server
labels:
  app: terrascan
spec:
replicas: 1
selector:
  matchLabels:
    app: terrascan
template:
  metadata:
    labels:
      app: terrascan
  spec:
    containers:
    - name: terrascan
      image: <TERRASCAN LATEST IMAGE>
      resources:
        limits:
          memory: "256Mi"
          cpu: "1"
      ports:
        - containerPort: 443
      livenessProbe:
        initialDelaySeconds: 30
        periodSeconds: 10
        timeoutSeconds: 5
        httpGet:
          path: /health
          port: 443
          scheme: HTTPS
      env:
        - name: K8S_WEBHOOK_API_KEY
          value: yoursecretapikey
      volumeMounts:
        - mountPath: /data/certs
          name: terrascan-certs-secret  
          readOnly: true
        - mountPath: /data/config
          name: terrascan-config
          readOnly: true
        - mountPath: /etc/secret-volume
          name: ssh-key-secret
          readOnly: true
        - mountPath: /etc/ssh-config-secret
          name: ssh-config
          readOnly: true
        - mountPath: /etc/ssh-known-hosts-secret
          name: ssh-known-hosts
          readOnly: true
      command: ["/bin/sh", "-c"]
      args:
        - >
          cp /etc/secret-volume/ssh-private-key /home/terrascan/.ssh/id_ed25519_github &&
          cp /etc/ssh-config-secret/ssh-config /home/terrascan/.ssh/config &&
          cp /etc/ssh-known-hosts-secret/ssh-known-hosts /home/terrascan/.ssh/known_hosts &&
          chmod -R 400 /home/terrascan/.ssh/* &&
          terrascan server --cert-path /data/certs/server.crt --key-path /data/certs/server.key -p 443 -l debug -c /data/config/config.toml
    volumes:
      #add all required ssh keys need to clone your repos
      - name: ssh-key-secret
        secret:
          secretName: ssh-key-secret
      #add a secret for git config file   
      - name: ssh-config
        secret:
          secretName: ssh-config-secret
      #add a configmap for the ssh known_hosts file
      - name: ssh-known-hosts
        configMap:
          name: known-hosts-config
      #add a configmap for the terrascan config.toml file    
      - name: terrascan-config
        configMap:
          name: terrascan-config
      #add a secret for the tls certificates        
      - name: terrascan-certs-secret
        secret:
          secretName: terrascan-certs-secret    
```            
**Service example**

```yaml
apiVersion: v1
kind: Service
metadata:
  name: terrascan-service
spec:
  selector:
    app: terrascan
  ports:
  - port: 443
    targetPort: 443
```

For private repositories, the following private keys needs to be added as a kubernetes secret:

 - `private key` and ssh `config` as Secret
 - `known_hosts`as ConFigmap


```
kubectl create secret generic ssh-key-secret \
  --from-file=ssh-privatekey= < path to your private key > \
  --from-file=ssh-publickey=< path to your public key >
```

```
kubectl create secret generic terrascan-certs-secret \
  --from-file= < path to your .key file > \
  --from-file= < path to your .crt file >
```

**Config-map**:

```
kubectl create configmap ssh-known-hosts --from-file=< path to your known hosts file >
```

```
kubectl create configmap terrascan-config  --from-file=<path to your config.toml file >
```
**ssh config secret**

```
kubectl create secret generic ssh-config-secret \
  --from-file=< path to your ssh config file >
```

##### Example ssh config file

```
Host github.com
  HostName github.com
  IdentityFile ~/.ssh/id_ed25519_github
```

After making changes to the webhook deployment file, apply this yaml in your cluster.

You can also run terrascan admission controller server outside cluster, for more information and instructions on configuring terrascan as an admission controller webhook, see https://docs.accurics.com/projects/accurics-terrascan/en/latest/integrations/admission-controller-webhooks-usage.

#### Step 2: Create a Dockerfile

Create a Dockerfile for the container. This container will run the script that triggers the remote Terrascan API server. The template for the script is below, after the Dockerfile. Please fill the values in the template to match your environment.

```DockerFile
# Dockerfile with a script to use terrascan's validating webhook
# configured in the kubernetes cluster, to scan a repo for violations
FROM alpine:3.12.0

#curl to send request to terrascan validating webhook
RUN apk add --no-cache curl

WORKDIR /home/terrascan

RUN mkdir bin

COPY scripts/argocd-terrascan-remote-scan.sh  bin/terrascan-remote-scan.sh

# create non root user
RUN addgroup --gid 101 terrascan && \
  adduser -S --uid 101 --ingroup terrascan terrascan && \
  chown -R terrascan:terrascan bin && \
  chmod u+x bin/terrascan-remote-scan.sh

USER 101

CMD ["sh"]
```

##### The terrascan-remote-scan script

```sh
#!/bin/sh

set -o errexit


TERRASCAN_SERVER="https://${SERVICE_NAME}"
IAC=${IAC_TYPE:-"k8s"}
IAC_VERSION=${IAC_VERSION:-"v1"}
CLOUD_PROVIDER=${CLOUD_PROVIDER:-"all"}
REMOTE_TYPE=${REMOTE_TYPE:-"git"}

if [ -z ${SERVICE_NAME} ]; then
    echo "Service Name Not set"
    exit 6
fi

if [ -z ${REMOTE_URL} ]; then
    echo "Remote URL Not set"
    exit 6
fi

SCAN_URL="${TERRASCAN_SERVER}/v1/${IAC}/${IAC_VERSION}/${CLOUD_PROVIDER}/remote/dir/scan"

echo "Connecting to the service: ${SERVICE_NAME} to scan the remote url: ${REMOTE_URL} \
  with configurations { IAC type: ${IAC}, IAC version: ${IAC_VERSION},  remote type: ${REMOTE_TYPE} , cloud provider: ${CLOUD_PROVIDER}}"


RESPONSE=$(curl -s -w \\n%{http_code} --location -k  --request POST "$SCAN_URL" \
--header 'Content-Type: application/json' \
--data-raw '{
"remote_type":"'${REMOTE_TYPE}'",
"remote_url":"'${REMOTE_URL}'"
}')

echo "$RESPONSE"

HTTP_STATUS=$(printf '%s\n' "$RESPONSE" | tail -n1)

if [ "$HTTP_STATUS" -eq 403 ]; then
    exit 3
elif [ "$HTTP_STATUS" -eq 200 ]; then
    exit 0
else
    exit 1
fi
```


#### Step 3: Configure PreSync hook to use container created in step 2

The following example hook yaml is mostly ready to be added to an existing kubernetes configuration.

[comment]: <> (it says 'mostly', what is the pending thing for user to add/configure?)

```yaml
apiVersion: batch/v1
kind: Job
metadata:
generateName: terrascan-hook-
namespace: <YOUR APP NAMESPACE>
annotations:
  argocd.argoproj.io/hook: PreSync            
spec:
ttlSecondsAfterFinished: 3600
template:
  spec:
    securityContext:
      seccompProfile:
        type: RuntimeDefault
    containers:
    - name: terrascan-argocd
      image: <IMAGE FROM STEP TWO>
      resources:
        requests:
          cpu: "1"
          memory: "256Mi"
        limits:
          cpu: "1"
          memory: "256Mi"
      env:
        - name: SERVICE_NAME
          value: <Name of service exposed for terrascan controller pod>
        - name: REMOTE_URL
          value: <YOUR PRIVATE REPOSITORY PATH>
        - name: IAC_TYPE
          value: <IAC TYPE YOU WANT SCAN> # If not provided default value is 'k8s'
        - name: IAC_VERSION
          value: <VERSION OF IAC TYPE SELECTED> # If not provided default value is 'v1'
        - name: CLOUD_PROVIDER
          value: <TYPE OF CLOUD PROVIDER> #If not provided default value is 'all'
        - name: REMOTE_TYPE
          value: <TYPE OF REMOTE> #If not provided default value is 'git'         
      args:
      - sh
      - /home/terrascan/bin/terrascan-remote-scan.sh
      securityContext:
        seccompProfile:
          type: RuntimeDefault
        allowPrivilegeEscalation: false
        readOnlyRootFilesystem: true
        runAsNonRoot: true
        runAsUser: 101
      livenessProbe:
        exec:
          command:
          - cat
          - /home/terrascan/bin/terrascan-remote-scan.sh
        periodSeconds: 10
        initialDelaySeconds: 10
      readinessProbe:
        exec:
          command:
          - cat
          - /home/terrascan/bin/terrascan-remote-scan.sh
        periodSeconds: 10
        initialDelaySeconds: 10
    restartPolicy: Never
backoffLimit: 1
```

To allow users to check for violations in the web interface, configure the job to delete after the specified time, using the parameter `ttlSecondsAfterFinished`. In addition, violation can be reported as webhook notifications, as shown in Method 1.

After configuring the presynchook yaml file, add the file to the relevant repository folder to configure Argo CD.

> **Note**: All the example yaml configuration files present in documentation are tested with k8s 1.19.7 version.
