package model

type Runtime struct { // 基础环境
	Name    string `json:"name" binding:"required"`
	Command string `json:"command" binding:"required"`
}
