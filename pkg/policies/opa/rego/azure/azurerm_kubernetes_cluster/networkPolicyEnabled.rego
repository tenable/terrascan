package accurics

networkPolicyEnabled[api.id]{
    api := input.azurerm_kubernetes_cluster[_]
    var := api.config.network_profile[_]
    not var.network_policy == "azure"
    not var.network_policy == "calico"
}