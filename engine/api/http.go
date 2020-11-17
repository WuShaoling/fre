package api

import (
	"engine/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetContainerRouter(function *service.FunctionService, container *service.ContainerService, r *gin.Engine) {
	w := r.Group("/")

	w.GET("ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	w.POST("/function", function.Create)
	w.DELETE("/function/:name", function.Delete)

	w.POST("/container/run", container.Create)
	w.DELETE("/container/:name", container.Delete)
}
