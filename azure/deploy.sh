#!/bin/bash

# Variables

ACR_NAME="promptdefender01"
IMAGE_NAME="prompt-defender"
CONTAINER_NAME=$1
RESOURCE_GROUP=$2

# Check if container name is provided
if [ -z "$CONTAINER_NAME" ] || [ -z "$RESOURCE_GROUP" ]; then
	echo "Usage: $0 <container-name> <resource-group>"
	exit 1
fi

# Login to Azure
echo "Logging in to Azure..."

# Deploy Terraform
echo "Deploying Terraform..."
terraform init -backend-config="storage_account_name=$ARM_STORAGE_ACCOUNT_NAME"
terraform apply -auto-approve -var "resourceGroupName=$RESOURCE_GROUP" -var "containerRegistryName=$ACR_NAME" || exit 1

# Login to Azure Container Registry
echo "Logging in to Azure Container Registry..."
az acr login --name $ACR_NAME

# Tag the Docker image
echo "Tagging the Docker image... as $ACR_NAME.azurecr.io/$IMAGE_NAME:latest"
docker tag $CONTAINER_NAME $ACR_NAME.azurecr.io/$IMAGE_NAME:latest || exit 1

# Push the Docker image to the Azure Container Registry
echo "Pushing the Docker image to the Azure Container Registry... "
docker push $ACR_NAME.azurecr.io/$IMAGE_NAME:latest || exit 1

echo "Deployment completed successfully."