package service

import (
	"engine/model"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ContainerService struct {
	engine     *Engine
	dataMap    map[string]model.Container // 保存容器信息
	resultData map[string]gin.H           // 同步执行模式下，暂存执行结果
	resultChan map[string]chan gin.H      // 同步执行模式下，等待在这里
}

func NewContainerService(engine *Engine) *ContainerService {
	service := &ContainerService{
		engine:     engine,
		dataMap:    make(map[string]model.Container),
		resultData: make(map[string]gin.H),
		resultChan: make(map[string]chan gin.H),
	}
	loadDataFromFile(&service.dataMap, ContainerDataFileName)
	return service
}

func (service *ContainerService) Create(ctx *gin.Context) {

	// 获取请求参数
	requestBody := &model.RunContainerRequestBody{}
	if err := ctx.ShouldBindJSON(requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, "参数不完整")
		return
	}

	// 获取模版信息
	template, ok := service.engine.TemplateService.dataMap[requestBody.TemplateName]
	if !ok {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("模版 %s 不存在", requestBody.TemplateName))
		return
	}

	// 创建容器的过程
	id, errInfo := service.newContainer(requestBody, &template)
	if errInfo != "" { // 如果容器创建失败了，直接返回
		ctx.JSON(http.StatusInternalServerError, errInfo)
		return
	}

	// 到这里，容器创建成功了

	// 如果是异步的，直接返回创建成功
	if !requestBody.Synchronized {
		ctx.JSON(http.StatusOK, gin.H{"id": id})
		return
	}

	// 如果是同步的，等待结果
	response := <-service.resultChan[id]
	ctx.JSON(http.StatusOK, response)
}

func (service *ContainerService) List(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, service.dataMap)
}

func (service *ContainerService) Delete(ctx *gin.Context) {
	// TODO 把容器进程直接杀掉
	ctx.JSON(http.StatusOK, "ok")
}

func (service *ContainerService) OnContainerRun(ctx *gin.Context) {
	id := ctx.Param("id")
	pid := ctx.Param("pid")

	// 这里直接返回 200
	ctx.JSON(http.StatusOK, nil)

	if pidInt, err := strconv.Atoi(pid); err != nil {
		log.Error(fmt.Sprintf("on container(%s) run format pid(%s) error, ", id, pid), err)
	} else {
		service.containerRunHandler(id, pidInt)
	}
}

func (service *ContainerService) OnContainerExit(ctx *gin.Context) {
	id := ctx.Param("id")
	data, err := ioutil.ReadAll(ctx.Request.Body)

	// 这里直接返回 200
	ctx.JSON(http.StatusOK, nil)

	if err != nil {
		log.Error("OnContainerExit read body error, ", err)
		service.containerExitHandler(id, "")
	} else {
		service.containerExitHandler(id, string(data))
	}
}

func (service *ContainerService) Dump(ctx *gin.Context) {
	if err := dumpDataToFile(service.dataMap, ContainerDataFileName); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	} else {
		ctx.String(http.StatusOK, "ok")
	}
}
