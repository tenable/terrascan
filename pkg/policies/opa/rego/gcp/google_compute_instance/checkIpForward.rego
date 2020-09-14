package accurics

{{.prefix}}{{.name}}{{.suffix}}[api.id]
{
     api := input.google_compute_instance[_]
     api.config.can_ip_forward == true
}

