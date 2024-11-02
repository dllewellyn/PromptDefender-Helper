package main

import (
	cache2 "PromptDefender-Keep/cache"
	"PromptDefender-Keep/endpoints"
	"PromptDefender-Keep/improve"
	"PromptDefender-Keep/score"
	"context"
	"fmt"
	"github.com/firebase/genkit/go/genkit"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"log"
	"os"
)

func main() {
	// Check if the file 'service-account.json' exists on disk
	if _, err := os.Stat("service-account.json"); os.IsNotExist(err) {
		// Retrieve the service account key from the environment variable
		serviceAccountKey := os.Getenv("SERVICE_ACCOUNT_KEY")
		if serviceAccountKey == "" {
			log.Fatal("SERVICE_ACCOUNT_KEY environment variable not set")
		}

		// Write the service account key to disk
		file, err := os.Create("service-account.json")
		if err != nil {
			log.Fatal(err)
		}

		_, err = file.WriteString(serviceAccountKey)
	}

	r := gin.Default()

	log.Println("Serving directory ./public at /")
	r.Static("/", "./public")

	r.LoadHTMLGlob("templates/*.html")

	ctx := context.Background()

	initialiseGenkit(ctx)

	scoreLlmTag := fx.ResultTags(`name:"scoreLlm.prompt"`)
	improveLlmTag := fx.ResultTags(`name:"llmImprover.prompt"`)

	app := fx.New(
		fx.Provide(
			ProvideModel,
			ProvideReflector,
			fx.Annotate(ProvideScoringPrompt, scoreLlmTag),
			ProvideScorer,
			fx.Annotate(ProvideImprovePrompt, improveLlmTag),
			ProvideImprover,
		),
		fx.Invoke(func(scorer score.Scorer, improver improve.Improver) {
			cache := cache2.NewInMemoryCache()
			endpoints.AddScorer(ctx, r, scorer, cache)
			endpoints.AddImprover(ctx, r, improver, cache)

			if os.Getenv("TEST_MODE") == "true" {
				log.Println("Initialising genkit in test mode")
				if err := genkit.Init(ctx, nil); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Println("Starting server on port", os.Getenv("PORT"))

				err := r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))

				if err != nil {
					log.Fatal(err)
				}
			}
		}),
	)

	app.Run()
}
