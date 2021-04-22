# Helm chart for deploying terrascan in server mode

----
this is WIP!! 
----

This chart deploys terrascan as a server within your kubernetes cluster. By default it runs just terrascan by itself, but 
user creates namespace and secrets

## Usage

* set up tls certs
* run...
```
helm install releasename .
```
This will use your current namespace unless `-n <namespace>` is specified

## TODO:
(in rough order of priority)
 - [ ] Storage support - volume for db
 - [ ] Make toml config/configmap a little smoother
 - [ ] Load balancer types
 - [ ] Support for ingress?
 - [ ] moar docs
 - [ ] Flag for UI enable/disable
 - [ ] Publish to Artifact hub?

