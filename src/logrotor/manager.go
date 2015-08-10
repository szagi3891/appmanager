package logrotor

import (
    "../errorStack"
    "../utils"
    configModule "../config"
)

func Init(config *configModule.File) (*Manager, *errorStack.Error) {
    
    //TODO - wyszczególnij wszystkie pliki gz znajdujące się w katalogu z logami
    //przy robieniu nowego gz pliku trzeba dodać wynikowy rozmiar tego pliku do sumy zajętości całego katalogu
    
    //, logDir string, rotatesize, rotatetime int
    
    return &Manager{
        config : config,
    }, nil
}


type Manager struct {
    config *configModule.File
}


func (self *Manager) New(name string, stdout bool) *LogWriter {
    
    //stdout true  : .out
    //stdout false : .err
    
    ext := "err"
    
    if stdout {
        ext = "out"
    }
        //., config.GetRotatesize(), config.GetRotatetime())    logDir
    
    return newLogWriter(self.config.GetLogDir() + "/" + name + "_" + utils.GetCurrentTimeName() + "." + ext)
}
