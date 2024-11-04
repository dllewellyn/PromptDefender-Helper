package endpoints

import (
	"PromptDefender-Keep/cache"
	"PromptDefender-Keep/dependencies"
	"PromptDefender-Keep/score"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

func scorePromptToUiFriendlyResponse(inputPrompt string, scorePrompt *score.PromptScore, loadedDefences []dependencies.Defence) UiDefences {

	var defences []UiDefence

	for _, defence := range loadedDefences {

		isPresent := false

		switch defence.Id {
		case dependencies.InContext:
			isPresent = scorePrompt.Defenses.InContextDefense
		case dependencies.SystemModeSelfReminder:
			isPresent = scorePrompt.Defenses.SystemModeSelfReminder
		case dependencies.SandwichDefence:
			isPresent = scorePrompt.Defenses.SandwichDefense
		case dependencies.XmlEncapsulation:
			isPresent = scorePrompt.Defenses.XMLEncapsulation
		case dependencies.RandomSequenceEnclosure:
			isPresent = scorePrompt.Defenses.RandomSequenceEnclosure
		}

		defences = append(defences, UiDefence{
			Name:        defence.Name,
			Description: defence.Description,
			IsPresent:   isPresent,
			Link:        defence.Link,
		})
	}

	uiDefences := UiDefences{
		Score:       scorePrompt.Score,
		Explanation: scorePrompt.Explanation,
		Defences:    defences,
		Prompt:      inputPrompt,
	}

	return uiDefences
}

func AddScorer(ctx context.Context, engine *gin.Engine, scorer score.Scorer, promptCache cache.Cache, loadedDefences []dependencies.Defence) {
	engine.POST("/score", func(c *gin.Context) {

		prompt := c.PostForm("prompt")

		cachedResponse, err := promptCache.Get(ctx, prompt+"_score")

		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusOK, "/error")
			return
		}

		if cachedResponse != "" {
			// Convert cachedResponse to PromptScore
			var response score.PromptScore
			err = json.Unmarshal([]byte(cachedResponse), &response)
			if err != nil {
				log.Println(err)
				c.Redirect(http.StatusOK, "/error")
				return
			}

			c.HTML(http.StatusOK, "score.html", scorePromptToUiFriendlyResponse(prompt, &response, loadedDefences))
			return
		}

		response, err := scorer.Score(prompt)

		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusOK, "/error")
		}

		cacheResponse(response, promptCache, prompt+"_score")

		c.HTML(http.StatusOK, "score.html", scorePromptToUiFriendlyResponse(prompt, response, loadedDefences))
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
