terraform {
  backend "azurerm" {
    container_name        = "tfstate"
    key             = "terraform.tfstate"
  }
}

provider "azurerm" {
  features {}
  subscription_id = var.subscriptionId
}

data "azurerm_resource_group" "resource_group" {
  name     = var.resourceGroupName
}

resource "azurerm_container_registry" "container_registry" {
  name                = var.containerRegistryName
  resource_group_name = data.azurerm_resource_group.resource_group.name
  location            = data.azurerm_resource_group.resource_group.location
  sku                 = "Basic"
  admin_enabled       = true
}

resource "azurerm_service_plan" "app_service_plan" {
  name                = "pdappserviceplan"
  location            = data.azurerm_resource_group.resource_group.location
  resource_group_name = data.azurerm_resource_group.resource_group.name
  os_type             = "Linux"
  sku_name            = "B1"
}

resource "azurerm_application_insights" "app_insights" {
  name                = "pdappinsights"
  location            = data.azurerm_resource_group.resource_group.location
  resource_group_name = data.azurerm_resource_group.resource_group.name
  application_type    = "web"
}

resource "azurerm_linux_web_app" "app_service" {
  name                = "pdappservice"
  location            = data.azurerm_resource_group.resource_group.location
  resource_group_name = data.azurerm_resource_group.resource_group.name
  service_plan_id     = azurerm_service_plan.app_service_plan.id


  app_settings  ={
    WEBSITE_ENABLE_SYNC_UPDATE_SITE = "true"
  }


  identity {
    type = "SystemAssigned"
  }

  auth_settings_v2 {
    auth_enabled     = true

    default_provider = "github"
    require_authentication = false
    unauthenticated_action = "AllowAnonymous"

    github_v2 {
      client_id                  = var.clientId
      client_secret_setting_name = "GITHUB_CLIENT_SECRET"
    }

    login {

    }
  }

  logs {
    http_logs {
      file_system {
        retention_in_mb = 100
        retention_in_days = 7
      }
    }
  }

  site_config {
    application_stack {
      docker_registry_url      = "https://${azurerm_container_registry.container_registry.name}.azurecr.io"
      docker_image_name        = "${var.containerName}:latest"
      docker_registry_username = azurerm_container_registry.container_registry.admin_username
      docker_registry_password = azurerm_container_registry.container_registry.admin_password
    }
  }
}

# This is a null resource that will run the az cli command to set the app settings
# It's needed because we are not using slots, and because of a dependency chain that 
# is caused by needing the key vault to be created before the app service 
resource "null_resource" "configure_app_service" {
  depends_on = [
    azurerm_linux_web_app.app_service,
    azurerm_key_vault_secret.gcloud_location,
    azurerm_key_vault_secret.gcloud_project,
    azurerm_key_vault_secret.service_account_key,
    azurerm_key_vault_secret.github_client_secret
  ]

    triggers = {
        always = timestamp()
    }

  provisioner "local-exec" {
    command = <<EOT
    az webapp config appsettings set --name pdappservice \
      --resource-group ${data.azurerm_resource_group.resource_group.name} \
      --settings GITHUB_CLIENT_SECRET="@Microsoft.KeyVault(SecretUri=${azurerm_key_vault_secret.github_client_secret.id})" \
                GCLOUD_LOCATION="@Microsoft.KeyVault(SecretUri=${azurerm_key_vault_secret.gcloud_location.id})" \
                GCLOUD_PROJECT="@Microsoft.KeyVault(SecretUri=${azurerm_key_vault_secret.gcloud_project.id})" \
                SERVICE_ACCOUNT_KEY="@Microsoft.KeyVault(SecretUri=${azurerm_key_vault_secret.service_account_key.id})" 
    EOT
  }
}

resource "azurerm_key_vault" "prompt_defender_kv" { #tfsec:ignore:azure-keyvault-no-purge
  name                = "promptdefender-keyvault"
  location            = data.azurerm_resource_group.resource_group.location
  resource_group_name = data.azurerm_resource_group.resource_group.name
  tenant_id           = azurerm_linux_web_app.app_service.identity[0].tenant_id
  sku_name            = "standard"

  network_acls {
    default_action = "Deny"
    bypass         = "AzureServices"

    ip_rules = [
      "0.0.0.0/0"
    ]
  }

  access_policy {
    tenant_id = azurerm_linux_web_app.app_service.identity[0].tenant_id
    object_id = azurerm_linux_web_app.app_service.identity[0].principal_id

    secret_permissions = [
      "Get",
    ]
  }

    access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Get",
      "List",
      "Set",
    ]
  }

}