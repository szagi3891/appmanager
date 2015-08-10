package main


import (
    "os"
    "os/signal"
    "fmt"
    proxyModule   "./proxy"
    backendModule "./backend"
    configModule "./config"
    logrotorModule "./logrotor"
    "./wwwpanel"
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
    
    //panic("TODO - jest problem na wyłączaniu aplikacji - wyłączony backend trzeba usunać z mapy")
    
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
    
    
    
    managerBackend, errInitManager := backendModule.Init(config.GetPortMain(), logrotor, appStdout, appStderr, config.GetGoCmd(), config.GetAppDir(), config.GetBuildDir(), config.GetAppMain(), config.GetAppUser(), config.GetGopath(), config.GetPortFrom(), config.GetPortTo())
    
    
    if errInitManager != nil {
            
        fmt.Println(errInitManager)
        os.Exit(1)
    }
    
    
    
    //TODO - zrobić pingowanie w nowy beckend, gotowość dopiero ma zgłosić jeśli będzie odpowiadał na zadanym porcie
    
    //TODO - zrobić rotowanie i gzipowanie logów
    
    //TODO - trzeba pozbyć się logowania poprzez fmt
    
    //zrobić obsługę zmiennej : rotatetotalsize
    //stare pliki z logami będą kasowane automatycznie żeby nie zapchać dysku
    
    //obiekt menedżera backendu powinien być przykryty obiektem proxy
    //tylko obiektem proxy powinniśmy sterować żeby sobie wysterować to co trzeba
    
    
    
    proxy, errStart := proxyModule.New(appStderr, config.GetPortMain(), logrotor, managerBackend)
    
    if errStart != nil {
        panic(errStart)
    }
    
    
    
                        //start panelu do zarządzania konfiguracją proxy
    wwwpanel.Start(8889, appStderr, managerBackend, proxy)
    
    
    
    
    
            //nie zakańczaj się
    stop := make(chan bool)
    <- stop
    
    
}
