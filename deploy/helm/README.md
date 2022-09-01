# Helm charts for deploying terrascan

This guide deploys terrascan as a server within your kubernetes cluster. Additionally, you can deploy a
Validating Webhook as well, that'll use the terrascan server as its backend.

In server mode, terrascan will act both as an API server for
performing remote scans of IAC, as well as a validating admission
webhook for a Kubernetes cluster. Further details can be found in
the [main documentation](https://runterrascan.io/).

## Usage
### Set up TLS certificates
A requirement to run an admission controller is that communication
happens over TLS. This helm chart expects to find the certificate
at `data/server.crt` and key at `data/server.key`.

There's a `data/domain.cnf` file available for you to edit and generate key & certificate. You can use the following command:

  ```bash
  openssl req -x509 -sha256 -nodes -newkey rsa:2048 -keyout data/server.key -out data/server.crt -config data/domain.cnf
  ```

In the `data/domain.cnf` file, we have configured DNS names as `terrascan.terrascan.svc`, assuming the defaults that service
will be named `terrascan` and hosted in `terrascan` namespace. You'll have to manually change that as per your requirements.

### Terrascan configuration file
This chart will look for a [terrascan configuration
file](https://runterrascan.io/docs/usage/config_options/)
at `data/config.toml`. If that file exists before running `helm
install`, it's contents will be loaded into a configMap and provided
to the terrascan server.

### Set up SSH config for private remote repo scan
If you're opting to utilise the remote repo scan feature for ***private*** repositories,
terrascan will require ssh capabilities to do that.
This helm chart expects to find the your ssh private key at `.ssh/private_key`,and .ssh known_hosts file at `.ssh/known_hosts`.
Your ssh public key must setup at the code repository hosting service, such as github, bitbucket, etc.

You can use the below content to create your `.ssh/known_hosts` file.

```bash
# known_hosts
github.com,192.30.255.113 ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==
bitbucket.org,104.192.141.1 ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAubiN81eDcafrgMeLzaFPsw2kNvEcqTKl/VqLat/MaB33pZy0y3rJZtnqwR2qOOvbwKZYKiEO1O6VqNEBxKvJJelCq0dTXWT5pbO2gDXC6h6QDXCaHo6pOHGPUy+YBaGQRGuSusMEASYiWunYN0vCAI8QaXnWMXNMdFP3jHAJH0eDsoiGnLPBlBp4TNm6rYI74nMzgz3B9IikW4WVK+dc8KZJZWYjAuORU3jc1c/NPskD2ASinf8v3xnfXeukU0sJ5N6m5E8VLjObPEO+mN2t/FZTMZLiFqPWc/ALSqnMnnhwrNi2rbfg/rd/IpL8Le3pSBne8+seeFVBoGqzHM9yXw==
gitlab.com,172.65.251.78 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBFSMqzJeV9rUzU4kWitGjeR4PWSa29SPqJ1fVkhtj3Hw9xjLVXVYrU9QlYWrOLXBpQ6KWjbjTDTdDkoohFzgbEY=
```
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

### Deploy
Once your TLS certificate is generated and the values in the
`values.yaml` configuration file have been reviewed, you can install
the chart with the following command:

1. Deploying Terrascan Server.

   For just installing a terrascan server deployment and service,

    ```
    helm install <release-name> . -n <namespace>
    ```
   Where `<release-name>` is the name you want to assign to this installed chart.
   This value will be used in various resources to make them both distinct and identifiable.

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

   ###Deploying Validating Webhook.

   For installing the terrascan deployment and service along the validating webhook,

    ```
    helm install <release-name> . -n <namespace> --set webhook.mode=true
    ```
   This will use your current namespace unless `-n <namespace>` is specified.

   #### Verification
   Try creating a resource that's scanning by the webhook

   ```bash
    kubectl run test-pod --image=nginx
   ```

   #### Clean Up

    ```bash
    helm uninstall <release-name> -n <namespace>
    ```

## TODO:
This chart is a WIP - we intend to add the following functionality in the near future:
- [x] Storage support - volume for db
- [x] Add section for setting the validating-webhook up.
- [x] Add secrets to add ssh capabilities in the container, to enable remote repo scan feature.
- [x] Support TLS certificate/key in existing secrets
- [ ] Support more load balancer types
- [ ] Support for ingress
- [ ] Flag for UI enable/disable
- [ ] Publish to Artifact hub
