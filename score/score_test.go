package score

import (
	"errors"
	"testing"
)

func TestParseJSON(t *testing.T) {
	validJSON := `{
        "explanation": "This is a test explanation.",
        "defenses": {
            "in_context_defense": 0.2,
            "system_mode_self_reminder": 0.1,
            "sandwich_defense": 0.0,
            "xml_encapsulation": 0.2,
            "random_sequence_enclosure": 0.1
        }
    }`
	invalidJSON := `{
        "explanation": "This is a test explanation.",
        "defenses": {
            "in_context_defense": 0.2,
            "system_mode_self_reminder": 0.1,
            "sandwich_defense": 0.0,
            "xml_encapsulation": 0.2,
            "random_sequence_enclosure": 0.1
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
            "explanation": "This is a test explanation.",
            "defenses": {
                "in_context_defense": 0.2,
                "system_mode_self_reminder": 0.1,
                "sandwich_defense": 0.0,
                "xml_encapsulation": 0.2,
                "random_sequence_enclosure": 0.1
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
