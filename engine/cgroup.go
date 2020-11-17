package main

import (
	"fmt"
	"github.com/containerd/cgroups"
	"github.com/opencontainers/runtime-spec/specs-go"
	"log"
)

func main() {
	shares := uint64(100)
	control, err := cgroups.New(
		cgroups.V1,
		cgroups.StaticPath("/test1"),
		&specs.LinuxResources{CPU: &specs.LinuxCPU{Shares: &shares}},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(control)
	//fmt.Println(control.Delete())
}
