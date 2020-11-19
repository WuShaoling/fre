package service

import (
	"engine/config"
	"engine/model"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"path"
	"time"
)

// 容器启动成功
func (service *ContainerService) containerRunHandler(id string, pid int) {
	log.Infof("on container run: id=%s, pid=%d", id, pid)

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
	log.Infof("on container exit: id=%s, status=%s", id, status)

	// 检查基本信息
	container := service.getAndCheckContainer(id)
	if container == nil {
		return
	}

	// 修改状态
	container.Status = Exit
	container.EndAt = time.Now().UnixNano()

	// 如果是同步的请求，返回数据
	if container.Synchronized {
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
			delete(service.resultData, id)
		} else {
			log.Warn(fmt.Sprintf("error chan not found, container id=%s, status=%s", id, status))
		}
	}

	// 清除 rootfs
	service.cleanRootFS(path.Join(config.GetContainerFsPath(), id))

	// 归还 cgroup
	service.engine.CgroupService.GiveBack(container.CgroupId)
}

func (service *ContainerService) functionResultHandler(data map[string]interface{}) {
	id := data["id"].(string)
	log.Infof("on function result: id=%s, data=%s", data["id"])

	container := service.getAndCheckContainer(id)
	if container != nil && container.Synchronized {
		if execError, ok := data["error"]; ok { // 执行出错
			service.resultData[id] = gin.H{"error": execError}
		} else {
			service.resultData[id] = gin.H{"data": data["containerResult"]}
		}
	}
}

func (service *ContainerService) getAndCheckContainer(id string) *model.Container {
	container, ok := service.dataMap[id]
	if !ok {
		log.Warnf("container(id=%s) not found", id)
		return nil
	} else if container.Status == Exit {
		log.Warnf("container(id=%s) already exit", id)
		return nil
	}
	return container
}
