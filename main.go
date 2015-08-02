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
    
    mutex sync.Mutex
    isClose chan bool
    
    ord1Flag bool
    ord1Size int64
    ord1Err1 error
    ord1Err2 error
    
    ord2Flag bool
    ord2Size int64
    ord2Err1 error
    ord2Err2 error
}

func newConnectionFlag() *connectionFlag {
    
    return &connectionFlag{
        isClose : make(chan bool),
    }
}

func (self *connectionFlag) isDone() bool {
    
    self.mutex.Lock()
    defer self.mutex.Unlock()
    
    return self.ord1Flag == true && self.ord2Flag == true
}
    
func (self *connectionFlag) Set1(size int64, err1, err2 error) {
    
    self.mutex.Lock()
    defer self.mutex.Unlock()
    
    self.ord1Flag = true
    self.ord1Size = size
    self.ord1Err1 = err1
    self.ord1Err2 = err2
    
    if self.ord1Flag == true && self.ord2Flag == true {
        close(self.isClose)
    }
}

func (self *connectionFlag) Set2(size int64, err1, err2 error) {
    
    self.mutex.Lock()
    defer self.mutex.Unlock()
    
    self.ord2Flag = true
    self.ord2Size = size
    self.ord2Err1 = err1
    self.ord2Err2 = err2
    
    if self.ord1Flag == true && self.ord2Flag == true {
        close(self.isClose)
    }
}

func (self *connectionFlag) Wait() {
    
    <- self.isClose
    
    fmt.Println("1: size :", self.ord1Size)
    fmt.Println("1: err1 :", self.ord1Err1)
    fmt.Println("1: err2 :", self.ord1Err2)
    fmt.Println("2: size :", self.ord2Size)
    fmt.Println("2: err1 :", self.ord2Err1)
    fmt.Println("2: err2 :", self.ord2Err2)
    //<- isClose
}


func main(){
    
    interruptNotify()
    
    
        //81 -> 9999
    
    addProxy, err1 := net.ResolveTCPAddr("tcp", "localhost:8888")
    
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
    
    
    //fmt.Println("nawiązano nowe połączenie")
    
    addrServer, errParse := net.ResolveTCPAddr("tcp", "localhost:9999")
    //addrServer, errParse := net.ResolveTCPAddr("tcp", "onet.pl:80")
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
        
        connFlag.Set1(Copy(connFlag, connIn, connDest))
        fmt.Println("koniec kopiowania 1 ...", current)
    }()
    
    go func() {
        
        connFlag.Set2(Copy(connFlag, connDest, connIn))
        fmt.Println("koniec kopiowania 2 ...", current)
    }()
    
    connFlag.Wait()
    
    fmt.Println("koniec", current)
}


//http://golang.org/pkg/net/#example_Listener
//http://golang.org/src/io/io.go?s=12247:12307#L340
//http://www.badgerr.co.uk/2011/06/20/golang-away-tcp-chat-server/


func Copy(flag *connectionFlag, src *net.TCPConn, dst *net.TCPConn) (int64, error, error) {
    
    written := int64(0)
    buf     := make([]byte, 32*1024)
    
    src.SetReadDeadline(time.Now().Add(time.Second))
    dst.SetWriteDeadline(time.Now().Add(time.Second))
    
    for {
        
        if flag.isDone() {
            return written, nil, nil
        }
        
        nr, err1 := src.Read(buf)
        
        if nr > 0 {
            
            nw, err2 := dst.Write(buf[0:nr])
            
            if nw > 0 {
                written += int64(nw)
            }
            
            if err1 != nil || err2 != nil {
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


