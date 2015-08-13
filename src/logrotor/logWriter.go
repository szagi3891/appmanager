package logrotor

import (
    configModule "../config"
)

func newLogWriter(name string, ext string, config *configModule.File) *logWriter {
    
    pipe    := make(chan *[]byte)
    isClose := make(chan bool)
    
    go runLogGroup(pipe, isClose, name, ext, config)
    
    return &logWriter{
        pipe    : pipe,
        isClose : isClose,
    }
}

type logWriter struct {
    pipe    chan *[]byte
    isClose chan bool
}

func copyArr(src *[]byte) *[]byte {
    
    cop := []byte{}
    
    for _, char := range *src {
        cop = append(cop, char)
    }
    
    return &cop
}

func (self *logWriter) WriteStringLn(p string) () {
    
    p2 := []byte(p + "\n")
    
    p3 := copyArr(&p2)
    
    self.pipe <- p3
}


//to dla applikacji które przekazują nam swoje uformowane logi
func (self *logWriter) Write(p []byte) (n int, err error) {
    
    p3 := copyArr(&p)
    
    
    self.pipe <- p3
    
    
    return len(p), nil
}


func (self *logWriter) Stop() {
    
    self.pipe <- nil
    
            //blokująco trzeba wszystko z tym logiem związane pozamykać
    <- self.isClose
}

//func createCondision(time 

/*
    I  stopień - grupujemy dane co ileś tam bajtów, po to aby niepotrzebnie nie męczyć dysku ciągłymi zapisami
    II stopień - otrzymaną paczkę od razu zapisujemy do pliku
*/

func runLogGroup(pipe chan *[]byte, isClose chan bool, name string, ext string, config *configModule.File) {
    
    buf           := []*[]byte{}
    size          := 0
    
    sendToFile    := make(chan []*[]byte)
    isCloseWriter := make(chan bool)
    
    
    go SaveData(name, sendToFile, isCloseWriter, ext, config)
    
    
    reciveData := func(newData *[]byte) bool {
        
        if newData == nil {
            
                                    //wyślij resztki
            sendToFile <- buf
                                    //zakończ działanie
            sendToFile <- nil
            
            close(isClose)
            return true

        }
        
        buf = append(buf, newData)
        size = size + len(*newData)
        return false
    }
    
    
    for {
        
        if size > 16 * 2<<10 {
            
            select {
                
                case newData := <- pipe : {
                    
                    if reciveData(newData) {
                        return
                    }
                }

                case sendToFile <- buf : {

                    buf  = []*[]byte{}
                    size = 0
                }
            }
        
        } else {
            
            newData := <- pipe
            
            if reciveData(newData) {
                return
            }
        }
    }
}


func SaveData(name string, saveIn chan []*[]byte, isCloseWriter chan bool, ext string, config *configModule.File) {
    
    
    file := CreateFile(config.GetLogDir() + "/" + name, ext)
    
    for {
        
        newData := <- saveIn
        
        
        if newData == nil {
            
            file.Close()
            close(isCloseWriter)
            
            return
        }
        
        
        for _, chankData := range newData {
            file.Write(chankData)
        }
        
        
                                        //rotuj plik
        
        if file.Size() > config.GetRotatesize() || file.GetTimeExist() > config.GetRotatetime() {
            
            file.Close()
            file = CreateFile(config.GetLogDir() + "/" + name, ext)   
        }
    }
}



//type logPipe struct {
    //timestart int           //czas w którym trafił pierwszy log do tego strumienia
    //size      int           //rozmiar danych które siedzą w tym pliku
//}
