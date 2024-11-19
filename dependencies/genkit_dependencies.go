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
)

func InitialiseGenkit(ctx context.Context) {
	if err := vertexai.Init(ctx, &vertexai.Config{
		ProjectID: os.Getenv("GCLOUD_PROJECT"),
		Location:  os.Getenv("GCLOUD_LOCATION"),
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
