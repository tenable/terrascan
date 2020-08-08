# Terrascan documentation
- About Terrascan
- Getting Started
    - Installation
    - Scanning
    - Terrascan CLI

- Architecture
    - Runtime
    - Infrastructure as Code Language Providers
    - Cloud Providers
    - Policy Engine
- Provider Reference
    - IAC Language Providers
        - Terraform (HCL2)
        - CloudFormation (JSON)
        - CloudFormation (YAML)
        - Kubernetes (YAML)
    - Cloud Providers
        - AWS
        - Azure
        - GCP
- Policies
    - AWS
    - Azure
    - GCP
    - Kubernetes
- Learning
    - pre-commit
    - super-linter


    Introduction: This section covers a general overview of what Envoy is, an architecture overview, how it is typically deployed, etc.

    Getting Started: Quickly get started with Envoy using Docker.

    Installation: How to build/install Envoy using Docker.

    Configuration: Detailed configuration instructions for Envoy. Where relevant, the configuration guide also contains information on statistics, runtime configuration, and APIs.

    Operations: General information on how to operate Envoy including the command line interface, hot restart wrapper, administration interface, a general statistics overview, etc.

    Extending Envoy: Information on how to write custom filters for Envoy.

    API reference: Envoy API detailed reference.

    Envoy FAQ: Have questions? We have answers. Hopefully.

## Using as pre-commit

Terrascan can be used on pre-commit hooks to prevent accidental introduction of security weaknesses into your repository.
This requires having pre-commit_ installed. An example configuration is provided in the comments of the here_ file in this repository.

.. _pre-commit: https://pre-commit.com/
.. _here: .pre-commit-config.yaml
