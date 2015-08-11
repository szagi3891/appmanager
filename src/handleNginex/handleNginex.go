package handleNginex

import (
    "io/ioutil"
    "../errorStack"
    //"fmt"
    "strconv"
    "strings"
    "../utils"
)

func Switch(port int, confTpl, confDest, confRun string) *errorStack.Error {
    
    content, errRead := ioutil.ReadFile(confTpl)
    
    if errRead != nil {
        return errorStack.From(errRead)
    }
    
    contentOut, errConvert := replacePort(&content, port)
    
    if errConvert != nil {
        return errConvert
    }
    
    errSave := ioutil.WriteFile(confDest, []byte(contentOut), 0644)
    
    if errSave != nil {
        return errorStack.From(errSave)
    }
    
    _, err := utils.ExecCommand(confRun)
    
    return err
}


func replacePort(content *[]byte, port int) (string, *errorStack.Error) {
    
    contentStr := string(*content)
    
    index := strings.Index(contentStr, "##port##")
    
    if index >= 0 {
        
        return contentStr[0:index] + strconv.FormatInt(int64(port), 10) + contentStr[index+8:], nil
        
    } else {
        return "", errorStack.Create("W pliku szablonowym nie znaleziono sekwencji znaków: ##port##")
    }
}


