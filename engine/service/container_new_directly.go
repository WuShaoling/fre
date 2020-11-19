package service

import (
	"encoding/json"
	"engine/config"
	"engine/model"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
)

// 直接启动容器进程

func (service *ContainerService) newContainerProcessDirectly(container *model.Container, template *model.Template, rootPath string) error {
	// 1. new pipe
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		log.Error(fmt.Sprintf("newContainerProcessDirectly(id=%s) new pipe error, ", container.Id), err)
		return err
	}
	defer writePipe.Close()

	// 2. new process
	initCmd, err := os.Readlink("/proc/self/exe")
	if err != nil {
		log.Errorf("newContainerProcessDirectly get init process error %v", err)
		return err
	}
	containerProcess := exec.Command(initCmd, "exec")
	containerProcess.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC,
	}
	containerProcess.ExtraFiles = []*os.File{readPipe}
	containerProcess.Env = append(os.Environ(), template.Envs...)
	containerProcess.Dir = path.Join(rootPath, "merge")
	// 临时使用std, TODO 重定向 Stdout Stderr 到日志文件
	//containerProcess.Stdin = os.Stdin
	containerProcess.Stdout = os.Stdout
	containerProcess.Stderr = os.Stderr
	if err := containerProcess.Start(); err != nil {
		log.Error(fmt.Sprintf("newContainerProcessDirectly(id=%s) start process error, ", container.Id), err)
		service.onNewContainerProcessDirectlyError(readPipe, nil)
		return err
	}

	// 3. 加入 cgroup
	if err := service.engine.CgroupService.Set(container.CgroupId, containerProcess.Process.Pid); err != nil {
		service.onNewContainerProcessDirectlyError(readPipe, containerProcess)
		return err
	}

	// 4. 发送运行命令
	entrypoint := service.engine.RuntimeService.getRuntimeByName(container.Template).Command
	entrypointParam := service.buildCommandParams(container, template, false)
	command := strings.Join(entrypoint, " ") + "|" + entrypointParam // python3 bootstrap.py|param str, |为解析命令的分隔符

	if _, err := writePipe.Write([]byte(command)); err != nil {
		log.Errorf("sendInitCommand failed, command=%s, error=%v", command, err)
		service.onNewContainerProcessDirectlyError(readPipe, containerProcess)
		return err
	}

	// 5. 设置容器为运行状态
	service.containerRunHandler(container.Id, containerProcess.Process.Pid)

	// 6. 异步 wait 容器进程退出
	go func() {
		if err := containerProcess.Wait(); err != nil {
			log.Error(fmt.Sprintf("wait error, id=%s, pid=%d", container.Id, containerProcess.Process.Pid), err)
		} else {
			service.containerExitHandler(container.Id, "container exit")
		}
	}()

	return nil
}

func (service *ContainerService) onNewContainerProcessDirectlyError(readPipe *os.File, containerProcess *exec.Cmd) {
	if readPipe != nil {
		_ = readPipe.Close()
	}

	if containerProcess != nil && containerProcess.Process != nil {
		_ = containerProcess.Process.Kill()
		_ = containerProcess.Wait()
	}
}

func (service *ContainerService) buildCommandParams(container *model.Container, template *model.Template, zygote bool) string {
	param := map[string]interface{}{
		"id":            container.Id,
		"functionPath":  path.Join("/code", template.Name),
		"handler":       template.Handler,
		"functionParam": container.FunctionParam,
		"server":        "http://localhost:80/",
	}

	if zygote {
		param["rootPath"] = path.Join(config.GetContainerFsPath(), container.Id, "mount")
		param["cgroupFiles"] = []string{} // TODO 附加参数
	}

	data, _ := json.Marshal(param)
	return string(data)
}
