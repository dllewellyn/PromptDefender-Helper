// tests/integration/improve_test.go
//go:build integration
// +build integration

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

var (
	server       string
	ctx          = context.Background()
	response     *http.Response
	responseBody string
)

func TestFeatures(t *testing.T) {
	serverURL := os.Getenv("BASE_URL")
	if serverURL == "" {
		t.Fatal("BASE_URL environment variable is not set")
	}

	server = serverURL

	// Run the Godog tests
	opts := godog.Options{
		Format: "pretty",
		Paths:  []string{"../features/improve.feature"},
	}

	status := godog.TestSuite{
		Name:                "integration",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	if status != 0 {
		t.Fail()
	}
}

func InitializeScenario(sc *godog.ScenarioContext) {
	sc.Step(`^I have a JSON request with the prompt "([^"]*)"$`, iHaveAJSONRequestWithThePrompt)
	sc.Step(`^I have a form request with the prompt "([^"]*)"$`, iHaveAFormRequestWithThePrompt)
	sc.Step(`^I have an invalid JSON request$`, iHaveAnInvalidJSONRequest)
	sc.Step(`^I have a form request with missing prompt data$`, iHaveAFormRequestWithMissingPromptData)
	sc.Step(`^I send a POST request to "([^"]*)"$`, iSendAPOSTRequestTo)
	sc.Step(`^the response status should be (\d+)$`, func(status int) error {
		return theResponseStatusShouldBe(status)
	})
	sc.Step(`^the response content type should be "([^"]*)"$`, theResponseContentTypeShouldBe)
	sc.Step(`^the response should contain an improved prompt$`, theResponseShouldContainAnImprovedPrompt)
	sc.Step(`^the response should contain "([^"]*)"$`, theResponseShouldContain)
	sc.Step(`^the response should redirect to "([^"]*)"$`, theResponseShouldRedirectTo)
}

// Step definitions (implement these functions)
func iHaveAJSONRequestWithThePrompt(prompt string) error {
	payload := map[string]string{"prompt": prompt}
	jsonData, _ := json.Marshal(payload)
	ctx = context.WithValue(ctx, "requestBody", string(jsonData))
	ctx = context.WithValue(ctx, "contentType", "application/json")
	return nil
}

func iHaveAFormRequestWithThePrompt(prompt string) error {
	formData := "prompt=" + prompt
	ctx = context.WithValue(ctx, "requestBody", formData)
	ctx = context.WithValue(ctx, "contentType", "application/x-www-form-urlencoded")
	return nil
}

func iHaveAnInvalidJSONRequest() error {
	ctx = context.WithValue(ctx, "requestBody", "{invalid_json}")
	ctx = context.WithValue(ctx, "contentType", "application/json")
	return nil
}

func iHaveAFormRequestWithMissingPromptData() error {
	formData := "invalid=data"
	ctx = context.WithValue(ctx, "requestBody", formData)
	ctx = context.WithValue(ctx, "contentType", "application/x-www-form-urlencoded")
	return nil
}

func iSendAPOSTRequestTo(endpoint string) error {
	requestBody := ctx.Value("requestBody").(string)
	contentType := ctx.Value("contentType").(string)

	req, err := http.NewRequest("POST", server+endpoint, strings.NewReader(requestBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", contentType)

	client := &http.Client{}
	response, err = client.Do(req)
	if err != nil {
		return err
	}

	bodyBytes, _ := io.ReadAll(response.Body)
	responseBody = string(bodyBytes)
	return nil
}

func theResponseStatusShouldBe(expectedStatus int) error {
	if response.StatusCode != expectedStatus {
		return fmt.Errorf("expected status %d but got %d", expectedStatus, response.StatusCode)
	}
	return nil
}

func theResponseContentTypeShouldBe(expectedContentType string) error {
	contentType := response.Header.Get("Content-Type")
	if contentType != expectedContentType {
		return fmt.Errorf("expected content type %s but got %s", expectedContentType, contentType)
	}
	return nil
}

func theResponseShouldContainAnImprovedPrompt() error {
	if responseBody == "" {
		return fmt.Errorf("response body is empty")
	}

	contentType := response.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
		var jsonData map[string]string
		err := json.Unmarshal([]byte(responseBody), &jsonData)
		if err != nil {
			return err
		}

		if jsonData["response"] == "" {
			return fmt.Errorf("response does not contain expected content (%s)", responseBody)
		}
	} else {
		if !strings.Contains(responseBody, "<turbo-frame") {
			return fmt.Errorf("response does not contain expected content (%s)", responseBody)
		}
	}

	return nil

}

func theResponseShouldContain(expectedContent string) error {
	if !strings.Contains(responseBody, expectedContent) {
		return fmt.Errorf("response does not contain expected content")
	}
	return nil
}

func theResponseShouldRedirectTo(redirectURL string) error {
	location, err := response.Location()
	if err != nil {
		return fmt.Errorf("expected redirect but got error: %v", err)
	}
	if location.Path != redirectURL {
		return fmt.Errorf("expected redirect to %s but got %s", redirectURL, location.Path)
	}
	return nil
}
