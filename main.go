package main

import (
	cache2 "PromptDefender-Keep/cache"
	"PromptDefender-Keep/dependencies"
	"PromptDefender-Keep/endpoints"
	"PromptDefender-Keep/improve"
	"PromptDefender-Keep/logger"
	"PromptDefender-Keep/score"
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/firebase/genkit/go/genkit"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func main() {
	// Check if the file 'service-account.json' exists on disk
	// Retrieve the service account key from the environment variable
	// Write the service account key to disk
	WriteServiceAccountKeyToFile()

	r := gin.Default()

	logger.Log.Debug("Serving directory ./public at /")

	r.Static("/", "./public")

	r.LoadHTMLGlob("templates/*.html")

	r.Use(func(c *gin.Context) {
		for k, v := range c.Request.Header {
			logger.Log.Info("Header", zap.String("key", k), zap.Strings("values", v))
		}
		c.Next()
	})

	ctx := context.Background()

	logger.Log.Debug("Initialising genkit")

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
			if os.Getenv("TEST_MODE") == "true" {
				logger.Log.Info("Initialising genkit in test mode")
				if err := genkit.Init(ctx, nil); err != nil {
					logger.Log.Fatal("Error initializing genkit", zap.Error(err))
				}
			} else {
				logger.Log.Info("Starting server on port", zap.String("port", os.Getenv("PORT")))

				cache := cache2.NewInMemoryCache()
				endpoints.AddScorer(ctx, r, scorer, cache, loadedDefences)
				endpoints.AddImprover(ctx, r, improver, cache)

				err := r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))

				if err != nil {
					logger.Log.Fatal("Error starting server", zap.Error(err))
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
			logger.Log.Fatal("SERVICE_ACCOUNT_KEY environment variable not set")
		}

		// Base64 decode the service account key
		decodedServiceAccountKey, err := base64.StdEncoding.DecodeString(serviceAccountKey)

		if err != nil {
			logger.Log.Fatal("Error decoding service account key", zap.Error(err))
		}

		serviceAccountKey = string(decodedServiceAccountKey)

		file, err := os.Create("service-account.json")
		if err != nil {
			logger.Log.Fatal("Error creating account key", zap.Error(err))
		}

		_, err = file.WriteString(serviceAccountKey)
		if err != nil {
			logger.Log.Fatal("Error writing service account key to file", zap.Error(err))
		}
	}
}
