package accurics

####fixed the minimum set of allowed volumes, this may change as per the user####

#rule for pod_security_policy
{{.prefix}}{{.name}}{{.suffix}}[psp.id] {
    psp := input.kubernetes_pod_security_policy[_]
    secure_volumes := [{{range .secure_volumes}}{{- printf "%q" . }},{{end}}]
    volume_field := psp.config.spec.volumes[_]
    not input_volume_type_allowed(volume_field, secure_volumes)
}

#rule for pod
{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    secure_volumes := [{{range .secure_volumes}}{{- printf "%q" . }},{{end}}]
    volume_fields := {x | pod.config.spec.volumes[_][x]; x != "name"}
    field := volume_fields[_]
    not input_volume_type_allowed(field, secure_volumes)
}

#rule for deployment, daemonset, job, replica_set, stateful_set, replication_controller
{{.prefix}}{{.name}}{{.suffix}}[kind.id] {
    item_list := [
        object.get(input, "kubernetes_daemonset", "undefined"),
        object.get(input, "kubernetes_deployment", "undefined"),
        object.get(input, "kubernetes_job", "undefined"),
        object.get(input, "kubernetes_replica_set", "undefined"),
        object.get(input, "kubernetes_replication_controller", "undefined"),
        object.get(input, "kubernetes_stateful_set", "undefined")
    ]

    item = item_list[_]
    item != "undefined"

    kind := item[_]
    secure_volumes := [{{range .secure_volumes}}{{- printf "%q" . }},{{end}}]
    volume_fields := {x | kind.config.spec.template.spec.volumes[_][x]; x != "name"}
    field := volume_fields[_]
    not input_volume_type_allowed(field, secure_volumes)
}

#rule for cron_job
{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    secure_volumes := [{{range .secure_volumes}}{{- printf "%q" . }},{{end}}]
    volume_fields := {x | cron_job.config.spec.jobTemplate.spec.template.spec.volumes[_][x]; x != "name"}
    field := volume_fields[_]
    not input_volume_type_allowed(field, secure_volumes)
}

input_volume_type_allowed(field, secure_volumes) {
    secure_volumes[_] == "*"
}

input_volume_type_allowed(field, secure_volumes) {
    field == secure_volumes[_]
}