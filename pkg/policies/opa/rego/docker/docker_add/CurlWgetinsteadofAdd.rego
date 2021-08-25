package accurics

{{.prefix}}{{.name}}{{.suffix}}[add.id]{
    add := input.docker.add[_]
    is_string(add.config)
    config := add.config
    checkWget(config)
}

{{.prefix}}{{.name}}{{.suffix}}[add.id] {
    add := input.docker.add[_]
    is_array(add.config)
    config := add.config
    checkWgetList(config)
}

checkWget(config) {
   re_match("https?", config)
}

 checkWgetList(config) {
     re_match("https?", config)
 }