package logrotor

import (
    "fmt"
    "../errorStack"
)


func Init(logDir string, rotatesize, rotatetime int) (*Manager, *errorStack.Error) {
    
    //TODO - wyszczególnij wszystkie pliki gz znajdujące się w katalogu z logami
    //przy robieniu nowego gz pliku trzeba dodać wynikowy rozmiar tego pliku do sumy zajętości całego katalogu
    
    return &Manager{
        logDir     : logDir,
        rotatesize : rotatesize,
        rotatetime : rotatetime,
        
    }, nil
}


type Manager struct {
    logDir     string
    rotatesize int
    rotatetime int
}



func (self *Manager) New(name string, stdout bool) *LogWriter {
    
    //stdout true  : .out
    //stdout false : .err
    
    panic("TODO - trzeba osbłużyć tworzenie nowej rurki zbierającej dane i zapisującej dane do pliku z logiem")
    return &LogWriter{
    }
}


type LogWriter struct {
    
}


func (self *LogWriter) Write(p []byte) (n int, err error) {
    
    fmt.Println("logwriter", string(p))
    return len(p), nil
}


//Read(p []byte) (n int, err error)
//Write(p []byte) (n int, err error)



type logPipe struct {
    timestart int           //czas w którym trafił pierwszy log do tego strumienia
    size      int           //rozmiar danych które siedzą w tym pliku
}

