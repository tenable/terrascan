package accurics

{{.prefix}}{{.name}}{{.suffix}}[dockerFrom]{
	dockerFrom := input.docker_dockerfile[_]
    config := dockerFrom.config
    contains(config, "healthcheck") 
}