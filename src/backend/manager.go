package backend


import (
    "os/exec"
    "sync"
    "errors"
    "strconv"
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


func (self *Manager) New(buildName string) (*Backend, error) {
    
    
    newBackend, errCerate := self.createNewBackend()
    
    if errCerate != nil {
        return nil, errCerate
    }
    
    buildPath := self.buildDir + "/" + buildName
    
    cmd := exec.Command(buildPath, strconv.FormatInt(int64(newBackend.port), 10))
    
    cmd.Dir = self.pwd
    
    out1 := outData{"stdout"}
    out2 := outData{"stderr"}
    
    cmd.Stdout = &out1
    cmd.Stderr = &out2
    
	err := cmd.Start()
    
    if err != nil {
		return nil, err
	}
    
    
    return newBackend, nil
}

