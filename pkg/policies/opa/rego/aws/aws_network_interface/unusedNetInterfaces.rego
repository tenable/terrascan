package accurics

{{.prefix}}{{.suffix}}{{.name}}[netInterface.id]{
   netInterface := input.aws_network_interface[_] 
   config := netInterface.config	
   object.get(config, "attachment", "undefined") == "undefined"
}

{{.prefix}}{{.suffix}}{{.name}}[netInterface.id]{
   netInterface := input.aws_network_interface[_] 
   config := netInterface.config.attachment[_]
   config.instance == ""
}

