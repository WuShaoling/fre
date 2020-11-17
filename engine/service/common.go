package service

import (
	"encoding/json"
	"engine/config"
	"engine/util"
	"path"
)

const (
	RuntimeDataFileName   = "runtime.json"
	FunctionDataFileName  = "function.json"
	ContainerDataFileName = "container.json"
)

func loadDataFromFile(v interface{}, filename string) {
	if data, err := util.ReadFromFile(getDataFilePath(filename)); err == nil {
		_ = json.Unmarshal(data, v)
	}
}

func getDataFilePath(filename string) string {
	return path.Join(config.GetDataPath(), filename)
}
