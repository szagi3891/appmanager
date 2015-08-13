package logrotor

import (
    "os"
    "io"
    "compress/gzip"
    "../utils"
    "../errorStack"
    "../applog"
)


type fileLog struct {
    filePath string
    file     *os.File
    isClose  bool
}

func createFile(path, ext string) *fileLog {
    
    return &fileLog {
        filePath : path + "_" + utils.GetCurrentTimeName() + "." + ext,
        file     : nil,
        isClose  : false,
    }
}

func (self *fileLog) write(chankData *[]byte) {
    
    //jeśli dyskryptor nie otworzony to dopiero teraz go otwórz
    
    if self.file == nil {
        
        fileNew := openFile(self.filePath)
        
        if fileNew == nil {
            return
        }
        
        self.file = fileNew
    }
    
    n, err := self.file.Write(*chankData)
    
    if err != nil {
        
        applog.WriteErrLn(errorStack.From(err).String())
        return
    }
    
    if n != len(*chankData) {
        applog.WriteErrLn(errorStack.Create("nieprawidłowa ilość zapisanych znaków do pliku").String())
    }
}

func (self *fileLog) close() {
    
    if self.file != nil {
        
        if self.isClose == false {
            
            self.isClose = true
            self.file.Close()
            
            //przy zamykani odpalana jest kompresja do gz
            
            errCompress := compress(self.filePath, self.filePath + ".gz")
            
            if errCompress == nil {
                
                errRemove := os.Remove(self.filePath)
                
                if errRemove != nil {
                    applog.WriteErrLn(errorStack.From(errRemove).String())
                }
                
            } else {
                
                applog.WriteErrLn(errCompress.String())
            }
        }
    }
}


func openFile(pathFile string) *os.File {
    
    fileNew, errCreate := os.OpenFile(pathFile, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0600)

    if errCreate != nil {

        applog.WriteErrLn(errorStack.From(errCreate).String())
        return nil
    }
    
    return fileNew
}




/*
na docelowym serwerze

-rw-r--r--  1 wm wm       629 sie 13 08:26 compress.go
-rw-r-----  1 wm wm 189252169 sie 13 08:10 dane.txt
-rw-------  1 wm wm  13817446 sie 13 08:26 dane.txt.gz

wm@wolnemedia:~/_logtest$ time go run compress.go 

real	0m7.376s
user	0m5.744s
sys	0m0.320s
*/


func compress(pathIn string, pathOut string) *errorStack.Error {
    
    
    fileIn, errFileIn := os.Open(pathIn)
    
    if errFileIn != nil {
        return errorStack.From(errFileIn)
    }
    
    defer fileIn.Close()
    
    
    fileOut, errCreate := os.OpenFile(pathOut, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0600)
    
    if errCreate != nil {
        return errorStack.From(errCreate)
    }
    
    defer fileOut.Close()
    
    
    writer := gzip.NewWriter(fileOut)
    defer writer.Close()
    
    _, errCopy := io.Copy(writer, fileIn)
    
    if errCopy != nil {
        return errorStack.From(errCopy)
    }
    
    return nil
}






