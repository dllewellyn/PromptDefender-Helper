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

type DefenceLevel = int

const (
	NoDefence DefenceLevel = iota
	PartialDefence
	FullDefence
)

type UiDefence struct {
	Name         string
	Description  string
	DefenceLevel DefenceLevel
	Link         string
}

type UiDefences struct {
	Defences    []UiDefence
	Score       float64
	Explanation string
	Prompt      string
}

// Function to convert a float between 0.0 and 0.2 to defence level
// 0.0 -> NoDefence
// 0.1 -> PartialDefence
// 0.2 -> FullDefence
func floatToDefenceLevel(input float64) DefenceLevel {
	if input == 0.0 {
		return NoDefence
	} else if input == 0.1 {
		return PartialDefence
	} else {
		return FullDefence
	}
}

func scorePromptToUiFriendlyResponse(inputPrompt string, scorePrompt *score.PromptScore, loadedDefences []dependencies.Defence) UiDefences {

	var defences []UiDefence

	for _, defence := range loadedDefences {

		defenceLevel := NoDefence

		switch defence.Id {
		case dependencies.InContext:
			defenceLevel = floatToDefenceLevel(scorePrompt.Defenses.InContextDefense)
		case dependencies.SystemModeSelfReminder:
			defenceLevel = floatToDefenceLevel(scorePrompt.Defenses.SystemModeSelfReminder)
		case dependencies.SandwichDefence:
			defenceLevel = floatToDefenceLevel(scorePrompt.Defenses.SandwichDefense)
		case dependencies.XmlEncapsulation:
			defenceLevel = floatToDefenceLevel(scorePrompt.Defenses.XMLEncapsulation)
		case dependencies.RandomSequenceEnclosure:
			defenceLevel = floatToDefenceLevel(scorePrompt.Defenses.RandomSequenceEnclosure)
		}

		defences = append(defences, UiDefence{
			Name:         defence.Name,
			Description:  defence.Description,
			DefenceLevel: defenceLevel,
			Link:         defence.Link,
		})
	}

	uiDefences := UiDefences{
		Score:       *scorePrompt.Score,
		Explanation: scorePrompt.Explanation,
		Defences:    defences,
		Prompt:      inputPrompt,
	}

	return uiDefences
}

func AddScorer(ctx context.Context, engine *gin.Engine, scorer score.Scorer, promptCache cache.Cache, loadedDefences []dependencies.Defence) {
	engine.POST("/score", func(c *gin.Context) {

		var prompt string
		var requestBody struct {
			Prompt string `json:"prompt"`
		}
		if c.Request.Header.Get("Content-Type") == "application/json" {
			if err := c.ShouldBindJSON(&requestBody); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
				return
			}
			prompt = requestBody.Prompt
		} else {
			prompt = c.PostForm("prompt")
		}

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

			// Check Content-Type header for desired response format
			if c.Request.Header.Get("Content-Type") == "application/json" {
				c.JSON(http.StatusOK, gin.H{
					"response": response,
				})
			} else {
				c.HTML(http.StatusOK, "score.html", scorePromptToUiFriendlyResponse(prompt, &response, loadedDefences))
			}

			return
		}

		response, err := scorer.Score(prompt)

		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusOK, "/error")
		}

		cacheResponse(response, promptCache, prompt+"_score")

		// Check Content-Type header for desired response format
		if c.Request.Header.Get("Content-Type") == "application/json" {
			c.JSON(http.StatusOK, gin.H{
				"response": response,
			})
		} else {
			c.HTML(http.StatusOK, "score.html", scorePromptToUiFriendlyResponse(prompt, response, loadedDefences))
		}
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
