package main

import (
    "os"
    "os/signal"
    "net"
    "io"
    "fmt"
    "sync"
    "time"
)


func interruptNotify() {
    
    signalInterrupt := make(chan os.Signal, 1)
    signal.Notify(signalInterrupt, os.Interrupt)
    
    go func(){
        
        for sig := range signalInterrupt {
            // sig is a ^C, handle it
            
            fmt.Println(sig)
            panic("signalInterrupt and stop")
        }
    }()
}

//var once sync.Once

var mutex sync.Mutex
var count int

func inc() {
    
    mutex.Lock()
    fmt.Println("nowe połączenie: ", count)
    count++
	mutex.Unlock()
}

func sub() {
    mutex.Lock()
    count--
    fmt.Println("zamykam        : ", count)
	mutex.Unlock()
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


func main(){
    
    interruptNotify()
    
    
        //81 -> 9999
    
    addProxy, err1 := net.ResolveTCPAddr("tcp", "localhost:8888")
    //addProxy, err1 := net.ResolveTCPAddr("tcp", "localhost:80")
    
    if err1 != nil {
        panic(err1)
    }
    
    fmt.Println("....", addProxy)
    
	listener, err := net.ListenTCP("tcp", addProxy)
    
	if err != nil {
		panic(err)
	}
    
	for {
        
		conn, err := listener.AcceptTCP()
        
		if err != nil {
			panic(err)
		}
        
        go handleConn(conn)
	}
}


//<-chan

func handleConn(connIn *net.TCPConn) {
    
    current := count
    
    inc()
    
    defer sub()
    
    
    defer connIn.Close()
    
    
    fmt.Println("nawiązano nowe połączenie")
    
    addrServer, errParse := net.ResolveTCPAddr("tcp", "localhost:9999")
    //addrServer, errParse := net.ResolveTCPAddr("tcp", "213.180.141.140:80")     //onet
    //addrServer, errParse := net.ResolveTCPAddr("tcp", "pl.wikipedia.org:80")
    
    
    if errParse != nil {
        panic(errParse)
    }
    
	connDest, err := net.DialTCP("tcp", nil, addrServer)
    
	if err != nil {
		panic(err)
	}
    
	defer connDest.Close()
    
    
    connFlag := newConnectionFlag()
    
    
    go func() {
        
        connFlag.Set1(Copy(connFlag, connIn, connDest, current))
        fmt.Println("koniec kopiowania 1 ...", current)
    }()
    
    go func() {
        
        connFlag.Set2(Copy(connFlag, connDest, connIn, current))
        fmt.Println("koniec kopiowania 2 ...", current)
    }()
    
    connFlag.Wait()
    
    fmt.Println("koniec", current)
}


//http://golang.org/pkg/net/#example_Listener
//http://golang.org/src/io/io.go?s=12247:12307#L340
//http://www.badgerr.co.uk/2011/06/20/golang-away-tcp-chat-server/


func Copy(flag *connectionFlag, src *net.TCPConn, dst *net.TCPConn, current int) (int64, error, error) {
    
    
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
            
            err2 = filterTimeout(err2)
            
            if nw > 0 {
                written += int64(nw)
            }
            
            if err1 != nil || err2 != nil {
                return written, err1, err2
            }
            
            if nr != nw {
                return written, err1, io.ErrShortWrite
            }
            
            continue
        }
        
        if err1 != nil {
            return written, err1, nil
        }
    }
}

func filterTimeout(err error) error {
    
    if err == io.EOF {
        return err
    }
    
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

