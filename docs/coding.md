# Development Guidelines and NFRs

## Key Development Considerations

1. **Code Structure**
   - Write Go code using clean and modular principles.
   - Maintain a separate package for reusable components.
   - Aim for simplicity: avoid over-complicating designs.

2. **Frontend Development**
   - Keep the frontend in the Hugo web app and try wherever possible to keep UI inside the ui folder.
   - Ensure seamless integration with the backend APIs using the hotwired libraries.
   - Use hotwired stimulus for any logic on the frontend, rather than using javascript
   - Avoid overloading the frontend with business logic; keep it lightweight.

3. **Deployment**
   - Use Azure App Services for deployment.
   - Automate all deployments through CI/CD pipelines.
   - Ensure infrastructure as code (IaC) configurations are up-to-date and version-controlled.

4. **Unit Testing**
   - Write unit tests for all appropriate classes.
   - Include tests for edge cases and negative scenarios.
   - Organise BDD test cases in the `tests` folder.

5. **Logging**
   - Use the `zap` library for structured logging.
   - Always initialise loggers using `logging.GetLogger`.
   - **When to log:**
     - Log at the start and end of major processes.
     - Log important decisions or branches in business logic.
     - Log all error and exception cases.
     - Use `INFO` level for routine processes and `DEBUG` for detailed insights during development.
     - Use `ERROR` or `FATAL` for critical errors requiring immediate attention.
   - **How to log:**
     - Use structured logging formats to include relevant metadata (e.g., trace IDs, timestamps, and user IDs).
     - Avoid logging sensitive data.
     - Consistently format log messages for readability.

6. **Error Handling**
   - Use Go's error interface consistently.
   - Return meaningful error messages and stack traces where applicable.
   - Avoid panics; recover gracefully from unexpected situations.

7. **Security**
   - Implement secure coding practices, especially for input sanitisation and validation.
   - Use Azure OpenAI Content Safety features to moderate and shield content where applicable.
   - Monitor for vulnerabilities and patch dependencies regularly.

8. **Keep it Simple**
   - Avoid overengineering.
   - Focus on delivering a maintainable, performant, and robust application.

---

## Non-Functional Requirements (NFRs)

1. **Scalability**
   - Support scaling up/down with Azure App Services.
   - Design the system to handle future increases in load without significant changes.

2. **Reliability**
   - Achieve 99.9% uptime (excluding planned maintenance).
   - Implement retry mechanisms for transient errors.

3. **Maintainability**
   - Follow clear coding conventions (e.g., `golint` or `gofmt`). You can use `make format` to run these.
   - Maintain documentation in Markdown format in the `/docs` folder
   - Ensure each module has clear ownership and responsibility and follows the best practices for go packages
      of one function per package
   - Use clean code standards 

4. **Observability**
   - Implement structured logging and monitoring for key application metrics.
   - Correlate logs with trace IDs to improve debugging efficiency.
   - Set up alerting for critical events (e.g., service downtime, high latency).

5. **Security**
   - Use HTTPS for all communication.
   - Store secrets securely using Azure Key Vault.
   - Regularly perform vulnerability scans and penetration tests.

6. **Portability**
   - Ensure the application can be deployed to other environments with minimal changes.
   - Document all dependencies clearly.

---

## Testing Strategy

1. **Unit Tests**
   - Create unit tests for all appropriate classes and business logic.
   - Include test cases for:
     - Expected behaviours.
     - Edge cases.
     - Negative scenarios.
   - Aim for 80%+ code coverage.

2. **Integration Tests**
   - Focus on LLM-based code integrations.
   - Ensure tests cover data handling, API interaction, and prompt handling workflows.
   - Validate the end-to-end flow of interactions between integrated services.

3. **API Tests**
   - Test API endpoints for:
     - Correct responses (200s, 400s, 500s, etc.).
     -


### Deployment scripts 

- The application contains terraform configuration which deploys to azure in the azure folder 
- The application uses makefiles to make application deployment, build and testing simplier
- The application uses Docker which is contained in the Dockerfile - which should be used for both local testing and deployment.