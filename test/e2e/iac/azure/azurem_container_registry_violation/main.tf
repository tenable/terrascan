provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "resource_group" {
  name     = "ptshggag1"
  location = "West US"
}

resource "azurerm_container_registry" "sample_container_registry" {
  name                     = "ptshggacr1"
  resource_group_name      = azurerm_resource_group.resource_group.name
  location                 = azurerm_resource_group.resource_group.location
  sku                      = "Premium"
  admin_enabled            = true
  georeplication_locations = ["East US", "West Europe"]
}

# Commented out on purpose
//resource "azurerm_management_lock" "resource-group-level" {
//  name       = "azurerm_container_registry.sample_container_registry_locked"
//  scope      = azurerm_container_registry.sample_container_registry_locked.id
//  lock_level = "CanNotDelete"
//  notes      = "${azurerm_container_registry.sample_container_registry_locked.name}"
//}
