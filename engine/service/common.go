package service

import (
	"encoding/json"
	"engine/config"
	"engine/util"
	"fmt"
	log "github.com/sirupsen/logrus"
	"path"
)

const (
	RuntimeDataFileName   = "runtime.json"
	TemplateDataFileName  = "template.json"
	ContainerDataFileName = "container.json"
)

const (
	Running = "running"
	Exit    = "exit"
)

func dumpDataToFile(v interface{}, filename string) error {
	data, err := json.Marshal(v)
	if err != nil {
		log.Error(fmt.Sprintf("json.Marshal error, filename=%s, data=%+v, error=%+v", filename, v, err))
	}
	return util.WriteToFile(filename, data)
}

func loadDataFromFile(v interface{}, filename string) {
	if data, err := util.ReadFromFile(getDataFilePath(filename)); err == nil {
		_ = json.Unmarshal(data, v)
	}
}

func getDataFilePath(filename string) string {
	return path.Join(config.GetDataPath(), filename)
}
