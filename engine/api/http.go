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

	runtimeGroup := r.Group("/runtime")
	runtimeGroup.GET("/", runtime.List)
	runtimeGroup.GET("/dump", runtime.Dump)

	templateGroup := r.Group("/template")
	templateGroup.GET("/", template.List)
	templateGroup.POST("/", template.Create)
	templateGroup.DELETE("/:name", template.Delete)
	templateGroup.GET("/dump", template.Dump)

	containerGroup := r.Group("/container")
	containerGroup.GET("/", container.List)
	containerGroup.POST("/", container.Create)
	containerGroup.DELETE("/:id", container.Delete)
	containerGroup.PUT("/callback/run/:id/:pid", container.OnContainerRun)
	containerGroup.PUT("/callback/exit/:id", container.OnContainerExit)
	containerGroup.GET("/Dump", container.Dump)
}
