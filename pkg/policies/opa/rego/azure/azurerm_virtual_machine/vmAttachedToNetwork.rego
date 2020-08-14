package accurics

{{.prefix}}vmAttachedToNetwork[vm.id] {
  vm := input.azurerm_virtual_machine[_]
  vm.type == "azurerm_virtual_machine"
  count(object.get(vm.config, "network_interface_ids", "undefined")) <= 0
}
