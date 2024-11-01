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

		prompt := c.PostForm("prompt")

		log.Println(prompt)

		cachedResponse, err := promptCache.Get(ctx, prompt)

		if err != nil {
			c.Redirect(http.StatusInternalServerError, "/error")
			return
		}

		if cachedResponse != "" {
			c.HTML(http.StatusOK, "improve.html", cachedResponse)
			return
		}

		response, err := improver.Improve(prompt)

		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusInternalServerError, "/error")
		}

		c.HTML(http.StatusOK, "improve.html", response)
	})
}
