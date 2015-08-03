package main


import (
    "os"
    "os/signal"
    "fmt"
    proxyModule   "./src/proxy"
    backendModule "./src/backend"
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



func main(){
    
    interruptNotify()
    
    
    backend1 := backendModule.New("127.0.0.1:9997")
    
    
    proxy, errStart := proxyModule.New("127.0.0.1:8888", "../appmanager_log", "../appmanager_build", backend1)
    
    if errStart != nil {
        panic(errStart)
    }
    
    
    fmt.Println("start proxy: ", proxy)
    
    
    time.Sleep(time.Second * 10)
    
    backend2 := backendModule.New("127.0.0.1:9998")
    
    fmt.Println("przełączam backend")
    
    proxy.Switch(backend2)
    
    
            //nie zakańczaj się
    stop := make(chan bool)
    <- stop
    
    
    //start i stop nowego backandu
    
    
    /*
        build_3_20150801_143212     - numer kolejny i data utworzenia


        makeBuild
            funkcja robi nowego builda i zwraca jego nazwę
        
        runBuild
            nazwa buildu, numer portu
    */
}


//<-chan
