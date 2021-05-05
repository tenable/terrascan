##### PRE-REQUISITE

1. Make sure you replace `<TERRASCAN_NAMESPACE>` placeholder with your target namespace where you to want to deploy the
terrascan server. The string replacement will be required in the `server/kustomization.yaml`,  `webhook/validating-webhook.yaml`
and the `certs/domain.cnf` file that is generated in step 1 of `Generating TLS Certificates` section.


2. Generate SSH keys and copy `~/.ssh/config`, `~/.ssh/known-hosts` and `~/.ssh/<generated_private_key>` to `server/.ssh/` directory.
   replace `<SSH_KEY_NAME>` with your private ssh key's name in `server/kustomization.yaml` and setup the generated public ssh key on
   your respective code repository hosting service, like github or bitbucket.

   You may also use this shell command :

   Let's assume your private key is `~/.ssh/github_rsa`

  ```bash
  mkdir server/.ssh
  cp ~/.ssh/config ~/.ssh/known_hosts <pwd>/server/.ssh/
  cp ~/.ssh/github_rsa ~/.ssh/github_rsa.pub <pwd>/server/.ssh/
  sed s/<SSH_KEY_NAME>/github_rsa/g server/kustomization.yaml
  ```

#### Generating TLS Certificates

1. Create a domain.cnf file.

  ```bash
  mkdir server/certs
  touch server/certs/domain.cnf
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
  emailAddress = me@email.com
  CN = terrascan
  [v3_req]
  subjectAltName = @alt_names
  [alt_names]
  DNS.1 = terrascan-validating-hook.<TERRASCAN_NAMESPACE>.svc
  >EOF
  ```

  Note: Please replace the placeholders like `<My_Country>`, etc as per your requirements.

2. Generate `server.key` and `server.crt`.

  ```bash
  openssl req -x509 -sha256 -nodes -newkey rsa:2048 -keyout server/certs/server.key -out server/certs/server.crt -config server/certs/domain.cnf
  ```

### Deploying Terrascan Server

1. Generate your TLS certificate files and place them in the `server/certs/` directory as `server.key` and `server.crt`.
   Please refer to the `Generating TLS Certificates` section above.

2. Place your terrascan `config.toml` in the `server/config/` directory or edit the existing one.

3. In `server/deployment.yaml` Replace `<TERRASCAN_SERVER_KEY_PLACEHOLDER>` with the string that
   you want your terrascan server key to be.

   You may also use this shell command :

   Let's assume we want the string `t3rrascan` as the server key.

  ```bash
  sed s/<TERRASCAN_SERVER_KEY_PLACEHOLDER>/t3rrascan/g server/deployment.yaml
  ```

4. Deploy the terrascan server. Skip this step if you're aiming to setup terrascan webhook.

   Note: Before running the command, please verify once that the `server/kustomization.yaml` is set with the desired parameters.

  ```bash
  kubectl apply -k server/
  ```

### Setting Up Terrascan Webhook

1. Please follow steps 1 to 3 from the `Deploying Terrascan Server` section above.

2. In `webhook/validating-webhook.yaml` Replace `<TERRASCAN_SERVER_KEY_PLACEHOLDER>` with the string that
   you want your terrascan server key to be.

   You may also use this shell command :

   Let's assume we want the string `t3rrascan` as the server key.

  ```bash
  sed s/<TERRASCAN_SERVER_KEY_PLACEHOLDER>/t3rrascan/g webhook/validating-webhook.yaml
  ```

  Note: Please ensure that the terrascan key value you use here is same as used
  in Step 3 of `Deploying Terrascan Server` section.

3. In `webhook/validating-webhook.yaml`, replace `<CA_BUNDLE_PLACEHOLDER>` with the base64 encoded value of the
   `server/certs/server.crt` that was setup in Step 1 of `Deploying Terrascan Server` section.

   You may also use this shell command :

  ```bash
  $CA_BUNDLE=(cat server/certs/server.crt | base64)
  sed s/<CA_BUNDLE_PLACEHOLDER>/$CA_BUNDLE/g webhook/validating-webhook.yaml
  ```

4. Deploy.

   Note: Before running the command, please verify once that the `server/kustomization.yaml` & `webhook/kustomization.yaml`
         are set with the desired parameters.

  ```bash
  kubectl apply -k webhook/
  ```
