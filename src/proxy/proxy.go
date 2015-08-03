package proxy


import (
    "net"
    "fmt"
    "../handleConn"
    backendModule "../backend"
)


type Proxy struct {
    
    addr     string
    log      string
    build    string
    listener *net.TCPListener
    logCh    chan *string
    backend  *backendModule.Backend
}


func New(addr, log, build string, backend *backendModule.Backend) (*Proxy, error) {
    
    
    addProxy, err1 := net.ResolveTCPAddr("tcp", addr)
    
    if err1 != nil {
        return nil, err1
    }
        
    
	listener, err2 := net.ListenTCP("tcp", addProxy)
    
	if err2 != nil {
		return nil, err2
	}
    
    
    logCh := make(chan *string)
    
    
    proxy := &Proxy{
        addr     : addr,
        log      : log,
        build    : build,
        listener : listener,
        logCh    : logCh,
        backend  : backend,
    }
    
    proxy.start()
    
    return proxy, nil
}


func (self *Proxy) Switch(backend *backendModule.Backend) {
    
    //old := self.serverBackend
    
    self.backend = backend
    
    //old.WaitWhenAllClose()
}


func printLog(logCh chan *string) {
    
    go func(){
        
        for logItem := range logCh {
            
            fmt.Println(logItem)
        }
        
    }()
}


func (self *Proxy) start() {
    
    go func(){
        
        for {

            conn, errAccept := self.listener.AcceptTCP()

            if errAccept != nil {
                
                errStr := "err: " + errAccept.Error()
                self.logCh <- &errStr
                
            } else {
                
                errConnectect := handleConn.HandleConn(self.backend, conn)
                
                if errConnectect != nil {
                    
                    errStr := "err: " + errConnectect.Error()
                    self.logCh <- &errStr
                    
                } else {
                    
                    //trzeba jakoś logować prawidłowe połączenia
                }
            }
        }
    }()
}