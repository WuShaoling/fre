package service

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/mydocker/core"
	"os"
)

func Remove(containerName string) {
	containerInfo, err := getContainerInfoByName(containerName)
	if err != nil {
		log.Errorf("Get container %s info error %v", containerName, err)
		return
	}
	if containerInfo.Status != core.Stop {
		log.Errorf("Couldn't remove running container")
		return
	}
	dirURL := fmt.Sprintf(core.DefaultInfoLocation, containerName)
	if err := os.RemoveAll(dirURL); err != nil {
		log.Errorf("Remove file %s error %v", dirURL, err)
		return
	}
	core.DeleteWorkSpace(containerInfo.Volume, containerName)
}
