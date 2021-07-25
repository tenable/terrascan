package accurics

{{.prefix}}{{.name}}{{.suffix}}[run.id]{

    run := input.run[_]
    UnpinnedVersion := run.config
    aptGet := regex.find_n("apt-get (-(-)?[a-zA-Z]+ *)*install", UnpinnedVersion, -1)
    aptGet != null

}