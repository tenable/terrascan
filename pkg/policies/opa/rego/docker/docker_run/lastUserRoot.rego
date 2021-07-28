package accurics

{{.prefix}}{{.name}}{{.suffix}}[cmd.id]{
    cmd := input.user[count(input.user) - 1]
    cmd.config == "root"
}