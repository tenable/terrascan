package accurics

{{.prefix}}noSecurityGroupAssociated[retVal] {
  vn := input.azurerm_virtual_network[_]
  vn.type = "azurerm_virtual_network"
  object.get(vn.config, "subnet", "undefined") != "undefined"
  not sgExists(vn.config)
  
  traverse = "subnet[0].security_group"
  retVal := { "Id": vn.id, "ReplaceType": "add", "CodeType": "block", "Traverse": traverse, "Attribute": "subnet.security_group", "AttributeDataType": "string", "Expected": "${azurerm_network_security_group.<security_group_name>.id}", "Actual": "" }
}

{{.prefix}}noSecurityGroupAssociated[retVal] {
  vn := input.azurerm_virtual_network[_]
  vn.type = "azurerm_virtual_network"
  object.get(vn.config, "subnet", "undefined") == "undefined"
  rc = "ewogICJzdWJuZXQiOiB7CiAgICAibmFtZSI6ICJzdWJuZXQzIiwKICAgICJhZGRyZXNzX3ByZWZpeCI6ICI8Y2lkcj4iLAogICAgInNlY3VyaXR5X2dyb3VwIjogIiR7YXp1cmVybV9uZXR3b3JrX3NlY3VyaXR5X2dyb3VwLjxzZWN1cml0eV9ncm91cF9uYW1lPi5pZH0iCiAgfQp9"
  traverse = ""
  retVal := { "Id": vn.id, "ReplaceType": "add", "CodeType": "block", "Traverse": traverse, "Attribute": "subnet", "AttributeDataType": "base64", "Expected": rc, "Actual": null }
}

sgExists(cfg) = true {
	subs = cfg.subnet[_]
    subs.security_group != ""
}

sgExists(cfg) = true {
	subs = cfg.subnet[_]
    object.get(subs, "security_group", "undefined") == "undefined"
}