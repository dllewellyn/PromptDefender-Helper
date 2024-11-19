#!/bin/bash

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <resource_group_name>"
    exit 1
fi

# $1 is the name of the resource group
RESOURCE_GROUP_NAME=$1
STORAGE_ACCOUNT_NAME="pdtfstate$RANDOM" 
CONTAINER_NAME="tfstate"

# Create storage account
az storage account create --name $STORAGE_ACCOUNT_NAME --resource-group $RESOURCE_GROUP_NAME --location eastus --sku Standard_LRS || exit 1

# Get storage account key
ACCOUNT_KEY=$(az storage account keys list --resource-group $RESOURCE_GROUP_NAME --account-name $STORAGE_ACCOUNT_NAME --query '[0].value' --output tsv) || exit 1

# Create blob container
az storage container create --name $CONTAINER_NAME --account-name $STORAGE_ACCOUNT_NAME --account-key $ACCOUNT_KEY || exit 1

# Output the backend configuration
echo "storage_account_name: $STORAGE_ACCOUNT_NAME"
echo "container_name: $CONTAINER_NAME"
echo "access_key: $ACCOUNT_KEY"