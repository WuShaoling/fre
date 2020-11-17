package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

var parallel bool
var count int

func init() {
	flag.IntVar(&count, "n", 10, "启动的容器数量")
	flag.BoolVar(&parallel, "p", false, "并发启动")
}

func runContainer() {
	cmd := exec.Command("python3", "./core.py")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		print(err)
	}
}

func main() {
	fmt.Printf("count=%d, parallel=%v\n", count, parallel)
	if parallel {
		wg := sync.WaitGroup{}
		wg.Add(count)

		for i := 0; i < count; i++ {
			go func() {
				runContainer()
				wg.Done()
			}()
		}

		wg.Wait()
	} else {
		for i := 0; i < count; i++ {
			runContainer()
		}
	}
}
