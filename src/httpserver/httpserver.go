package httpserver

import (
    "net"
    "net/http"
    "strconv"
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
    self.handler.ServeHTTP(out, req)
    self.wg.Done()
}

func Start(port int64, funcHandle func(out http.ResponseWriter, req *http.Request)) (func(), *errorStack.Error) {
    
    
    addr := ":" + strconv.FormatInt(port, 10)
    
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


func Start2(port int64, funcHandle func(out http.ResponseWriter, req *http.Request)) *errorStack.Error {
    
    servInst := serv {
        handler : http.HandlerFunc(funcHandle),
    }
    
    addr := ":" + strconv.FormatInt(port, 10)
    
	s := &http.Server{
		Addr           : addr,
		Handler        : &servInst,
		ReadTimeout    : 10 * time.Second,
		WriteTimeout   : 10 * time.Second,
		MaxHeaderBytes : 1 << 20,
	}
    
	errStart := s.ListenAndServe()
    
    if errStart != nil {
        return errorStack.From(errStart)
    }
    
    return nil
}

