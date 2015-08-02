package serverBackend


import (
    "sync"
    "fmt"
)


type ServerBackend struct {
    addr   string
    mutex  sync.Mutex
    count  int
    active int
}

func New(addr string) *ServerBackend {
    
    return &ServerBackend{
        addr : addr,
        count : 0,
    }
}


func (self *ServerBackend) Active() int {
    
    return self.active
}

func (self *ServerBackend) GetAddr() string {
    
    return self.addr
}


//func (self *ServerBackend) Inc() {

func (self *ServerBackend) Inc() {
    
    self.mutex.Lock()
    self.active++
    self.count++
    fmt.Println("nowe połączenie: ", self.count)
	self.mutex.Unlock()
}

func (self *ServerBackend) Sub() {
    
    self.mutex.Lock()
    self.active--
    fmt.Println("zamykam        : ", self.count)
	self.mutex.Unlock()
}