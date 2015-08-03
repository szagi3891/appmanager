package serverBackend


import (
    "os/exec"
    "sync"
    "fmt"
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

//Read(p []byte) (n int, err error)
//Write(p []byte) (n int, err error)

type outData struct {
    name string
}

func (self *outData) Write(p []byte) (n int, err error) {
    
    fmt.Println(self.name + " Write " + string(p))
    return len(p), nil
}

func (self *Manager) New(buildName string) (*Backend, error) {
    
    
    newBackend, errCerate := self.createNewBackend()
    
    if errCerate != nil {
        return nil, errCerate
    }
    
    
    //panic("TODO - tutaj trzeba uruchomić tą aplikację backendową")
    
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



type Backend struct {
    addr   string
    port   int
    mutex  sync.Mutex
    count  int
    active int
}


func (self *Backend) Active() int {
    
    return self.active
}

func (self *Backend) GetAddr() string {
    
    return self.addr + ":" + strconv.FormatInt(int64(self.port), 10)
}



func (self *Backend) Inc() {
    
    self.mutex.Lock()
    self.active++
    self.count++
    fmt.Println("nowe połączenie: ", self.addr, " count: ", self.count, " active: ", self.active)
	self.mutex.Unlock()
}

func (self *Backend) Sub() {
    
    self.mutex.Lock()
    self.active--
    fmt.Println("zamykam        : ", self.addr, " count: ", self.count, " active: ", self.active)
	self.mutex.Unlock()
}