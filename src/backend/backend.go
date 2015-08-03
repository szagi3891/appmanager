package serverBackend


import (
    "sync"
    "fmt"
)


type Backend struct {
    addr   string
    mutex  sync.Mutex
    count  int
    active int
}


func New(addr string) *Backend {
    
    return &Backend{
        addr : addr,
        count : 0,
    }
}


func (self *Backend) Active() int {
    
    return self.active
}

func (self *Backend) GetAddr() string {
    
    return self.addr
}


//func (self *ServerBackend) Inc() {

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