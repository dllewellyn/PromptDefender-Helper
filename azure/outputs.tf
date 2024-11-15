output "azure_container_registry" {
  value       = azurerm_container_registry.container_registry.name
    description = "The name of the Azure Container Registry"
}
