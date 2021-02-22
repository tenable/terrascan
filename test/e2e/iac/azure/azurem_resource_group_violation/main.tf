provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_management_lock" "example1" {
  name       = "azurerm_resource_group.example1"
  scope      = azurerm_resource_group.example.id
  lock_level = "CanNotDelete"
  # azurerm_management_lock does not contain tags, and we cannot match them not unless the resource is deployed in the cloud.
  notes      = "Cannot Delete Resource"
}