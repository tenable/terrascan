### this policy depends on the parameters specified by the user/client. Here we are considering that no hostPath are allowed###
package accurics

#rule for pod
{{.prefix}}{{.name}}{{.suffix}}[pod.id] {
    pod := input.kubernetes_pod[_]
    vols := pod.config.spec.volumes[_]
    parameters := {}
    has_field(vols, "hostPath")
    allowedPaths := get_allowed_paths(parameters)
    input_hostpath_violation(allowedPaths, vols)
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
    vols := kind.config.spec.template.spec.volumes[_]
    #parameters := {  'allowedHostPath' :[{ 'readOnly': true, 'pathPrefix': '/foo' }] }
    parameters := {}
    has_field(vols, "hostPath")
    allowedPaths := get_allowed_paths(parameters)
    input_hostpath_violation(allowedPaths, vols)
}

#rule for cron_job
{{.prefix}}{{.name}}{{.suffix}}[cron_job.id] {
    cron_job := input.kubernetes_cron_job[_]
    vols := cron_job.config.spec.jobTemplate.spec.template.spec.volumes[_]
    #parameters := {  'allowedHostPath' :[{ 'readOnly': true, 'pathPrefix': '/foo' }] }
    parameters := {}
    has_field(vols, "hostPath")
    allowedPaths := get_allowed_paths(parameters)
    input_hostpath_violation(allowedPaths, vols)
}

#function for all KINDs
has_field(object, field) = true {
    object[field]
}

#now allowed paths are null, this function will run##
get_allowed_paths(params) = out {
    not params.allowedHostPath == "undefined"
    out = []
}

input_hostpath_violation(allowedPaths, volume) {
    allowedPaths == []
}

### below functions are for violation when user has specified the hostPath, for testing uncomment the parameter array of objects at top####

get_allowed_paths(params) = out {
    out = params.allowedHostPath
}

input_hostpath_violation(allowedPaths, volume) {
    not input_hostpath_allowed(allowedPaths, volume)
}

input_hostpath_allowed(allowedPaths, volume) {
    allowedHostPath := allowedPaths[_]
    path_matches(allowedHostPath.pathPrefix, volume.hostPath.path)
    not allowedHostPath.readOnly == true
}

input_hostpath_allowed(allowedPaths, volume) {
    allowedHostPath := allowedPaths[_]
    path_matches(allowedHostPath.pathPrefix, volume.hostPath.path)
    allowedHostPath.readOnly
    not writeable_input_volume_mounts(volume.name)
}

writeable_input_volume_mounts(volume_name) {
	containers := input.kubernetes_pod[_].config.spec.containers[_]
    mount := containers.volumeMounts[_]
    mount.name == volume_name
    not mount.readOnly
}

path_matches(prefix, path) {
    a := split(trim(prefix, "/"), "/")
    b := split(trim(path, "/"), "/")
    prefix_matches(a, b)
}

prefix_matches(a, b) {
    count(a) <= count(b)
    not any_not_equal_upto(a, b, count(a))
}

any_not_equal_upto(a, b, n) {
    a[i] != b[i]
    i < n
}
