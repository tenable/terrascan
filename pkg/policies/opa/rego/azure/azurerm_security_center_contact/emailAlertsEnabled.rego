package accurics

emailAlertsEnabled[api.id]{
    api := input.azurerm_security_center_contact[_]
    not api.config.alert_notifications == true
}