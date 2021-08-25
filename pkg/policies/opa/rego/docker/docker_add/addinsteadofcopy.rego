package accurics

{{.prefix}}{{.name}}{{.suffix}}[add.id]{
	add := input.docker.add[_]
    is_string(add.config)
    config := add.config
    checkCopy(config)
}

{{.prefix}}{{.name}}{{.suffix}}[add.id] {
    add := input.docker.add[_]
    is_array(add.config)
    config := add.config
    checkCopyList(config)
}

checkCopy(config) {
    contains(config, [".tar", ".tar."][_])
}

checkCopyList(config) {
    arrayList := [".tar", ".tar."]
    some i
    checkcopyList := arrayList[i]
    contains(config[_], checkcopyList)
}