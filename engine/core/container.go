package core

import (
	"engine/model"
	"engine/service"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// 只有创建成功才加入 map
func NewContainer(engine *service.Engine, template *model.Template, param map[string]interface{}) int {


}

//const (
//	Running string = "running"
//	Stop    string = "stopped"
//	Exit    string = "exited"
//
//	DefaultInfoLocation string = "/var/run/mydocker/%s/"
//	ConfigName          string = "config.json"
//	ContainerLogFile    string = "container.log"
//	RootUrl             string = "/root"
//	MntUrl              string = "/root/mnt/%s"
//	WriteLayerUrl       string = "/root/writeLayer/%s"
//)
//
//type ContainerInfo struct {
//	Pid         string `json:"pid"`        // 容器的init进程在宿主机上的 PID
//	Id          string `json:"id"`         // 容器Id
//	Name        string `json:"name"`       // 容器名
//	Command     string `json:"command"`    // 容器内init运行命令
//	CreatedTime string `json:"createTime"` // 创建时间
//	Status      string `json:"status"`     // 容器的状态
//	Volume      string `json:"volume"`     // 容器的数据卷
//}

