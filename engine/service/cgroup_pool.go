package service

import (
	"engine/config"
	"engine/model"
	"github.com/containerd/cgroups"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/ventu-io/go-shortid"
	"log"
	"sync"
)

const cgroupPrefix = "/fre_"

type CgroupPoolService struct {
	pool  []*CgroupInfo
	mutex sync.Mutex
}

type CgroupInfo struct {
	Id     string
	Cgroup *cgroups.Cgroup
}

func NewCgroupPoolService() *CgroupPoolService {
	service := &CgroupPoolService{
		pool: make([]*CgroupInfo, 0, config.SysConfigInstance.CgroupPoolSize),
	}

	for i := 0; i < config.SysConfigInstance.CgroupPoolSize; i++ {
		cgroupInfo := service.newCgroupInfo(nil)
		if cgroupInfo == nil {
			log.Fatal("[error] NewCgroupPoolService error")
		}
		service.pool = append(service.pool, cgroupInfo)
	}

	return service
}

func (service *CgroupPoolService) newCgroupInfo(limit *model.ResourceLimit) *CgroupInfo {
	id := shortid.MustGenerate()

	linuxResource := &specs.LinuxResources{
		Memory: &specs.LinuxMemory{},
		CPU:    &specs.LinuxCPU{},
	}
	if limit != nil {
		linuxResource.Memory.Limit = &limit.Memory
		linuxResource.CPU.Shares = &limit.CpuShare
	}

	cgroup, err := cgroups.New(cgroups.V1, cgroups.StaticPath(cgroupPrefix+id), linuxResource)
	if err != nil {
		log.Println("[error] NewCgroupPoolService error, ", err)
		return nil
	}

	return &CgroupInfo{
		Id:     id,
		Cgroup: &cgroup,
	}
}

func (service *CgroupPoolService) Get(limit *model.ResourceLimit) *CgroupInfo {
	service.mutex.Lock()

	n := len(service.pool)

	if n > 0 {
		c := service.pool[n-1]
		service.pool = service.pool[0 : n-1]
		service.mutex.Unlock()
		return c
	}

	service.mutex.Unlock()
	return service.newCgroupInfo(limit)
}

func (service *CgroupPoolService) GiveBack(id string) {
	cgroup, err := cgroups.Load(cgroups.V1, cgroups.StaticPath(cgroupPrefix+id))
	if err != nil {
		log.Println("[error] cgroups.Load error: ", err)
		return
	}

	if len(service.pool) < config.SysConfigInstance.CgroupPoolSize {
		service.mutex.Lock()
		service.pool = append(service.pool, &CgroupInfo{
			Id:     id,
			Cgroup: &cgroup,
		})
		service.mutex.Unlock()
	} else {
		_ = cgroup.Delete()
	}
}
