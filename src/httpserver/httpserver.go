package httpserver

import (
    "net"
    "net/http"
    "time"
    "sync"
    "../errorStack"
)

type serv struct {
    wg      sync.WaitGroup
    handler http.Handler
}

func (self *serv) ServeHTTP(out http.ResponseWriter, req *http.Request) {
    self.wg.Add(1)
    defer self.wg.Done()
    self.handler.ServeHTTP(out, req)
}

func Start(addr string, funcHandle func(out http.ResponseWriter, req *http.Request)) (func(), *errorStack.Error) {
    
    lisiner, errStart := net.Listen("tcp", addr)
	
    if errStart != nil {
        return func(){}, errorStack.From(errStart)
    }
    
    servInst := serv {
        handler : http.HandlerFunc(funcHandle),
    }
    
	s := &http.Server{
		Handler        : &servInst,
		ReadTimeout    : 10 * time.Second,
		WriteTimeout   : 10 * time.Second,
		MaxHeaderBytes : 1 << 20,
	}
    
	go s.Serve(lisiner)
    
    return func(){
        lisiner.Close()
        servInst.wg.Wait()
    }, nil
}
