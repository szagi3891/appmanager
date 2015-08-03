package backend


import (
    "os/exec"
    "sync"
    "errors"
    "strconv"
    "bytes"
    "time"
    "fmt"
)


func Init(pwd string, buildDir string, portStart, portEnd int) *Manager {
    
    ports := map[int]*Backend{}
    
    for i:=portStart; i<= portEnd; i++ {
        ports[i] = nil
    }
    
    return &Manager{
        pwd      : pwd,
        buildDir : buildDir,
        ports    : ports,
    }
}


type Manager struct {
    mutex    sync.Mutex
    pwd      string
    buildDir string
    ports    map[int]*Backend
}


func (self *Manager) createNewBackend() (*Backend, error) {
    
    self.mutex.Lock()
    self.mutex.Unlock()
    
    for portIndex, value := range self.ports {
        
        if value == nil {
            
            newBeckend := Backend{addr : "127.0.0.1", port : portIndex}
            
            self.ports[portIndex] = &newBeckend
            
            return &newBeckend, nil
        }
    }
    
    return nil, errors.New("Error alocation ports")
}


func (self *Manager) getSha1Repo() (string, error) {
    
    cmd := exec.Command("git", "rev-parse", "HEAD")
    
    cmd.Dir = self.pwd
    
    var stderr bytes.Buffer
	var stdout bytes.Buffer
	
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
    
	err := cmd.Run()
	
    if err != nil {
        return "", err
    }
    
    out := stdout.String()
    
    if len(out) >= 40 {
        return out[0:40], nil
    }
    
    return "", errors.New("Zbyt mała odpowiedź")
}


func (self *Manager) MakeBuild() error {
    
    repoSha1, errRepo := self.getSha1Repo()
    
    if errRepo != nil {
        return errRepo
    }
    
    current := time.Now()
    
    year, montch, day := current.Date()
    
    hour   := current.Hour()
    minute := current.Minute()
    second := current.Second()
    
    name := "build_" + frm(year, 4) + frm(int(montch), 2) + frm(day, 2) + frm(hour, 2) + frm(minute, 2) + frm(second, 2) + "_" + repoSha1
    
    fmt.Println(name)
    
    //go build ./src/main.go
    
    return nil
}


func frm(liczba int, digit int) string {
    
    out := strconv.FormatInt(int64(liczba), 10)
    
    for len(out) < digit {
        out = "0" + out
    }
    
    return out
}


func (self *Manager) New(buildName string) (*Backend, error) {
    
    
    newBackend, errCerate := self.createNewBackend()
    
    if errCerate != nil {
        return nil, errCerate
    }
    
    buildPath := self.buildDir + "/" + buildName
    
    cmd := exec.Command(buildPath, strconv.FormatInt(int64(newBackend.port), 10))
    
    cmd.Dir = self.pwd
    
    out1 := outData{}
    out2 := outData{}
    
    cmd.Stdout = &out1
    cmd.Stderr = &out2
    
	err := cmd.Start()
    
    if err != nil {
		return nil, err
	}
    
    newBackend.process = cmd.Process
    
    return newBackend, nil
}

