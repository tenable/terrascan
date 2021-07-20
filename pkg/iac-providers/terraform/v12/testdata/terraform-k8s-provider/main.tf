terraform {
 required_providers {
   kubernetes-beta = {
     source  = "hashicorp/kubernetes"
     version = "= 2.3.2"
   }
 }
}

provider "kubernetes-beta" {
  config_path = "~/.kube/config"
}
resource "kubernetes_namespace" "test" {
  metadata {
    name = "nginx"
  }
}

resource "kubernetes_service" "test" {
  metadata {
    name      = "nginx"
    namespace = kubernetes_namespace.test.metadata.0.name
  }
  spec {
    selector = {
      app = kubernetes_deployment.example.spec.0.template.0.metadata.0.labels.test
    }
    type = "NodePort"
    port {
      node_port   = 30201
      port        = 80
      target_port = 80
    }
  }
}

resource "kubernetes_deployment" "example" {
  metadata {
    name = "terraform-example"
    namespace = kubernetes_namespace.test.metadata.0.name
    labels = {
      test = "MyExampleApp"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        test = "MyExampleApp"
      }
    }

    template {
      metadata {
        labels = {
          test = "MyExampleApp"
        }
      }

      spec {
        container {
          image = "nginx:1.7.8"
          name  = "example1"
        }

        container {
          image = "nginx:1.7.8"
          name  = "example2"
        }

        container {
          image = "nginx:1.7.8"
          name  = "example3"
        }

        service_account_name            = "default"
        automount_service_account_token = false
      }
    }
  }
}
