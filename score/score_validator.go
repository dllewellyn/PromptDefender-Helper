package score

import (
	"encoding/json"
	"strings"
)

type ValidatorPromptInput struct {
	Score string `json:"model_result"`
}

type DecodedValidatorResult struct {
	Consistent  bool    `json:"consistent"`
	Explanation *string `json:"explanation"`
}
type ValidateScore = func(score string) (string, error)

type Validator struct {
	ValidateScore ValidateScore
}

func NewValidator(ValidateScore ValidateScore) Validator {
	return Validator{
		ValidateScore: ValidateScore,
	}
}

// Function to parse the JSON string into the struct
func parseJSONToValidatorResult(input string) (DecodedValidatorResult, error) {
	cleaned := strings.TrimPrefix(input, "```json")
	cleaned = strings.ReplaceAll(cleaned, "\n", "")
	cleaned = strings.TrimSuffix(cleaned, "```")

	// Unmarshal the JSON into the struct
	var result DecodedValidatorResult
	err := json.Unmarshal([]byte(cleaned), &result)
	return result, err
}

func (v *Validator) Validate(score string) (string, bool, error) {

	result, err := v.ValidateScore(score)

	if err != nil {
		return "", false, err
	}

	if result == "" {
		return "", false, nil
	}

	parsedResult, err := parseJSONToValidatorResult(result)

	if err != nil {
		return "", false, err
	}

	if parsedResult.Consistent {
		return "", parsedResult.Consistent, nil
	}

	return *parsedResult.Explanation, parsedResult.Consistent, nil
}
