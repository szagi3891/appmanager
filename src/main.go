package main


import (
    "os"
    "os/signal"
    backendModule "./backend"
    configModule "./config"
    logrotorModule "./logrotor"
    "./wwwpanel"
    "./applog"
)


func main(){
    
    code := run()
    
    os.Exit(code)
}

func run() int {
    
    if len(os.Args) != 2 {
        applog.WriteErrLn("Spodziewano się dokładnie dwóch parametrów")
        return 1
    }
    
    
    config, errParse := configModule.Parse(os.Args[1])
    
    if errParse != nil {
        applog.WriteErrLn(errParse.String())
        return 1
    }
    
    
    logrotor, errLogrotor := logrotorModule.Init(config)
    
    if errLogrotor != nil {
        
        applog.WriteErrLn(errLogrotor.String())
        return 1
    }
    
    
                                                //logi aplikacji    
    logs := logrotor.NewLogs("appmanager")
    
    defer logs.Stop()
    
    
    managerBackend, errInitManager := backendModule.Init(config, logrotor, logs)
    
    
    if errInitManager != nil {
        
        applog.WriteErrLn(errInitManager.String())
        return 1
    }
    
    defer managerBackend.Stop()
    
    
    //TODO
    //zrobić obsługę zmiennej : rotatetotalsize
    //stare pliki z logami będą kasowane automatycznie żeby nie zapchać dysku
    
    //TODO
    //zdarzyło się tak że przy wyłącznaiu programu, nie zdążył się skompresować główny log
    //obserwować
    
    /*
        //https://www.youtube.com/watch?v=InG72scKPd4
        //debuger w go napisany do go


        //listowanie wszystkich aplikacji słuchających na porcie
        //netstat -tulpn
        //lsof -i - coś podobnego
        
    
                //łapać jeszcze jedno zdarzenie
        // Handle SIGINT and SIGTERM.
        ch := make(chan os.Signal)
        signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
        log.Println(<-ch)
        
        kill -l - wyświetla dopuszczalne nazwy sygnałów
        kill -kill 1234 - wysłanie sygnału SIGKILL do procesu z pid=1234
        kill -9 1234 - to samo wyżej, bo SIGKILL ma nr 9
        kill -int 1234 - wysłanie sygnału SIGINT do procesu z pid=1234
        kill -2 1234 - to samo wyżej, bo SIGINT ma nr 2
        
        //sposób na fajne proxowanie http
        //http://siberianlaika.ru/node/29/
    */
    
    
                        //start panelu do zarządzania konfiguracją proxy
    closePanel, errStartPanel := wwwpanel.Start(8889, logs, managerBackend)
    
    if errStartPanel != nil {
        applog.WriteErrLn(errStartPanel.String())
        return 1
    }
    
    defer closePanel()
    
    
    
    signalInterrupt := make(chan os.Signal, 1)
    signal.Notify(signalInterrupt, os.Interrupt)
    <- signalInterrupt
    
    logs.Std.WriteStringLn("wyłączam proxy z kodem wyjścia 0")
    
    return 0
}