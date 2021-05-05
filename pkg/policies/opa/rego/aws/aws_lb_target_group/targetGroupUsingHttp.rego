package accurics

{{.prefix}}targetGroupUsingHttp[tg_group.id] {
    tg_group = input.aws_lb_target_group[_]
    upper(tg_group.config.protocol) == "HTTP"
}