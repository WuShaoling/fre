package service

import (
	"encoding/json"
	"engine/model"
	"engine/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TemplateService struct {
	engine  *Engine
	dataMap map[string]*model.Template
}

func NewTemplateService(engine *Engine) *TemplateService {
	service := &TemplateService{
		engine:  engine,
		dataMap: make(map[string]*model.Template),
	}
	loadDataFromFile(&service.dataMap, TemplateDataFileName)
	return service
}

func (service *TemplateService) Create(ctx *gin.Context) {
	functionInfo := &model.Template{}
	if err := ctx.ShouldBindJSON(functionInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, "参数不完整")
		return
	}

	if _, ok := service.dataMap[functionInfo.Metadata.Name]; ok {
		ctx.JSON(http.StatusBadRequest, functionInfo.Metadata.Name+" 已存在，请更换")
		return
	}

	// TODO 拉取代码
	// TODO 加载额外所需的共享库和依赖包

	// 更新缓存
	service.dataMap[functionInfo.Metadata.Name] = functionInfo

	// 写回文件
	data, err := json.Marshal(service.dataMap)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	if err := util.WriteToFile(getDataFilePath(TemplateDataFileName), data); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, "ok")
}

func (service *TemplateService) List(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, service.dataMap)
}

func (service *TemplateService) Delete(ctx *gin.Context) {
	// TODO
	ctx.JSON(http.StatusOK, "ok")
}

func (service *TemplateService) Dump(ctx *gin.Context) {
	if err := dumpDataToFile(service.dataMap, TemplateDataFileName); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	} else {
		ctx.String(http.StatusOK, "ok")
	}
}

func (service *TemplateService) getTemplateByName(name string) *model.Template {
	return service.dataMap[name]
}

func (service *TemplateService) checkIsTemplateExist(name string) bool {
	_, ok := service.dataMap[name]
	return ok
}
