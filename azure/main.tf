provider "azurerm" {
    features {}
    subscription_id = var.subscriptionId
}

resource "azurerm_resource_group" "resource_group" {
    name     = var.resourceGroupName
    location =  var.location
}

resource "azurerm_container_registry" "container_registry" {
    name                = var.containerRegistryName
    resource_group_name = azurerm_resource_group.resource_group.name
    location            = azurerm_resource_group.resource_group.location
    sku                 = "Basic"
    admin_enabled       = true
}

resource "azurerm_app_service_plan" "app_service_plan" {
    name                = "pdappserviceplan"
    location            = azurerm_resource_group.resource_group.location
    resource_group_name = azurerm_resource_group.resource_group.name
    kind                = "Linux"
    reserved            = true

    sku {
        tier = "Basic"
        size = "B1"
    }
}

resource "azurerm_linux_web_app" "app_service" {
    name                = "pdappservice"
    location            = azurerm_resource_group.resource_group.location
    resource_group_name = azurerm_resource_group.resource_group.name
    service_plan_id = azurerm_app_service_plan.app_service_plan.id
    
     app_settings = {
        "GITHUB_CLIENT_SECRET" = var.clientSecret
     }

      auth_settings_v2 {
        auth_enabled = true
        default_provider = "github"

        github_v2 {
            client_id     = var.clientId
            client_secret_setting_name = "GITHUB_CLIENT_SECRET"
        }

        login {
            
        }
    }

    site_config {
        application_stack {
            docker_registry_url = "https://${azurerm_container_registry.container_registry.name}.azurecr.io"
            docker_image_name      = "${var.containerName}:latest"
        }
        # linux_fx_version = "DOCKER|${azurerm_container_registry.container_registry.name}.azurecr.io/${var.containerName}:latest"
    }
}

