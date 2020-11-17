package main

import (
	"engine/api"
	"engine/config"
	"engine/service"
	"flag"
	"github.com/gin-gonic/gin"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "", "config path")
}

func main() {
	flag.Parse()
	config.InitSysConfig(configPath)
	gin.SetMode(gin.ReleaseMode)

	// new service
	service.NewRuntimeService()
	zygoteService := service.NewZygoteService()
	cgroupService := service.NewCgroupPoolService()
	functionService := service.NewFunctionService()
	containerService := service.NewContainerService(cgroupService, functionService, zygoteService)

	r := gin.Default()
	api.SetContainerRouter(functionService, containerService, r)

	_ = r.Run(":80")
}
