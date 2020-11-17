package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetContainerRouter(r *gin.Engine) {
	w := r.Group("/")

	w.GET("ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

}
