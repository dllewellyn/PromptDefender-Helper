package dependencies

import (
	"PromptDefender-Keep/score"
	"context"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/dotprompt"
	"github.com/invopop/jsonschema"
	"go.uber.org/fx"
	"log"
)

func ProvideScoringPrompt(model ai.Model, reflector *jsonschema.Reflector) *dotprompt.Prompt {
	prompt, err := dotprompt.Open("scoring_prompt")
	if err != nil {
		log.Printf("Error opening scoring prompt: %v", err)
		log.Fatal(err)
	}

	scoreLlmPrompt, err := dotprompt.Define("scoreLlm.prompt", prompt.TemplateText,
		dotprompt.Config{
			Model:        model,
			InputSchema:  reflector.Reflect(score.LlmScoringPromptInput{}),
			OutputFormat: ai.OutputFormatText,
		},
	)
	if err != nil {
		log.Printf("Error defining scoreLlm.prompt: %v", err)
		log.Fatal(err)
	}
	return scoreLlmPrompt
}

func ProvideValidateScorePrompt(model ai.Model, reflector *jsonschema.Reflector) *dotprompt.Prompt {
	prompt, err := dotprompt.Open("validate_prompt_consistency")

	if err != nil {
		log.Fatal(err)
	}

	validateScorePrompt, err := dotprompt.Define("validateScore.prompt", prompt.TemplateText,
		dotprompt.Config{
			Model:        model,
			InputSchema:  reflector.Reflect(score.ValidatorPromptInput{}),
			OutputFormat: ai.OutputFormatText,
		},
	)

	if err != nil {
		log.Fatal(err)
	}
	return validateScorePrompt
}

func ProvideValidator(params struct {
	fx.In
	ValidateScorePrompt *dotprompt.Prompt `name:"validateScore.prompt"`
}) score.Validator {

	validateScoreFlow := genkit.DefineFlow("validateScore", func(ctx context.Context, input string) (string, error) {
		response, err := params.ValidateScorePrompt.Generate(context.Background(), &dotprompt.PromptRequest{
			Variables: score.ValidatorPromptInput{Score: input},
		}, nil)

		if err != nil {
			return "", err
		}

		return response.Text(), nil
	})

	invokeRequest := func(score string) (string, error) {
		return validateScoreFlow.Run(context.Background(), score)
	}

	return score.NewValidator(invokeRequest)
}

func ProvideScorer(params struct {
	fx.In
	ScoreLlmPrompt *dotprompt.Prompt `name:"scoreLlm.prompt"`
}, ScoreValidator score.Validator) score.Scorer {

	scorePromptFlow := genkit.DefineFlow("scorePromptSecurity", func(ctx context.Context, input string) (string, error) {
		response, err := params.ScoreLlmPrompt.Generate(ctx, &dotprompt.PromptRequest{
			Variables: score.LlmScoringPromptInput{
				StartingPrompt: input,
			},
		}, nil)

		if err != nil {
			return "", err
		}

		return response.Text(), nil
	})

	invokeRequest := func(prompt string) (string, error) {
		return scorePromptFlow.Run(context.Background(), prompt)
	}

	return score.NewLlmScorer(invokeRequest, ScoreValidator)
}
