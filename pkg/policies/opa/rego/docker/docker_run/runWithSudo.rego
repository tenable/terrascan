package accurics

{{.prefix}}{{.name}}{{.suffix}}[run.id]{
    
    run := input.run[_]
    ContainsSudo := run.config{
    re_match("^( )*sudo", ContainsSudo)
    }
    ContainsSudo := run.config{
    re_match("( )*&& sudo", ContainsSudo)
    }

}