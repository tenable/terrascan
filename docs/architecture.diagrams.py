from diagrams import Cluster, Diagram
from diagrams.onprem import iac
from diagrams.onprem.client import Client
from diagrams.onprem.compute import Server
from diagrams.aws.compute import ECS
from diagrams.azure.compute import VM
from diagrams.gcp.compute import GCE

with Diagram("Terrascan architecture", show=False):
    cli = Client("CLI")
    server = Server("API server")

    with Cluster("Runtime"):
        runtime = [
            ECS("Input Validate"),
            ECS("Process"),
            ECS("Output")
        ]

    with Cluster("IaC Providers"):
        tf = iac.Terraform("Terraform")
        ansible = iac.Ansible("Ansible")
        iac.Awx("Awx")


    with Cluster("Policy Engine"):
        policy = [
            ECS("AWS"),
            VM("Azure"),
            GCE("GCP")
        ]

    server >> runtime >> tf >> policy
    cli >> runtime >> ansible >> policy

