package model

type Container struct { // 容器实例
	Id      string `json:"id"`      // 唯一标识
	StartAt int64  `json:"startAt"` // 启动容器时间
	RunAt   int64  `json:"runAt"`   // 函数开始运行时间
	EndAt   int64  `json:"endAt"`   // 运行结束时间
}
