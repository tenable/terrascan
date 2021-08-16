package accurics

{{.prefix}}{{.name}}{{.suffix}}[command.id]
{
   item_list := [
        object.get(input, "docker_add", "undefined"),
        object.get(input, "docker_copy", "undefined"),
        object.get(input, "docker_dockerfile", "undefined"),
        object.get(input, "docker_from", "undefined"),
        object.get(input, "docker_run", "undefined"),
    ]

    item = item_list[_]
    item != "undefined"
    command := item[_]
    
    instructions := {"copy", "add", "run"}
	some i
	check := [x | item[i].config == instructions[y]; x := item[i]]
    
    some j, k
	Counter := [x | check[j].line - check[k].line == -1; x := check[j]]
}