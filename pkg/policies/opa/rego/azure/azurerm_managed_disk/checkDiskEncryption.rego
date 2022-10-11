package accurics

{{.prefix}}{{.name}}{{.suffix}}[managed_disk.id] {
    managed_disk := input.azurerm_managed_disk[_]
    checkEncryption(managed_disk.config)
}

checkEncryption(inputConfig) {
    inputConfig.encryption_settings[_].enabled != true
}

checkEncryption(inputConfig) {    
    object.get(inputConfig, "encryption_settings", "undefined") == [null, false, "undefined", [], {}][_]
}
