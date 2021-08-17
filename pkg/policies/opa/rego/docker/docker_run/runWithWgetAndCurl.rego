package accurics

{{.prefix}}{{.name}}{{.suffix}}[run.id] {
    run := input.docker_run[_]    
    
    configArray := [config | config := input.docker_run[_].config] 
    configString := concat(" | ", configArray)
    
    contains(configString, "wget")
    contains(configString, "curl")
}