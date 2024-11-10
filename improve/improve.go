package improve

import (
	"PromptDefender-Keep/logger"
	"golang.org/x/exp/rand"
)

type Improver interface {
	Improve(input string) (string, error)
}

type LlmImprover struct {
	ImproveFunc func(input string, randomSequence string) (string, error)
}

type LlmPromptImproverInput struct {
	StartingPrompt string `json:"starting_prompt"`
	RandomSequence string `json:"random_sequence"`
}

func NewLlmImprover(improveFunc func(input string, randomSequence string) (string, error)) *LlmImprover {
	if improveFunc == nil {
		logger.Log.Error("Improve func is nil")
		return nil
	}
	
	return &LlmImprover{ImproveFunc: improveFunc}
}

func (l *LlmImprover) Improve(input string) (string, error) {
	return l.ImproveFunc(input, RandomString(6))
}

// Generate a random string of alphanumeric characters of length n
func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
