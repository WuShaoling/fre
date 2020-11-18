package service

type ZygoteService struct {
	engine *Engine
}

func NewZygoteService(engine *Engine) *ZygoteService {
	return &ZygoteService{
		engine: engine,
	}
}
