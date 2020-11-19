package service
//
//import (
//	"encoding/json"
//	"engine/model"
//	"errors"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	log "github.com/sirupsen/logrus"
//	"github.com/ventu-io/go-shortid"
//	"net/url"
//	"os"
//	"os/exec"
//	"path"
//	"strings"
//	"syscall"
//	"time"
//)
//
//// 创建过程中，如果哪一步失败了，需要清理掉之前的结果
//func (service *ContainerService) newContainer(template *model.Template, requestBody *model.RunContainerRequestBody) (string, error) {
//	id := shortid.MustGenerate()
//
//	// 1. 对于同步函数，注册回调chan
//	if requestBody.Synchronized {
//		service.resultChan[id] = make(chan gin.H)
//	}
//
//	// 2. 获取 cgroup
//	cgroupId, err := service.engine.CgroupService.Get(&template.ResourceLimit)
//	if cgroupId == "" {
//		service.onNewContainerError(id, "", "")
//		return "", errors.New("获取Cgroup失败")
//	}
//
//	// 3. new root fs
//	rootPath, err := service.newRootFS(id, template.Runtime)
//	if err != nil {
//		service.onNewContainerError(id, rootPath, "")
//		return "", errors.New("创建根目录失败")
//	}
//
//	container := &model.Container{
//		Id:            shortid.MustGenerate(),
//		Template:      template.Name,
//		FunctionParam: requestBody.FunctionParam,
//		Synchronized:  requestBody.Synchronized,
//		CreateAt:      time.Now().UnixNano(),
//	}
//	service.dataMap[container.Id] = container
//
//	// 3. 基于 zygote 创建或者直接启动容器
//	zygoteProcess := service.engine.ZygoteService.SearchMatchProcess(template.Runtime)
//	if zygoteProcess != nil {
//		err = service.newContainerProcessByZygote(container, template, rootPath, zygoteProcess)
//		// TODO send message to zygoteProcess
//	} else {
//		err = service.newContainerProcessDirectly(container, template, rootPath)
//	}
//
//	if err != nil {
//		service.onNewContainerError(id, rootPath, cgroupId)
//		return "", err
//	}
//
//	return id, nil
//}
//
//func (service *ContainerService) newContainerProcessByZygote(container *model.Container, template *model.Template,
//	rootPath string, zygoteProcess *model.ZygoteProcessNode) error {
//	// TODO
//	return nil
//}
//
//func (service *ContainerService) newContainerProcessDirectly(container *model.Container, template *model.Template, rootPath string) error {
//	// 1. new pipe
//	readPipe, writePipe, err := os.Pipe()
//	if err != nil {
//		log.Error(fmt.Sprintf("newContainerProcessDirectly(id=%s) new pipe error, ", container.Id), err)
//		return err
//	}
//	defer writePipe.Close()
//
//	// 2. new process
//	initCmd, err := os.Readlink("/proc/self/exe")
//	if err != nil {
//		log.Errorf("newContainerProcessDirectly get init process error %v", err)
//		return err
//	}
//	containerProcess := exec.Command(initCmd, "exec")
//	containerProcess.SysProcAttr = &syscall.SysProcAttr{
//		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC,
//	}
//	containerProcess.ExtraFiles = []*os.File{readPipe}
//	containerProcess.Env = append(os.Environ(), template.Envs...)
//	containerProcess.Dir = path.Join(rootPath, "merge")
//	// 临时使用std, TODO 重定向 Stdout Stderr 到日志文件
//	//containerProcess.Stdin = os.Stdin
//	containerProcess.Stdout = os.Stdout
//	containerProcess.Stderr = os.Stderr
//	if err := containerProcess.Start(); err != nil {
//		log.Error(fmt.Sprintf("newContainerProcessDirectly(id=%s) start process error, ", container.Id), err)
//		service.onNewContainerProcessDirectlyError(readPipe, nil)
//		return err
//	}
//
//	// 3. 加入 cgroup
//	if err := service.engine.CgroupService.Set(container.CgroupId, containerProcess.Process.Pid); err != nil {
//		service.onNewContainerProcessDirectlyError(readPipe, containerProcess)
//		return err
//	}
//
//	// 4. 发送运行命令
//	ctx := url.Values{}
//	ctx.Add("id", container.Id)
//	ctx.Add("functionPath", path.Join("/code", template.Name))
//	ctx.Add("handler", template.Handler)
//	functionParam, _ := json.Marshal(container.FunctionParam)
//	ctx.Add("functionParam", string(functionParam))
//
//	command := service.engine.RuntimeService.getRuntimeByName(container.Template).Command
//	command = append(command, ctx.Encode())
//
//	if _, err := writePipe.Write([]byte(strings.Join(command, " "))); err != nil {
//		log.Error("sendInitCommand failed, ", err)
//		service.onNewContainerProcessDirectlyError(readPipe, containerProcess)
//		return err
//	}
//
//	// 5. 设置容器为运行状态
//	service.containerRunHandler(container.Id, containerProcess.Process.Pid)
//
//	// 6. 异步 wait 容器进程退出
//	go func() {
//		if err := containerProcess.Wait(); err != nil {
//			log.Error(fmt.Sprintf("wait error, id=%s, pid=%d", container.Id, containerProcess.Process.Pid), err)
//		} else {
//			service.containerExitHandler(container.Id, "container exit")
//		}
//	}()
//
//	return nil
//}
//
//func (service *ContainerService) onNewContainerError(id, rootPath, cgroupId string) {
//	close(service.resultChan[id])
//	delete(service.resultChan, id)
//
//	if rootPath != "" {
//		service.cleanRootFS(rootPath)
//	}
//	if cgroupId != "" {
//		service.engine.CgroupService.GiveBack(cgroupId)
//	}
//
//	delete(service.dataMap, id)
//}
//
//func (service *ContainerService) onNewContainerProcessDirectlyError(readPipe *os.File, containerProcess *exec.Cmd) {
//	if readPipe != nil {
//		_ = readPipe.Close()
//	}
//
//	if containerProcess != nil && containerProcess.Process != nil {
//		_ = containerProcess.Process.Kill()
//		_ = containerProcess.Wait()
//	}
//}
//
//func (service *ContainerService) onNewContainerProcessByZygoteError() {
//	// TODO
//}
