package accurics

{{.prefix}}{{.name}}{{.suffix}}[managed_disk.id] {
    managed_disk := input.azurerm_managed_disk[_]
    encryption_settings := managed_disk.config.encryption_settings[_]
    encryption_settings.enabled == false
}

{{.prefix}}{{.name}}{{.suffix}}[managed_disk.id] {
    managed_disk := input.azurerm_managed_disk[_]
    managed_disk.config.encryption_settings == []
}