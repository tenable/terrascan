package accurics

{{.prefix}}{{.name}}{{.suffix}}[run.id]{

    run := input.run[_]
    aptGetInstall := run.config
    aptGet := regex.find_n("apt-get (-(-)?[a-zA-Z]+ *)*install", aptGetInstall, -1)
    aptGet != null

}