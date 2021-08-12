package accurics

{{.prefix}}{{.name}}{{.suffix}}[run]{
    run := input.docker_run[_]
    config := run.config

    yum := regex.find_n("yum (-(-)?[a-zA-Z]+ *)*(group|local)?install", config, -1)
	yum != null
}