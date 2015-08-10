package proxy


import (
    "strconv"
    "../handleConn"
    "../errorStack"
    backendModule "../backend"
    logrotorModule "../logrotor"
)


type Proxy struct {
    manager   *backendModule.Manager
    backend   *backendModule.Backend
}


func New(appStderr *logrotorModule.LogWriter, mainPort int, logrotor *logrotorModule.Manager, manager *backendModule.Manager) (*Proxy, *errorStack.Error) {
    
    
    backend, errStartLastBuild := manager.StartLastBuild()
    
    if errStartLastBuild != nil {
        return nil, errStartLastBuild
    }
    
    
    addr := "127.0.0.1:" + strconv.FormatInt(int64(mainPort), 10)
    
    
    proxy := &Proxy{
        manager   : manager,
        backend   : backend,
    }
    
    
    errStart := handleConn.Start(addr, appStderr, func() *backendModule.Backend{
        return proxy.backend
    })
    
    if errStart != nil {
        return nil, errStart
    }
    
    return proxy, nil
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


