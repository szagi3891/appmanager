package main

import (
    "os"
    "os/signal"
    "net"
    "io"
    "fmt"
    "sync"
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
    
    fmt.Println("nawiązano nowe połączenie")
    
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
    
    
    
    setClose := make(chan bool)
    
    go func() {
        
        io.Copy(connIn, connDest)
        
        fmt.Println("koniec kopiowania 1 ...", current)
        //close(isClose)
        
        setClose <- false
    }()
    
    go func() {
        
        io.Copy(connDest, connIn)
        
        fmt.Println("koniec kopiowania 2 ...", current)
        //close(isClose)
        
        setClose <- false
    }()
    
    fmt.Println("czytam", current)
    
    <- setClose
    
    fmt.Println("koniec", current)
    
    //<- isClose
    
}


//http://golang.org/pkg/net/#example_Listener
//http://golang.org/src/io/io.go?s=12247:12307#L340
//http://www.badgerr.co.uk/2011/06/20/golang-away-tcp-chat-server/

/*
func Copy(dst Writer, src Reader) (int64, error) {
    
    written := int64(0)
    
    buf := make([]byte, 32*1024)
    
    for {
        
        nr, er := src.Read(buf)
        
        if nr > 0 {
            
            nw, ew := dst.Write(buf[0:nr])
            
            if nw > 0 {
                written += int64(nw)
            }
            
            if ew != nil {
                err = ew
                break
            }
            
            if nr != nw {
                return written, ErrShortWrite
            }
        }
        
        if er == EOF {
            break
        }
        
        if er != nil {
            return written, er
        }
    }
    
    return written, err
}
*/

/*

    c.SetReadDeadline(time.Now())
if _, err := c.Read(one); err == io.EOF {
  l.Printf(logger.LevelDebug, "%s detected closed LAN connection", id)
  c.Close()
  c = nil
} else {
  var zero time.Time
  c.SetReadDeadline(time.Time{})
}
*/

