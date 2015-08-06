package proxy


import (
    "net"
    "fmt"
    "strconv"
    "../handleConn"
    backendModule "../backend"
    logrotorModule "../logrotor"
)


type Proxy struct {
    
    listener *net.TCPListener
    logCh    chan *string
    backend  *backendModule.Backend
}


func New(mainPort int, logrotor *logrotorModule.Manager, backend *backendModule.Backend) (*Proxy, error) {
    
    //TODO
        //trzeba będzie utworzyć dwa strumienie na logi z obiektu proxy
    
    addr := "127.0.0.1:" + strconv.FormatInt(int64(mainPort), 10)
    
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
        listener : listener,
        logCh    : logCh,
        backend  : backend,
    }
    
    proxy.start()
    
    return proxy, nil
}


func (self *Proxy) Switch(backend *backendModule.Backend) {
    
    self.backend = backend
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