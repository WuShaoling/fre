package service

import (
	"engine/config"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"syscall"
)

func (service *ContainerService) newRootFS(id, runtime string) (string, error) {
	basePath := path.Join(config.GetContainerFsPath(), id)
	lowerPath := path.Join(config.GetRuntimePath(), runtime)
	mergePath := path.Join(basePath, "merge")
	upperPath := path.Join(basePath, "upper")
	workerPath := path.Join(basePath, "worker")

	// 创建目录
	if err := os.MkdirAll(mergePath, 0755); err != nil {
		log.Error(fmt.Sprintf("new mergePath %s for container %s error,", mergePath, id), err)
		return basePath, err
	}
	if err := os.MkdirAll(upperPath, 0755); err != nil {
		log.Error(fmt.Sprintf("new upperPath %s for container %s error,", upperPath, id), err)
		return basePath, err
	}
	if err := os.MkdirAll(workerPath, 0755); err != nil {
		log.Error(fmt.Sprintf("new workerPath %s for container %s error,", workerPath, id), err)
		return basePath, err
	}

	data := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerPath, upperPath, workerPath)
	if err := syscall.Mount("overlay", mergePath, "overlay", 0, data); err != nil {
		log.Error(fmt.Sprintf("overlay mount for container %s error", id), err)
		return basePath, err
	}

	return basePath, nil
}

func (service *ContainerService) cleanRootFS(rootPath string) {
	_ = syscall.Unmount(path.Join(rootPath, "merge"), 0)
	_ = os.RemoveAll(rootPath)
}
