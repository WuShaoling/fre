package service

import (
	"engine/model"
)

type RuntimeService struct {
	dataMap map[string]model.Runtime
}

func NewRuntimeService() *RuntimeService {
	service := &RuntimeService{
		dataMap: make(map[string]model.Runtime),
	}
	loadDataFromFile(service.dataMap, RuntimeDataFileName)
	return service
}
