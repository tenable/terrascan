package accurics

{{.prefix}}{{.name}}{{.suffix}}[dockerShell.id]{
	dockerShell := input.docker_dockerfile[_]
    config := dockerShell.config 
    contains(config, "shell") 
}