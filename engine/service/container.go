package service

import (
	"engine/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ContainerService struct {
	engine  *Engine
	dataMap map[string]model.Container
}

func NewContainerService(engine *Engine) *ContainerService {
	service := &ContainerService{
		engine:  engine,
		dataMap: make(map[string]model.Container),
	}
	loadDataFromFile(&service.dataMap, ContainerDataFileName)
	return service
}

func (service *ContainerService) Create(ctx *gin.Context) {

}

func (service *ContainerService) List(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, service.dataMap)
}

func (service *ContainerService) Delete(ctx *gin.Context) {
	// TODO
	ctx.JSON(http.StatusOK, "ok")
}
