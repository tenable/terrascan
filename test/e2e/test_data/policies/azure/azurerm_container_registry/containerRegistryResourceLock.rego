package accurics

{{.prefix}}containerRegistryResourceLock[retVal] {
    registry := input.azurerm_container_registry[_]
    registry_input := input
    registry.type == "azurerm_container_registry"

	not resourceLockExist(registry, registry_input)

    rc = "cmVzb3VyY2UgImF6dXJlcm1fbWFuYWdlbWVudF9sb2NrIiAiIyNyZXNvdXJjZV9uYW1lIyMiIHsKICBuYW1lICAgICAgID0gImF6dXJlcm1fbWFuYWdlbWVudF9sb2NrLiMjcmVzb3VyY2VfbmFtZSMjIgogIHNjb3BlICAgICAgPSBhenVyZXJtX2NvbnRhaW5lcl9yZWdpc3RyeS4jI3Jlc291cmNlX25hbWUjIy5pZAogIGxvY2tfbGV2ZWwgPSAiQ2FuTm90RGVsZXRlIgogICMgYXp1cmVybV9tYW5hZ2VtZW50X2xvY2sgZG9lcyBub3QgY29udGFpbiB0YWdzLCBhbmQgd2UgY2Fubm90IG1hdGNoIHRoZW0gbm90IHVubGVzcyB0aGUgcmVzb3VyY2UgaXMgZGVwbG95ZWQgaW4gdGhlIGNsb3VkLgogIG5vdGVzICAgICAgPSAiQ2Fubm90IERlbGV0ZSBSZXNvdXJjZSIKfQ=="
    decode_rc = base64.decode(rc)
    replaced_registry_id := replace(decode_rc, "##resource_name##", registry.name)

    traverse = ""
    retVal := { "Id": registry.id, "ReplaceType": "add", "CodeType": "resource", "Traverse": traverse, "Attribute": "", "AttributeDataType": "resource", "Expected": base64.encode(replaced_registry_id), "Actual": null }
}

resourceLockExist(registry, registry_input) = exists {
	resource_lock_exist_set := { resource_lock_id | input.azurerm_management_lock[i].type == "azurerm_management_lock"; resource_lock_id := input.azurerm_management_lock[i].config.scope }
	resource_lock_exist_set[registry.id]
    exists = true
} else = exists {
	resource_lock_exist_set := { resource_id | input.azurerm_management_lock[i].type == "azurerm_management_lock"; resource_id := input.azurerm_management_lock[i].config.name }
    registry_name := sprintf("azurerm_container_registry.%s", [registry.name])
	resource_lock_exist_set[registry_name]
    exists = true
}