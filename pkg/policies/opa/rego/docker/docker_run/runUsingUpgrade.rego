package accurics

{{.prefix}}{{.name}}{{.suffix}}[dockerRun.id] {
    dockerRun := input.docker_RUN[_]
    is_string(dockerRun.config)
    config := dockerRun.config
    checkupgrade(config)
}

{{.prefix}}{{.name}}{{.suffix}}[dockerRun.id] {
    dockerRun := input.docker_RUN[_]
    is_array(dockerRun.config)
    config := dockerRun.config
    checkupgradeList(config)
}

checkupgrade(config) {
    contains(config, ["apt-get upgrade", "apt-get dist-upgrade"][_])
}

checkupgradeList(config) {
    arrayList := ["apt-get upgrade", "apt-get dist-upgrade"]
    some i
    checkyumList := arrayList[i]
    contains(config[_], checkyumList)
}