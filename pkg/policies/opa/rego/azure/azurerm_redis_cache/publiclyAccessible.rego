package accurics

{{.prefix}}{{.name}}[retVal] {
    redis := input.azurerm_redis_firewall_rule[_]
    redis.config.start_ip == "0.0.0.0"
    redis.config.end_ip == "0.0.0.0"
    {{.isEntire}}
    retVal := { "Id": redis.id, "ReplaceType": "delete", "CodeType": "resource", "Traverse": "", "Attribute": "", "AttributeDataType": "resource", "Expected": null, "Actual": null }
}

{{.prefix}}{{.name}}[retVal] {
    redis := input.azurerm_redis_firewall_rule[_]
    redis.config.start_ip != "0.0.0.0"
    checkScopeIsPublic(redis.config.start_ip)
    redis.config.end_ip != "0.0.0.0"
    checkScopeIsPublic(redis.config.end_ip)
    not {{.isEntire}}
    retVal := { "Id": redis.id, "ReplaceType": "delete", "CodeType": "resource", "Traverse": "", "Attribute": "", "AttributeDataType": "resource", "Expected": null, "Actual": null }
}

scopeIsPrivate(scope) {
  private_ips = ["10.0.0.0/8", "192.168.0.0/16", "172.16.0.0/12"]
  net.cidr_contains(private_ips[_], scope)
}

checkScopeIsPublic(val) = true {
  glob.match("[0-9]*.[0-9]*.[0-9]*.*", [], val)
  not scopeIsPrivate(val)
}