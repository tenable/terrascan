package accurics

{{.prefix}}checkStorageContainerAccess[storage_container.id] {
  storage_container := input.azurerm_storage_container[_]
  storage_container.config.container_access_type != "private"
}
