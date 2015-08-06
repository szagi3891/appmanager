package applog

import (
    "os"
)

var pipeOut chan *string
var pipeErr chan *string

func WriteStdLn(newline string) {
    
    if pipeOut == nil {
        
        pipeOut = make(chan *string)
        
        go printer(os.Stdout, pipeOut)
    }
    
    pipeOut <- &newline
}

func WriteErrLn(newline string) {
    
    if pipeErr == nil {
        
        pipeErr = make(chan *string)
        
        go printer(os.Stderr, pipeErr)
    }
    
    pipeErr <- &newline
}

func printer(file *os.File, pipe chan *string) {
    
    for logLine := range pipe {
        
        lineOut := []byte(*logLine + "\n")
        count   := len(lineOut)
        
        n, err := file.Write(lineOut)
        
        if err != nil {
            panic(err)
        }
        
        if n != count {
            panic("różne długości")
        }
    }
}