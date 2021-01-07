package accurics

#rule for pod_security_policy
{{.prefix}}{{.name}}{{.suffix}}[psp.id] {
    psp := input.kubernetes_pod_security_policy[_]
    affected_volumes := ["glusterfs", "quobyte", "storageos", "scaleIO"]
    volume_type := psp.config.spec.volumes[_]
    volNotAllowed(volume_type, affected_volumes)
}

#rule for pod
{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    affected_volumes := ["glusterfs", "quobyte", "storageos", "scaleIO"]
    volume_types := {x | pod.config.spec.volumes[_][x]; x != "name"}
    vol:= volume_types[_]
    volNotAllowed(vol, affected_volumes)
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
    affected_volumes := ["glusterfs", "quobyte", "storageos", "scaleIO"]
    volume_types := {x | kind.config.spec.template.spec.volumes[_][x]; x != "name"}
    vol:= volume_types[_]
    volNotAllowed(vol, affected_volumes)
}

#rule for cron_job
{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    affected_volumes := ["glusterfs", "quobyte", "storageos", "scaleIO"]
    volume_types := {x | cron_job.config.spec.jobTemplate.spec.template.spec.volumes[_][x]; x != "name"}
    vol:= volume_types[_]
    volNotAllowed(vol, affected_volumes)
}

volNotAllowed(field, affected_volumes) {
    field == affected_volumes[_]
}