package accurics

{{.prefix}}{{.name}}{{.suffix}}[vio.id]{
	vio := input.docker_dockerfile[_]
    is_array(vio.config)
    config := vio.config 
    checkHealthCheck(config)
}

checkHealthCheck(config) {
      contains(config[_], "healthcheck") 
} 