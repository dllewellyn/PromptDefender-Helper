package main

import (
	cache2 "PromptDefender-Keep/cache"
	"PromptDefender-Keep/dependencies"
	"PromptDefender-Keep/endpoints"
	"PromptDefender-Keep/gh"
	"PromptDefender-Keep/improve"
	"PromptDefender-Keep/logger"
	"PromptDefender-Keep/score"
	"context"
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"

	"github.com/firebase/genkit/go/genkit"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func init() {
	// Check if the file 'service-account.json' exists on disk
	// Retrieve the service account key from the environment variable
	// Write the service account key to disk
	WriteServiceAccountKeyToFile()
}

func main() {

	ctx := context.Background()

	logger.GetLogger().Debug("Initialising genkit")

	dependencies.InitialiseGenkit(ctx)

	scoreLlmTag := fx.ResultTags(`name:"scoreLlm.prompt"`)
	improveLlmTag := fx.ResultTags(`name:"llmImprover.prompt"`)

	app := fx.New(
		fx.Provide(
			dependencies.ProvideModel,
			dependencies.ProvideReflector,
			fx.Annotate(dependencies.ProvideScoringPrompt, scoreLlmTag),
			dependencies.ProvideScorer,
			fx.Annotate(dependencies.ProvideImprovePrompt, improveLlmTag),
			dependencies.ProvideImprover,
			dependencies.ProvideDefences,
		),
		fx.Invoke(func(scorer score.Scorer, improver improve.Improver, loadedDefences []dependencies.Defence) {
			logger.GetLogger().Info("Invoke function called")

			logger.GetLogger().Info("Test mode", zap.String("test_mode", os.Getenv("TEST_MODE")))

			if os.Getenv("TEST_MODE") == "true" {
				logger.GetLogger().Info("Initialising genkit in test mode")
				if err := genkit.Init(ctx, nil); err != nil {
					logger.GetLogger().Fatal("Error initializing genkit", zap.Error(err))
				}
			} else {
				logger.GetLogger().Info("Starting server on port", zap.String("port", os.Getenv("PORT")))
				r := gin.Default()

				if logger.GetLogger() == nil {
					log.Fatal("Logger not initialised")
				}

				logger.GetLogger().Debug("Serving directory ./public at /")

				r.Static("/", "./public")

				r.LoadHTMLGlob("templates/*.html")

				cache := cache2.NewInMemoryCache()
				endpoints.AddScorer(ctx, r, scorer, cache, loadedDefences)
				endpoints.AddImprover(ctx, r, improver, cache)
				setupGitHubAppBackend(r, scorer)

				err := r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))

				if err != nil {
					logger.GetLogger().Fatal("Error starting server", zap.Error(err))
				}
			}
		}),
	)

	app.Run()
}

func WriteServiceAccountKeyToFile() {
	if _, err := os.Stat("service-account.json"); os.IsNotExist(err) {

		serviceAccountKey := os.Getenv("SERVICE_ACCOUNT_KEY")
		if serviceAccountKey == "" {
			logger.GetLogger().Fatal("SERVICE_ACCOUNT_KEY environment variable not set")
		}

		file, err := os.Create("service-account.json")
		if err != nil {
			logger.GetLogger().Fatal("Error creating account key", zap.Error(err))
		}

		_, err = file.WriteString(serviceAccountKey)
		if err != nil {
			logger.GetLogger().Fatal("Error writing service account key to file", zap.Error(err))
		}
	}
}
func setupGitHubAppBackend(r *gin.Engine, scorer score.Scorer) {
	logger.GetLogger().Info("Setting up GitHub App backend")
	r.POST("/github/callback", func(c *gin.Context) {
		logger.GetLogger().Info("Handling GitHub callback")
		gh.HandleWebhook(c, scorer)
	})
}
