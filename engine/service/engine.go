package service

type Engine struct {
	RuntimeService    *RuntimeService
	TemplateService   *TemplateService
	ContainerService  *ContainerService
	ZygoteService     *ZygoteService
	CgroupPoolService *CgroupPoolService
}

func NewEngine() *Engine {

	engine := &Engine{}

	engine.RuntimeService = NewRuntimeService(engine)
	engine.CgroupPoolService = NewCgroupPoolService()
	engine.TemplateService = NewTemplateService(engine)
	engine.ZygoteService = NewZygoteService(engine)
	engine.ContainerService = NewContainerService(engine)

	return engine
}
