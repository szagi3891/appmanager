package config


import (
    "fmt"
    "strings"
    "os"
    "bufio"
    "path/filepath"
    "strconv"
    "../errorStack"
)


type File struct {
    appDir   string
    buildDir string
    portFrom int
    portTo   int
    goCmd    string
}


func Parse(path string) (*File, *errorStack.Error) {
    
    
    path, errAbsPath := filepath.Abs(path)
    
    if errAbsPath != nil {
        return nil, errorStack.From(errAbsPath)
    }
    
    
    lines, errRead := readLines(path)
    
    if errRead != nil {
        return nil, errRead
    }
    
    
    mapConfig, errConvert := convertToMap(lines)
    
    if errConvert != nil {
        return nil, errConvert
    }
    
    
    pathBase := filepath.Dir(path)
    
    
    
    configFile := File{}
    
    
    configFile.buildDir = pathBase + "/build";
    
    
    
    
    appDir, errAppDir := getFromMap(&mapConfig, "appdir", pathBase)
    
    if errAppDir != nil {
        return nil, errAppDir
    }
    
    configFile.appDir = appDir
    
    
    portFrom, errPortFrom := getInt(&mapConfig, "portfrom")
    
    if errPortFrom != nil {
        return nil, errPortFrom
    }
    
    configFile.portFrom = portFrom
    
    
    portTo, errPortTo := getInt(&mapConfig, "portto")
    
    if errPortTo != nil {
        return nil, errPortTo
    }
    
    configFile.portTo = portTo
    
    
    goCmd, errGoCmd := getFromMap(&mapConfig, "gocmd", pathBase)
    
    if errGoCmd != nil {
        return nil, errGoCmd
    }
    
    configFile.goCmd = goCmd
    
    
    
    fmt.Println(configFile)
    
    /*
    GetBuildDir
    "../appmanager_build"
    */
    
    //configFile.appdir = 
    /*
    keys := []string{"port", "appmain"}
    
    for _, paramName := range keys {
        
        value, isSet := mapConfig[paramName]
        
        if isSet {
            outConfig[paramName] = value
        } else {
            return nil, errorStack.Create("Brak klucza: " + paramName)
        }
    }
    */
    
    return &configFile, nil
}


func (self *File) GetPortFrom() int {
    
    return self.portFrom
}

func (self *File) GetGoCmd() string {
    
    return self.goCmd
}

func (self *File) GetPortTo() int {
    
    return self.portTo
}

func (self *File) GetAppDir() string {
    
    return self.appDir
}

func (self *File) GetBuildDir() string {
    
    return self.buildDir
}

func getInt(mapConfig *map[string]string, propName string) (int, *errorStack.Error) {
    
    value, isValue := (*mapConfig)[propName]
    
    if isValue == false {
        return 0, errorStack.Create("Brak zmiennej : " + propName)
    }
    
    valueParse, errParse := strconv.ParseInt(value, 10, 64)
    
    if errParse != nil {
        return 0, errorStack.From(errParse)
    }
    
    return int(valueParse), nil
}

func getFromMap(mapConfig *map[string]string, propName, pathBase string) (string, *errorStack.Error) {
    
    value, isValue := (*mapConfig)[propName]
    
    if isValue == false {
        return "", errorStack.Create("Brak zmiennej : " + propName)
    }
    
    valueAbs, errApp := filepath.Abs(pathBase + "/" + value)
    
    if errApp != nil {
        return "", errorStack.From(errApp)
    }
    
    return valueAbs, nil
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

