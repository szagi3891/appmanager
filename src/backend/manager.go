package backend


import (
    "os/exec"
    "sync"
    "strconv"
    "bytes"
    "sort"
    userModule "os/user"
    "syscall"
    "../errorStack"
    logrotorModule "../logrotor"
    utils "../utils"
    "io/ioutil"
    "fmt"
    "../handleConn"
)


func Init(mainPort int, logrotor *logrotorModule.Manager,  appStdout, appStderr *logrotorModule.LogWriter, gocmd, pwd, buildDir, appMain, appUser, gopath string, portStart, portEnd int) (*Manager, *errorStack.Error) {
    
    
    ports := map[int]*Backend{}
    
    for i:=portStart; i<= portEnd; i++ {
        ports[i] = nil
    }
    
    uid, gid, errLookup := lookupUser(appUser)
    
    if errLookup != nil {
        return nil, errLookup
    }
    
    manager := Manager{
        mainPort  : mainPort,
        //backend   : backend,
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
    }
    
    
    backend, errStartLastBuild := manager.startLastBuild()
    
    if errStartLastBuild != nil {
        return nil, errStartLastBuild
    }
    
    manager.backend = backend
    
    
    
    addr := "127.0.0.1:" + strconv.FormatInt(int64(mainPort), 10)
    
    errStart := handleConn.Start(addr, appStderr, func() (string, func(), func()) {
        
        backend := manager.GetActiveBackend()
        
        return backend.GetAddr(), backend.Inc, backend.Sub
    })
    
    if errStart != nil {
        return nil, errStart
    }
    
    
    
    return &manager, nil
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
    mainPort  int
    backend   *Backend
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

func (self *Manager) GetActiveBackend() *Backend {
    return self.backend
}

func (self *Manager) SwitchByNameAndPort(name string, port int) bool {
    
    backend, isFind := self.ports[port]
    
    if isFind && backend != nil && backend.Name() == name {
        self.backend = backend
        return true
    }
    
    return false
}

func (self *Manager) DownByNameAndPort(name string, port int) bool {
    
    //trzeba wprowadzić mutex
    
    backend, isFind := self.ports[port]
    
    if isFind && backend != nil && backend.Name() == name {
        
        delete(self.ports, port)
        
        backend.Stop()
        return true
    }
    
    return false
}

func (self *Manager) GetMainPort() int {
    return self.mainPort
}

func (self *Manager) GetSha1Repo() (string, *errorStack.Error) {
    
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


func (self *Manager) MakeBuild() (string, *errorStack.Error) {
    
    repoSha1, errRepo := self.GetSha1Repo()
    
    if errRepo != nil {
        return "", errRepo
    }
    
    
    buildName := "build_" + utils.GetCurrentTimeName()  + "_" + repoSha1
    
    
    //cmd := exec.Command(self.gocmd, "build", "-race", "-o", self.buildDir + "/" + buildName, self.appMain)
    cmd := exec.Command(self.gocmd, "build", "-o", self.buildDir + "/" + buildName, self.appMain)
    
    cmd.Dir = self.pwd
    cmd.Env = []string{"GOPATH=" + self.gopath}
    
    var bufOut bytes.Buffer
    var bufErr bytes.Buffer
    
	cmd.Stdout = &bufOut
	cmd.Stderr = &bufErr
    
	err := cmd.Run()
    
    fmt.Println(bufOut.String())
    fmt.Println(bufErr.String())
    
    //panic("TODO - wyniki działania programu z bufora trzeba przepchnąć do loga")
    
    if err != nil {
        return "", errorStack.From(err)
	}
    
    return buildName, nil
}


func (self *Manager) createNewBackend(buildName string) (*Backend, *errorStack.Error) {
    
    self.mutex.Lock()
    self.mutex.Unlock()
    
    for portIndex, value := range self.ports {
        
        if value == nil {
            
            newBeckend := Backend{
                name    : buildName,
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
    
    cmd.Stdout = newBackend.stdout
    cmd.Stderr = newBackend.stderr
    
	err := cmd.Start()
    
    if err != nil {
        return nil, errorStack.From(err)
	}
    
    newBackend.process = cmd.Process
    
    return newBackend, nil
}


//sprawdza czy z obecnego komitu można zrobić nowego builda
func (self *Manager) IsAvailableNewCommit(listBuild *[]string, lastCommitRepo string) bool {
    
    for _, name := range *listBuild {
        if lastCommitRepo == name[21:] {
            return false
        }
    }
    
    return true
}

type AppInfo struct {
    Port    int
    Name    string
    Active  int
}

func (self *Manager) GetAppList() *[]*AppInfo {
    
    order  := []string{}
    preOut := map[string]*AppInfo{}
    
    for _, back := range self.ports {
        
        if back != nil {
            
            fullName := back.Name() + "_" + strconv.FormatInt(int64(back.Port()), 10)
            
            order = append(order, fullName)
            
            preOut[fullName] = &AppInfo{
                Port : back.Port(),
                Name : back.Name(),
                Active : back.Active(),
            }
        }
    }
    
    sort.Sort(sort.Reverse(sort.StringSlice(order)))
    
    
    out := []*AppInfo{}
    
    for _, name := range order {
        out = append(out, preOut[name])
    }
    
    return &out
}


func (self *Manager) GetListBuild() (*[]string, *errorStack.Error) {
    
    
    list, errList := ioutil.ReadDir(self.buildDir)
    
    if errList != nil {
        return nil, errorStack.From(errList)
    }
    
    out := []string{}
    
    for _, item := range list {
        
        name := item.Name()
        
        if isValidBuildName(name) {
            out = append(out, name)
        }
    }
    
    return &out, nil
}


//startuje ostatniego builda, jeśli nie ma nic w katalogu z buildami to sobie tworzy takowego builda
func (self *Manager) startLastBuild() (*Backend, *errorStack.Error) {
    
    list, errList := self.GetListBuild()
    
    if errList != nil {
        return nil, errList
    }
    
    lastName, isFind := findLast(list)
    
    if isFind == true {
        
        return self.New(lastName)
        
    } else {
        
        newBuildName, errBuild := self.MakeBuild()
        
        if errBuild != nil {
            return nil, errBuild
        }
        
        return self.New(newBuildName)   
    }
}


func findLast(list *[]string) (string, bool) {
    
    max := ""
    last := ""
    
    for _, item := range *list {
        
        data := item[6:20]
        
        if max < data {
            max  = data
            last = item
        }
    }
    
    if last == "" {
        return "", false
    } else {
        return last, true
    }
}


func isValidBuildName(name string) (bool) {
    
    //5 + 1 + 14 + 1 + 40
    //build_14cyfr_40cyfr
    
    if len(name) != 61 {
        return false
    }
    
    name1 := name[0:6]
    name2 := name[6:20]
    name3 := name[20:21]
    name4 := name[21:]
    
    return name1 == "build_" && isDigit(name2, false) && name3 == "_" && isDigit(name4, true)
}


func isDigit(name string, isHash bool) bool {
    
    for i:=0; i<len(name); i++ {    //, char := range name {
        
        if isHash {
            
            if "0"[0] <= name[i] && name[i] <= "9"[0] || "a"[0] <= name[i] && name[i] <= "f"[0] {
                //ok
            } else {
                return false
            }
        
        } else {
            if "0"[0] <= name[i] && name[i] <= "9"[0] {
                //ok
            } else {
                return false
            }
        }
    }
    
    return true
}