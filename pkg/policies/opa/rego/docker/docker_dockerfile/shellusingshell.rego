package accurics

{{.prefix}}{{.name}}{{.suffix}}[dockerShell.id]{
	dockerShell := input.docker_dockerfile[_]
    config := dockerShell.config 
    is_array(dockerFrom.config)
    checkShell(config)
}

checkShell(config) {
     contains(config, "shell") 
} 