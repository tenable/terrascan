package accurics

{{.prefix}}{{.name}}{{.suffix}}[endpoint_slice.id] {
   endpoint_slice = input.kubernetes_endpoint_slice[_]
   address := endpoint_slice.config.endpoints[_].addresses[_]

   not_allowed_addresses := ["127.0.0.0/8", "169.254.0.0/16"]
   net.cidr_contains(not_allowed_addresses[_], address)
}