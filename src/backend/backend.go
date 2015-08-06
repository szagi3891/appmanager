package backend


import (
    "os"
    "sync"
    "fmt"
    "strconv"
    logrotorModule "../logrotor"
)


type Backend struct {
    stop    bool
    addr    string
    port    int
    mutex   sync.Mutex
    active  int
    process *os.Process
    isClose chan bool
    stdout  *logrotorModule.LogWriter
    stderr  *logrotorModule.LogWriter
}


func (self *Backend) checkStrop() {
    
    self.mutex.Lock()
    self.mutex.Unlock()
    
    if self.stop == true && self.active == 0 {
        
        if self.process != nil {
            
            close(self.isClose)
            
            fmt.Println("zlecam morderstwo ...")
            
            errKill := self.process.Kill()
            
            if errKill != nil {
                fmt.Println(errKill)
            }
            
            self.process = nil
        }
    }
}


func (self *Backend) Stop() {
    
    self.mutex.Lock()    
    self.stop = true
    self.mutex.Unlock()
    
    self.checkStrop()
    
                        //czekaj aż wszystkie powiązane procesy z tym procesem się zakończą
    <- self.isClose
    
    self.stdout.Stop()
    self.stderr.Stop()
}


func (self *Backend) Active() int {
    
    return self.active
}


func (self *Backend) GetAddr() string {
    
    return self.addr + ":" + strconv.FormatInt(int64(self.port), 10)
}


func (self *Backend) Inc() {
    
    self.mutex.Lock()
    
    if self.stop == true {
        panic("TODO - blokada kolejnych połączeń na ten adres")
    }
    
    self.active++
    fmt.Println("nowe połączenie: ", self.addr, " active: ", self.active)
	self.mutex.Unlock()
}


func (self *Backend) Sub() {
    
    self.mutex.Lock()
    self.active--
    fmt.Println("zamykam        : ", self.addr, " active: ", self.active)
	self.mutex.Unlock()
    
    self.checkStrop()
}
