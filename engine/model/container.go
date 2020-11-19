package model

type Container struct { // 容器实例
	Id       string `json:"id"`       // 唯一标识
	Template string `json:"template"` // 函数模版
	Pid      int    `json:"pid"`      // 进程ID
	CgroupId string `json:"cgroupId"` // cgroupId

	FunctionParam map[string]interface{} `json:"functionParam"` // 运行参数
	Synchronized  bool                   `json:"synchronized"`  // 同步等待结果

	Status   string `json:"status"`   // 状态
	CreateAt int64  `json:"createAt"` // 启动创建时间
	RunAt    int64  `json:"runAt"`    // 函数开始运行时间
	EndAt    int64  `json:"endAt"`    // 运行结束时间
}
