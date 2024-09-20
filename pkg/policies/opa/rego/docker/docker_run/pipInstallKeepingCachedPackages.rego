package accurics

{{.prefix}}{{.name}}{{.suffix}}[run.id]{
    run := input.docker_RUN[_]
    config := run.config
    hasCacheFlag(config)
}

hasCacheFlag(values) {
	commands = split(values, "&&")
	some i
	instruction := commands[i]
	re_match("pip(3)? (-(-)?[a-zA-Z]+ *)*install", instruction) == true
	not contains(instruction, "--no-cache-dir")
}
