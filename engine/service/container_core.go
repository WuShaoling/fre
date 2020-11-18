package service

import (
	"encoding/json"
	"engine/config"
	"engine/model"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/ventu-io/go-shortid"
	"os"
	"os/exec"
	"path"
	"syscall"
	"time"
)

// 创建过程中，如果哪一步失败了，需要清理掉之前的结果
func (service *ContainerService) newContainer(requestBody *model.RunContainerRequestBody, template *model.Template) (id, errInfo string) {

	id = shortid.MustGenerate()
	errInfo = ""

	containerInfo := model.Container{
		Id:           id,
		Template:     requestBody.TemplateName,
		Synchronized: requestBody.Synchronized,
		CreateAt:     time.Now().UnixNano(),
	}

	// 1. 对于同步函数，注册回调chan
	if requestBody.Synchronized {
		service.resultChan[containerInfo.Id] = make(chan gin.H)
	}

	// 2. 创建容器根目录
	rootPath := path.Join(config.GetContainerFsPath(), id, "merge")
	err := service.newRootfs(rootPath)
	if err != nil {
		errInfo = "创建根目录失败"
		return
	}

	// 3. 获取 cgroup
	containerInfo.CgroupId = service.engine.CgroupService.Get(&template.ResourceLimit)
	if containerInfo.CgroupId == "" {
		errInfo = "获取Cgroup失败"
		return
	}

	// 4. 基于 zygote 创建或者直接启动容器
	zygoteProcess := service.engine.ZygoteService.SearchMatchProcess(template.Environment.Runtime)
	if zygoteProcess != nil {
		err = service.newContainerProcessByZygote()
		// TODO send message to zygoteProcess
	} else {
		err = service.newContainerProcessDirectly(&containerInfo, template, requestBody.Param, rootPath)
	}

	if err != nil {
		// 清理
		errInfo = err.Error()
	}

	service.dataMap[id] = containerInfo
	return
}

func (service *ContainerService) newContainerProcessByZygote() error {

}

func (service *ContainerService) newContainerProcessDirectly(containerInfo *model.Container, template *model.Template,
	param map[string]interface{}, rootPath string) error {

	// 1. new pipe
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		log.Error(fmt.Sprintf("newContainerProcessDirectly(id=%s) new pipe error, ", containerInfo.Id), err)
		return err
	}

	// 2. new process
	initCmd, err := os.Readlink("/proc/self/exe")
	if err != nil {
		log.Errorf("newContainerProcessDirectly get init process error %v", err)
		return err
	}
	cmd := exec.Command(initCmd, "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC,
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	cmd.Env = append(os.Environ(), template.Environment.Envs...)
	cmd.Dir = rootPath
	// TODO 重定向 Stdout Stderr 到日志文件

	// 3. start process
	if err := cmd.Start(); err != nil {
		log.Error(fmt.Sprintf("newContainerProcessDirectly(id=%s) start process error, ", containerInfo.Id), err)
		return err
	}

	// 4. 加入 cgroup
	if err := service.engine.CgroupService.Set(containerInfo.CgroupId, cmd.Process.Pid); err != nil {
		return err
	}

	// 5. 构造并发送 exec 参数
	commandParam := make(map[string]interface{})
	if err := service.sendInitCommand(commandParam, writePipe); err != nil {
		return err
	}

	return nil
}

func (service *ContainerService) sendInitCommand(containerInfo *model.Container, template *model.Template,
	param map[string]interface{}, writePipe *os.File) error {

	data, err := json.Marshal(commandParam)
	if err != nil {
		log.Error("sendInitCommand json.Marshal error, ", err)
		return err
	}

	defer writePipe.Close()
	if _, err := writePipe.Write(data); err != nil {
		log.Error("sendInitCommand failed,", err)
		return err
	}

	return nil
}

func (service *ContainerService) newRootfs(rootPath string) error {
	return nil
}
