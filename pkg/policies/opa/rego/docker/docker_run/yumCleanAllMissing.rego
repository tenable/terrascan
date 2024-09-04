package accurics

{{.prefix}}{{.name}}{{.suffix}}[run.id]{

    run := input.docker_RUN[_]
    yumCleanAllMissing := run.config
    startswith(yumCleanAllMissing, "yum") 
    not contains(yumCleanAllMissing, "yum clean all")

}