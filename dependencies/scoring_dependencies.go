package dependencies

import (
	"PromptDefender-Keep/logger"
	"PromptDefender-Keep/score"
	"context"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/dotprompt"
	"github.com/invopop/jsonschema"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func ProvideScoringPrompt(model ai.Model, reflector *jsonschema.Reflector) *dotprompt.Prompt {
	prompt, err := dotprompt.Open("scoring_prompt")
	if err != nil {
		logger.Log.Error("Error opening scoring prompt", zap.Error(err))
		logger.Log.Fatal("Error opening scoring prompt", zap.Error(err))
	}

	scoreLlmPrompt, err := dotprompt.Define("scoreLlm.prompt", prompt.TemplateText,
		dotprompt.Config{
			Model:        model,
			InputSchema:  reflector.Reflect(score.LlmScoringPromptInput{}),
			OutputFormat: ai.OutputFormatText,
		},
	)
	if err != nil {
		logger.Log.Error("Error defining scoreLlm.prompt", zap.Error(err))
		logger.Log.Fatal("Error defining scoreLlm.prompt", zap.Error(err))
	}
	return scoreLlmPrompt
}

func ProvideScorer(params struct {
	fx.In
	ScoreLlmPrompt *dotprompt.Prompt `name:"scoreLlm.prompt"`
}) score.Scorer {

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

	return score.NewLlmScorer(invokeRequest)
}
