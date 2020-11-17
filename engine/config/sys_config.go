package config

import "github.com/gin-gonic/gin"

type sysConfig struct {
}

var SysConfigInstance *sysConfig

func InitSysConfig() {

	gin.SetMode(gin.ReleaseMode)

	SysConfigInstance = &sysConfig{

	}
}
