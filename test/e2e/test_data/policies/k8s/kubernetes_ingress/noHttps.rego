package accurics

{{.prefix}}{{.name}}{{.suffix}}[ingress.id] {
    ingress = input.kubernetes_ingress[_]
    re_match("^(extensions|networking.k8s.io)", ingress.config.apiVersion) #can be from two apis "extensions", "networking.k8s.io"
    not https_complete(ingress.config)
}
##two conditions ingress spec should have a tls key map and annotation kubernetes.io/ingress.allow-http = false
https_complete(arg) = true {
    object.get(arg.spec, "tls", "undefined") != "undefined"
    arg.metadata.annotations["kubernetes.io/ingress.allow-http"] == "false"
}