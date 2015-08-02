package main


import (
    "os"
    "os/signal"
    "net"
    "fmt"
    "./src/handleConn"
    "./src/serverBackend"
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



func main(){
    
    interruptNotify()
    
    
    addProxy, err1 := net.ResolveTCPAddr("tcp", "localhost:8888")
    
    if err1 != nil {
        panic(err1)
    }
    
    
    fmt.Println("start proxy: ", addProxy)
    
    
	listener, err := net.ListenTCP("tcp", addProxy)
    
	if err != nil {
		panic(err)
	}
    
    serverBackend := serverBackend.New("localhost:9999")
    
	for {
        
		conn, err := listener.AcceptTCP()
        
		if err != nil {
			panic(err)
		}
        
        errConnectect := handleConn.HandleConn(serverBackend, conn)
        
        if errConnectect != nil {
            panic(errConnectect)
        }
	}
}


//<-chan
