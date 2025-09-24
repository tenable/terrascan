package accurics

{{.prefix}}{{.name}}{{.suffix}}[dockerRun]{
	dockerRun := input.docker_RUN[_]
    config := dockerRun.config
    regex.match(`^cd\s`, config)
}
