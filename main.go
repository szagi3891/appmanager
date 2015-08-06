package main


import (
    "os"
    "os/signal"
    "fmt"
    proxyModule   "./src/proxy"
    backendModule "./src/backend"
    "time"
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
    
    
    
    managerBackend, errInitManager := backendModule.Init(logrotor, config.GetGoCmd(), config.GetAppDir(), config.GetBuildDir(), config.GetAppMain(), config.GetAppUser(), config.GetGopath(), config.GetPortFrom(), config.GetPortTo())
    
    
    if errInitManager != nil {
        
        fmt.Println(errInitManager)
        os.Exit(1)
    }
    
    
    //fmt.Println(config.GetRotatesize(), config.GetRotatetime())
    //os.Exit(1)
    
    /*
    errMake := managerBackend.MakeBuild()
    
    fmt.Println("errMake:", errMake)
    
    return
    */
    
    
    
    
    /*
        format logów ...
        
        build_czasutworzenia_sha1

        build_czasutworzenia_sha1_czasuruchomienia.out
        build_czasutworzenia_sha1_czasuruchomienia.err

        build_czasutworzenia_sha1_czasuruchomienia.out.gz

        appmanager_czasuruchomienia.out
        appmanager_czasuruchomienia.err
    */
    
    
    /*
        nowy proces backendowy
            inicjuje dwa rotatory w ramach strumieni wyjściowych z procesu
        przy zabijaniu procesu trzeba wysłać sygnały do tych dwóch obiektów odnośnie tego że mają pozamykać ładnie
        pliki z logami na których operowały ...
        
        przy zrotowaniu pliku
            odpallić nową gorutinę która będzie kompresowałą stary plik i zrobi z niego gzip-a
    */
    
    /*

        zrobić obsługę zmiennej : rotatetotalsize
        stare pliki z logami będą kasowane automatycznie żeby nie zapchać dysku
    */
    
    backend1, errCreate1 := managerBackend.New("build_20150804140526_381cd491cd49208d4a667912ef55fb78ab8469b1")
    
    
    if errCreate1 != nil {
        
        fmt.Println(errCreate1)
        os.Exit(1)
    }
    
    
    
    proxy, errStart := proxyModule.New(config.GetPortMain(), logrotor, backend1)
    
    if errStart != nil {
        panic(errStart)
    }
    
    
    fmt.Println("start proxy: ", proxy)
    
    
    
    time.Sleep(time.Second * 10)
    
    
    
    fmt.Println("przełączam backend")
    
    
    backend2, errBackend2 := managerBackend.New("build_20150805153050_381cd491cd49208d4a667912ef55fb78ab8469b1")
    
    if errBackend2 != nil {
        
        fmt.Println(errBackend2)
        os.Exit(1)
    }
    
    
    proxy.Switch(backend2)
    
    
    
    
            //nie zakańczaj się
    stop := make(chan bool)
    <- stop
    
    
    //start i stop nowego backandu
    
    
    //TODO - nowy byt, struktura reprezentująca listę buildów ...
    
}
