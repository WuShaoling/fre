package model

type Function struct { // 函数模版
	Runtime       string            `json:"runtime" binding:"required"` // 基础环境
	Name          string            `json:"name" binding:"required"`    // 函数名, 唯一
	Handler       string            `json:"handler" binding:"required"` // 函数入口文件
	Volume        string            `json:"volume"`                     // 数据卷挂载目录
	Envs          map[string]string `json:"envs"`                       // 环境变量
	ResourceLimit ResourceLimit     `json:"resourceLimit"`              // 资源限制
}
