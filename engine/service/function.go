package service

import (
	"encoding/json"
	"engine/model"
	"engine/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FunctionService struct {
	dataMap map[string]model.Function
}

func NewFunctionService() *FunctionService {
	service := &FunctionService{
		dataMap: make(map[string]model.Function),
	}
	loadDataFromFile(&service.dataMap, FunctionDataFileName)
	return service
}

func (service *FunctionService) Create(ctx *gin.Context) {
	functionInfo := &model.Function{}
	if err := ctx.ShouldBindJSON(functionInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, "请求失败，参数不完整")
		return
	}

	if _, ok := service.dataMap[functionInfo.Name]; ok {
		ctx.JSON(http.StatusBadRequest, functionInfo.Name+"已存在，请更换")
		return
	}

	// 更新缓存
	service.dataMap[functionInfo.Name] = *functionInfo

	// 写回文件
	data, err := json.Marshal(service.dataMap)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	if err := util.WriteToFile(getDataFilePath(FunctionDataFileName), data); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, "ok")

	// TODO 上传代码
}

func (service *FunctionService) Delete(ctx *gin.Context) {
	// TODO 删除代码
	// TODO 删除信息
	// TODO 更新文件
	ctx.JSON(http.StatusOK, "ok")
}
