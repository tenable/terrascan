package accurics

{{.prefix}}{{.suffix}}{{.name}}[netInterface.id]{
   netInterface := input.aws_network_interface[_]
   config := netInterface.config
   some i
   secGroups := config.security_groups[i]
   name = split(secGroups, ".")[1] 
   secGroupCheck(name)
}  

secGroupCheck(a) {
    secConfig := input.aws_security_group[_]
    groupId := secConfig.name
    groupId == a
    secConfig.config.ingress[_].cidr_blocks[_] == "0.0.0.0/0"
}