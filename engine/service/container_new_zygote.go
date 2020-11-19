package service

import (
	"engine/model"
)

// 通过 zygote 启动容器进程

func (service *ContainerService) newContainerProcessByZygote(container *model.Container, template *model.Template,
	rootPath string, zygoteProcess *model.ZygoteProcessNode) error {
	// TODO
	return nil
}

func (service *ContainerService) onNewContainerProcessByZygoteError() {
	// TODO
}
