package config

import (
	"engine/util"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
)

type SysConfig struct {
	RootPath        string `yaml:"rootPath"`        // 应用根目录
	ZygoteMaxMemory int    `yaml:"zygoteMaxMemory"` // zygote 池最大内存，MB
	CgroupPoolSize  int    `yaml:"cgroupPoolSize"`  // cgroup 缓存池大小
}

const (
	LogPath         = "log"       // 日志根目录
	DataPath        = "data"      // 数据目录
	RuntimePath     = "runtime"   // 运行时环境根目录
	VolumeHostPath  = "volume"    // 数据卷主机端根目录
	ContainerFsPath = "container" // 容器文件系统根目录

	RuntimeFunctionPath = "/code" // 运行时环境里面代码根目录
)

var SysConfigInstance *SysConfig

func InitSysConfig(configPath string) {
	if configPath != "" { // use config from file
		log.Println("[info] load config from", configPath)
		SysConfigInstance = &SysConfig{}

		f, err := os.Open(configPath)
		if err != nil {
			log.Fatal("[error] init error, ", err)
		}

		if err := yaml.NewDecoder(f).Decode(SysConfigInstance); err != nil {
			log.Fatal("[error] load yaml config error, ", err)
		}
	} else { // use default config
		log.Println("[info] use default config")
		SysConfigInstance = &SysConfig{
			RootPath:        "./fre",
			ZygoteMaxMemory: 1024,
			CgroupPoolSize:  2,
		}
	}

	if err := util.MkdirIfNotExist(GetLogPath()); err != nil {
		log.Fatal("[error]", err)
	}
	if err := util.MkdirIfNotExist(GetDataPath()); err != nil {
		log.Fatal("[error]", err)
	}
	if err := util.MkdirIfNotExist(GetRuntimePath()); err != nil {
		log.Fatal("[error]", err)
	}
	if err := util.MkdirIfNotExist(GetVolumeHostPath()); err != nil {
		log.Fatal("[error]", err)
	}
	if err := util.MkdirIfNotExist(GetContainerFsPath()); err != nil {
		log.Fatal("[error]", err)
	}
}

func GetLogPath() string         { return path.Join(SysConfigInstance.RootPath, LogPath) }
func GetDataPath() string        { return path.Join(SysConfigInstance.RootPath, DataPath) }
func GetRuntimePath() string     { return path.Join(SysConfigInstance.RootPath, RuntimePath) }
func GetVolumeHostPath() string  { return path.Join(SysConfigInstance.RootPath, VolumeHostPath) }
func GetContainerFsPath() string { return path.Join(SysConfigInstance.RootPath, ContainerFsPath) }
