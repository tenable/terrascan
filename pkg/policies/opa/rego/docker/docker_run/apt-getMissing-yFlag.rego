package accurics

{{.prefix}}{{.name}}{{.suffix}}[dockerRun.id] {
    dockerRun := input.run[_]
    is_string(dockerRun.config)
    config := dockerRun.config
    checkAptGetUse(config)
    not checkManualInput(config)
}
{{.prefix}}{{.name}}{{.suffix}}[dockerRun.id] {
    dockerRun := input.run[_]
    is_array(dockerRun.config)
    config := dockerRun.config
    checkAptGetUseList(config)
    not checkManualInputList(config)
}
checkAptGetUse(config) {
    contains(config, "apt-get")
}
checkAptGetUseList(config) {
    contains(config[_], "apt-get")
}
checkManualInput(config) {
    flags := ["-y", "yes", "assumeyes", "-qy"]
    contains(config, flags[_])
}
checkManualInputList(config) {
    flags := ["-y", "yes", "assumeyes", "-qy"]
    some i
    checkFlag := flags[i]
    contains(config[_], checkFlag)
}