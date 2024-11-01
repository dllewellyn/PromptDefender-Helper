# PromptDefender-Keep

## Overview
PromptDefender-Keep is a project designed to handle scoring and improving prompts using machine learning models. It leverages the Go programming language and integrates with various libraries and frameworks such as `gin`, `fx`, and `genkit` to provide a robust and scalable solution.

## Project Structure
- `main.go`: The entry point of the application. It sets up the HTTP server, initializes dependencies, and starts the application.
- `score/`: Contains the logic for scoring prompts.
- `improve/`: Contains the logic for improving prompts.
- `cache/`: Provides caching mechanisms to store and retrieve data efficiently.
- `endpoints/`: Defines the HTTP endpoints for scoring and improving prompts.

## Dependencies
- `gin`: A web framework for building HTTP servers.
- `fx`: A dependency injection framework for Go.
- `genkit`: A toolkit for building AI-powered applications.
- `dotprompt`: A library for defining and generating prompts.

## Installation
1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/PromptDefender-Keep.git
    cd PromptDefender-Keep
    ```
2. Install dependencies:
    ```sh
    go mod tidy
    ```
3. Set up environment variables:
    ```sh
    export PORT=8080
    export TEST_MODE=false
    ```

## Usage
1. Run the application:
    ```sh
    go run main.go
    ```
2. Access the application at [http://localhost:8080](http://localhost:8080).

## Key Components
### Scoring Prompts
The scoring functionality is implemented in the `score` package. It uses an `LlmScorer` to evaluate prompts and return a score.

### Improving Prompts
The improvement functionality is implemented in the `improve` package. It uses an `LlmImprover` to suggest improvements for prompts.

### Caching
The `cache` package provides an in-memory caching mechanism to store and retrieve data efficiently.

### Endpoints
The `endpoints` package defines the HTTP endpoints for scoring and improving prompts. These endpoints are registered in the `main.go` file.

## Example
To score a prompt, send a POST request to `/score` with the prompt text. To improve a prompt, send a POST request to `/improve` with the prompt text.

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

## License
This project is licensed under the MIT License. See the `LICENSE` file for details.