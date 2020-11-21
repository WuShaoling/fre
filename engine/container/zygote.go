package container

import (
	"engine/runtime"
)

type ZygoteProcessNode struct {
	Id     int
	Socket int // socket
}

type ZygoteService struct {
	radixTreeMap map[string]*ZygoteProcessNode
}

func NewZygoteService(runtimeList map[string]runtime.Runtime) *ZygoteService {
	service := &ZygoteService{
		radixTreeMap: make(map[string]*ZygoteProcessNode),
	}
	for name, r := range runtimeList {
		if r.CanZygote {
			if root, err := service.buildRadixTree(); err != nil {
				service.radixTreeMap[name] = root
			}
		}
	}
	return service
}

func (service *ZygoteService) buildRadixTree() (*ZygoteProcessNode, error) {
	return &ZygoteProcessNode{}, nil
}

func (service *ZygoteService) SearchMatchProcess(runtimeName string) *ZygoteProcessNode {
	return nil
}
