from diagrams import Cluster, Diagram
from diagrams.onprem.iac import Terraform


with Diagram("Terrascan architecture"):
    Terraform("IaC Provider")

