package accurics

{{.prefix}}{{.name}}{{.suffix}}[apt]{
	apt := input.workdir[_]
	conval := apt.config
    
    not re_match("(^/[A-z0-9-_+]*)|(^[A-z0-9-_+]:\\\\.*)|(^\\$[{}A-z0-9-_+].*)", conval)
    
}
