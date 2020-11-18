package model

type Container struct { // 容器实例
	Id       string `json:"id"`       // 唯一标识
	Pid      int    `json:"pid"`      // 进程ID
	Template string `json:"template"` // 函数模版
	CgroupId string `json:"cgroupId"` // cgroupId

	Status       string `json:"status"`       // 状态
	Synchronized bool   `json:"synchronized"` // 等待结果

	CreateAt int64 `json:"createAt"` // 启动创建时间
	RunAt    int64 `json:"runAt"`    // 函数开始运行时间
	EndAt    int64 `json:"endAt"`    // 运行结束时间
}
