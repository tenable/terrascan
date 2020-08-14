package accurics

KubeDashboardDisabled[api.id]{
    api := input.azurerm_kubernetes_cluster[_]
    var := api.config.addon_profile[_]
    data := var.kube_dashboard[_]
    not data.enabled == false
}