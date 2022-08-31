package accurics

{{.prefix}}{{.name}}{{.suffix}}[retval] {
    security_center_settings := input.azurerm_security_center_setting[_]
    upper(security_center_settings.config.setting_name) != "{{.value}}"
    security_center_settings.config.enabled != true

    retval := {
        "Id": security_center_settings.id,
        "ReplaceType": "edit",
        "CodeType": "attribute",
        "Traverse": "enabled",
        "Attribute": "enabled",
        "AttributeDataType": "bool",
        "Expected": true,
        "Actual": security_center_settings.config.enabled
    }
}