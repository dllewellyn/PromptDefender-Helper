package main

import (
	"PromptDefender-Keep/score"
	"context"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/dotprompt"
	"github.com/firebase/genkit/go/plugins/vertexai"
	"github.com/gin-gonic/gin"
	"github.com/invopop/jsonschema"
	"google.golang.org/api/option"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()

	ctx := context.Background()

	initialiseGenkit(ctx)

	reflector, model := createReflectorAndModel()

	prompt, err := dotprompt.Open("extract_data")

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

	if err != nil {
		log.Fatal(err)
	}

	// Load HTML template from file
	tmpl, err := template.ParseFiles("templates/template.html")
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(tmpl)

	// Define a route and handler
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "template.html", "gopher")
	})

	r.POST("/score", func(c *gin.Context) {
		// Get the prompt from the text body
		prompt := c.PostForm("prompt")

		scorer := score.NewLlmScorer(func(prompt string) (string, error) {
			flow := genkit.DefineFlow("scorePromptSecurity", func(ctx context.Context, input string) (string, error) {
				response, err := scoreLlmPrompt.Generate(ctx, &dotprompt.PromptRequest{
					Variables: score.LlmScoringPromptInput{
						input,
					},
				}, nil)

				if err != nil {
					return "", err
				}

				return response.Text(), nil
			})

			return flow.Run(ctx, prompt)
		})

		response, err := scorer.Score(prompt)

		if err != nil {
			// Send an error response
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.HTML(http.StatusOK, "score.html", response)
	})

	// Start the server
	r.Run(":8080")
}

func initialiseGenkit(ctx context.Context) {
	if err := vertexai.Init(ctx, &vertexai.Config{
		ProjectID:     os.Getenv("GCLOUD_PROJECT"),
		Location:      "us-central1",
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
