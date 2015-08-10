package logrotor

import (
    "fmt"
    "os"
    "../errorStack"
    "../applog"
)

func newLogWriter(pathFile string) *logWriter {
    
    pipe    := make(chan *[]byte)
    isClose := make(chan bool)
    
    go runLogGroup(pipe, isClose, pathFile)
    
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

func (self *logWriter) WriteString(p string) () {
    
    //fmt.Println(p)
    
    p2 := []byte(p)
    
    p3 := copyArr(&p2)
    
    self.pipe <- p3
}

func (self *logWriter) Write(p []byte) (n int, err error) {
    
    //fmt.Println(string(p))
    
    
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

func runLogGroup(pipe chan *[]byte, isClose chan bool, pathFile string) {
    
    buf           := []*[]byte{}
    size          := 0
    
    sendToFile    := make(chan []*[]byte)
    isCloseWriter := make(chan bool)
    
    
    go SaveData(pathFile, sendToFile, isCloseWriter)
    
    
    reciveData := func(newData *[]byte) bool {

        if newData == nil {
            
            fmt.Println("Otrzymałem nil-a - trzeba zamknąć tego loga - wysyłam resztki z bufora")
            
            sendToFile <- buf
            
            fmt.Println("Otrzymałem nil-a - trzeba zamknąć tego loga - teraz wysyłam nila")
            
            sendToFile <- nil       //zakończ działanie

            fmt.Println("Otrzymałem nil-a - trzeba zamknąć tego loga - się wreszcie zamykam")
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
    
    var file *os.File
    
    
    for {
        
        newData := <- saveIn
        
        
        if file == nil {
            
            fileNew, errCreate := os.OpenFile(pathFile, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0600)
            
            if errCreate != nil {
                
                applog.WriteErrLn(errorStack.From(errCreate).String())
                return
            }
            
            file = fileNew
        }
        
        
        if newData == nil {
            
            if file != nil {
                file.Close()
            }
            
            close(isCloseWriter)
            
            return
        }
        
        
        for _, chankData := range newData {
            
            fmt.Print(string(*chankData))
            n, err := file.Write(*chankData)
            
            if err != nil {
                
                applog.WriteErrLn(errorStack.From(err).String())
                continue
            }
            
            if n != len(*chankData) {
                
                applog.WriteErrLn(errorStack.Create("nieprawidłowa ilość zapisanych znaków do pliku").String())
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

