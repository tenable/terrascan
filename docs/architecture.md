# Architecture

Terrascan's architecture is built to be modular to facilitate adding IaC languages and policies. At a high level Terrascan is composed of the following architectural components: a command line interface, API server, runtime, pluggable IaC proviers, pluggable policy engine, notifier, and writter.

* Command Line Interface = Provides CLI input to Terrascan.
* API Server = Provider input to Terrascan through an API.
* Runtime = Performs input validation and process inputs
* IaC Providers = Converts IaC language into normalized JSON
* Policy Engine = Applies policies against normalized JSON
* Notifier = Providers webhooks for results of Terrascan scans.
* Writter = Writes results into various formats like JSON, YAML, or XML.

![Terrascan architecture](terrascan_architecture.png)


