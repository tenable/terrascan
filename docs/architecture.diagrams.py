from diagrams import Cluster, Diagram
from diagrams.onprem.iac import Terraform
from diagrams.onprem.client import Client

with Diagram("Terrascan architecture", show=False):
    Client("CLI") >> Terraform("IaC Provider")

