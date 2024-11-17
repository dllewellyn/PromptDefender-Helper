resource "azurerm_key_vault_secret" "gcloud_location" {
  name         = "GCLOUD-LOCATION"
  value        = var.gcloudLocation
  key_vault_id = azurerm_key_vault.prompt_defender_kv.id
}

resource "azurerm_key_vault_secret" "gcloud_project" {
  name         = "GCLOUD-PROJECT"
  value        = var.gcloudProject
  key_vault_id = azurerm_key_vault.prompt_defender_kv.id
}

resource "azurerm_key_vault_secret" "service_account_key" {
  name         = "SERVICE-ACCOUNT-KEY"
  value        = var.serviceAccountKey
  key_vault_id = azurerm_key_vault.prompt_defender_kv.id
}

resource "azurerm_key_vault_secret" "github_client_secret" {
  name         = "GITHUB-CLIENT-SECRET"
  value        = var.clientSecret
  key_vault_id = azurerm_key_vault.prompt_defender_kv.id
}