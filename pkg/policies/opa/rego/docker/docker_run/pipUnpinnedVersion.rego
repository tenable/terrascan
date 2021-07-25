package accurics

{{.prefix}}{{.name}}{{.suffix}}[run.id]
{
    run := input.run[_]
    config := run.config
    version := regex.find_n("pip(3)? (-(-)?[a-zA-Z]+ *)*install", config, -1)
	version != null
}