package subsystems

import(
	"testing"
)

func TestFindCgroupMountpoint(t *testing.T) {
	t.Logf("cpu subsystem merge point %v\n", FindCgroupMountpoint("cpu"))
	t.Logf("cpuset subsystem merge point %v\n", FindCgroupMountpoint("cpuset"))
	t.Logf("memory subsystem merge point %v\n", FindCgroupMountpoint("memory"))
}