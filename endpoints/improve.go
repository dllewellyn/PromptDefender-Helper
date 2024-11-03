package endpoints

import (
	"PromptDefender-Keep/cache"
	"PromptDefender-Keep/improve"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// AddImprover adds an improve endpoint to the gin engine
func AddImprover(ctx context.Context, engine *gin.Engine, improver improve.Improver, promptCache cache.Cache) {
	engine.POST("/api/improve", func(c *gin.Context) {

		var prompt struct {
			Prompt string `json:"prompt"`
		}

		if err := c.BindJSON(&prompt); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		cachedResponse, err := promptCache.Get(ctx, prompt.Prompt)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		if cachedResponse != "" {
			c.JSON(http.StatusOK, gin.H{"result": cachedResponse})
			return
		}

		response, err := improver.Improve(prompt.Prompt)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		err = promptCache.Set(ctx, prompt.Prompt, response)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": response})
	})
}
