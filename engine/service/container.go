package service

import (
	"engine/model"
	"engine/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (engine *Engine) CreateContainer(ctx *gin.Context) {
	requestBody := &model.ContainerRunRequest{}

	if err := ctx.ShouldBindJSON(requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "BadParameter"})
		return
	}

	template, ok := engine.templateService.Get(requestBody.TemplateName)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "TemplateNotExit"})
		return
	}

	runtime, ok := engine.runtimeService.Get(template.Runtime)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "RuntimeNotExit"})
		return
	}

	requestId := util.UniqueId()
	id, err := engine.containerService.Create(requestId, runtime, template, requestBody.FunctionParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if requestBody.Synchronized {
		c := make(chan gin.H, 1)
		engine.functionResultWaitChanMap[requestId] = c
		response := <-c
		fmt.Println(response)
		ctx.JSON(http.StatusOK, response)
	} else {
		ctx.JSON(http.StatusOK, gin.H{"id": id})
	}
}

func (engine *Engine) ListContainer(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, engine.containerService.List())
}

func (engine *Engine) StopContainer(ctx *gin.Context) {
	// TODO 把容器进程直接杀掉
	ctx.JSON(http.StatusOK, "ok")
}

func (engine *Engine) DumpContainer(ctx *gin.Context) {
	if err := engine.containerService.Dump(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"data": "ok"})
	}
}
