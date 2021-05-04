# Integration of Terrascan with Argo CD

Terrascan can be configured as an Argo CD job during the application sync process using argocd’s resource hook. The PreSync resource hook is the best way to evaluate the kubernetes deployment configuration and report any violations.

## Terrascan can be integrated with Argo CD in two ways
___
1. Use terrascan as a pre-sync hook to scan remote repositories
2. Use terrascan’s k8s admission controller along with a pre-sync that scans a configured repository with the admission controller webhook


### 1. Configure terrascan as a PreSync hook and scan the remote repository. 
___

#### Configure a PreSync hook

The following example hook yaml is mostly ready to be added to an existing kubernetes configuration. Just make sure that the secrets,  known_hosts  and ssh_config volume are relevant and specify a terrascan image. You can also map a slack notification script to the container which will send notifications to your Slack webhook endpoint after the embedded script scans the repo.
 
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
     # if you want to add slack notification script add one more volume here
     Volumes:
       - name: notification-scripts
         hostPath:
           path: <PATH TO YOUR SCRIPT>
       #add all required ssh keys need to clone your repos
       - name: ssh-keys-github
         secret:
           secretName: github-secret
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
         cp /etc/github-secret/ssh-private-key /home/terrascan/.ssh/id_ed25519_github &&
         cp /etc/ssh-config-secret/ssh-config /home/terrascan/.ssh/config &&
         cp /etc/ssh-known-hosts-secret/ssh-known-hosts /home/terrascan/.ssh/known_hosts &&
         chmod -R 400 /home/terrascan/.ssh/* &&
         /go/bin/terrascan scan -r git -u <YOUR REPOSITORY PATH>-i k8s -t k8s | /data/notify_slack.sh webhook-tests argo-cd https://hooks.slack.com/services/TXXXXXXXX/XXXXXXXXXXX/0XXXXXXXXXXXXXXXXXX
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
         - mountPath: /etc/github-secret
           name: ssh-keys-github
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
 
As shown, the PreSync requires access to the repository where IaC is stored, using the same branch (default) as the Argo CD application pipeline.
Configuring the job to delete only after the specified time see ttlSecondsAfterFinished will allow users to check for violations in the User Interface, the alternative is through notifications.

Example slack notification script

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
 
For non-public repositories, the private key, known_hosts and ssh config needs to be added as a kubernetes secret, configmap and secret respectively.
 
```sh 
 kubectl create secret generic github-secret \
   --from-file=ssh-privatekey= < path to your private key > \
    --from-file=ssh-publickey=< path to your public key >
```

Config-map: 

``` 
  kubectl  create configmap ssh-known-hosts --from-file=
   < path to you known hosts file >
```   
 
ssh config secret

```
 kubectl create secret generic ssh-config-secret \
   --from-file=< path to your ssh config file >
```   
 
Example ssh config file

``` 
 Host github.com
  HostName github.com
  IdentityFile ~/.ssh/id_ed25519_github
```

Once the presynchook yaml file is completely configured, add this file to your repository folder for which Argo CD pipeline is configured.


### 2. Use Terrascan Admission Webhook from PreSyncHook
___
You can use the already deployed terrascan’s k8s admission controller webhook to scan the remote repository from Argo CD PreSync hook.
To configure, follow below steps


#### Step 1: Configure terrascan admission controller webhook deployment yaml file with required keys and volumes.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
name: terra-controller-webhook
labels:
  app: terrascan-admission-webhook
spec:
replicas: 1
selector:
  matchLabels:
    app: terrascan-admission-webhook
template:
  metadata:
    labels:
      app: terrascan-admission-webhook
  spec:
    containers:
    - name: terrascan-admission-webhook
      image: <TERRASCAN IMAGE>
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
        - mountPath: /data
          name: certspath
        - mountPath: /etc/github-secret
          name: ssh-keys-github
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
          cp /etc/github-secret/ssh-private-key /home/terrascan/.ssh/id_ed25519_github &&
          cp /etc/ssh-config-secret/ssh-config /home/terrascan/.ssh/config &&
          cp /etc/ssh-known-hosts-secret/ssh-known-hosts /home/terrascan/.ssh/known_hosts &&
          chmod -R 400 /home/terrascan/.ssh/* &&
          terrascan server --cert-path /data/server.crt --key-path /data/server.key -p 443 -l debug -c /data/config.toml
    volumes:
      - name: certspath
        hostPath:
          path: <YOUR CERTIFICATE FOLDER PATH>
      #add all required ssh keys need to clone your repos
      - name: ssh-keys-github
        secret:
          secretName: github-secret
      #add a secret for git config file   
      - name: ssh-config
        secret:
          secretName: ssh-config-secret
      #add a configmap for the ssh known_hosts file
      - name: ssh-known-hosts
        configMap:
          name: known-hosts-config
```            

For non-public repositories, the private key, known hosts and ssh config needs to be added as a kubernetes secret, configmap and secret respectively.
  
```
kubectl create secret generic github-secret \
  --from-file=ssh-privatekey= < path to your private key > \
  --from-file=ssh-publickey=< path to your public key >
``` 

Config-map: 

``` 
kubectl create configmap ssh-known-hosts --from-file=< path to you known hosts file >
``` 

ssh config secret

``` 
kubectl create secret generic ssh-config-secret \
  --from-file=< path to your ssh config file >
``` 

Example ssh config file

``` 
Host github.com
  HostName github.com
  IdentityFile ~/.ssh/id_ed25519_github
``` 

After making changes to the webhook deployment file, apply this yaml in your cluster.

You can also run terrascan admission controller server outside cluster, for more information on configuring terrascan as an admission controller webhook, follow https://docs.accurics.com/projects/accurics-terrascan/en/latest/integrations/admission-controller-webhooks-usage 


#### Step 2: Create a Dockerfile for the container which has the terrascan script to run the remote scan against the terrascan’s admission controller webhook.

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
  adduser -S --uid 101 --ingroup terrascan terrascan

RUN chown -R terrascan:terrascan bin && \
  chmod u+x bin/terrascan-remote-scan.sh

USER terrascan

CMD ["sh"]
```

terrascan-remote-scan script

```sh
#!/bin/sh
TERRASCAN_SERVER="https://${SERVICE_NAME}"
IAC="k8s"
IAC_VERSION="v1"
CLOUD_PROVIDER="all"
REMOTE_TYPE="git"

SCAN_URL="${TERRASCAN_SERVER}/v1/${IAC}/${IAC_VERSION}/${CLOUD_PROVIDER}/remote/dir/scan"

RESPONSE=$(curl -s -w \\n%{http_code} --location -k  --request POST "$SCAN_URL" \
--header 'Content-Type: application/json' \
--data-raw '{
"remote_type":"'${REMOTE_TYPE}'",
"remote_url":"'${REMOTE_URL}'"
}')

echo "$RESPONSE"

HTTP_STATUS=$(printf '%s\n' "$RESPONSE" | tail -n1)

if [ "$HTTP_STATUS" -ne 200 ]; then
  exit 3
fi
```


#### Step 3: Configure PreSync hook to use container created in step 2

The following example hook yaml is mostly ready to be added to an existing kubernetes configuration.

```yaml
apiVersion: batch/v1
kind: Job
metadata:
generateName: terrascan-hook-
namespace: test
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
      image: <IMAGE FROM STEP 2>
      resources:
        requests:
          cpu: "1"
          memory: "256Mi"
        limits:
          cpu: "1"
          memory: "256Mi"
      env:
        - name: SERVICE_NAME
          value: terra-controller-service.default.svc
        - name: REMOTE_URL
          value: <YOUR PRIVATE REPOSITORY PATH>
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

Configuring the job to delete only after the specified time see ttlSecondsAfterFinished will allow users to check for violations in the User Interface, the alternative is through notifications.
Once the presynchook yaml file is completely configured add this file to your Repository folder which you want to configure for Argo CD.

` All the example yaml configuration files present in documentation are tested with k8s 1.19.7 version. 
` 
