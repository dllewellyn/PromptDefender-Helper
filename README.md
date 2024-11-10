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