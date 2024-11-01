package improve

type Improver interface {
	Improve(input string) (string, error)
}

type LlmImprover struct {
	ImproveFunc func(input string) (string, error)
}

type LlmPromptImproverInput struct {
	StartingPrompt string `json:"starting_prompt"`
}

func NewLlmImprover(improveFunc func(input string) (string, error)) *LlmImprover {
	return &LlmImprover{ImproveFunc: improveFunc}
}

func (l *LlmImprover) Improve(input string) (string, error) {
	return l.ImproveFunc(input)
}
