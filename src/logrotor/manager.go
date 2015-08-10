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

type Logs struct {
    Std *logWriter
    Err *logWriter
}

func (self *Logs) Stop() {
    
    self.Std.Stop()
    self.Err.Stop()
}

func (self *Manager) NewLogs(name string) *Logs {
    
    return &Logs {
        Std : self.newSingleLog(name, true),
        Err : self.newSingleLog(name, false),
    }
}

func (self *Manager) newSingleLog(name string, stdout bool) *logWriter {
    
    //stdout true  : .out
    //stdout false : .err
    
    ext := "err"
    
    if stdout {
        ext = "out"
    }
        //., config.GetRotatesize(), config.GetRotatetime())    logDir
    
    return newLogWriter(self.config.GetLogDir() + "/" + name + "_" + utils.GetCurrentTimeName() + "." + ext)
}
