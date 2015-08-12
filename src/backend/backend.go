package backend


import (
    //"os"
    "os/exec"
    "strconv"
    "fmt"
    "syscall"
    "../errorStack"
    logrotorModule "../logrotor"
)


type Backend struct {
    name string
    addr string
    port int
    cmd  *exec.Cmd
    logs *logrotorModule.Logs
    //process *os.Process
}

func newBackend(buildDir, buildName, appDir string, uid, gid uint32, port int, logs *logrotorModule.Logs) (*Backend, *errorStack.Error) {
    
    newBackend := Backend{
        name    : buildName,
        addr    : "127.0.0.1",
        port    : port,
        logs    : logs,
    }
    
    
    buildPath := buildDir + "/" + buildName
    
    cmd := exec.Command(buildPath, strconv.FormatInt(int64(newBackend.port), 10))
    
    cmd.Dir = appDir
    
                                                    //uruchomienie na koncie określonego użytkownika
    cmd.SysProcAttr = &syscall.SysProcAttr{}
    cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uid, Gid: gid}
    
    cmd.Stdout = newBackend.logs.Std
    cmd.Stderr = newBackend.logs.Err
    
	err := cmd.Start()
    
    if err != nil {
        return nil, errorStack.From(err)
	}
    
    //newBackend.process = cmd.Process
    newBackend.cmd = cmd
    
    
    return &newBackend, nil
}

func (self *Backend) Name() string {
    return self.name
}

func (self *Backend) Stop() {
    
    errKill := self.cmd.Process.Kill()
    
    
    //syscall.Kill(self.process.Pid, syscall.SIGINT)
    
    
    if errKill != nil {
        fmt.Println(errKill)
    }
    
                        //czekaj aż się zakończy ten proces
    self.cmd.Wait()
    
    self.logs.Stop()
    
}


func (self *Backend) Port() int {
    return self.port
}


func (self *Backend) GetAddr() string {
    
    return self.addr + ":" + strconv.FormatInt(int64(self.port), 10)
}


