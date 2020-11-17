package model

type RunConfig struct {
	Image         string
	Cmd           []string
	IT            bool // tty, interactive
	Detach        bool
	ResourceLimit *CGroupConfig
	Volume        map[string]string
	Envs          map[string]string
}

type CGroupConfig struct {
	Memory   string
	CpuShare string
	CpuSet   string
}
