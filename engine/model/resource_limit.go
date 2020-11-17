package model

type ResourceLimit struct { // 资源限制
	Memory   int64  `json:"memory"`
	CpuShare uint64 `json:"cpuShares"`
	Timeout  int    `json:"timeout"` // 超时时间
}
