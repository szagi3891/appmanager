package logrotor

import (
    "os"
    "fmt"
)

func newLogWriter(pathFile string) *LogWriter {
    
    pipe    := make(chan *[]byte)
    isClose := make(chan bool)
    
    go runLogGroup(pipe, isClose, pathFile)
    
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
    
    //fmt.Println(string(p))
    self.pipe <- &p
    return len(p), nil
}


func (self *LogWriter) Stop() {
    
    self.pipe <- nil
    
            //blokująco trzeba wszystko z tym logiem związane pozamykać
    <- self.isClose
}

//func createCondision(time 

/*
    I  stopień - grupujemy dane co ileś tam bajtów, po to aby niepotrzebnie nie męczyć dysku ciągłymi zapisami
    II stopień - otrzymaną paczkę od razu zapisujemy do pliku
*/

func runLogGroup(pipe chan *[]byte, isClose chan bool, pathFile string) {
    
    buf           := []*[]byte{}
    size          := 0
    
    sendToFile    := make(chan []*[]byte)
    isCloseWriter := make(chan bool)
    
    
    go SaveData(pathFile, sendToFile, isCloseWriter)
    
    
    reciveData := func(newData *[]byte) bool {

        if newData == nil {
            
            sendToFile <- buf
            sendToFile <- nil       //zakończ działanie

            close(isClose)
            return true

        }
        
        buf = append(buf, newData)
        size = size + len(*newData)
        return false
    }
    
    
    for {
        
        if size > 5000 {
            
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


func SaveData(pathFile string, saveIn chan []*[]byte, isCloseWriter chan bool) {
    
    
    file, errCreate := os.OpenFile(pathFile, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0600)
    
    if errCreate != nil {
        fmt.Println(errCreate)
        panic("DAsdAS")
    }
    
    
    for {
        
        newData := <- saveIn

        if newData == nil {
            
            file.Close()
            close(isCloseWriter)

            return
        }
        
        
        for _, chankData := range newData {
            
            n, err := file.Write(*chankData)
            fmt.Println(string(*chankData))
            if err != nil {
                
                fmt.Println(err)
                panic("dasd")
            }
            
            if n != len(*chankData) {
                
                panic("nieprawidłowa ilość zapisanych znaków do pliku")
            }
        }
        
        

        //jeśli rotowanie

            //tak ->
            //zamknij plik
            //zgzipuj nową zawartość
            //usuń poprzednią wersję
            //otwórz nowy dyskryptor
            //zgłóś gotowość że zadanie wykonane
            //saveResult <- true

            //nie ->
            //zgłość gotowość że zadanie wykonane
            //saveResult <- true
    }
}


//type logPipe struct {
    //timestart int           //czas w którym trafił pierwszy log do tego strumienia
    //size      int           //rozmiar danych które siedzą w tym pliku
//}

