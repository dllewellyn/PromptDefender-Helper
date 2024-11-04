Feature: Improve endpoint
  As a user
  I want to send a prompt to the /improve endpoint
  So that I can receive an improved version of the prompt

  Scenario: Successfully improve a prompt with JSON request
    Given I have a JSON request with the prompt "This is a test prompt"
    When I send a POST request to "/improve"
    Then the response status should be 200
    And the response content type should be "application/json; charset=utf-8"
    And the response should contain an improved prompt

  Scenario: Successfully improve a prompt with form data
    Given I have a form request with the prompt "This is a test prompt"
    When I send a POST request to "/improve"
    Then the response status should be 200
    And the response content type should be "text/html; charset=utf-8"
    And the response should contain an improved prompt

  Scenario: Handle invalid JSON request
    Given I have an invalid JSON request
    When I send a POST request to "/improve"
    Then the response status should be 400

  Scenario: Handle missing form data
    Given I have a form request with missing prompt data
    When I send a POST request to "/improve"
    Then the response status should be 400