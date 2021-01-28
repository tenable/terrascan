package accurics

{{.prefix}}checkStorageContainerAccess[storage_container.id] {
  storage_container := input.azurerm_storage_container[_]
  not checkAccessType(storage_container.config.container_access_type)
}

checkAccessType(accesstype) {
  contains(accesstype, "private")
}

checkAccessType(accesstype) {
  contains(accesstype, "PRIVATE")
}