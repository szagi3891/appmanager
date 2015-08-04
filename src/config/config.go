package config

import (
    //"fmt"
    "strings"
    "os"
    "bufio"
    "../errorStack"
)


type File struct {
    config  map[string]string
}


func Parse(path string, keys *[]string) (*File, *errorStack.Error) {
    
    lines, errRead := readLines(path)
    
    if errRead != nil {
        return nil, errRead
    }
    
    mapConfig, errConvert := convertToMap(lines)
    
    if errConvert != nil {
        return nil, errConvert
    }
    
    outConfig := map[string]string{}
    
    for _, paramName := range (*keys) {
        
        value, isSet := mapConfig[paramName]
        
        if isSet {
            outConfig[paramName] = value
        } else {
            return nil, errorStack.Create("Brak klucza: " + paramName)
        }
    }
    
    return &File{config : outConfig}, nil
}


func convertToMap(lines *[]string) (map[string]string, *errorStack.Error) {
    
    configMap := map[string]string{}
    
    for _, line := range *lines {
        key, value, isParsed, errParse := parseLine(line)
        
        if errParse != nil {
            return nil, errParse
        }
        
        if isParsed {
            
            if _, isSet := configMap[key]; isSet {
                return nil, errorStack.Create("Zduplikowany klucz: " + key)
            }
            
            configMap[key] = value
        }
    }
    
    return configMap, nil
}


func parseLine(line string) (string, string, bool, *errorStack.Error) {
    
    if line[0] == "#"[0] {
        return "", "", false, nil
    }
    
    lineTrim := strings.TrimSpace(line)
    
    index := strings.Index(lineTrim, ":")
    
    if index >= 0 {
        
        key   := strings.TrimSpace(lineTrim[0:index])
        value := strings.TrimSpace(lineTrim[index+1:])
        
        if key == "" {
            return "", "", false, errorStack.Create("Pusty klucz: " + line)
        }
        
        if value == "" {
            return "", "", false, errorStack.Create("Pusty klucz: " + line)
        }
        
        return key, value, true, nil
    
    } else {
        
        return "", "", false, errorStack.Create("Błąd parsowania: " + lineTrim)
    }
}

func readLines(path string) (*[]string, *errorStack.Error) {
    
    file, err := os.Open(path)
    
    if err != nil {
        return nil, errorStack.From(err)
    }
    
    defer file.Close()
    
    var lines []string
    
    scanner := bufio.NewScanner(file)
    
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    
    errScaner := scanner.Err()
    
    if errScaner != nil {
        return nil, errorStack.From(errScaner)
    }
    
    return &lines, nil
}

