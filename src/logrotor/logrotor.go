package logrotor

import (
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



func (self *Manager) New(name string) *LogWriter {
    
    return &LogWriter{
    }
}


type LogWriter struct {
    
}



type logPipe struct {
    timestart int           //czas w którym trafił pierwszy log do tego strumienia
    size      int           //rozmiar danych które siedzą w tym pliku
}

