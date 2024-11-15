variable containerName {
  type        = string
  default     = "prompt-defender"
  description = "description"
}

variable resourceGroupName {
  type        = string
  default     = "pdresourcegroup"
  description = "description"
}

variable location {
  type        = string
  default     = "West Europe"
  description = "description"
}

variable containerRegistryName {
  type        = string
  default     = "promptdefenderregistry"
  description = "description"
}

variable subscriptionId {
  type        = string
  description = "description"
}

// Github OAuth App

variable clientId {
  type        = string
  description = "description"
}

variable clientSecret {
  type        = string
  description = "description"
}