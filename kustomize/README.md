## INSTALLING TERRASCAN IN A KUBERNETES CLUSTER USING KUSTOMIZE

This guide will help you install terrascan server inside your kubernetes cluster, as per your specific use case.
We have covered four use cases in the sections below:

  - Deploying Terrascan Server

      terrascan operating in http server mode.

  - Deploying Terrascan Server in TLS Mode

      terrascan operating in https server mode.

  - Deploying Terrascan Server for Remote Repository Scan

  - Setting Up Terrascan Webhook

### PRE-REQUISITE
1. Make sure you have required access on the kubernetes cluster to create and update the following resources:

  - Secrets
  - Configmaps
  - Deployments
  - Services
  - ValidatingWebhookConfiguration (only if you're aiming to deploy the webhook as well)

  **_If it is not a production level cluster, you probably do have the required access._**

2. Make sure you have `kubectl`, `kustomize` and `openssh` installed on your local machine.

3. Make sure you replace `<TERRASCAN_NAMESPACE>` placeholder with your target namespace where you to want to deploy the
terrascan server. The string replacement will be required in the following files:

  - `base/kustomization.yaml`
  - `server/kustomization.yaml`
  - `server-tls/kustomization.yaml`
  - `certs/domain.cnf` (that is generated in step 1 of `Deploying Terrascan Server in TLS Mode` section)
  - `webhook/kustomization.yaml` (only if you're aiming to deploy the webhook as well)
  - `webhook/validating-webhook.yaml` (only if you're aiming to deploy the webhook as well)

### Deploying Terrascan Server

Deploy terrascan in server mode operating in plain HTTP mode.

1. Place your terrascan `config.toml` in the `base/config/` directory or edit the existing one.

2. Deploy the terrascan server. Skip this step if you're aiming to setup terrascan in tls mode, terrascan webhook or
   terrascan server for remote repository scan.

   Note: Before running the command, please verify once that the `server/kustomization.yaml` is set with the desired parameters.

    ```bash
    kubectl apply -k server/
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
    CN = terrascan.<TERRASCAN_NAMESPACE>.svc
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
    kubectl apply -k server-tls/
    ```

### Deploying Terrascan Server For Remote Repository Scan

For scanning remote IaC file repositories, Terrascan must be provided with the required SSH keys to connect and clone the
repository locally to scan it. The following steps will help in setting up for that.

1. Follow steps 1-3 of the `Deploying TerraScan Server in TLS mode` section.

2. Generate SSH keys and copy `~/.ssh/config`, `~/.ssh/known-hosts` and `~/.ssh/<generated_private_key>` to
   `server-remote-repo-scan/.ssh/` directory. Replace `<SSH_KEY_NAME>` with your private ssh key's name in
   `server-remote-repo-scan/kustomization.yaml` and setup the generated public ssh key on your respective code repository
   hosting service, like github or bitbucket.

   **You may also use this shell command:**

   _Let's assume your private key is `~/.ssh/github_rsa`_

    ```bash
    mkdir server-remote-repo-scan/.ssh
    cp ~/.ssh/config ~/.ssh/known_hosts server-remote-repo-scan/.ssh/
    cp ~/.ssh/github_rsa ~/.ssh/github_rsa.pub server-remote-repo-scan/.ssh/
    sed s/<SSH_KEY_NAME>/github_rsa/g server-remote-repo-scan/kustomization.yaml
    ```

3. Deploy. Skip this step if you're aiming to setup terrascan webhook.

   **Note:** Before running the command, please verify once that the `server-remote-repo-scan/kustomization.yaml` is set
   with the desired parameters.

   ```bash
    kubectl apply -k server-remote-repo-scan/
    ```

### Setting Up Terrascan Webhook
If you want to setup a Validating Webhook that scans your incoming kubernetes resources using terrascan,
follow the steps below.

1. If you aim to use the deployed terrascan server solely by the validating webhook, follow steps 1 to 3 from the
   `Deploying Terrascan Server in TLS mode` section above.

   **OR**

   If you aim to use the deployed terrascan server both by the validating webhook and argocd pre-sync hooks, follow steps 1 to 2 from the
   `Deploying Terrascan Server For Remote Repository Scan` section above.

2. In `webhook/validating-webhook.yaml` and `webhook/deployment-env.yaml` file, Replace `<WEBHOOK_API_KEY>`with the string that
   you want your terrascan server key to be.

   **You may also use this shell command:**

   *Let's assume we want the string `t3rrascan` as the server key.*

    ```bash
    sed s/<WEBHOOK_API_KEY>/t3rrascan/g webhook/validating-webhook.yaml
    sed s/<WEBHOOK_API_KEY>/t3rrascan/g webhook/deployment-env.yaml
    ```

3. In `webhook/validating-webhook.yaml`, replace `<CA_BUNDLE>` with the base64 encoded value of the
   `server/certs/server.crt` that was setup in Step 2 of `Deploying Terrascan Server in TLS Mode` section.

   *You may also use this shell command:*

    ```bash
    $CA_BUNDLE=(cat server-tls/certs/server.crt | base64)
    sed s/<CA_BUNDLE>/$CA_BUNDLE/g webhook/validating-webhook.yaml
    ```

4. In the `webhook/validating-webhook.yaml` file, set the `webhooks.rules` section as per your requirement. By default,
   we have setup a rule to block possibly all the resources from being created or updated. This might not be correct for
   your use case, refer the kubernetes admission webhook docs for the same.

   *The following command might help as well.*
    ```bash
    kubectl explain ValidatingWebhookConfiguration.webhooks.rules
    ```

5. Deploy.
   Note: Before running the command, please verify once that the `server/kustomization.yaml` & `webhook/kustomization.yaml`
         are set with the desired parameters.

    ```bash
    kubectl apply -k webhook/
    ```
