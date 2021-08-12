package accurics

{{.prefix}}{{.name}}{{.suffix}}[dockerRun.id]{
	dockerRun := input.docker_from[_]
    config := dockerRun.config
    not contains(config, ":latest")    
}