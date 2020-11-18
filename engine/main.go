package main

import (
	"engine/api"
	"engine/config"
	"engine/service"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "", "config path")
}

func main() {
	// init
	flag.Parse()
	config.InitSysConfig(configPath)
	gin.SetMode(gin.ReleaseMode)

	freEngine := service.NewEngine()
	r := gin.Default()
	api.SetContainerRouter(freEngine.RuntimeService, freEngine.TemplateService, freEngine.ContainerService, r)

	log.Println("listen on :80")
	_ = r.Run(":80")
}
