# Overview

Terrascan can be integrated into many tools in the development pipeline. When integrated into a tool, vulnerability scanning is automated as part of the commit or build process.
It can run on a developer's laptop, a SCM (e.g. GitHub), and CI\CD servers (e.g. ArgoCD and Jenkins). It also has a built in Admission Controller for Kubernetes. 

Please see the following guides for integrating Terrascan in different use cases. If the product you want to integrate with is not listed, do not fret. Terrascan supports many output formats (**YAML**, **JSON**, **XML**, **JUNIT-XML** and **SARIF**) to suit the variety of tools in the ecosystem. For example, it's straightforward to integrate with **Jenkins** using the **JUNIT-XML** format.

Go to the [Usage](../usage/command_line_mode.md#configuring-the-output-format-for-a-scan) page for more details.

### Integration Guides:

1. [Kubernetes (K8s) Admissions webhooks](admission-controller-webhooks-usage.md)
2. [ArgoCD](argocd-integration.md)
3. [Atlantis](atlantis-integration.md)
4. [Github and GitLab](cicd.md)

### Community Guides and Blogs:
* [Azure DevOps](https://lgulliver.github.io/terrascan-in-azure-devops/) Credit to  [@lrgulliver](https://twitter.com/lrgulliver) (Liam Gulliver)
