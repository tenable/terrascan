package accurics

{{.prefix}}{{.name}}{{.suffix}}[dockerRun]{
	dockerRun := input.docker_run[_]
    config := dockerRun.config
    re_match(".*apk.*upgrade", config)  
}