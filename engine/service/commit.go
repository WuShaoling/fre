package service

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/mydocker/core"
	"os/exec"
)

func Commit(containerName, imageName string) {
	mntURL := fmt.Sprintf(core.MntUrl, containerName)
	mntURL += "/"

	imageTar := core.RootUrl + "/" + imageName + ".tar"

	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntURL, ".").CombinedOutput(); err != nil {
		log.Errorf("Tar folder %s error %v", mntURL, err)
	}
}
