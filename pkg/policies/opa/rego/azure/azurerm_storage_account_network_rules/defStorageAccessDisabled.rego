package accurics

defStorageAccessDisabled[san_rules.id] {
	san_rules := input.azurerm_storage_account_network_rules[_]
	lower(san_rules.config.default_action) != "deny"
}
