package accurics

{{.prefix}}{{.name}}{{.suffix}}[managed_disk.id] {
    managed_disk := input.azurerm_managed_disk[_]
    checkEncryption(managed_disk.config)
}

checkEncrytion(inputConfig) {
    inputConfig.encryption_settings.enabled != true
}

checkEncryption(inputConfig) {
    count(inputConfig.encryption_settings) == 0
}

checkEncryption(inputConfig) {
    not inputConfig.encryption_settings
}