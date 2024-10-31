package main

import (
	"PromptDefender-Keep/cache"
	"PromptDefender-Keep/score"
	"context"
	"fmt"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/dotprompt"
	"github.com/firebase/genkit/go/plugins/vertexai"
	"github.com/gin-gonic/gin"
	"github.com/invopop/jsonschema"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()
	r.Static("/", "./public")
	r.LoadHTMLGlob("templates/*.html")

	ctx := context.Background()

	initialiseGenkit(ctx)

	reflector, model := createReflectorAndModel()

	err, scoreLlmPrompt := retrievePrompt(model, reflector)

	if err != nil {
		log.Fatal(err)
	}

	promptScorerFlow := genkit.DefineFlow("scorePromptSecurity", func(ctx context.Context, input string) (string, error) {
		response, err := scoreLlmPrompt.Generate(ctx, &dotprompt.PromptRequest{
			Variables: score.LlmScoringPromptInput{
				StartingPrompt: input,
			},
		}, nil)

		if err != nil {
			return "", err
		}

		return response.Text(), nil
	})

	r.POST("/score", func(c *gin.Context) {
		// Get the prompt from the text body
		prompt := c.PostForm("prompt")
		promptCache := cache.NewInMemoryCache()

		// Check if the prompt is in the cache
		cachedResponse, err := promptCache.Get(ctx, prompt)

		if err != nil {
			c.HTML(http.StatusOK, "error.html", gin.H{})
			return
		}

		if cachedResponse != "" {
			c.HTML(http.StatusOK, "error.html", gin.H{"response": cachedResponse})
			return
		}

		scorer := score.NewLlmScorer(func(prompt string) (string, error) {
			return promptScorerFlow.Run(ctx, prompt)
		})

		response, err := scorer.Score(prompt)

		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{})
		}

		// Cache the response
		err = promptCache.Set(ctx, prompt, response.Explanation)

		if err != nil {
			// Log the error but ignore - we don't want to fail the request
			// if the cache fails
			log.Println(err)
		}

		c.HTML(http.StatusOK, "score.html", *response)
	})

	// Start the server
	err = r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))

	if err != nil {
		log.Fatal(err)
	}
}

func retrievePrompt(model ai.Model, reflector *jsonschema.Reflector) (error, *dotprompt.Prompt) {
	prompt, err := dotprompt.Open("scoring_prompt")

	if err != nil {
		log.Fatal(err)
	}

	scoreLlmPrompt, err := dotprompt.Define("score_llm.prompt", prompt.TemplateText,
		dotprompt.Config{
			Model:        model,
			InputSchema:  reflector.Reflect(score.LlmScoringPromptInput{}),
			OutputFormat: ai.OutputFormatText,
		},
	)
	return err, scoreLlmPrompt
}

func initialiseGenkit(ctx context.Context) {
	if err := vertexai.Init(ctx, &vertexai.Config{
		ProjectID:     os.Getenv("GCLOUD_PROJECT"),
		Location:      os.Getenv("GCLOUD_LOCATION"),
		ClientOptions: []option.ClientOption{option.WithCredentialsFile("service-account.json")},
	}); err != nil {
		log.Fatal(err)
	}
}

func createReflectorAndModel() (*jsonschema.Reflector, ai.Model) {
	r := &jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}

	g := vertexai.Model("gemini-1.5-pro")
	if g == nil {
		log.Fatal("Model is nil")
	}

	dotprompt.SetDirectory("prompts")

	return r, g
}
