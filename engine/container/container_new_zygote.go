package container

import (
	"engine/runtime"
	"engine/template"
)

// 通过 zygote 启动容器进程

func (service *Service) newContainerProcessByZygote(runtime runtime.Runtime, template template.Template,
	container *Container, zygoteProcess *ZygoteProcessNode) error {
	// TODO
	return nil
}

func (service *Service) onNewContainerProcessByZygoteError() {
	// TODO
}
