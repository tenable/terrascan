package accurics

kubeDashboardDisabled[api.id] {
    api := input.azurerm_kubernetes_cluster[_]
    profile := api.config.addon_profile[_]
    dashboard := profile.kube_dashboard[_]
    dashboard.enabled == true
}
