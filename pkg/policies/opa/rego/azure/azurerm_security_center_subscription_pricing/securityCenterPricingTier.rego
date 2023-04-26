package accurics

{{.prefix}}securityCenterPricingTier[subscription.id] {
  subscription := input.azurerm_security_center_subscription_pricing[_]
  subscription.config.tier != "Standard"
}