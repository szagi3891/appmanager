package proxy


import (
    "net"
    "strconv"
    "../handleConn"
    "../errorStack"
    backendModule "../backend"
    logrotorModule "../logrotor"
)


type Proxy struct {
    
    appStderr *logrotorModule.LogWriter
    listener  *net.TCPListener
    backend   *backendModule.Backend
}


func New(appStderr *logrotorModule.LogWriter, mainPort int, logrotor *logrotorModule.Manager, backend *backendModule.Backend) (*Proxy, error) {
    
    
    addr := "127.0.0.1:" + strconv.FormatInt(int64(mainPort), 10)
    
    addProxy, err1 := net.ResolveTCPAddr("tcp", addr)
    
    if err1 != nil {
        return nil, err1
    }
        
    
	listener, err2 := net.ListenTCP("tcp", addProxy)
    
	if err2 != nil {
		return nil, err2
	}
    
    
    proxy := &Proxy{
        appStderr : appStderr,
        listener  : listener,
        backend   : backend,
    }
    
    proxy.start()
    
    return proxy, nil
}


func (self *Proxy) Switch(backend *backendModule.Backend) {
    
    self.backend = backend
}



func (self *Proxy) start() {
    
    go func(){
        
        for {

            conn, errAccept := self.listener.AcceptTCP()

            if errAccept != nil {
                
                self.appStderr.WriteString(errorStack.From(errAccept).String())
                
            } else {
                
                errConnectect := handleConn.HandleConn(self.backend, conn)
                
                if errConnectect != nil {
                    
                    self.appStderr.WriteString(errorStack.From(errConnectect).String())
                    
                } else {
                    
                    //TODO
                    //trzeba jakoś logować prawidłowe połączenia
                }
            }
        }
    }()
}