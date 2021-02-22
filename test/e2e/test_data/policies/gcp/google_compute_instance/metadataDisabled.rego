package accurics

{{.name}}[api.id] {
    api := array.concat(input.google_compute_instance, input.google_compute_project_metadata)[_]
    api.config.metadata != null
    lower(object.get(api.config.metadata, "{{.metaKey}}", "undefined")) != "true"
}

{{.name}}[api.id] {
    api := array.concat(input.google_compute_instance, input.google_compute_project_metadata)[_]
    propUndefinedOrNull(api.config, "metadata")
}

propUndefinedOrNull(obj, prop) = true {
	obj[prop] == null
}

propUndefinedOrNull(obj, prop) = true {
	object.get(obj, prop, "undefined") == "undefined"
}