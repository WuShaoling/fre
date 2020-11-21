package container

import (
	"engine/runtime"
	"engine/template"
	"errors"
)

// 通过 zygote 启动容器进程

func (service *Service) newContainerProcessByZygote(runtime runtime.Runtime, template template.Template, container *Container) error {
	// 构建参数
	functionExecContext := service.buildFunctionExecContext(template, container)
	if functionExecContext == "" {
		return errors.New("BuildFunctionExecContextError")
	}

	return service.zygoteService.NewContainerByZygoteProcess(runtime.Name, template.Name, functionExecContext)
}
