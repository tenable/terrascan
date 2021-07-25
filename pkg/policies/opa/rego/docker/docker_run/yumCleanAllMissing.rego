package accurics

{{.prefix}}{{.name}}{{.suffix}}[yumCleanAllMissing]{

    run := input.run[_]
    yumCleanAllMissing := run.config
    output := regex.match("yum install", yumCleanAllMissing)
    output == true

}