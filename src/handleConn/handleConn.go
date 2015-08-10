package handleConn


import (
    "net"
    "time"
    "io"
    "fmt"
    "../errorStack"
    logrotorModule "../logrotor"
)


func Start(addr string, logs *logrotorModule.Logs, getBackend func() (string, func(), func())) *errorStack.Error {
    
    addProxy, err1 := net.ResolveTCPAddr("tcp", addr)
    
    if err1 != nil {
        return errorStack.From(err1)
    }
        
    
	listener, err2 := net.ListenTCP("tcp", addProxy)
    
	if err2 != nil {
        return errorStack.From(err2)
	}
    
    go func(){
        
        for {

            conn, errAccept := listener.AcceptTCP()

            if errAccept != nil {
                
                logs.Err.WriteString(errorStack.From(errAccept).String())
                
            } else {
                
                errConnectect := handleConn(getBackend, conn)
                
                if errConnectect != nil {
                    
                    logs.Err.WriteString(errorStack.From(errConnectect).String())
                    
                } else {
                    
                    //TODO
                    //trzeba jakoś logować prawidłowe połączenia
                }
            }
        }
    }()
    
    return nil
}

type connectionFlag struct {
    
    isClose1 chan bool
    ord1Size int64
    ord1Err1 error
    ord1Err2 error
    
    isClose2 chan bool
    ord2Size int64
    ord2Err1 error
    ord2Err2 error
}

func newConnectionFlag() *connectionFlag {
    
    return &connectionFlag{
        isClose1 : make(chan bool),
        isClose2 : make(chan bool),
    }
}

func (self *connectionFlag) isDone() bool {
    
    select {
        case <- self.isClose1 : {
            return true
        }
        default : {
        }
    }
    
    select {
        case <- self.isClose2 : {
            return true
        }
        default : {
        }
    }
    
    return false
}
    
func (self *connectionFlag) Set1(size int64, err1, err2 error) {
    
    self.ord1Size = size
    self.ord1Err1 = err1
    self.ord1Err2 = err2
    close(self.isClose1)
}

func (self *connectionFlag) Set2(size int64, err1, err2 error) {
    
    self.ord2Size = size
    self.ord2Err1 = err1
    self.ord2Err2 = err2
    close(self.isClose2)
}

func (self *connectionFlag) Wait() {
    
    <- self.isClose1
    <- self.isClose2
    
    fmt.Println("1: size :", self.ord1Size)
    fmt.Println("1: err1 :", self.ord1Err1)
    fmt.Println("1: err2 :", self.ord1Err2)
    fmt.Println("2: size :", self.ord2Size)
    fmt.Println("2: err1 :", self.ord2Err1)
    fmt.Println("2: err2 :", self.ord2Err2)
}




func handleConn(getBackend func() (string, func(), func()), connIn *net.TCPConn) error {
    
    addr, incCounter, subCounter := getBackend()
    
    addrServer, errParse := net.ResolveTCPAddr("tcp", addr)   //backend.GetAddr()
    
    if errParse != nil {
        return errParse
    }
    
	connDest, err := net.DialTCP("tcp", nil, addrServer)
    
	if err != nil {
		return err
	}
    
    
    go func(){
        
        defer connIn.Close()
        defer connDest.Close()

        incCounter()
        defer subCounter()
        
        
        connFlag := newConnectionFlag()
        
        go func() {
            connFlag.Set1(Copy(connFlag, connIn, connDest))
        }()
        
        go func() {
            connFlag.Set2(Copy(connFlag, connDest, connIn))
        }()
        
        connFlag.Wait()
    }()
    
    return nil
}


//http://golang.org/pkg/net/#example_Listener
//http://golang.org/src/io/io.go?s=12247:12307#L340
//http://www.badgerr.co.uk/2011/06/20/golang-away-tcp-chat-server/


func Copy(flag *connectionFlag, src *net.TCPConn, dst *net.TCPConn) (int64, error, error) {
    
    
    written := int64(0)
    buf     := make([]byte, 32*1024)
    
    
    for {
        
        if flag.isDone() {
            return written, nil, nil
        }
        
        src.SetReadDeadline(time.Now().Add(time.Second))
        
        nr, err1 := src.Read(buf)
        
        err1 = filterTimeout(err1)
        
        if nr > 0 {
            
            dst.SetWriteDeadline(time.Now().Add(time.Second))
            
            nw, err2 := dst.Write(buf[0:nr])
            
            if nw > 0 {
                written += int64(nw)
            }
            
            
            err2 = filterTimeout(err2)
            
            if err2 != nil {
                return written, err1, err2
            }
            
            if nr != nw {
                return written, err1, io.ErrShortWrite
            }
        }
        
        if err1 != nil {
            return written, err1, nil
        }
    }
}


    /*if err == io.EOF {
        return err
    }*/
    

func filterTimeout(err error) error {
    
    if isTimeout(err) {
        return nil
    }
    
    return err
}

func isTimeout(err error) bool {
    
    if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
        return true
    }
    
    return false
}

