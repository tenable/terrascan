package accurics

{{.prefix}}listenerNotHttps[listener.id] {
    listener = input.aws_lb_listener[_]
    upper(listener.config.protocol) == "HTTP"
    not listener.default_action.redirect.protocol
}

{{.prefix}}listenerNotHttps[listener.id] {
    listener = input.aws_lb_listener[_]
    upper(listener.config.protocol) == "HTTP"
    upper(listener.default_action.redirect.protocol) != "HTTPS"
}

{{.prefix}}listenerNotHttps[listener.id] {
    listener = input.aws_lb_listener[_]
    listener.config.port == 80
}