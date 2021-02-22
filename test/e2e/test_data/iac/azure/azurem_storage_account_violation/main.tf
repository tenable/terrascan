provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "a_rg" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "a_vn" {
  name                = "virtnetname"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.a_rg.location
  resource_group_name = azurerm_resource_group.a_rg.name
}

resource "azurerm_subnet" "a_subnet" {
  name                 = "subnetname"
  resource_group_name  = azurerm_resource_group.a_rg.name
  virtual_network_name = azurerm_virtual_network.a_vn.name
  address_prefixes       = ["10.0.2.0/24"]
  service_endpoints    = ["Microsoft.Sql", "Microsoft.Storage"]
}

resource "azurerm_storage_account" "storageRulesCheck" {
  name                      = "storageaccountname"
  resource_group_name       = "some-group"
  location                  = "westus"
  account_tier              = "Standard"
  account_replication_type  = "GRS"
  enable_https_traffic_only = false
  blob_properties {

  }
  tags = {
    environment = "staging"
  }

  network_rules {
    default_action = "Allow"
    bypass         = ["None", "Logging", "Metrics"]
    ip_rules       = ["0.0.0.0/0"]
  }
}

resource "azurerm_storage_account_network_rules" "defStorageAccessDisabled" {
  resource_group_name  = "resourcename"
  storage_account_name = "myacc"
  default_action       = "Allow"
}

resource "azurerm_storage_account_network_rules" "noAzureServices" {
  resource_group_name  = "resourceName"
  storage_account_name = "myacc"

  default_action = "Allow"
  ip_rules       = ["127.0.0.1"]
  bypass         = ["Metrics"]
}
