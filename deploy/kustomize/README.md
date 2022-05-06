## Installing terrascan in a Kubernetes cluster using Kustomize

This guide will help you install terrascan server inside your kubernetes cluster.
We have covered the following use cases in the sections below.

  - [Deploying Terrascan Server](#deploying-terrascan-server)  
    Terrascan operating in http server mode.

  - [Deploying Terrascan Server in TLS Mode](#deploying-terrascan-server-in-tls-mode)  
    Terrascan operating in https server mode. This deployment is also a foundation for the terrascan webhook setup.

  - [Deploying Terrascan Server for Remote Repository Scan](#deploying-terrascan-server-for-private-remote-repository-scan)  
    Terrascan in https server mode installed with ssh capabilities, to scan ***private*** remote repositories. For remote
    scanning public repos, deploying `Terrascan Server in TLS Mode` is sufficient.
    This deployment can be handy for use-cases like an argocd pre-sync hook that sends remote repository scan requests to the server.

  - [Setting Up Terrascan Webhook](#setting-up-terrascan-webhook)  
    A Kubernetes Validating Webhook, that safeguards your cluster by denying the creation of kubernetes resources that
    can cause potential security violations.

  - [Clean Up](#clean-up)

### Pre-requisite
1. Make sure you have required access on the kubernetes cluster to create and update the following resources:

  - Secrets
  - Configmaps
  - Deployments
  - Services
  - ValidatingWebhookConfiguration (only if you're aiming to deploy the webhook as well)

  **If it is not a production level cluster, you probably do have the required access.**

2. Make sure you have `kubectl`, `kustomize` and `openssh` installed on your local machine.

3. Make sure you replace `<TERRASCAN_NAMESPACE>` placeholder with your target namespace where you to want to deploy the
terrascan server. The string replacement will be required in the following files:

  - `base/kustomization.yaml`
  - `server/kustomization.yaml`
  - `server-tls/kustomization.yaml`
  - `server-remote-repo-scan/kustomization.yaml`
  - `server-tls/certs/domain.cnf` (that is generated in step 1 of `Deploying Terrascan Server in TLS Mode` section)
  - `webhook/kustomization.yaml` (only if you're aiming to deploy the webhook as well)
  - `webhook/validating-webhook.yaml` (only if you're aiming to deploy the webhook as well)

  *Make sure your pwd is same as this README.md file*

  Let's assume that the desired namespace is 'terrascan'.
  ```bash
  sed -i "" "s/<TERRASCAN_NAMESPACE>/terrascan/g" base/kustomization.yaml
  sed -i "" "s/<TERRASCAN_NAMESPACE>/terrascan/g" server/kustomization.yaml
  sed -i "" "s/<TERRASCAN_NAMESPACE>/terrascan/g" server-tls/kustomization.yaml
  sed -i "" "s/<TERRASCAN_NAMESPACE>/terrascan/g" server-remote-repo-scan/kustomization.yaml
  sed -i "" "s/<TERRASCAN_NAMESPACE>/terrascan/g" server-tls/certs/domain.cnf
  sed -i "" "s/<TERRASCAN_NAMESPACE>/terrascan/g" webhook/kustomization.yaml
  sed -i "" "s/<TERRASCAN_NAMESPACE>/terrascan/g" webhook/validating-webhook.yaml
  ```
4. Ensure that your desired namespace exist.

  Let's assume that the desired namespace is 'terrascan'.
  ```bash
  kubectl create namespace terrascan
  ```

### Deploying Terrascan Server

Deploy terrascan in server mode operating in plain HTTP mode.

1. Place your terrascan `config.toml` in the `base/config/` directory or edit the existing one.

2. Deploy the terrascan server. Skip this step if you're aiming to setup terrascan in tls mode, terrascan webhook or
   terrascan server for remote repository scan.

   Note: Before running the command, please verify once that the `server/kustomization.yaml` is set with the desired parameters.

    ```bash
    kustomize build server/ | kubectl apply -f -
    ```

### Deploying Terrascan Server in TLS Mode

Deploy terrascan in server mode operating in HTTPS mode.

1. Follow Step 1 from `Deploying Terrascan Server` section

2. Create a domain.cnf file.

    ```bash
    mkdir server-tls/certs
    touch server-tls/certs/domain.cnf
    cat << EOF > certs/domain.cnf
    [req]
    default_bits = 2048
    prompt = no
    default_md = sha256
    x509_extensions = v3_req
    distinguished_name = dn
    [dn]
    C = <My_Country>
    ST = <My_State>
    L = <My_Location>
    O = <My_Organization>
    emailAddress = <My_Email>
    CN = terrascan.<TERRASCAN_NAMESPACE>.svc.cluster.local
    [v3_req]
    subjectAltName = @alt_names
    [alt_names]
    DNS.1 = terrascan.<TERRASCAN_NAMESPACE>.svc.cluster.local
    >EOF
    ```

    **Note:** Please replace the placeholders like `<My_Country>`,`<My_State>` etc as per your requirements.

3. Generate `server.key` and `server.crt`.

    ```bash
    openssl req -x509 -sha256 -nodes -newkey rsa:2048 -keyout server-tls/certs/server.key -out server-tls/certs/server.crt -config server-tls/certs/domain.cnf
    ```

4. Deploy the terrascan server. Skip this step if you're aiming to setup terrascan webhook or terrascan server for remote repository scan.

   **Note:** Before running the command, please verify once that the `server-tls/kustomization.yaml` is set with the desired parameters.

    ```bash
    kustomize build server-tls/ | kubectl apply -f -
    ```

### Deploying Terrascan Server For Private Remote Repository Scan

For scanning ***Private*** remote IaC file repositories, Terrascan must be provided with the required SSH keys to connect and clone the
repository locally to scan it. The following steps will help in setting up for that.

1. Follow steps 1-3 of the `Deploying TerraScan Server in TLS mode` section.

2. Generate SSH keys and copy your generated private key to
   `server-remote-repo-scan/.ssh/` directory. Replace `<SSH_KEY_NAME>` with your private ssh key's name in
   `server-remote-repo-scan/kustomization.yaml` and setup the generated public ssh key on your respective code repository
   hosting service, like github, gitlab or bitbucket.


   *Let's assume your private key is `~/.ssh/github_rsa`*

  ```bash
    mkdir server-remote-repo-scan/.ssh
    cp ~/.ssh/github_rsa server-remote-repo-scan/.ssh/
    sed -i "" "s/<SSH_KEY_NAME>/github_rsa/g" server-remote-repo-scan/kustomization.yaml
  ```

   Apart from the ssh key, we also require a `known_hosts` file. You can create `server-remote-repo-scan/.ssh/known_hosts` using the below command.
   If you're dealing with some other code repository service than github, bitbucket or gitlab, please modify the known_host file accordingly.


  ```bash
    cat << EOF >> server-remote-repo-scan/.ssh/known_hosts
    # known_hosts
    github.com,192.30.255.113 ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==
    bitbucket.org,104.192.141.1 ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAubiN81eDcafrgMeLzaFPsw2kNvEcqTKl/VqLat/MaB33pZy0y3rJZtnqwR2qOOvbwKZYKiEO1O6VqNEBxKvJJelCq0dTXWT5pbO2gDXC6h6QDXCaHo6pOHGPUy+YBaGQRGuSusMEASYiWunYN0vCAI8QaXnWMXNMdFP3jHAJH0eDsoiGnLPBlBp4TNm6rYI74nMzgz3B9IikW4WVK+dc8KZJZWYjAuORU3jc1c/NPskD2ASinf8v3xnfXeukU0sJ5N6m5E8VLjObPEO+mN2t/FZTMZLiFqPWc/ALSqnMnnhwrNi2rbfg/rd/IpL8Le3pSBne8+seeFVBoGqzHM9yXw==
    gitlab.com,172.65.251.78 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBFSMqzJeV9rUzU4kWitGjeR4PWSa29SPqJ1fVkhtj3Hw9xjLVXVYrU9QlYWrOLXBpQ6KWjbjTDTdDkoohFzgbEY=
  ```

  After pasting the command, press RETURN key, type `EOF` on the prompt and press RETURN key again.

3. Deploy.

   **Note:** Before running the command, please verify once that the `server-remote-repo-scan/kustomization.yaml` is set
   with the desired parameters.

   ```bash
   kustomize build server-remote-repo-scan/ | kubectl apply -f -
    ```

### Setting Up Terrascan Webhook
If you want to setup a Validating Webhook that scans your incoming kubernetes resources using terrascan,
follow the steps below.

1. If you aim to use the deployed terrascan server solely by the validating webhook, follow steps 1 to 3 from the
   `Deploying Terrascan Server in TLS mode` section above.

   **OR**

   If you aim to use the deployed terrascan server both by the validating webhook and for remote repository scans, follow
   steps 1 to 2 from the `Deploying Terrascan Server For Remote Repository Scan` section above and
   replace `- ../server-tls` with `- ../server-remote-repo-scan` in the `webhook/kustomization.yaml` file.

2. In `webhook/validating-webhook.yaml` and `webhook/deployment-env.yaml` file, Replace `<WEBHOOK_API_KEY>`with the string that
   you want your terrascan server key to be.

   **You may also use this shell command:**

   *Let's assume we want the string `t3rrascan` as the server key.*

    ```bash
    sed -i "" "s/<WEBHOOK_API_KEY>/t3rrascan/g" webhook/validating-webhook.yaml
    sed -i "" "s/<WEBHOOK_API_KEY>/t3rrascan/g" webhook/deployment-env.yaml
    ```

3. In `webhook/validating-webhook.yaml`, replace `<CA_BUNDLE>` with the base64 encoded value of the
   `server/certs/server.crt` that was setup in Step 2 of `Deploying Terrascan Server in TLS Mode` section.

   *You may also use this shell command:*

    ```bash
    CA_BUNDLE=$(cat server-tls/certs/server.crt | base64)
    sed -i "" "s/<CA_BUNDLE>/$CA_BUNDLE/g" webhook/validating-webhook.yaml
    ```

4. In the `webhook/validating-webhook.yaml` file, set the `webhooks.rules` section as per your requirement. By default,
   we have setup a rule to block possibly all the resources from being created or updated. This might not be correct for
   your use case, refer the kubernetes admission webhook docs for the same.

   *The following command might help as well.*
    ```bash
    kubectl explain ValidatingWebhookConfiguration.webhooks.rules
    ```

5. Deploy.

   5.1 Deploy the webhook's backend: terrascan deployment and service.

   **Note:** Before running the command, please verify once that the `server/kustomization.yaml` & `webhook/kustomization.yaml`
   are set with the desired parameters.

   ```bash
   kustomize build webhook/ | kubectl apply -f -
   ```

    5.2 Verify that the terrascan server pod is up and ready to server.

    ```bash
    kubectl -n <TERRASCAN_NAMESPACE> get pods -w
    ```

    When the pod is in running state, verify the logs.

    ```bash
    kubectl -n terrascan logs <pod-name> -f
    ```

    When there is a log message that says "server listening at port 9010", proceed to the next step.

    5.3 Deploy the webhook.

    ```bash
    kubectl apply -f webhook/validating-webhook.yaml
    ```

### Clean Up

  Deleting the namespace that you used, will delete all the resources itself.
  ```bash
  kubectl delete ns <TERRASCAN_NAMESPACE>
  ```

  If for some reason you don't want to delete the namespace, for example, in case you deployed to the `default` namespace
  and deleting the default namespace will make you lose the default kubernetes service as well, which may cause further trouble.

  ```bash
  kustomize build <section-directory>/ | kubectl delete -f
  ```

