package accurics

{{.prefix}}checkKeyExpirationIsSet[retVal] {
  vault_key := input.azurerm_key_vault_key[_]
  vault_key.config.expiration_date == null
  traverse = "expiration_date"
  expected := getExpiryRfc3339(time.now_ns())
  retVal := { "Id": vault_key.id, "ReplaceType": "add", "CodeType": "attribute", "Traverse": traverse, "Attribute": "expiration_date", "AttributeDataType": "string", "Expected": expected, "Actual": null }
}

{{.prefix}}checkKeyExpirationIsSet[retVal] {
  vault_key := input.azurerm_key_vault_key[_]
  vault_key.config.expiration_date != null
  now := time.now_ns()
  expiration := time.parse_rfc3339_ns(vault_key.config.expiration_date)
  (expiration - now) > (2 * 365 * 24 * 60 * 60 * 1000000000) # 2 years
  traverse = "expiration_date"
  expected := getExpiryRfc3339(now)
  retVal := { "Id": vault_key.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "expiration_date", "AttributeDataType": "string", "Expected": expected, "Actual": null }
}

getExpiryRfc3339(curtime) = expiry {
  expiryNs := time.add_date(curtime, 1, 0, 1)
  dateAr := time.date(expiryNs)
  timeAr := time.clock(expiryNs)
  expiry := sprintf("%d-%d-%dT%d:%d:%dZ", array.concat(dateAr, timeAr))
}