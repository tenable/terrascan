package accurics

{{.prefix}}{{.name}}{{.suffix}}[run.id]
{
    run := input.run[_]
    config := run.config
    
    contains(config, ["sudo apt-get install*", "apt-get install*"][_])
    contains(config, "=")
    cleanString := trim_space(split(config, "=")[1])
    packages := split(cleanString, " ")
    pack := packages[_]
    not withVersion(pack)  
}

withVersion(pack) {
	re_match("[A-Za-z0-9_-]+[-:][$](.+)", pack)
}

withVersion(pack) {
	regex.match("[A-Za-z0-9_-]+[:-]([0-9]+.)+[0-9]+", pack)
}

withVersion(pack) {
	regex.match("[A-Za-z0-9_-]+=(.+)", pack)
}

