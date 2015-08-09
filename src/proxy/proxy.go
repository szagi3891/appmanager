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
    mainPort  int
    appStderr *logrotorModule.LogWriter
    listener  *net.TCPListener
    manager   *backendModule.Manager
    backend   *backendModule.Backend
}


func New(appStderr *logrotorModule.LogWriter, mainPort int, logrotor *logrotorModule.Manager, manager *backendModule.Manager) (*Proxy, *errorStack.Error) {
    
    
    backend, errStartLastBuild := manager.StartLastBuild()
    
    if errStartLastBuild != nil {
        return nil, errStartLastBuild
    }
    
    
    addr := "127.0.0.1:" + strconv.FormatInt(int64(mainPort), 10)
    
    addProxy, err1 := net.ResolveTCPAddr("tcp", addr)
    
    if err1 != nil {
        return nil, errorStack.From(err1)
    }
        
    
	listener, err2 := net.ListenTCP("tcp", addProxy)
    
	if err2 != nil {
        return nil, errorStack.From(err2)
	}
    
    
    proxy := &Proxy{
        mainPort  : mainPort,
        appStderr : appStderr,
        listener  : listener,
        manager   : manager,
        backend   : backend,
    }
    
    proxy.start()
    
    return proxy, nil
}

func (self *Proxy) GetMainPort() int {
    return self.mainPort
}

func (self *Proxy) GetActive() *backendModule.Backend {
    return self.backend
}

func (self *Proxy) SwitchByNameAndPort(name string, port int) bool {
    
    backend, isSwitch := self.manager.SwitchByNameAndPort(name, port)
    
    if isSwitch {
        self.backend = backend
    }
    
    return isSwitch
}


func (self *Proxy) DownByNameAndPort(name string, port int) bool {
    
    isDown := self.manager.DownByNameAndPort(name, port)
    
    return isDown
}


/*
func (self *Proxy) Switch(backend *backendModule.Backend) {
    
    self.backend = backend
}
*/


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