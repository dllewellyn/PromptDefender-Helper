resource "random_string" "suffix" {
  length  = 8
  special = false
}

resource "azurerm_key_vault" "prompt_defender_kv" { #tfsec:ignore:azure-keyvault-no-purge
  name                = "promptdefender-keyvault-${random_string.suffix.result}"
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
resource "azurerm_key_vault_secret" "gcloud_location" {
  name           = "GCLOUD-LOCATION"
  value          = var.gcloudLocation
  key_vault_id   = azurerm_key_vault.prompt_defender_kv.id
  content_type   = "text/plain"
  expiration_date = timeadd(timestamp(), "8760h") // 1 year from now
}

resource "azurerm_key_vault_secret" "gcloud_project" {
  name           = "GCLOUD-PROJECT"
  value          = var.gcloudProject
  key_vault_id   = azurerm_key_vault.prompt_defender_kv.id
  content_type   = "text/plain"
  expiration_date = timeadd(timestamp(), "8760h") // 1 year from now
}

resource "azurerm_key_vault_secret" "service_account_key" {
  name           = "SERVICE-ACCOUNT-KEY"
  value          = var.serviceAccountKey
  key_vault_id   = azurerm_key_vault.prompt_defender_kv.id
  content_type   = "application/json"
  expiration_date = timeadd(timestamp(), "8760h") // 1 year from now
}

resource "azurerm_key_vault_secret" "github_client_secret" {
  name           = "GITHUB-CLIENT-SECRET"
  value          = var.clientSecret
  key_vault_id   = azurerm_key_vault.prompt_defender_kv.id
  content_type   = "text/plain"
  expiration_date = timeadd(timestamp(), "8760h") // 1 year from now
}