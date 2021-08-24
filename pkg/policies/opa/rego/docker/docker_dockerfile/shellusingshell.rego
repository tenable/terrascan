package accurics

{{.prefix}}{{.name}}{{.suffix}}[dockerShell.id]{
	dockerShell := input.docker_dockerfile[_]
    config := dockerShell.config 
    is_array(dockerShell.config)
    checkShell(config)
}

checkShell(config) {
     contains(config[_], "shell") 
} 