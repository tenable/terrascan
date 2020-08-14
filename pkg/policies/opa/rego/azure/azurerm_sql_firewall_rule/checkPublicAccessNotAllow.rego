package accurics

{{.prefix}}{{.name}}[retVal] {
    sql_rule := input.azurerm_sql_firewall_rule[_]
    sql_rule.config.start_ip_address == "0.0.0.0"
    sql_rule.config.end_ip_address == "0.0.0.0"
    {{.isEntire}}
    retVal := { "Id": sql_rule.id, "ReplaceType": "delete", "CodeType": "resource", "Traverse": "", "Attribute": "", "AttributeDataType": "resource", "Expected": null, "Actual": null }
}

{{.prefix}}{{.name}}[retVal] {
    sql_rule := input.azurerm_sql_firewall_rule[_]
    sql_rule.config.start_ip_address == "0.0.0.0"
    sql_rule.config.end_ip_address == "255.255.255.255"
    not {{.isEntire}}
    retVal := { "Id": sql_rule.id, "ReplaceType": "delete", "CodeType": "resource", "Traverse": "", "Attribute": "", "AttributeDataType": "resource", "Expected": null, "Actual": null }
}