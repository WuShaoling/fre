package service

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/mydocker/core"
	"io/ioutil"
	"os"
)

func Log(containerName string) {
	dirURL := fmt.Sprintf(core.DefaultInfoLocation, containerName)
	logFileLocation := dirURL + core.ContainerLogFile
	file, err := os.Open(logFileLocation)
	defer file.Close()
	if err != nil {
		log.Errorf("Log container open file %s error %v", logFileLocation, err)
		return
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Errorf("Log container read file %s error %v", logFileLocation, err)
		return
	}
	fmt.Fprint(os.Stdout, string(content))
}
