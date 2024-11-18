# Deployment 

## Setting up a Resource Group in Azure

You can set up a resource group in Azure using either the Azure Portal or the Azure CLI.

### Using the Azure Portal

1. Go to the [Azure Portal](https://portal.azure.com/).
2. In the left-hand menu, select **Resource groups**.
3. Click on **+ Create**.
4. Fill in the required details:
    - **Subscription**: Select your Azure subscription.
    - **Resource group**: Enter a unique name for the resource group.
    - **Region**: Select the region where you want the resource group to be located.
5. Click **Review + create** and then **Create**.

### Using the Azure CLI

1. Open your terminal or command prompt.
2. Use the following command to create a resource group:
    ```sh
    az group create --name <ResourceGroupName> --location <Location>
    ```
    Replace `<ResourceGroupName>` with your desired resource group name and `<Location>` with the Azure region (e.g., `eastus`).

Example:
```sh
az group create --name MyResourceGroup --location eastus
```

## One of setup 

To use the `one-off-setup.sh` script, follow these steps:

Run the Script: Execute the script by providing the resource group name created in the previous step. For example, if your resource group name is MyResourceGroup, you would run the script as follows:

## Script Explanation

Variable Initialization: It initializes variables for the resource group name, a randomly generated storage account name, and a container name.

Create Storage Account: It creates a storage account in the specified resource group and location (eastus).

Retrieve Storage Account Key: It retrieves the storage account key.

Create Blob Container: It creates a blob container in the storage account.

Output Configuration: Finally, it outputs the storage account name, container name, and access key.

Example Output
If you run the script with MyResourceGroup as the resource group name, the output will look something like this:

```sh
storage_account_name: pdtfstate12345
container_name: tfstate
access_key: <your_storage_account_key>
```

For the next step, set the following environment variables:

```sh
export ARM_ACCESS_KEY=<key from previous step>
export ARM_STORAGE_ACCOUNT_NAME=<storage_account_name>
```

### Required Environment Variables

To execute the `deploy.sh` script, you need to set the following environment variables:

- `ARM_STORAGE_ACCOUNT_NAME`: The name of the Azure storage account used for Terraform state.
- `ARM_ACCESS_KEY`: The access key for the Azure storage account.

Example:
```sh
export ARM_STORAGE_ACCOUNT_NAME=<your_storage_account_name>
export ARM_ACCESS_KEY=<your_storage_account_key>
```

### Terraform Variables

The following `TF_VAR` environment variables need to be set for the Terraform configuration in `variables.tf`:

- `TF_VAR_subscriptionId`: Your Azure subscription ID.
- `TF_VAR_clientId`: The client ID of your GitHub OAuth app.
- `TF_VAR_clientSecret`: The client secret of your GitHub OAuth app.
- `TF_VAR_gcloudLocation`: The location for your Google Cloud resources.
- `TF_VAR_gcloudProject`: Your Google Cloud project ID.
- `TF_VAR_serviceAccountKey`: The service account key for Google Cloud.

Example:
```sh
export TF_VAR_subscriptionId=<your_subscription_id>
export TF_VAR_clientId=<your_client_id>
export TF_VAR_clientSecret=<your_client_secret>
export TF_VAR_gcloudLocation=<your_gcloud_location>
export TF_VAR_gcloudProject=<your_gcloud_project>
export TF_VAR_serviceAccountKey=<your_service_account_key>
```


### Destruction 

```sh
az keyvault set-policy --name promptdefender-keyvault --object-id 06717f30-12b0-4e52-9960-5d5e757e51c9 --secret-permissions get list set delete purge --key-permissions get list --certificate-permissions get list
```