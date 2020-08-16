from diagrams import Cluster, Diagram
from diagrams.aws.compute import ECS
from diagrams.aws.management import Cloudformation
from diagrams.aws.integration import ConsoleMobileApplication
from diagrams.azure.compute import VM
from diagrams.gcp.compute import GCE
from diagrams.programming.language import Bash
from diagrams.onprem import iac
from diagrams.onprem.compute import Server


with Diagram("Terrascan architecture", show=False):
    cli = Bash("CLI")
    server = Server("API server")
    notifier = ConsoleMobileApplication("Notifier (Webhook)")
    writer = Bash("Writer (JSON, YAML, XML)")

    with Cluster("Runtime"):
        ECS("Input Validate")
        ECS("Process")
        output = ECS("Output")

    with Cluster("IaC Providers"):
        tf = iac.Terraform("Terraform")
        ansible = iac.Ansible("Ansible")
        cft = Cloudformation("CloudFormation")


    with Cluster("Policy Engine"):
        policy = [
            ECS("AWS"),
            VM("Azure"),
            GCE("GCP")
        ]

    server >> output >> tf >> policy >> notifier
    cli >> output >> ansible >> policy >> writer
    cft >> policy

