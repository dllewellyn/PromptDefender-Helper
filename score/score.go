package score

import "strconv"

type InvokeRequest = func(prompt string) (string, error)

type LlmScoringPromptInput struct {
	StartingPrompt string `json:"startingPrompt"`
}
type LlmScorer struct {
	InvokeRequest InvokeRequest
}

type Scorer interface {
	Score(prompt string) (int, error)
}

func NewLlmScorer(InvokeRequest InvokeRequest) Scorer {

	return &LlmScorer{
		InvokeRequest: InvokeRequest,
	}
}

func (s *LlmScorer) Score(prompt string) (int, error) {
	response, err := s.InvokeRequest(prompt)

	if err != nil {
		return 0, err
	}

	return strconv.Atoi(response)
}
