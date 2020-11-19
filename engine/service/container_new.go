package service

import (
	"engine/model"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ventu-io/go-shortid"
	"time"
)

// 创建过程中，如果哪一步失败了，需要清理掉之前的结果
func (service *ContainerService) newContainer(template *model.Template, requestBody *model.RunContainerRequestBody) (string, error) {
	id := shortid.MustGenerate()

	// 1. 对于同步函数，注册回调chan
	if requestBody.Synchronized {
		service.resultChan[id] = make(chan gin.H)
	}

	// 2. 获取 cgroup
	cgroupId, err := service.engine.CgroupService.Get(&template.ResourceLimit)
	if cgroupId == "" {
		service.onNewContainerError(id, "", "")
		return "", errors.New("获取Cgroup失败")
	}

	// 3. new root fs
	rootPath, err := service.newRootFS(id, template.Runtime)
	if err != nil {
		service.onNewContainerError(id, rootPath, "")
		return "", errors.New("创建根目录失败")
	}

	container := &model.Container{
		Id:            id,
		Template:      template.Name,
		FunctionParam: requestBody.FunctionParam,
		Synchronized:  requestBody.Synchronized,
		CreateAt:      time.Now().UnixNano(),
	}
	service.dataMap[container.Id] = container

	// 3. 基于 zygote 创建或者直接启动容器
	zygoteProcess := service.engine.ZygoteService.SearchMatchProcess(template.Runtime)
	if zygoteProcess != nil {
		err = service.newContainerProcessByZygote(container, template, rootPath, zygoteProcess)
	} else {
		err = service.newContainerProcessDirectly(container, template, rootPath)
	}

	if err != nil {
		service.onNewContainerError(id, rootPath, cgroupId)
		return "", err
	}

	return id, nil
}

func (service *ContainerService) onNewContainerError(id, rootPath, cgroupId string) {
	close(service.resultChan[id])
	delete(service.resultChan, id)

	if rootPath != "" {
		service.cleanRootFS(rootPath)
	}
	if cgroupId != "" {
		service.engine.CgroupService.GiveBack(cgroupId)
	}

	delete(service.dataMap, id)
}
