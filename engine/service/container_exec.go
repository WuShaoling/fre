package core

import (
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func Exec() error {
	// read command from pipe
	command, err := readFunctionContextFromPipe()
	if err != nil {
		return err
	}

	// chroot
	if err := setUpMount(); err != nil {
		return err
	}

	// 查找可执行文件
	path, err := exec.LookPath(command[0])
	if err != nil {
		log.Errorf("Exec loop path error %v", err)
		return err
	}
	log.Infof("Find path %s", path)

	if err := syscall.Exec(path, command[0:], os.Environ()); err != nil {
		log.Errorf("syscall.Exec failed, error=%v", err.Error())
		return err
	}
	return nil
}

func readFunctionContextFromPipe() ([]string, error) {
	pipe := os.NewFile(uintptr(3), "pipe")
	defer pipe.Close()

	data, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("exec read pipe error %v", err)
		return nil, err
	}

	command := string(data)

	separatorIndex := strings.Index(command, "|")
	entrypoint := command[:separatorIndex]
	entrypointParam := command[separatorIndex+1:]

	// ["python3", "bootstrap.py", "json化后的参数，可以包含空格"]
	return append(strings.Split(entrypoint, " "), entrypointParam), nil
}

func setUpMount() error {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("Get current location error %v", err)
		return err
	}
	log.Infof("Current location is %s", pwd)

	err = syscall.Chroot(pwd)
	if err != nil {
		log.Errorf("chroot to %s error %v", pwd, err)
	}
	return err
	//err = pivotRoot(pwd)

	//merge proc
	//defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	//syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	//syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}

//func pivotRoot(root string) error {
//	/**
//	  为了使当前root的老 root 和新 root 不在同一个文件系统下，我们把root重新mount了一次
//	  bind mount是把相同的内容换了一个挂载点的挂载方法
//	*/
//	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
//		return fmt.Errorf("Mount rootfs to itself error: %v", err)
//	}
//
//	// 创建 rootfs/.pivot_root 存储 old_root
//	pivotDir := filepath.Join(root, ".pivot_root")
//	if err := os.Mkdir(pivotDir, 0777); err != nil {
//		return err
//	}
//
//	// pivot_root 到新的rootfs, 现在老的 old_root 是挂载在rootfs/.pivot_root
//	// 挂载点现在依然可以在mount命令中看到
//	if err := syscall.PivotRoot(root, pivotDir); err != nil {
//		return fmt.Errorf("pivot_root %v", err)
//	}
//
//	// 修改当前的工作目录到根目录
//	if err := syscall.Chdir("/"); err != nil {
//		return fmt.Errorf("chdir / %v", err)
//	}
//
//	// umount rootfs/.pivot_root
//	pivotDir = filepath.Join("/", ".pivot_root")
//	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
//		return fmt.Errorf("unmount pivot_root dir %v", err)
//	}
//
//	// 删除临时文件夹
//	return os.Remove(pivotDir)
//}
