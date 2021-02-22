package accurics

{{.prefix}}storageAccountOpenToPublic[retVal] {
  storage_account := input.azurerm_storage_account[_]
  some i
  network_rule := storage_account.config.network_rules[i]
  arrayContains(network_rule.ip_rules, "0.0.0.0/0") == true
  expected := [ item | item := replace_cidr(storage_account.config.network_rules[i].ip_rules[_]) ]
  traverse := sprintf("network_rules[%d].ip_rules", [i])
  retVal := { "Id": storage_account.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": traverse, "AttributeDataType": "list", "Expected": expected, "Actual": storage_account.config.network_rules[i].ip_rules }
}

arrayContains(items, elem) = true {
  items[_] = elem
} else = false {
	true
}

replace_cidr(cidr) = value {
	cidr == "0.0.0.0/0"
    value := "<cidr>"
}

replace_cidr(cidr) = value {
	cidr != "0.0.0.0/0"
    value := cidr
}