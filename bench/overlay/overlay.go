package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"syscall"
	"time"
)

var parallel bool
var count int

func init() {
	flag.IntVar(&count, "n", 1, "启动的容器数量")
	flag.BoolVar(&parallel, "p", false, "并发启动")
}

func mount(id int) {
	t1 := time.Now().UnixNano()

	// make home path
	homePath := fmt.Sprintf("container_%d", id)
	if err := os.Mkdir(homePath, 0777); err != nil {
		log.Println("mkdir "+homePath, err)
	}

	// make upper path
	upperPath := homePath + "/upper"
	if err := os.Mkdir(upperPath, 0777); err != nil {
		log.Println("mkdir "+upperPath, err)
	}

	// make worker path
	workerPath := homePath + "/worker"
	if err := os.Mkdir(workerPath, 0777); err != nil {
		log.Println("mkdir "+workerPath, err)
	}

	// make mount path
	mountPath := homePath + "/merge"
	if err := os.Mkdir(mountPath, 0777); err != nil {
		log.Println("mkdir "+mountPath, err)
	}

	t2 := time.Now().UnixNano()
	data := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", "rootfs", upperPath, workerPath)
	fmt.Println(data)
	fmt.Println(mountPath)
	if err := syscall.Mount("overlay", mountPath, "overlay", 0, data); err != nil {
		log.Println("syscall.Mount", err)
	}

	t3 := time.Now().UnixNano()
	fmt.Printf("%d, %d, %d, %d\n", id, (t2-t1)/1e3, (t3-t2)/1e3, (t3-t1)/1e3)
}

func clean() {
	for i := 0; i < count; i++ {
		homePath := fmt.Sprintf("container_%d", i)
		mountPath := homePath + "/merge"
		if err := syscall.Unmount(mountPath, 0); err != nil {
			log.Println(err)
		}
		if err := os.RemoveAll(homePath); err != nil {
			log.Println(err)
		}
	}
}

func main() {
	flag.Parse()
	fmt.Printf("count=%d, parallel=%v\n", count, parallel)
	if parallel {
		wg := sync.WaitGroup{}
		wg.Add(count)

		for i := 0; i < count; i++ {
			go func(id int) {
				mount(id)
				wg.Done()
			}(i)
		}

		wg.Wait()
	} else {
		for i := 0; i < count; i++ {
			mount(i)
		}
	}

	clean()
}
