//go:build integration
// +build integration

// tests/integration/score_integration_test.go
package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type PromptTest struct {
	Prompt           string  `json:"prompt"`
	ExpectedScore    float64 `json:"expected_score"`
	ExpectedDefenses struct {
		InContextDefense        float64 `json:"in_context_defense"`
		SystemModeSelfReminder  float64 `json:"system_mode_self_reminder"`
		SandwichDefense         float64 `json:"sandwich_defense"`
		XMLEncapsulation        float64 `json:"xml_encapsulation"`
		RandomSequenceEnclosure float64 `json:"random_sequence_enclosure"`
	} `json:"expected_defenses"`
}

func TestScoreEndpoint(t *testing.T) {
	serverURL := os.Getenv("BASE_URL")
	if serverURL == "" {
		t.Fatal("BASE_URL environment variable is not set")
	}

	promptFiles := []string{
		"../prompts/valid_prompt.json",
		"../prompts/invalid_prompt.json",
	}

	for _, file := range promptFiles {
		t.Run(file, func(t *testing.T) {
			promptTest, err := loadPromptTest(file)
			if err != nil {
				t.Fatalf("Failed to load prompt test: %v", err)
			}

			response, err := sendScoreRequest(serverURL, promptTest.Prompt)
			if err != nil {
				t.Fatalf("Failed to send score request: %v", err)
			}

			var result map[string]interface{}
			err = json.Unmarshal(response, &result)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			score := result["response"].(map[string]interface{})["score"].(float64)
			assert.Equal(t, promptTest.ExpectedScore, score, "Expected score does not match")

			defenses := result["response"].(map[string]interface{})["defenses"].(map[string]interface{})
			assert.Equal(t, promptTest.ExpectedDefenses.InContextDefense, defenses["in_context_defense"].(float64), "InContextDefense does not match")
			assert.Equal(t, promptTest.ExpectedDefenses.SystemModeSelfReminder, defenses["system_mode_self_reminder"].(float64), "SystemModeSelfReminder does not match")
			assert.Equal(t, promptTest.ExpectedDefenses.SandwichDefense, defenses["sandwich_defense"].(float64), "SandwichDefense does not match")
			assert.Equal(t, promptTest.ExpectedDefenses.XMLEncapsulation, defenses["xml_encapsulation"].(float64), "XMLEncapsulation does not match")
			assert.Equal(t, promptTest.ExpectedDefenses.RandomSequenceEnclosure, defenses["random_sequence_enclosure"].(float64), "RandomSequenceEnclosure does not match")
		})
	}
}

func loadPromptTest(file string) (*PromptTest, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var promptTest PromptTest
	err = json.Unmarshal(data, &promptTest)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &promptTest, nil
}

func sendScoreRequest(serverURL, prompt string) ([]byte, error) {
	requestBody := map[string]string{"prompt": prompt}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", serverURL+"/score", ioutil.NopCloser(bytes.NewReader(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return responseData, nil
}
