package accurics

{{.prefix}}{{.name}}{{.suffix}}[dockerFrom]{
	dockerFrom := input.docker_from[_]
    config := dockerFrom.config
    not contains(config, "healthcheck") 
}