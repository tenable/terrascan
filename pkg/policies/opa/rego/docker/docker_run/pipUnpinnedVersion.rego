package accurics

{{.prefix}}{{.name}}{{.suffix}}[run.id]
{
    run := input.docker_RUN[_]
    config := run.config
    
    contains(config, ["pip install", "pip3 install"][_])
    contains(config, "--upgrade")
    cleanString := trim_space(split(config, "--upgrade")[1])
    packages := split(cleanString, " ")
    pack := packages[_]
    not withVersion(pack)  
}

{{.prefix}}{{.name}}{{.suffix}}[run.id]
{
    run := input.docker_RUN[_]
    config := run.config  
    contains(config, ["pip install", "pip3 install"][_])
    not contains(config, "--")
    cleanString := trim_space(split(config, "install")[1])
    packages := split(cleanString, " ")
    pack := packages[_]
}

withVersion(pack) {
	re_match("[A-Za-z0-9_-]+[-:][$](.+)", pack)
}

withVersion(pack) {
	re_match("[A-Za-z0-9_-]+[:-]([0-9]+.)+[0-9]+", pack)
}

withVersion(pack) {
	re_match("[A-Za-z0-9_-]+=(.+)", pack)
}

