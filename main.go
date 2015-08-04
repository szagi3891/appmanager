package main


import (
    "os"
    "os/signal"
    "fmt"
    proxyModule   "./src/proxy"
    backendModule "./src/backend"
    "time"
    "./src/config"
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
    
    if len(os.Args) != 2 {
        fmt.Println("Spodziewano się dokładnie dwóch parametrów")
    }
    
    configField := []string{"port", "gocmd", "appdir", "appmain", "portfrom", "portto"}
    
    configFile, errParse := config.Parse(os.Args[1], &configField)
    
    if errParse != nil {
        
        fmt.Println(errParse)
        os.Exit(1)
    }
    
    
    //TODO - sprawdzić czy takie polecenie się wykona :
    //GOPATH=/home/wm/GOPATH /usr/local/go/bin/go
    
    
    
    //git rev-parse HEAD
    //pobranie aktualnego hasza komita z katalogu
    
    // struktura z listą buildów
    //, "../appmanager_build"
    
    
    //testowanie czy katalog z logami istnieje,
    //testowanie czy katalog z buildami istnieje
    
    
    /*
        lista

            nie można przełączyć builda, jeśli trwa aktualnie przełączanie nowej wersji
            przy uruchomionej aplikacji pokazywać ile aktualnie trwa połączeń

            ewentualnie można zrobić strumieniowanie danych na temat ilości połączeń przełączanych z jednej wersji aplikacji na drugą
    */
    
    //TODO - trzeba w pierwszej kolejności zrobić startowanie builda z konkretnej binarki
    
    //TODO - proxy będzie zarządzał numerami portów na których mogą pracować aplikacje
    
    
    
    
    //GOPATH=/home/wm/GOPATH /usr/local/go/bin/go build -o testowy_build ./cms/src/main.go
    //po zbudowaniu binarki powinien zostać zmieniony włąściciel oraz grupa
    
    //odpalanie procesu powinno się odbywać na określonym użytkowniku
    //https://golang.org/pkg/os/exec/#Cmd               -> SysProcAttr
    //https://golang.org/pkg/syscall/#SysProcAttr       -> Credential
    //https://golang.org/pkg/syscall/#Credential        -> Credential struct
    
    
    
    //"/usr/local/go/bin/go", 
    
    
    managerBackend := backendModule.Init("../wolnemedia", "../appmanager_build", 9990, 9999)
    
    
    managerBackend.MakeBuild()
    
    return
    
                            //build w kontekście tego katalogu będzie odpalany
    
    
    backend1, errCreate1 := managerBackend.New("build_20150803111303_3128586f693bb8005253ae17eb0d95ea25573b94")       //127.0.0.1
    
    
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
    
    
    backend2, errBackend2 := managerBackend.New("build_20150803111303_3128586f693bb8005253ae17eb0d95ea25573b94")
    
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
