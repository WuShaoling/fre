package container

import (
	"bufio"
	"encoding/json"
	"engine/config"
	"engine/runtime"
	"engine/template"
	"errors"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"os/exec"
	"path"
)

type ZygoteProcess struct {
	id         string
	pid        int
	packageSet []string
}

type ZygoteProcessTree map[string]*ZygoteProcess // key 暂时用 templateName，实际应该为 ZygoteProcess.id

type ZygoteService struct {
	zygoteProcessUnixSocket  map[string]*net.UnixConn     // key 为 id
	runtimeZygoteProcessTree map[string]ZygoteProcessTree // key 为 runtimeName
}

func NewZygoteService(runtimeSet map[string]runtime.Runtime, templateSet map[string]template.Template) *ZygoteService {
	log.Info("NewZygoteService")
	service := &ZygoteService{
		zygoteProcessUnixSocket:  make(map[string]*net.UnixConn),
		runtimeZygoteProcessTree: make(map[string]ZygoteProcessTree),
	}

	if err := service.startUnixSocketServer(); err != nil {
		// 启动失败直接退出
		log.Fatal("NewZygoteService startUnixSocketServer error, ", err)
	}

	// 对 templateSet 按 runtime 分组
	templateGroup := make(map[string][]template.Template)
	for _, v := range templateSet {
		templateGroup[v.Runtime] = append(templateGroup[v.Runtime], v)
	}

	// 对于每一种 runtime，构造 ZygoteProcessRadixTree
	for runtimeName, templateList := range templateGroup {
		if r, ok := runtimeSet[runtimeName]; ok && r.ZygoteCommand != nil && len(r.ZygoteCommand) != 0 {
			service.runtimeZygoteProcessTree[runtimeName] = service.buildRuntimeZygoteProcessTree(runtimeSet[runtimeName], templateList)
		}
	}

	return service
}

func (service *ZygoteService) NewContainerByZygoteProcess(runtimeName, templateName string, functionExecCtx string) error {
	radixTree, ok := service.runtimeZygoteProcessTree[runtimeName]
	if !ok {
		return errors.New("RuntimeNotSupportZygote")
	}

	process, ok := radixTree[templateName]
	if !ok {
		return errors.New("NoMatchZygoteProcessFound")
	}

	unixSocket, ok := service.zygoteProcessUnixSocket[process.id]
	if !ok {
		return errors.New("ZygoteProcessUnixSocketNotFound")
	}
	_, err := unixSocket.Write([]byte(functionExecCtx))

	if err != nil {
		log.Errorf("send new container command to zygote(id=%s, pid=%d) process error, %v", process.id, process.pid, err)
	}
	return err
}

func (service *ZygoteService) buildRuntimeZygoteProcessTree(r runtime.Runtime, templateList []template.Template) ZygoteProcessTree {
	zygoteProcessTree := ZygoteProcessTree{}

	// 对于每一个 template
	for _, t := range templateList {
		// 构造进程参数
		param := map[string]interface{}{
			"id":               t.Name,
			"packageSet":       t.Packages,
			"serverSocketFile": config.SysConfigInstance.ZygoteUnixSocketFile,
		}
		data, err := json.Marshal(param)
		if err != nil {
			log.Error("newZygoteProcess error, ", err)
			continue
		}
		// 启动进程
		cmd := exec.Command(r.ZygoteCommand[0], append(r.ZygoteCommand[1:], string(data))...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = path.Join(config.GetZygoteCodePath(), r.Name)
		if err := cmd.Start(); err != nil {
			log.Errorf("newZygoteProcess start process error, runtime=%s, template=%s, error=%+v", r.Name, t.Name, err)
			continue
		}
		// 加入 tree 中
		zygoteProcessTree[t.Name] = &ZygoteProcess{
			id:         t.Name,
			pid:        cmd.Process.Pid,
			packageSet: t.Packages,
		}
	}
	return zygoteProcessTree
}

func (service *ZygoteService) startUnixSocketServer() error {
	unixAddr, err := net.ResolveUnixAddr("unix", config.SysConfigInstance.ZygoteUnixSocketFile)
	if err != nil {
		log.Error("ZygoteService ResolveUnixAddr error, ", err)
		return err
	}

	unixListener, err := net.ListenUnix("unix", unixAddr)
	if err != nil {
		log.Error("ZygoteService ListenUnix error, ", err)
		return err
	}
	defer unixListener.Close()

	go func() {
		for {
			unixConn, err := unixListener.AcceptUnix()
			if err != nil {
				log.Error("ZygoteService AcceptUnix error, ", err)
				continue
			}
			log.Info("zygote process connected : " + unixConn.RemoteAddr().String())
			go func() {
				reader := bufio.NewReader(unixConn)
				id, err := reader.ReadString('\n')
				if err != nil {
					log.Error("startUnixSocketServer ReadString error, ", err)
					return
				}
				service.zygoteProcessUnixSocket[id] = unixConn
			}()
		}
	}()
	return nil
}
