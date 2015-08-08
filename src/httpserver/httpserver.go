package httpserver

import (
    "net/http"
    "strconv"
    "time"
    "../errorStack"
)

type serv struct {
    handler http.Handler
}

func (self *serv) ServeHTTP(out http.ResponseWriter, req *http.Request) {
    self.handler.ServeHTTP(out, req)
}

func Start(port int64, funcHandle func(out http.ResponseWriter, req *http.Request)) *errorStack.Error {
    
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

