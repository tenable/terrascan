package accurics

{{.prefix}}{{.name}}{{.suffix}}[run.id]{

    run := input.run[_]
    yumCleanAllMissing := run.config
    startswith(yumCleanAllMissing, "yum") 
    not contains(yumCleanAllMissing, "yum clean all")

}