package accurics

{{.prefix}}{{.name}}{{.suffix}}[dockerFrom]{
	dockerFrom := input.docker_dockerfile[_]
    is_array(dockerFrom.config)
    config := dockerFrom.config 
    checkHealthCheck(config)
}

checkHealthCheck(config) {
      contains(config, "healthcheck") 
} 