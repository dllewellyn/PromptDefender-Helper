package endpoints

import (
	"PromptDefender-Keep/cache"
	"PromptDefender-Keep/improve"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func AddImprover(ctx context.Context, engine *gin.Engine, improver improve.Improver, promptCache cache.Cache) {
	engine.POST("/improve", func(c *gin.Context) {
		prompt, err := getPromptFromRequest(c)

		if err != nil {
			handleError(c, http.StatusBadRequest, "Sorry, something went wrong. Please check your request try again.")
			return
		}

		if prompt == "" {
			handleError(c, http.StatusBadRequest, "Sorry, something went wrong. Please check your request try again.")
			return
		}

		cachedResponse, err := promptCache.Get(ctx, prompt+"_improve")
		if err != nil {
			log.Println("Error getting cached response:", err)
			handleError(c, http.StatusInternalServerError, "Sorry, something went wrong. Please try again.")
			return
		}

		if cachedResponse != "" {
			handleResponse(c, cachedResponse)
			return
		}

		response, err := improver.Improve(prompt)

		if err != nil {
			log.Println(err)
			handleError(c, http.StatusInternalServerError, "Sorry, something went wrong. Please try again.")
			return
		}

		err = promptCache.Set(ctx, prompt+"_improve", response)
		if err != nil {
			log.Println(err)
			handleError(c, http.StatusInternalServerError, "Sorry, something went wrong. Please try again.")
			return
		}

		handleResponse(c, response)
	})
}

func handleResponse(c *gin.Context, response string) {
	if c.GetHeader("Content-Type") == "application/json" {
		c.JSON(http.StatusOK, gin.H{"response": response})
	} else {
		c.HTML(http.StatusOK, "improve.html", response)
	}
}

func getPromptFromRequest(c *gin.Context) (string, error) {
	contentType := c.GetHeader("Content-Type")
	if contentType == "application/json" {
		var jsonData map[string]string
		if err := c.ShouldBindJSON(&jsonData); err != nil {
			return "", err
		}
		return jsonData["prompt"], nil
	}
	return c.PostForm("prompt"), nil
}

func handleError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"message": message})
}
