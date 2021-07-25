package accurics

{{.prefix}}{{.name}}{{.suffix}}[run.id]{
    
    run := input.run[_]
    ContainsSudo := run.config
    regex.match("^( )*sudo", ContainsSudo) == true

}