package accurics


{{.prefix}}{{.name}}{{.suffix}}[id] {
    add := input.docker[_]
    add.config.add != null
    id = add.id
}
