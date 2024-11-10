package score

import (
	"PromptDefender-Keep/logger"
	"encoding/json"
	"math"
	"strings"

	"go.uber.org/zap"
)

type Defenses struct {
	InContextDefense        float64 `json:"in_context_defense"`
	SystemModeSelfReminder  float64 `json:"system_mode_self_reminder"`
	SandwichDefense         float64 `json:"sandwich_defense"`
	XMLEncapsulation        float64 `json:"xml_encapsulation"`
	RandomSequenceEnclosure float64 `json:"random_sequence_enclosure"`
}

type PromptScore struct {
	Score       *float64 `json:"score"`
	Explanation string   `json:"explanation"`
	Defenses    Defenses `json:"defenses"`
}

// Function to parse the JSON string into the struct
func parseJSON(input string) (PromptScore, error) {
	cleaned := strings.TrimPrefix(input, "```json")
	cleaned = strings.ReplaceAll(cleaned, "\n", "")
	cleaned = strings.TrimSuffix(cleaned, "```")

	// Unmarshal the JSON into the struct
	var result PromptScore
	err := json.Unmarshal([]byte(cleaned), &result)
	if err != nil {
		logger.Log.Error("Error unmarshalling JSON", zap.Error(err))
		return result, err
	}

	// Set score to total of all defences
	score := 0.0
	score += result.Defenses.InContextDefense
	score += result.Defenses.SystemModeSelfReminder
	score += result.Defenses.SandwichDefense
	score += result.Defenses.XMLEncapsulation
	score += result.Defenses.RandomSequenceEnclosure

	// Ensure the score has only 2 decimal points at most
	score = math.Round(score*100) / 100
	result.Score = &score

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
