package service

import "engine/model"

type ZygoteService struct {
	engine       *Engine
	radixTreeMap map[string]*model.ZygoteProcessNode
}

func NewZygoteService(engine *Engine) *ZygoteService {
	service := &ZygoteService{
		engine:       engine,
		radixTreeMap: make(map[string]*model.ZygoteProcessNode),
	}

	for name, runtime := range service.engine.RuntimeService.dataMap {
		if runtime.CanZygote {
			if root, err := service.buildRadixTree(); err != nil {
				service.radixTreeMap[name] = root
			}
		}
	}

	return service
}

func (service *ZygoteService) buildRadixTree() (*model.ZygoteProcessNode, error) {
	return &model.ZygoteProcessNode{}, nil
}

func (service *ZygoteService) SearchMatchProcess(runtimeName string) *model.ZygoteProcessNode {
	return nil
}
