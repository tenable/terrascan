package accurics


{{.prefix}}{{.name}}{{.suffix}}[id] {
    user := input.docker[_]
    user.config.user == null
    id = user.id
}
