package accurics

networkPolicyEnabled[api.id] {
	api := input.azurerm_kubernetes_cluster[_]
	profile := api.config.network_profile[_]
	not checkNetworkPolicy(profile.network_policy)
}

networkPolicyEnabled[api.id] {
	api := input.azurerm_kubernetes_cluster[_]
	object.get(api.config, "network_profile", "undefined") == "undefined"
}

checkNetworkPolicy(policy) {
	policy == "azure"
}

checkNetworkPolicy(policy) {
	policy == "calico"
}
