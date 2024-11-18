## Deployment

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