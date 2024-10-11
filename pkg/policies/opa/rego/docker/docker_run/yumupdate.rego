package accurics

{{.prefix}}{{.name}}{{.suffix}}[dockerRun.id] {
    dockerRun := input.docker_RUN[_]
    is_string(dockerRun.config)
    config := dockerRun.config
    checkyumUpdate(config)
}

{{.prefix}}{{.name}}{{.suffix}}[dockerRun.id] {
    dockerRun := input.docker_RUN[_]
    is_array(dockerRun.config)
    config := dockerRun.config
    checkyumUpdateArray(config)
}

checkyumUpdate(config) {
   contains(config, ["yum update", "yum update-to", "yum upgrade", "yum upgrade-to"][_])
}

checkyumUpdateArray(config) {
    arrayList := ["yum update", "yum update-to", "yum upgrade", "yum upgrade-to"]
    some i
    checkyumList := arrayList[i]
    contains(config[_], checkyumList)
} 