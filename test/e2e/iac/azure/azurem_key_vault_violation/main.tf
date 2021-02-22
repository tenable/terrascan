provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "resource_group_details" {
  name     = "my-resource-group"
  location = "West US"
}

resource "azurerm_key_vault" "vault_details" {
  name                = "keyvaultkeyexample"
  location            = azurerm_resource_group.resource_group_details.location
  resource_group_name = azurerm_resource_group.resource_group_details.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  soft_delete_enabled = false

  sku_name = "premium"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "create",
      "get",
    ]

    secret_permissions = [
      "set",
    ]
  }
}

resource "azurerm_monitor_diagnostic_setting" "key_vault_logging" {
  name               = "azurerm_key_vault_vault_details_logs"
  target_resource_id = azurerm_key_vault.vault_details.id
  storage_account_id = azurerm_storage_account.storage_account.id
  log {
    category = "AuditEvent"
    enabled  = true
  }
}

resource "azurerm_storage_account" "storage_account" {
  name                      = "examplestorageaccount987"
  resource_group_name       = azurerm_resource_group.resource_group_details.name
  location                  = azurerm_resource_group.resource_group_details.location
  account_tier              = "Standard"
  account_replication_type  = "LRS"
  enable_https_traffic_only = false

  network_rules {
    default_action = "Allow"
    bypass         = ["None", "Logging", "Metrics"]
    ip_rules       = ["0.0.0.0/0"]
  }
}
