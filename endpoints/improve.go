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

		cachedResponse, err := promptCache.Get(ctx, prompt+"_improve")

		if err != nil {
			c.Redirect(http.StatusOK, "/error")
			return
		}

		if cachedResponse != "" {
			c.HTML(http.StatusOK, "improve.html", cachedResponse)
			return
		}

		response, err := improver.Improve(prompt)

		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusOK, "/error")
		}

		err = promptCache.Set(ctx, prompt+"_improve", response)

		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusOK, "/error")
		}
		c.HTML(http.StatusOK, "improve.html", response)
	})
}
