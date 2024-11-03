package main

import (
	"PromptDefender-Keep/improve"
	"context"
	"log"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/dotprompt"
	"github.com/invopop/jsonschema"
	"go.uber.org/fx"
)

func ProvideImprovePrompt(model ai.Model, reflector *jsonschema.Reflector) *dotprompt.Prompt {
	prompt, err := dotprompt.Open("suggest_improvements")

	if err != nil {
		log.Fatal(err)
	}

	scoreLlmPrompt, err := dotprompt.Define("suggest_improvements.prompt", prompt.TemplateText,
		dotprompt.Config{
			Model:        model,
			InputSchema:  reflector.Reflect(improve.LlmPromptImproverInput{}),
			OutputFormat: ai.OutputFormatText,
		},
	)

	if err != nil {
		log.Fatal(err)
	}
	return scoreLlmPrompt
}

func ProvideImprover(params struct {
	fx.In
	LlmImprover *dotprompt.Prompt `name:"llmImprover.prompt"`
}) improve.Improver {

	llmImproverPromptFlow := genkit.DefineFlow("llmImprover", func(ctx context.Context, input improve.LlmPromptImproverInput) (string, error) {
		response, err := params.LlmImprover.Generate(ctx, &dotprompt.PromptRequest{
			Variables: input,
		}, nil)

		if err != nil {
			return "", err
		}

		return response.Text(), nil
	})

	return improve.NewLlmImprover(func(prompt string, randomSequence string) (string, error) {
		return llmImproverPromptFlow.Run(context.Background(), improve.LlmPromptImproverInput{
			StartingPrompt: prompt, RandomSequence: randomSequence,
		})
	})
}
