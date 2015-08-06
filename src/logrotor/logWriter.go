package logrotor

import (
    "fmt"
)

func newLogWriter(pathFile string) *LogWriter {
    
    pipe    := make(chan *[]byte)
    isClose := make(chan bool)
    
    go runLogPipe(pipe, isClose, pathFile)
    
    return &LogWriter{
        pipe    : pipe,
        isClose : isClose,
    }
}

type LogWriter struct {
    pipe    chan *[]byte
    isClose chan bool
}

func (self *LogWriter) Write(p []byte) (n int, err error) {
    
    self.pipe <- &p
    return len(p), nil
}


func (self *LogWriter) Stop() {
    
    close(self.pipe)
    
            //blokująco trzeba wszystko z tym logiem związane pozamykać
    <- self.isClose
}


func runLogPipe(pipe chan *[]byte, isClose chan bool, pathFile string) {
    
    buf       := []*[]byte{}
    size      := 0
    pipeWrite := make(chan []*[]byte)
    
    go runWritePipe(pipeWrite, isClose, pathFile)
    
    for newData := range pipe {
        
        buf = append(buf, newData)
        size = size + len(*newData)
        
        if size > 1024 {
            
            pipeWrite <- buf
            
            buf  = []*[]byte{}
            size = 0
        }
    }
    
    pipeWrite <- buf
    close(pipeWrite)
}


func runWritePipe(pipeWrite chan []*[]byte, isClose chan bool, pathFile string) {
    
    for newPack := range pipeWrite {
        
        fmt.Println("nowa paczka danych: ", newPack, len(newPack))
    }
    
    
        //daj sygnał światu że zakończyliśmy swoją egzystencję i żeby o nas pamietał
    close(isClose)
}


//type logPipe struct {
    //timestart int           //czas w którym trafił pierwszy log do tego strumienia
    //size      int           //rozmiar danych które siedzą w tym pliku
//}

