# PromptDefender-Keep

## LLM security

Want to find out more about securing LLM applicaitons?

Check out  [The LLM security framework](https://llmsecurity.safetorun.com/)

## Overview

PromptDefender-Helper is a project designed to handle scoring and improving the security of prompts using GenAI.

We use genkit to provide two key functionalities to a user:

* Prompt scoring - to determine the security of a prompt and its resistance to attacks (particularly prompt injection
  attacks)
* Prompt improvement - to automatically apply improvements to a prompt to make it more secure

## Quickstart

1. Clone the repository:
    ```sh
    git clone https://github.com/dllewellyn/PromptDefender-Helper.git
    cd PromptDefender-Helper
    ```
2. Install dependencies:
    ```sh
    go mod tidy
    ```
3. Set up environment variables:
    ```sh
    export PORT=8080
    export TEST_MODE=false
   make run 
    ```
   
## Project Structure

- `main.go`: The entry point of the application. It sets up the HTTP server, initializes dependencies, and starts the
  application.
- `score/`: Contains the logic for scoring prompts.
- `improve/`: Contains the logic for improving prompts.
- `cache/`: Provides caching mechanisms to store and retrieve data efficiently.
- `endpoints/`: Defines the HTTP endpoints for scoring and improving prompts.

## Dependencies

- `gin`: A web framework for building HTTP servers.
- `fx`: A dependency injection framework for Go.
- `genkit`: A toolkit for building AI-powered applications.
- `dotprompt`: A library for defining and generating prompts.

For the frontend, we use:

- `hotwired/turbo`: A framework for building modern web applications.
- `hotwired/stimulus`: A JavaScript framework for building web applications.
- `shepherd`: A library for guiding users through a series of steps.

### Content-Types 

The API accepts and returns JSON data if the content-type is set to `application/json`. Otherwise, it returns HTMl wrapped in 
hotwired turbo.

## Usage

1. Run the application:
    ```sh
    go run main.go
    ```
2. Access the application at [http://localhost:8080](http://localhost:8080).

## Key Components

### Scoring Prompts

The scoring functionality is implemented in the `score` package. It uses an `LlmScorer` to evaluate prompts and return a
score.

### Improving Prompts

The improvement functionality is implemented in the `improve` package. It uses an `LlmImprover` to suggest improvements
for prompts.

### Caching

The `cache` package provides an in-memory caching mechanism to store and retrieve data efficiently.

### Endpoints

The `endpoints` package defines the HTTP endpoints for scoring and improving prompts. These endpoints are registered in
the `main.go` file.

## Example

To score a prompt, send a POST request to `/score` with the prompt text. To improve a prompt, send a POST request to
`/improve` with the prompt text.

# Testing

## Overview

To run the unit tests, run

```sh 
make test 
```

To run the integration tests, run

```sh
make integration-test
```

### Manual API runs

In vscode, you can use the `REST Client` extension to run the API calls. These are in the folder `/test/http`

### Testing prompts

To test the prompt separately, you can run the app using genkit.

First ensure you have genkit installed, you can follow the instructions here:

[https://firebase.google.com/docs/genkit-go/get-started-go](https://firebase.google.com/docs/genkit-go/get-started-go)

```sh
make genkit_mode
```

### Running Integration Tests for the `/score` Endpoint

To run the integration tests for the `/score` endpoint, follow these steps:

1. Ensure the application is running:
    ```sh
    make run
    ```

2. Run the integration tests:
    ```sh
    make integration-test
    ```

The integration tests will load prompt inputs from files in the `tests/prompts/` directory and verify that the defense matches the expected score for each prompt.

## Contributing

1. Fork the repository.
2. Create a new branch:
    ```sh
    git checkout -b feature-branch
    ```
3. Make your changes.
4. Commit your changes:
    ```sh
    git commit -am 'Add new feature'
    ```
5. Push to the branch:
    ```sh
    git push origin feature-branch
    ```
6. Create a new Pull Request.


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