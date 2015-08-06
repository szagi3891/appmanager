package backend


import (
    "os/exec"
    "sync"
    "strconv"
    "bytes"
    //"fmt"
    userModule "os/user"
    "syscall"
    "../errorStack"
    logrotorModule "../logrotor"
    utils "../utils"
)


func Init(logrotor *logrotorModule.Manager,  appStdout, appStderr *logrotorModule.LogWriter, gocmd, pwd, buildDir, appMain, appUser, gopath string, portStart, portEnd int) (*Manager, *errorStack.Error) {
    
    ports := map[int]*Backend{}
    
    for i:=portStart; i<= portEnd; i++ {
        ports[i] = nil
    }
    
    uid, gid, errLookup := lookupUser(appUser)
    
    if errLookup != nil {
        return nil, errLookup
    }
    
    return &Manager{
        logrotor  : logrotor,
        appStdout : appStdout,
        appStderr : appStderr,
        gocmd     : gocmd,
        pwd       : pwd,
        buildDir  : buildDir,
        appMain   : appMain,
        appUser   : appUser,
        ports     : ports,
        uid       : uid,
        gid       : gid,
        gopath    : gopath,
    }, nil
}


func lookupUser(appUser string) (uint32, uint32, *errorStack.Error) {
    
    userDetail, errLookup := userModule.Lookup(appUser)
    
    if errLookup != nil {
        return 0, 0, errorStack.From(errLookup)
    }
    
    uid, errUid := strconv.ParseInt(userDetail.Uid, 10, 64)
    
    if errUid != nil {
        return 0, 0, errorStack.From(errUid)
    }
    
    gid, errGid := strconv.ParseInt(userDetail.Gid, 10, 64)
    
    if errGid != nil {
        return 0, 0, errorStack.From(errGid)
    }
    
    return uint32(uid), uint32(gid), nil
}
    


type Manager struct {
    logrotor  *logrotorModule.Manager
    appStdout *logrotorModule.LogWriter
    appStderr *logrotorModule.LogWriter
    gocmd     string
    mutex     sync.Mutex
    pwd       string
    buildDir  string
    appMain   string
    appUser   string
    ports     map[int]*Backend
    uid       uint32
    gid       uint32
    gopath    string
}



func (self *Manager) getSha1Repo() (string, *errorStack.Error) {
    
    cmd := exec.Command("git", "rev-parse", "HEAD")
    
    cmd.Dir = self.pwd
    
    var stderr bytes.Buffer
	var stdout bytes.Buffer
	
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
    
	err := cmd.Run()
	
    if err != nil {
        return "", errorStack.From(err)
    }
    
    out := stdout.String()
    
    if len(out) >= 40 {
        return out[0:40], nil
    }
    
    return "", errorStack.Create("Zbyt mała odpowiedź")
}


func (self *Manager) MakeBuild() *errorStack.Error {
    
    repoSha1, errRepo := self.getSha1Repo()
    
    if errRepo != nil {
        return errRepo
    }
    
    
    buildName := "build_" + utils.GetCurrentTimeName()  + "_" + repoSha1
    
    
    cmd := exec.Command(self.gocmd, "build", "-o", self.buildDir + "/" + buildName, self.appMain)
    
    cmd.Dir = self.pwd
    cmd.Env = []string{"GOPATH=" + self.gopath}
    
    var bufOut bytes.Buffer
    var bufErr bytes.Buffer
    
	cmd.Stdout = &bufOut
	cmd.Stderr = &bufErr
    
	err := cmd.Run()
    
    panic("TODO - wyniki działania programu z bufora trzeba przepchnąć do loga")
    
    if err != nil {
        return errorStack.From(err)
	}
    
    return nil
}


func (self *Manager) createNewBackend(buildName string) (*Backend, *errorStack.Error) {
    
    self.mutex.Lock()
    self.mutex.Unlock()
    
    for portIndex, value := range self.ports {
        
        if value == nil {
            
            newBeckend := Backend{
                addr    : "127.0.0.1",
                port    : portIndex,
                isClose : make(chan bool),
                stdout  : self.logrotor.New(buildName, true),
                stderr  : self.logrotor.New(buildName, false),
                
            }
            
            self.ports[portIndex] = &newBeckend
            
            return &newBeckend, nil
        }
    }
    
    return nil, errorStack.Create("Error alocation ports")
}

//uruchomienie konkretnego builda
func (self *Manager) New(buildName string) (*Backend, *errorStack.Error) {
    
    newBackend, errCerate := self.createNewBackend(buildName)
    
    if errCerate != nil {
        return nil, errCerate
    }
    
    buildPath := self.buildDir + "/" + buildName
    
    cmd := exec.Command(buildPath, strconv.FormatInt(int64(newBackend.port), 10))
    
    cmd.Dir = self.pwd
    
                                                    //uruchomienie na koncie określonego użytkownika
    cmd.SysProcAttr = &syscall.SysProcAttr{}
    cmd.SysProcAttr.Credential = &syscall.Credential{Uid: self.uid, Gid: self.gid}
    
    cmd.Stdout = newBackend.stdout      //self.logrotor.New(buildName, true)
    cmd.Stderr = newBackend.stderr      //self.logrotor.New(buildName, false)
    
	err := cmd.Start()
    
    if err != nil {
        return nil, errorStack.From(err)
	}
    
    newBackend.process = cmd.Process
    
    return newBackend, nil
}

