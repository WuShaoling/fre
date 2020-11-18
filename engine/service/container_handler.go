package service

import (
	"engine/model"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

// 容器启动成功
func (service *ContainerService) containerRunHandler(id string, pid int) {
	log.Info(fmt.Sprintf("OnContainerRun id=%s, pid=%d", id, pid))

	// 检查基本信息
	containerInfo := service.getAndCheckContainer(id)
	if containerInfo == nil {
		return
	}

	// 修改状态
	containerInfo.Status = Running
	containerInfo.RunAt = time.Now().UnixNano()
	containerInfo.Pid = pid
}

// 容器退出
// 直接创建的容器，wait 返回结果时调用此方法
// zygote创建的容器，zygote父进程wait返回后通过http调用此方法
func (service *ContainerService) containerExitHandler(id string, status string) {
	log.Info(fmt.Sprintf("ContainerService id=%s, status=%s", id, status))

	// 检查基本信息
	containerInfo := service.getAndCheckContainer(id)
	if containerInfo == nil {
		return
	}

	// 修改状态
	containerInfo.Status = Exit
	containerInfo.EndAt = time.Now().UnixNano()

	// 如果是同步的请求，返回数据
	if containerInfo.Synchronized {
		result := gin.H{"id": id}
		if data, ok := service.resultData[id]; !ok { // 数据不存在，说明函数非正常退出
			result["error"] = status
		} else {
			result["data"] = data
		}

		if dataChan, ok := service.resultChan[id]; ok {
			dataChan <- result
			close(service.resultChan[id])
			delete(service.resultChan, id)
		} else {
			log.Warn(fmt.Sprintf("error chan not found, container id=%s, status=%s", id, status))
		}
	}

	// TODO 执行清理工作
}

func (service *ContainerService) getAndCheckContainer(id string) *model.Container {
	containerInfo, ok := service.dataMap[id]
	if !ok {
		log.Warn(fmt.Sprintf("container(%s) not found", id))
		return nil
	} else if containerInfo.Status == Exit {
		log.Warn(fmt.Sprintf("container(%s) already exit", id))
		return nil
	}
	return &containerInfo
}
