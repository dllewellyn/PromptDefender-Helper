package score

import (
	"encoding/json"
	"os"
	"strings"
)

type Defenses struct {
	InContextDefense        bool `json:"in_context_defense"`
	SystemModeSelfReminder  bool `json:"system_mode_self_reminder"`
	SandwichDefense         bool `json:"sandwich_defense"`
	XMLEncapsulation        bool `json:"xml_encapsulation"`
	RandomSequenceEnclosure bool `json:"random_sequence_enclosure"`
}

type PromptScore struct {
	Score       float64  `json:"score"`
	Explanation string   `json:"explanation"`
	Defenses    Defenses `json:"defenses"`
}

// Function to parse the JSON string into the struct
func parseJSON(input string) (PromptScore, error) {

	// Write the string to file
	err := os.WriteFile("response.json", []byte(input), 0644)

	if err != nil {
		return PromptScore{}, err
	}

	cleaned := strings.TrimPrefix(input, "```json")
	cleaned = strings.ReplaceAll(cleaned, "\n", "")
	cleaned = strings.TrimSuffix(cleaned, "```")

	// Unmarshal the JSON into the struct
	var result PromptScore
	err = json.Unmarshal([]byte(cleaned), &result)
	return result, err
}

type InvokeRequest = func(prompt string) (string, error)

type LlmScoringPromptInput struct {
	StartingPrompt string `json:"startingPrompt"`
}

type LlmScorer struct {
	InvokeRequest InvokeRequest
}

type Scorer interface {
	Score(prompt string) (*PromptScore, error)
}

func NewLlmScorer(InvokeRequest InvokeRequest) Scorer {

	return &LlmScorer{
		InvokeRequest: InvokeRequest,
	}
}

func (s *LlmScorer) Score(prompt string) (*PromptScore, error) {
	response, err := s.InvokeRequest(prompt)

	if err != nil {
		return nil, err
	}

	score, err := parseJSON(response)

	if err != nil {
		return nil, err
	}

	return &score, nil
}
