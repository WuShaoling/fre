package service

type Engine struct {
	RuntimeService   *RuntimeService
	TemplateService  *TemplateService
	ContainerService *ContainerService
	ZygoteService    *ZygoteService
	CgroupService    *CgroupService
}

func NewEngine() *Engine {

	engine := &Engine{}

	engine.CgroupService = NewCgroupService()
	engine.RuntimeService = NewRuntimeService(engine)
	engine.TemplateService = NewTemplateService(engine)
	engine.ZygoteService = NewZygoteService(engine)
	engine.ContainerService = NewContainerService(engine)

	return engine
}
