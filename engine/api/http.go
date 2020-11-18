package api

import (
	"engine/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetContainerRouter(runtime *service.RuntimeService, template *service.TemplateService, container *service.ContainerService, r *gin.Engine) {
	w := r.Group("/")

	w.GET("ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	w.GET("/runtime", runtime.List)

	w.GET("/template", template.List)
	w.POST("/template", template.Create)
	w.DELETE("/template/:name", template.Delete)

	w.GET("/container", container.List)
	w.POST("/container/run", container.Create)
	w.DELETE("/container/:name", container.Delete)
}
