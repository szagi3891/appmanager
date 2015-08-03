package backend

import (
    "fmt"
)

//Read(p []byte) (n int, err error)
//Write(p []byte) (n int, err error)

type outData struct {
    name string
}

func (self *outData) Write(p []byte) (n int, err error) {
    
    fmt.Println(self.name + " Write " + string(p))
    return len(p), nil
}
