package accurics

{{.prefix}}{{.name}}{{.suffix}}[run.name] {
    run := input.docker_RUN[_]    
    config := run.config
    config == ["apt-get update", "yum update", "sudo apt-get update", "sudo yum update"][_]
}