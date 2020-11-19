package main

import (
	"engine/api"
	"engine/config"
	"engine/core"
	"engine/service"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "", "config path")
}

func main() {
	// init
	flag.Parse()

	if len(os.Args) >= 2 && os.Args[1] == "exec" { // 容器进程
		_ = core.Exec()
		return
	}

	config.InitSysConfig(configPath)
	gin.SetMode(gin.ReleaseMode)

	freEngine := service.NewEngine()
	r := gin.Default()
	api.SetContainerRouter(freEngine.RuntimeService, freEngine.TemplateService, freEngine.ContainerService, r)

	log.Println("listen on :80")
	_ = r.Run(":80")
}
