package accurics

{{.prefix}}{{.name}}{{.suffix}}[run.id]{
    run := input.docker_RUN[_]
    config := run.config
    
    
    commands := split(config, "&&")
    command := commands[_]
    contains(command, "yum install")
    len := count(regex.split("yum (group|local)?install ?(-(-)?[a-zA-Z]+ *)*", command))
    packages_array := array.slice(regex.split("yum (group|local)?install ?(-(-)?[a-zA-Z]+ *)*", command), 1, len)[0]
    packages := split(packages_array, " ")
    not checkVersion(packages)
}

checkVersion(arg) {
    pack := arg[_]
	re_match("[A-Za-z0-9_-]+[-:][$](.+)", pack)
}

checkVersion(arg) {
    pack := arg[_]
	re_match("[A-Za-z0-9_-]+[:-]([0-9]+.)+[0-9]+", pack)
}

checkVersion(arg) {
    pack := arg[_]
	re_match("[A-Za-z0-9_-]+=(.+)", pack)
}