package dependencies

import (
	"context"
	"log"
	"os"

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
		log.Fatal(err)
	}

	dotprompt.SetDirectory("prompts")

	if err := googlecloud.Init(
		ctx,
		googlecloud.Config{ProjectID: os.Getenv("GCLOUD_PROJECT")},
	); err != nil {
		log.Fatal(err)
	}

}

func ProvideModel() ai.Model {
	g := vertexai.Model("gemini-1.5-pro")
	if g == nil {
		log.Fatal("Model is nil")
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
