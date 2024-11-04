package score

import (
	"errors"
	"testing"
)

func TestParseJSON(t *testing.T) {
	validJSON := `{
		"score": 0.85,
		"explanation": "This is a test explanation.",
		"defenses": {
			"in_context_defense": true,
			"system_mode_self_reminder": false,
			"sandwich_defense": true,
			"xml_encapsulation": false,
			"random_sequence_enclosure": true
		}
	}`
	invalidJSON := `{
		"score": 0.85,
		"explanation": "This is a test explanation.",
		"defenses": {
			"in_context_defense": true,
			"system_mode_self_reminder": false,
			"sandwich_defense": true,
			"xml_encapsulation": false,
			"random_sequence_enclosure": true
		`

	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{"Valid JSON", validJSON, false},
		{"Invalid JSON", invalidJSON, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseJSON(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("parseJSON() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestLlmScorer_Score(t *testing.T) {
	mockInvokeRequest := func(prompt string) (string, error) {
		if prompt == "error" {
			return "", errors.New("mock error")
		}
		return `{
			"score": 0.85,
			"explanation": "This is a test explanation.",
			"defenses": {
				"in_context_defense": true,
				"system_mode_self_reminder": false,
				"sandwich_defense": true,
				"xml_encapsulation": false,
				"random_sequence_enclosure": true
			}
		}`, nil
	}

	scorer := NewLlmScorer(mockInvokeRequest)

	tests := []struct {
		name      string
		prompt    string
		expectErr bool
	}{
		{"Valid prompt", "test prompt", false},
		{"Error prompt", "error", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := scorer.Score(tt.prompt)
			if (err != nil) != tt.expectErr {
				t.Errorf("Score() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}
