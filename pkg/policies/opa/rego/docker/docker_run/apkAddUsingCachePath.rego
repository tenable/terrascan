package accurics

{{.prefix}}{{.name}}{{.suffix}}[run.id]
{
run := input.run[_]
config := run.config
containsApkAddWithoutNoCache(config)
}

 containsApkAddWithoutNoCache(config) {
	command := trim_space(config)
	startswith(command, "apk ")
	contains(command, " add ")
	not contains(command, "--no-cache")
}