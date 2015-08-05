package main


import (
    "os"
    "os/signal"
    "fmt"
    proxyModule   "./src/proxy"
    backendModule "./src/backend"
    "time"
    configModule "./src/config"
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
    
    
    //testowanie czy katalog z logami istnieje,
    //testowanie czy katalog z buildami istnieje
    
    
    /*
        lista

            nie można przełączyć builda, jeśli trwa aktualnie przełączanie nowej wersji
            przy uruchomionej aplikacji pokazywać ile aktualnie trwa połączeń

            ewentualnie można zrobić strumieniowanie danych na temat ilości połączeń przełączanych z jednej wersji aplikacji na drugą
    */
    
    managerBackend, errInitManager := backendModule.Init(config.GetGoCmd(), config.GetAppDir(), config.GetBuildDir(), config.GetAppMain(), config.GetAppUser(), config.GetGopath(), config.GetPortFrom(), config.GetPortTo())
    
    
    if errInitManager != nil {
        
        fmt.Println(errInitManager)
        os.Exit(1)
    }
    
    
    /*
    errMake := managerBackend.MakeBuild()
    
    fmt.Println("errMake:", errMake)
    
    return
    */
    
    
    
    
    /*
        todo - zrobić obsługę parametrów dotyczących rotowania logów
        rotatesize : granica w bajtach po której przekroczeniu rozpoczynamy nowy plik
        rotatetime : granica w sekundach po przekroczeniu której rozpoczynamy nowy plik
        
        nowy proces backendowy
            inicjuje dwa rotatory w ramach strumieni wyjściowych z procesu
        przy zabijaniu procesu trzeba wysłać sygnały do tych dwóch obiektów odnośnie tego że mają pozamykać ładnie
        pliki z logami na których operowały ...
        
        przy zrotowaniu pliku
            odpallić nową gorutinę która będzie kompresowałą stary plik i zrobi z niego gzip-a
    
    */
                            //build w kontekście tego katalogu będzie odpalany
    
    
    /*
    keys := []string{"port"}
    
    for _, paramName := range keys {
        
        value, isSet := mapConfig[paramName]
        
        if isSet {
            outConfig[paramName] = value
        } else {
            return nil, errorStack.Create("Brak klucza: " + paramName)
        }
    }
    */
    
    
    
    backend1, errCreate1 := managerBackend.New("build_20150804140526_381cd491cd49208d4a667912ef55fb78ab8469b1")       //127.0.0.1
    
    
    if errCreate1 != nil {
        
        fmt.Println(errCreate1)
        os.Exit(1)
    }
    
    
    
    proxy, errStart := proxyModule.New("127.0.0.1:8888", "../appmanager_log", backend1)
    
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
    
    
    /*
        build_3_20150801_143212     - numer kolejny i data utworzenia
        

        makeBuild
            funkcja robi nowego builda i zwraca jego nazwę
        
        proxy.runBuild("nazwa buildu", port)
            nazwa buildu, numer portu
    */
}


//<-chan
