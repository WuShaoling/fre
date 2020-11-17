package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mydocker/api"
	"github.com/mydocker/config"
)

func main() {

	config.InitSysConfig()

	r := gin.Default()
	api.SetContainerRouter(r)

	_ = r.Run(":8080")

}
