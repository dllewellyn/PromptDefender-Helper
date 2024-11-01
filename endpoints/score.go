package endpoints

import (
	"PromptDefender-Keep/cache"
	"PromptDefender-Keep/score"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type UiDefence struct {
	Name        string
	Description string
	IsPresent   bool
	Link        string
}

type UiDefences struct {
	Defences    []UiDefence
	Score       float64
	Explanation string
	Prompt      string
}

func scorePromptToUiFriendlyResponse(inputPrompt string, scorePrompt *score.PromptScore) UiDefences {
	inContextDefence := UiDefence{
		Name:        "In Context Defence",
		Description: "This defence is used to...",
		IsPresent:   scorePrompt.Defenses.InContextDefense,
	}

	systemModeSelfReminder := UiDefence{
		Name:        "System Mode Self Reminder",
		Description: "This defence is used to...",
		IsPresent:   scorePrompt.Defenses.SystemModeSelfReminder,
	}

	sandwichDefence := UiDefence{
		Name:        "Sandwich Defence",
		Description: "This defence is used to...",
		IsPresent:   scorePrompt.Defenses.SandwichDefense,
	}

	xmlEncapsulation := UiDefence{
		Name:        "XML Encapsulation",
		Description: "This defence is used to...",
		IsPresent:   scorePrompt.Defenses.XMLEncapsulation,
	}

	randomSequenceEnclosure := UiDefence{
		Name:        "Random Sequence Enclosure",
		Description: "This defence is used to...",
		IsPresent:   scorePrompt.Defenses.RandomSequenceEnclosure,
	}

	defences := []UiDefence{
		inContextDefence,
		systemModeSelfReminder,
		sandwichDefence,
		xmlEncapsulation,
		randomSequenceEnclosure,
	}

	uiDefences := UiDefences{
		Score:       scorePrompt.Score,
		Explanation: scorePrompt.Explanation,
		Defences:    defences,
		Prompt:      inputPrompt,
	}

	return uiDefences
}

func AddScorer(ctx context.Context, engine *gin.Engine, scorer score.Scorer, promptCache cache.Cache) {
	engine.POST("/score", func(c *gin.Context) {

		prompt := c.PostForm("prompt")

		cachedResponse, err := promptCache.Get(ctx, prompt)

		if err != nil {
			c.Redirect(http.StatusInternalServerError, "/error")
			return
		}

		if cachedResponse != "" {
			// Convert cachedResponse to PromptScore
			var response score.PromptScore
			err = json.Unmarshal([]byte(cachedResponse), &response)
			if err != nil {
				c.Redirect(http.StatusInternalServerError, "/error")
				return
			}

			c.HTML(http.StatusOK, "score.html", scorePromptToUiFriendlyResponse(prompt, &response))
			return
		}

		response, err := scorer.Score(prompt)

		if err != nil {
			c.Redirect(http.StatusInternalServerError, "/error")
		}

		cacheResponse(response, promptCache, prompt)

		c.HTML(http.StatusOK, "score.html", scorePromptToUiFriendlyResponse(prompt, response))
	})
}

func cacheResponse(response *score.PromptScore, promptCache cache.Cache, prompt string) {
	responseJson, err := json.Marshal(response)

	if err != nil {
		log.Println(err)
	}

	err = promptCache.Set(context.Background(), prompt, string(responseJson))

	if err != nil {
		// Log the error but ignore - we don't want to fail the request
		// if the cache fails
		log.Println(err)
	}
}
