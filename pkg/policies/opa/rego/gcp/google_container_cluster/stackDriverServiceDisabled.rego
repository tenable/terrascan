package accurics

{{.name}}[container_cluster.id] {
  container_cluster := input.google_container_cluster[_]
  container_cluster.config["{{.service}}_service"] != "{{.service}}.googleapis.com/kubernetes"
}
