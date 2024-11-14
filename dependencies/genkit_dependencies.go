package dependencies

import (
	"context"
	"os"

	"PromptDefender-Keep/logger"

	"go.uber.org/zap"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/plugins/dotprompt"
	"github.com/firebase/genkit/go/plugins/googlecloud"
	"github.com/firebase/genkit/go/plugins/vertexai"
	"github.com/invopop/jsonschema"
	"google.golang.org/api/option"
)

func InitialiseGenkit(ctx context.Context) {
	if err := vertexai.Init(ctx, &vertexai.Config{
		ProjectID:     os.Getenv("GCLOUD_PROJECT"),
		Location:      os.Getenv("GCLOUD_LOCATION"),
		ClientOptions: []option.ClientOption{option.WithCredentialsFile("service-account.json")},
	}); err != nil {
		logger.Log.Fatal("Error initializing Vertex AI", zap.Error(err))
	}

	dotprompt.SetDirectory("prompts")

	if err := googlecloud.Init(
		ctx,
		googlecloud.Config{ProjectID: os.Getenv("GCLOUD_PROJECT")},
	); err != nil {
		logger.Log.Fatal("Error initializing Google Cloud", zap.Error(err))
	}

}

func ProvideModel() ai.Model {
	g := vertexai.Model("gemini-1.5-pro")
	if g == nil {
		logger.Log.Fatal("Model is nil")
	}

	return g
}

func ProvideReflector() *jsonschema.Reflector {
	r := &jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}

	return r
}

func init() {
	// Check if the file 'service-account.json' exists on disk
	// Check if the environment variables are set
	// Log and exit if the environment variables are not set

	if _, err := os.Stat("service-account.json"); os.IsNotExist(err) {
		logger.Log.Fatal("service-account.json not found")
	}

	requiredEnvVars := []string{"GCLOUD_PROJECT", "GCLOUD_LOCATION"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			logger.Log.Fatal("Environment variable not set", zap.String("env var", envVar))
		}
	}
}
