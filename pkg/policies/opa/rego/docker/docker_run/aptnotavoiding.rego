package accurics

{{.prefix}}{{.name}}{{.suffix}}[dockerRun.id] {
    dockerRun := input.docker_run[_]
    is_string(dockerRun.config)
    config := dockerRun.config
    avoidAdditionalPackages(config)
}

avoidAdditionalPackages(config) {
	flags := ["--no-install-recommends", "apt::install-recommends=false"]
	contains(config, flags[_])
}

avoidAdditionalPackages(config) {
    arrayList := ["--no-install-recommends", "apt::install-recommends=false"]
    some i
    checkyumList := arrayList[i]
    contains(config[_], checkyumList)
}