# Helm chart for deploying terrascan

This chart deploys terrascan as a server within your kubernetes cluster. By default it runs just terrascan by itself, but,
user creates namespace and secrets.

In server mode, terrascan will act both as an API server for
performing remote scans of IAC, as well as a validating admission
webhook for a Kubernetes cluster. Further details can be found in
the [main documentation](https://docs.accurics.com/projects/accurics-terrascan/en/latest/).
There are two helm charts:

1. In the `server/` directory : to deploy terrascan in server mode.
2. In the `webhook/` directory : to setup a validating webhook that uses the deployed terrascan server from step 1, as its backend.

## Usage
### Set up TLS certificates
A requirement to run an admission controller is that communication
happens over TLS. This helm chart expects to find the certificate
at `server/data/server.crt` and key at `server/data/server.key`.
If you opt to deploy the webhook as well, please copy `server/data/server.crt` at `webhook/data/server.crt`

### Set up SSH config for remote repo scan
If you're opting to utilise the remote repo scan feature,
terrascan will require ssh capabilities to do that.
This helm chart expects to find the your ssh private key at `.ssh/private_key`,
ssh config file at `.ssh/config` and .ssh known_hosts file at `.ssh/known_hosts`.
Your ssh public key must setup at the code repository hosting service, such as github, bitbucket, etc.

**Note:** This is an optional feature and not a requirement.

### Persistent storage
By default, this chart will deploy terrascan with a `emptyDir`
volume - basically a temporary volume. If you intend to use the
admission controller functionality, then you may want to store the
admission controller database on a persistent volume. This chart
supports specifying a [persistent volume claim](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) for
the database - as storage, PVs, and PVCs are a wide topic within
Kubernetes ecosystem, the details of the PV/PVC creation are left
to the individual.

To specify the use of a PVC, set `persistence.enable` to `true`, and then specify the name of an existing PVC:

```
persistence:
  enabled: false
  existingclaim: pvcClaimName
```

### Terrascan configuration file
This chart will look for a [terrascan configuration
file](https://docs.accurics.com/projects/accurics-terrascan/en/latest/usage/#config-file)
at `server/data/config.toml`. If that file exists before running `helm
install`, it's contents will be loaded into a configMap and provided
to the terrascan server.

### Deploy
Once your TLS certificate is generated and the values in the
`values.yaml` configuration file have been reviewed, you can install
the chart with the following command:

1. Deploying Terrascan Server.

    *Ensure that your current working directory is `server/`.*
    ```
    helm install <releasename-for-server> .
    ```
    Where `<releasename-for-server>` is the name you want to assign to this installed chart.
    This value will be used in various resources to make them both distinct and identifiable.

    This will use your current namespace unless `-n <namespace>` is specified

    #### Verification

    You can query for the pod using the following command.
    ```
    kubectl get pod -n <namespace> -w
    ```
    Watch the pod until it attains the `Running` state.

    Verify the logs of the terrascan pod using the following command.
    ```
    kubectl -n <namespace> logs <pod-name>
    ```
   If you see a log that goes like `server listening on port : <port-name>`, the deployment went smooth.

2. Deploying Validating Webhook.

    *Ensure that your current working directory is `webhook/`.*
    ```
    helm install <releasename-for-webhook> .
    ```
   This will use your current namespace unless `-n <namespace>` is specified.
   ***Ensure that you provide the exact same <namespace> value as you did to deploy the `server/` chart in step 1.***



## TODO:
This chart is a WIP - we intend to add the following functionality in the near future:
 - [x] Storage support - volume for db
 - [x] Add section for setting the validating-webhook up.
 - [x] Add secrets to add ssh capabilities in the container, to enable remote repo scan feature.
 - [ ] Support more load balancer types
 - [ ] Support for ingress
 - [ ] Flag for UI enable/disable
 - [ ] Publish to Artifact hub
 - [ ] Support TLS certificate/key in existing secrets
