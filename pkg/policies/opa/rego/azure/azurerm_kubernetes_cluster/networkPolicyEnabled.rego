package accurics

networkPolicyEnabled[api.id]{
    api := input.azurerm_kubernetes_cluster[_]
    profile := api.config.network_profile[_]
    profile.network_policy != "azure"
}

networkPolicyEnabled[api.id]{
    api := input.azurerm_kubernetes_cluster[_]
    object.get(api.config, "network_profile", "undefined") == "undefined"
}