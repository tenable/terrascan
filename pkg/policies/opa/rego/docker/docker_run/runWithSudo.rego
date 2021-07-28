package accurics

{{.prefix}}{{.name}}{{.suffix}}[run.id]{    
    run := input.run[_]    
    checkSudo(run.config)
}
checkSudo(config) {
	startswith(config, "sudo")
}
checkSudo(config) {
	contains(config, "&& sudo")
}