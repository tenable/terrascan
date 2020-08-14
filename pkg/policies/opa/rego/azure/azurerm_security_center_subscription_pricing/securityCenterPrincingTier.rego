package accurics

{{.prefix}}securityCenterPrincingTier[subscription.id] {
  subscription := input.azurerm_security_center_subscription_pricing[_]
  subscription.config.tier != "Standard"
}