package accurics

clientCertificateEnabled[container_cluster.id] {
  container_cluster := input.google_container_cluster[_]
  master := container_cluster.config.master_auth[_]
  master.client_certificate_config[_].issue_client_certificate == true
}
