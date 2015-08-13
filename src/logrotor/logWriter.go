package logrotor


func newLogWriter(pathFile string, ext string) *logWriter {
    
    pipe    := make(chan *[]byte)
    isClose := make(chan bool)
    
    go runLogGroup(pipe, isClose, pathFile, ext)
    
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

func runLogGroup(pipe chan *[]byte, isClose chan bool, pathFile string, ext string) {
    
    buf           := []*[]byte{}
    size          := 0
    
    sendToFile    := make(chan []*[]byte)
    isCloseWriter := make(chan bool)
    
    
    go SaveData(pathFile, sendToFile, isCloseWriter, ext)
    
    
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


func SaveData(pathFile string, saveIn chan []*[]byte, isCloseWriter chan bool, ext string) {
    
    
    file := createFile(pathFile, ext)
    
    for {
        
        newData := <- saveIn
        
        
        if newData == nil {
            
            file.close()
            close(isCloseWriter)
            
            return
        }
        
        
        for _, chankData := range newData {
            file.write(chankData)
        }
        
        //jeśli przekroczono rozmiar,
            //zamknij plik
            //utwórz nowy
    }
}



//type logPipe struct {
    //timestart int           //czas w którym trafił pierwszy log do tego strumienia
    //size      int           //rozmiar danych które siedzą w tym pliku
//}
