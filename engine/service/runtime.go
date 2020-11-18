package service

import (
	"engine/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RuntimeService struct {
	engine  *Engine
	dataMap map[string]model.Runtime
}

func NewRuntimeService(engine *Engine) *RuntimeService {
	service := &RuntimeService{
		engine:  engine,
		dataMap: make(map[string]model.Runtime),
	}
	loadDataFromFile(&service.dataMap, RuntimeDataFileName)
	return service
}

func (service *RuntimeService) List(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, service.dataMap)
}
