# Using Terrascan as a Kubernetes Admission Controller

## Overview
Terrascan can be integrated with K8s [admissions webhooks](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/).
It can be used as one of the validating webhooks to be used and scan new configurations.

In this guide, we'll demonstrate how Terrascan can be configured to:
* Scan configuration changes policies when an object is being created or updated
* Allow / reject the request in case a violation is detected


## Installation Guide

### Create an instance
Your Terrascan instance has the following requirements for being able to scan K8s configurations.

1. Be accessible via HTTPS. Make sure your cloud firewall is configured to allow this.
1. Have a valid SSL certificate for the served domain name. To do that, choose one of our suggested methods:
  1. Use a subdomain of your choosing (e.g dev-terrascan-k8s.accurics.com) and create a valid certificate for this subdomain through your SSL certificate provider. [Let's Encrypt](https://letsencrypt.org/) is a free, simple to use certificate authority you can use.
  1. Use a reverse-proxy to serve SSL requests; for example, use Cloudflare Flexible to get a certificate by a trusted-CA to your [self-signed certificate](https://www.digitalocean.com/community/tutorials/openssl-essentials-working-with-ssl-certificates-private-keys-and-csrs).
  1. Generate a self-signed certificate and have your K8s cluster trust it. To add a trusted CA to ca-pemstore, as demonstrated in [paraspatidar's blog post](https://medium.com/@paraspatidar/add-ssl-tls-certificate-or-pem-file-to-kubernetes-pod-s-trusted-root-ca-store-7bed5cd683d).
1. Use the Terrascan docker as demonstrated in this document, or run it from the sources.

### Run Terrascan webhook service
Run Terrascan docker image in your server using the following command:
 ```bash
  sudo docker run -p 443:9443 -v <DATA_PATH>:/data -u root -e K8S_WEBHOOK_API_KEY=<API_KEY>> accurics/terrascan server --cert-path /data/cert.pem --key-path /data/key.pem
 ```
`<API_KEY>` is a key used for authentication between your K8s environment and  the Terrascan server. Generate your preferred key and use it here.

`<DATA_PATH>` is a directory path in your server where both the certificate and the private key .pem files are stored.
In addition, this directory is used to save the webhook logs. (An SQLite file)

You can specify a config file that specifies which policies to use in the scan and which violations should lead to rejection.

A config file example: ```my_terrscan_config.toml```
  ```bash
[severity]
level = "medium"
[rules]
    skip-rules = [
        "accurics.kubernetes.IAM.107"
    ]

[k8s-deny-rules]
  denied-categories = [
      "Network Ports Security"
  ]
  denied-severity = "high"
  ```

You can specify the following configurations:
*  **scan-rules** - one or more rules to scan
*  **skip-rules** - one or more rules to skip while scanning
*  **severity** - the minimal level of severity of the policies to be scanned


* **k8s-deny-rules** - specify the rules that should cause a rejection of the admission request
  *  **denied-categories** - one or more policy categories that are not allowed in the detected violations
  *  **denied-severity** - the minimal level of severity that should cause a rejection

In order to use a configuration file, add it as a command line argument:

```<terrascan> -c /data/my_terrscan_config.toml```


### Configure K8s to send webhooks
Configure a new ```ValidatingWebhookConfiguration``` in your Kubernetes environment and specify your Terrascan server endpoint.

Example:
   ```bash
    cat <<EOF | kubectl apply -f -
    apiVersion: admissionregistration.k8s.io/v1
    kind: ValidatingWebhookConfiguration
    metadata:
      name: my.validation.example.check
    webhooks:
      - name: my.validation.example.check
        rules:
          - apiGroups:
              - ""
            apiVersions:
              - v1
            operations:
              - CREATE
              - UPDATE
            resources:
              - pods
              - services
        failurePolicy: Fail
        clientConfig:
          url: https://<SERVER_ADDRESS>/v1/k8s/webhooks/<API_KEY>/scan
        sideEffects: None
        admissionReviewVersions: ["v1"]
    EOF
  ```

* You can modify the `rules` that trigger the webhook according to your preferences.
* Update the ```clientConfig``` URL with your terrascan server address and the API key you generated before.


### Test your settings
Try to run a new pod / service. For example:
``` Bash
  kubectl run mynginx --image=nginx
```

Go to ```https://<SERVER_ADDRESS>/k8s/webhooks/<API_KEY>/logs``` and verify your request is logged.
