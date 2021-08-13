package accurics

{{.prefix}}{{.name}}{{.suffix}}[dockerFrom.id]{
	dockerFrom := input.docker_from[_]
    config := dockerFrom.config
    not contains(config, ":latest")    
}