package main


import (
    "os"
    "os/signal"
    "fmt"
    proxyModule   "./src/proxy"
    backendModule "./src/backend"
    //"time"
    configModule "./src/config"
    logrotorModule "./src/logrotor"
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
    
    if len(os.Args) != 2 {
        fmt.Println("Spodziewano się dokładnie dwóch parametrów")
        os.Exit(1)
    }
    
    
    interruptNotify()
    
    
    config, errParse := configModule.Parse(os.Args[1])
    
    if errParse != nil {
        
        fmt.Println(errParse)
        os.Exit(1)
    }
    
    
    
    //listowanie wszystkich aplikacji słuchających na porcie
    //netstat -tulpn
    //lsof -i - coś podobnego
    
    
    logrotor, errLogrotor := logrotorModule.Init(config.GetLogDir(), config.GetRotatesize(), config.GetRotatetime())
    
    if errLogrotor != nil {
        
        fmt.Println(errLogrotor)
        os.Exit(1)
    }
    
    
    //TODO - Trzeba zrobić jakaś strukturkę zbierającą te dwie struktury
    
                                                //logi aplikacji
    appStdout := logrotor.New("appmanager", true)
    appStderr := logrotor.New("appmanager", false)
    
    
    
    managerBackend, errInitManager := backendModule.Init(logrotor, appStdout, appStderr, config.GetGoCmd(), config.GetAppDir(), config.GetBuildDir(), config.GetAppMain(), config.GetAppUser(), config.GetGopath(), config.GetPortFrom(), config.GetPortTo())
    
    
    if errInitManager != nil {
            
        fmt.Println(errInitManager)
        os.Exit(1)
    }
    
    
    /*
    errMake := managerBackend.MakeBuild()
    fmt.Println("errMake:", errMake)
    return
    */
    
    
    //TODO - zrobić pingowanie w nowy beckend, gotowość dopiero ma zgłosić jeśli będzie odpowiadał na zadanym porcie
    
    //TODO - zrobić rotowanie i gzipowanie logów
    
    
    /*

        zrobić obsługę zmiennej : rotatetotalsize
        stare pliki z logami będą kasowane automatycznie żeby nie zapchać dysku
    */
    
    backend1, errCreate1 := managerBackend.New("build_20150806232318_b82180a227005c7db9d49f054970c5d0dc8df7d1")
    
    
    if errCreate1 != nil {
        
        fmt.Println(errCreate1)
        os.Exit(1)
    }
    
    
    
    proxy, errStart := proxyModule.New(appStderr, config.GetPortMain(), logrotor, backend1)
    
    if errStart != nil {
        panic(errStart)
    }
    
    
    fmt.Println("start proxy: ", proxy)
    
    
    /*
    time.Sleep(time.Second * 10)
    
    fmt.Println("przełączam backend")
    
    backend2, errBackend2 := managerBackend.New("build_20150806221324_f1c4c02114226c90a4a202dcbeec65366970fb55")
    
    if errBackend2 != nil {
        
        fmt.Println(errBackend2)
        os.Exit(1)
    }
    
    proxy.Switch(backend2)
    */
    
    
    
            //nie zakańczaj się
    stop := make(chan bool)
    <- stop
    
    
    //start i stop nowego backandu
    
    
    //TODO - nowy byt, struktura reprezentująca listę buildów ...
    
}
