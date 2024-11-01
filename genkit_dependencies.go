package main

import (
	"context"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/plugins/dotprompt"
	"github.com/firebase/genkit/go/plugins/vertexai"
	"github.com/invopop/jsonschema"
	"google.golang.org/api/option"
	"log"
	"os"
)

func initialiseGenkit(ctx context.Context) {
	if err := vertexai.Init(ctx, &vertexai.Config{
		ProjectID:     os.Getenv("GCLOUD_PROJECT"),
		Location:      os.Getenv("GCLOUD_LOCATION"),
		ClientOptions: []option.ClientOption{option.WithCredentialsFile("service-account.json")},
	}); err != nil {
		log.Fatal(err)
	}

	dotprompt.SetDirectory("prompts")

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
