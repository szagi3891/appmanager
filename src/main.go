package main


import (
    "os"
    "os/signal"
    "fmt"
    backendModule "./backend"
    configModule "./config"
    logrotorModule "./logrotor"
    "./wwwpanel"
)


func main(){
    
    code := run()
    
    os.Exit(code)
}

func run() int {
    fmt.Println("DEBUG: start app")
    if len(os.Args) != 2 {
        fmt.Println("Spodziewano się dokładnie dwóch parametrów")
        return 1
    }
    
    
    config, errParse := configModule.Parse(os.Args[1])
    
    if errParse != nil {
        
        fmt.Println(errParse)
        return 1
    }
    
    
    
    //listowanie wszystkich aplikacji słuchających na porcie
    //netstat -tulpn
    //lsof -i - coś podobnego
    
    
    logrotor, errLogrotor := logrotorModule.Init(config)
    
    if errLogrotor != nil {
        
        fmt.Println(errLogrotor)
        return 1
    }
    
    
                                                //logi aplikacji    
    logs := logrotor.NewLogs("appmanager")
    
    defer logs.Stop()
    
    
    managerBackend, errInitManager := backendModule.Init(config, logrotor, logs)
    
    
    if errInitManager != nil {
        
        fmt.Println(errInitManager)
        return 1
    }
    
    defer managerBackend.Stop()
    
    
    
    //TODO - zrobić pingowanie w nowy beckend, gotowość dopiero ma zgłosić jeśli będzie odpowiadał na zadanym porcie
    
    //TODO - zrobić rotowanie i gzipowanie logów
    
    //TODO - trzeba pozbyć się logowania poprzez fmt
    
    //zrobić obsługę zmiennej : rotatetotalsize
    //stare pliki z logami będą kasowane automatycznie żeby nie zapchać dysku
    
    
                //TODO - ten sposób proxowania trzeba będzie zaimplementować
    //sposób na fajne proxowanie http
    //http://siberianlaika.ru/node/29/
    
    //wypaśnie o serwerach http ...
    //http://jan.newmarch.name/go/http/chapter-http.html
    
    
                        //start panelu do zarządzania konfiguracją proxy
    go wwwpanel.Start(8889, logs, managerBackend)
    
    
    
    
    signalInterrupt := make(chan os.Signal, 1)
    signal.Notify(signalInterrupt, os.Interrupt)
    
    
    for sig := range signalInterrupt {
        // sig is a ^C, handle it
        
                                        //zakończ tylko wtedy jeśli przyjdzie Ctrl + C        
        if sig == os.Interrupt {
            
            logs.Std.WriteStringLn("Przyszedł sygnał ctrl + c - wyłączam proxy")
            close(signalInterrupt)
            
        } else {
            
            logs.Std.WriteStringLn("Przyszedł nieznany sygnał - ignoruję")
        }
    }
    
    logs.Std.WriteStringLn("wyłączam proxy z kodem wyjścia 0")
    
    return 0
}