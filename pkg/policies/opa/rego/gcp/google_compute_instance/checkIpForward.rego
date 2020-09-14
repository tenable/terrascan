package accurics

checkIpForward[api.id]
{
     api := input.google_compute_instance[_]
     not api.config.can_ip_forward == true
}