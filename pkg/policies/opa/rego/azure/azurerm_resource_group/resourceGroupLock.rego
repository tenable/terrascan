package accurics

{{.prefix}}resourceGroupLock[retVal] {
    resource_group := input.azurerm_resource_group[_]
    registry_input := input
    resource_group.type == "azurerm_resource_group"

	not resourceLockExist(resource_group)

    rc = "cmVzb3VyY2UgImF6dXJlcm1fbWFuYWdlbWVudF9sb2NrIiAiIyNyZXNvdXJjZV9uYW1lIyMiIHsKICBuYW1lICAgICAgID0gImF6dXJlcm1fcmVzb3VyY2VfZ3JvdXAuIyNyZXNvdXJjZV9uYW1lIyMiCiAgc2NvcGUgICAgICA9IGF6dXJlcm1fcmVzb3VyY2VfZ3JvdXAuIyNyZXNvdXJjZV9uYW1lIyMuaWQKICBsb2NrX2xldmVsID0gIkNhbk5vdERlbGV0ZSIKICAjIGF6dXJlcm1fbWFuYWdlbWVudF9sb2NrIGRvZXMgbm90IGNvbnRhaW4gdGFncywgYW5kIHdlIGNhbm5vdCBtYXRjaCB0aGVtIG5vdCB1bmxlc3MgdGhlIHJlc291cmNlIGlzIGRlcGxveWVkIGluIHRoZSBjbG91ZC4KICBub3RlcyAgICAgID0gIkNhbm5vdCBEZWxldGUgUmVzb3VyY2UiCn0="
    decode_rc = base64.decode(rc)
    replaced_resource_group_id := replace(decode_rc, "##resource_name##", resource_group.name)

    traverse = ""
    retVal := { "Id": resource_group.id, "ReplaceType": "add", "CodeType": "resource", "Traverse": traverse, "Attribute": "", "AttributeDataType": "resource", "Expected": base64.encode(replaced_resource_group_id), "Actual": null }
}

resourceLockExist(resource_group) = exists {
	resource_lock_exist_set := { resource_lock_id | resource_lock_id := input.azurerm_management_lock[i].config.scope }
	resource_lock_exist_set[resource_group.id]
    exists = true
} else = exists {
	resource_lock_exist_set := { resource_id | resource_id := input.azurerm_management_lock[i].config.name }
    resource_group_name := sprintf("azurerm_resource_group.%s", [resource_group.name])
	resource_lock_exist_set[resource_group_name]
    exists = true
}