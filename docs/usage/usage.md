# Usage Overview
Terrascan is a static code analyzer for Infrastructure as Code (IaC). It can executed with the native binary/executable or by using the [`docker`](#using-docker) container.

## Installing Terrascan

For steps to install or run Terrascan from docker, see the section [Getting Started](getting-started.md).

## Building Terrascan
Terrascan is a Go binary that you can build locally. This is useful if you want to be on the latest version or when developing Terrascan.

``` Bash
$ git clone git@github.com:accurics/terrascan.git
$ cd terrascan
$ make build
$ ./bin/terrascan
```

## Using Terrascan

This section provides an overview of the different ways you can use Terrascan:

1. [Command line mode](command_line_mode.md) provides list of Terrascan commands with descriptions.
2. [Server mode](server_mode.md) using Terrascan as API server

See the [Configuring Terrascan](config_options) for steps to use a TOML file to configure webhook notifications.


## Integrations

Terrascan can be integrated into various platforms and configured to validate policies to provide run time security. Currently Terrascan supports the following integrations:

1. [Kubernetes (K8s) Admissions webhooks](integrations/admission-controller-webhooks-usage.md)
2. [ArgoCD](integrations/argocd-integration.md)
3. [Atlantis](integrations/atlantis-integration.md)
4. [Github and GitLab or CI/CD pipelines](integrations/cicd.md)