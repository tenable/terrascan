package accurics

{{.prefix}}allowLessHosts[retVal] {
  redis := input.azurerm_redis_firewall_rule[_]
  sHosts := calculateHosts(redis.config.start_ip)
  eHosts := calculateHosts(redis.config.end_ip)
  abs(eHosts - sHosts) >= 256
  rc := "ewogICJzdGFydF9pcCI6ICI8c3RhcnRfaXA+IiwKICAiZW5kX2lwIjogIjxlbmRfaXA+Igp9"
  retVal := { "Id": redis.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": "", "Attribute": "", "AttributeDataType": "base64", "Expected": rc, "Actual": { "start_ip": redis.config.start_ip, "end_ip": redis.config.end_ip } }
}

calculateHosts(val) = ans {
  ipVals := split(val, ".")
  # 2^24 = 16777216, 2^16 = 65536, 2^8 = 256
  # no of hosts in IP p.q.r.s = (p * 2^24) + (q * 2^16) + (r * 2^8) + s
  ans = (to_number(ipVals[0]) * 16777216) + (to_number(ipVals[1]) * 65536) + (to_number(ipVals[2]) * 256) + to_number(ipVals[3])
}