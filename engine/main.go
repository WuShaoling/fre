package main

import (
	"engine/api"
	"engine/config"
	"engine/container"
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
		_ = container.Exec()
		return
	}

	config.InitSysConfig(configPath)
	gin.SetMode(gin.ReleaseMode)

	freEngine := service.NewEngine()
	r := gin.Default()
	api.SetContainerRouter(freEngine, r)

	log.Println("listen on :" + config.SysConfigInstance.ServePort)
	_ = r.Run(":" + config.SysConfigInstance.ServePort)
}
