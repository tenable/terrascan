package accurics

defStorageAccessDisabled[api.id]{
    api := input.azurerm_storage_account_network_rules[_]
    not api.config.default_action == "Deny"
}