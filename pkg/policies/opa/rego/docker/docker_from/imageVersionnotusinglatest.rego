package accurics

{{.prefix}}{{.name}}{{.suffix}}[dockerFrom.id]{
	dockerFrom := input.docker_FROM[_]
    config := dockerFrom.config
    contains(config, ":latest")
}

{{.prefix}}{{.name}}{{.suffix}}[dockerFrom.id]{
	dockerFrom := input.docker_FROM[_]
    config := dockerFrom.config
    config != "scratch"
    not contains(config, ":")
}
