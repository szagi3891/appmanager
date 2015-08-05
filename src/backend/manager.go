package backend


import (
    "os/exec"
    "sync"
    "strconv"
    "bytes"
    "time"
    //"fmt"
    userModule "os/user"
    "syscall"
    "../errorStack"
)


func Init(gocmd, pwd, buildDir, appMain, appUser string, portStart, portEnd int) (*Manager, *errorStack.Error) {
    
    ports := map[int]*Backend{}
    
    for i:=portStart; i<= portEnd; i++ {
        ports[i] = nil
    }
    
    uid, gid, errLookup := lookupUser(appUser)
    
    if errLookup != nil {
        return nil, errLookup
    }
    
    return &Manager{
        gocmd    : gocmd,
        pwd      : pwd,
        buildDir : buildDir,
        appMain  : appMain,
        appUser  : appUser,
        ports    : ports,
        uid      : uid,
        gid      : gid,
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
    gocmd    string
    mutex    sync.Mutex
    pwd      string
    buildDir string
    appMain  string
    appUser  string
    ports    map[int]*Backend
    uid      uint32
    gid      uint32
}


func (self *Manager) createNewBackend() (*Backend, *errorStack.Error) {
    
    self.mutex.Lock()
    self.mutex.Unlock()
    
    for portIndex, value := range self.ports {
        
        if value == nil {
            
            newBeckend := Backend{addr : "127.0.0.1", port : portIndex}
            
            self.ports[portIndex] = &newBeckend
            
            return &newBeckend, nil
        }
    }
    
    return nil, errorStack.Create("Error alocation ports")
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
    
    current := time.Now()
    
    year, montch, day := current.Date()
    
    hour   := current.Hour()
    minute := current.Minute()
    second := current.Second()
    
    buildName := "build_" + frm(year, 4) + frm(int(montch), 2) + frm(day, 2) + frm(hour, 2) + frm(minute, 2) + frm(second, 2) + "_" + repoSha1
    
    //fmt.Println(buildName)
    
    //go build -o ../appmanager_build/nowy ../wolnemedia/src/main.go
    //go build ./src/main.go
    
    cmd := exec.Command(self.gocmd, "build", "-o", self.buildDir + "/" + buildName, self.appMain)
    
    cmd.Dir = self.pwd
    
    out1 := outData{}
    out2 := outData{}
    
    cmd.Stdout = &out1
    cmd.Stderr = &out2
    
	err := cmd.Run()
    
    if err != nil {
        return errorStack.From(err)
	}
    
    return nil
}


func frm(liczba int, digit int) string {
    
    out := strconv.FormatInt(int64(liczba), 10)
    
    for len(out) < digit {
        out = "0" + out
    }
    
    return out
}


//uruchomienie konkretnego builda
func (self *Manager) New(buildName string) (*Backend, *errorStack.Error) {
    
    
    newBackend, errCerate := self.createNewBackend()
    
    if errCerate != nil {
        return nil, errCerate
    }
    
    buildPath := self.buildDir + "/" + buildName
    
    cmd := exec.Command(buildPath, strconv.FormatInt(int64(newBackend.port), 10))
    
    cmd.Dir = self.pwd
    
    cmd.SysProcAttr = &syscall.SysProcAttr{}
    cmd.SysProcAttr.Credential = &syscall.Credential{Uid: self.uid, Gid: self.gid}
    
    out1 := outData{}
    out2 := outData{}
    
    cmd.Stdout = &out1
    cmd.Stderr = &out2
    
	err := cmd.Start()
    
    if err != nil {
        return nil, errorStack.From(err)
	}
    
    newBackend.process = cmd.Process
    
    return newBackend, nil
}

