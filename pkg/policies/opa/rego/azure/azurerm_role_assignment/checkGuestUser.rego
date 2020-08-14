package accurics

{{.prefix}}checkGuestUser[role_assignment.id] {
  role_assignment := input.azurerm_role_assignment[_]
  role_assignment.config.role_definition_name == "Guest"
}
