package logrotor

import (
    "../errorStack"
    "../utils"
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
    
    ext := "err"
    
    if stdout {
        ext = "out"
    }
    
    return newLogWriter(self.logDir + "/" + name + "_" + utils.GetCurrentTimeName() + "." + ext)
}
