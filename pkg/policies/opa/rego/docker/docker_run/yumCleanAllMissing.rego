package accurics

{{.prefix}}{{.name}}{{.suffix}}[yumCleanAllMissing]{

    run := input.run[_]
    yumCleanAllMissing := run.config
    output := re_match("yum install", yumCleanAllMissing)
    output == true

}