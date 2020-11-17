package service

import (
	"engine/model"
	"github.com/gin-gonic/gin"
)

type ContainerService struct {
	dataMap map[string]model.Container

	zygoteService     *ZygoteService
	functionService   *FunctionService
	cgroupPoolService *CgroupPoolService
}

func NewContainerService(cgroupPoolService *CgroupPoolService, functionService *FunctionService, zygoteService *ZygoteService) *ContainerService {
	service := &ContainerService{
		zygoteService:     zygoteService,
		functionService:   functionService,
		cgroupPoolService: cgroupPoolService,
		dataMap:           make(map[string]model.Container),
	}
	loadDataFromFile(&service.dataMap, ContainerDataFileName)
	return service
}

func (service *ContainerService) Create(context *gin.Context) {

}

func (service *ContainerService) Delete(context *gin.Context) {

}
