package accurics

{{.prefix}}{{.name}}{{.suffix}}[dockerRun]{
	dockerRun := input.docker_run[_]
    config := dockerRun.config
    contains(config, "cd")
}