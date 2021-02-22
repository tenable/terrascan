provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "ptshggarg1"
  location = "West US"
}

resource "azurerm_virtual_network" "accurics_virtual_network" {
  name                = "ptshggavn1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "accurics_gateway_subnet" {
  name                 = "ptshggags1"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.accurics_virtual_network.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_public_ip" "accurics_public_ip" {
  name                = "ptshggapi1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  allocation_method   = "Static"
  sku                 = "Standard"
  domain_name_label   = "accurics"
}

resource "azurerm_public_ip" "accurics_public_ip2" {
  name                = "ptshggapi2"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  allocation_method   = "Static"
  sku                 = "Standard"
  domain_name_label   = "accurics2"
}

resource "azurerm_application_gateway" "accurics_application_gateway" {
  name                = "ptshggaag1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku {
    name     = "WAF_v2"
    tier     = "WAF_v2"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "accurics_gateway_config"
    subnet_id = azurerm_subnet.accurics_gateway_subnet.id
  }

  frontend_port {
    name = "accurics_gateway_frontend_port"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "accurics_gateway_frontend_ip"
    public_ip_address_id = azurerm_public_ip.accurics_public_ip.id
  }

  backend_address_pool {
    name         = "accurics_gateway_backend_pool"
    ip_addresses = ["10.0.0.5"]
  }

  backend_http_settings {
    name                  = "backend_http_settings_name"
    cookie_based_affinity = "Disabled"
    path                  = "/"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 60
  }

  http_listener {
    name                           = "accurics_gateway_listener"
    frontend_ip_configuration_name = "accurics_gateway_frontend_ip"
    frontend_port_name             = "accurics_gateway_frontend_port"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "accurics_gateway_rule"
    rule_type                  = "Basic"
    http_listener_name         = "accurics_gateway_listener"
    backend_address_pool_name  = "accurics_gateway_backend_pool"
    backend_http_settings_name = "backend_http_settings_name"
  }
}

resource "azurerm_application_gateway" "accurics_application_gateway2" {
  name                = "ptshggaag2"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku {
    name     = "WAF_v2"
    tier     = "WAF_v2"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "accurics_gateway_config"
    subnet_id = azurerm_subnet.accurics_gateway_subnet.id
  }

  frontend_port {
    name = "accurics_gateway_frontend_port"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "accurics_gateway_frontend_ip"
    public_ip_address_id = azurerm_public_ip.accurics_public_ip2.id
  }

  backend_address_pool {
    name         = "accurics_gateway_backend_pool"
    ip_addresses = ["10.0.0.6"]
  }

  backend_http_settings {
    name                  = "backend_http_settings_name"
    cookie_based_affinity = "Disabled"
    path                  = "/"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 60
  }

  http_listener {
    name                           = "accurics_gateway_listener"
    frontend_ip_configuration_name = "accurics_gateway_frontend_ip"
    frontend_port_name             = "accurics_gateway_frontend_port"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = "accurics_gateway_rule"
    rule_type                  = "Basic"
    http_listener_name         = "accurics_gateway_listener"
    backend_address_pool_name  = "accurics_gateway_backend_pool"
    backend_http_settings_name = "backend_http_settings_name"
  }

  waf_configuration {
    enabled          = false
    firewall_mode    = "Detection"
    rule_set_type    = "OWASP"
    rule_set_version = "3.1"
  }
}
