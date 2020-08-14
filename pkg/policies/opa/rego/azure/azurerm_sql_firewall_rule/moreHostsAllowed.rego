package accurics

{{.prefix}}moreHostsAllowed[retVal] {
  sql_rule := input.azurerm_sql_firewall_rule[_]
  sHosts := calculateHosts(sql_rule.config.start_ip_address)
  eHosts := calculateHosts(sql_rule.config.end_ip_address)
  abs(eHosts - sHosts) >= 256
  rc := "ewogICJzdGFydF9pcF9hZGRyZXNzIjogIjxzdGFydF9pcF9hZGRyZXNzPiIsCiAgImVuZF9pcF9hZGRyZXNzIjogIjxlbmRfaXBfYWRkcmVzcz4iCn0="
  retVal := { "Id": sql_rule.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": "", "Attribute": "", "AttributeDataType": "base64", "Expected": rc, "Actual": { "start_ip_address": sql_rule.config.start_ip_address, "end_ip_address": sql_rule.config.end_ip_address } }
}

calculateHosts(val) = ans {
  ipVals := split(val, ".")
  # 2^24 = 16777216, 2^16 = 65536, 2^8 = 256
  # no of hosts in IP p.q.r.s = (p * 2^24) + (q * 2^16) + (r * 2^8) + s
  ans = (to_number(ipVals[0]) * 16777216) + (to_number(ipVals[1]) * 65536) + (to_number(ipVals[2]) * 256) + to_number(ipVals[3])
}