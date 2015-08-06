package config


import (
    //"fmt"
    "strings"
    "os"
    "bufio"
    "path/filepath"
    "strconv"
    "../errorStack"
)


type File struct {
    appDir     string
    buildDir   string
    logDir     string
    portMain   int
    portFrom   int
    portTo     int
    goCmd      string
    appMain    string
    appUser    string
    gopath     string
    rotatesize int
    rotatetime int
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
    
    
    configFile.buildDir = pathBase + "/build"
    configFile.logDir   = pathBase + "/log"
    
    errCheckBuildDir := checkDirectory(configFile.buildDir)
    
    if errCheckBuildDir != nil {
        return nil, errCheckBuildDir
    }
    
    errCheckLogDir := checkDirectory(configFile.logDir)
    
    if errCheckLogDir != nil {
        return nil, errCheckLogDir
    }
    
    
    
    appDir, errAppDir := getFromMap(&mapConfig, "appdir", pathBase)
    
    if errAppDir != nil {
        return nil, errAppDir
    }
    
    configFile.appDir = appDir
    
    
    portMain, errPortMain := getInt(&mapConfig, "portmain")
    
    if errPortMain != nil {
        return nil, errPortMain
    }
    
    configFile.portMain = portMain
    
    
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
    
    
 
    
    rotatesize, errRotatesize := getInt(&mapConfig, "rotatesize")
    
    if errRotatesize != nil {
        return nil, errRotatesize
    }
    
    configFile.rotatesize = rotatesize
    
    
    rotatetime, errRotatetime := getInt(&mapConfig, "rotatetime")
    
    if errRotatetime != nil {
        return nil, errRotatetime
    }
    
    configFile.rotatetime = rotatetime
    
    
    
    
    
    goCmd, errGoCmd := getFromMap(&mapConfig, "gocmd", pathBase)
    
    if errGoCmd != nil {
        return nil, errGoCmd
    }
    
    configFile.goCmd = goCmd
    
    
    
    appMain, errAppMain := getFromMap(&mapConfig, "appmain", pathBase)
    
    if errAppMain != nil {
        return nil, errAppMain
    }
    
    configFile.appMain = appMain
    
    
    
    appUser, errAppUser := getFromMap(&mapConfig, "appuser", pathBase)
    
    if errAppUser != nil {
        return nil, errAppUser
    }
    
    configFile.appUser = appUser
    
    
    gopath, errGopath := getFromMap(&mapConfig, "gopath", pathBase)
    
    if errGopath != nil {
        return nil, errGopath
    }
    
    configFile.gopath = gopath
   
    
    
    
    return &configFile, nil
}


func checkDirectory(path string) *errorStack.Error {

	info, err := os.Stat(path)

	if err != nil {
        
        return errorStack.From(err)
	}
    
	if info.IsDir() == false {
        
        return errorStack.Create("it's a directory: " + path)
	}
    
    return nil
}

func (self *File) GetRotatesize() int {
    
    return self.rotatesize
}
    
func (self *File) GetRotatetime() int {
    
    return self.rotatetime
}

func (self *File) GetPortMain() int {
    return self.portMain
}

func (self *File) GetGopath() string {
    
    return self.gopath
}

func (self *File) GetAppUser() string {
    
    return self.appUser
}

func (self *File) GetPortFrom() int {
    
    return self.portFrom
}



func (self *File) GetAppMain() string {
    
    return self.appMain
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

func (self *File) GetLogDir() string {
    
    return self.logDir
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

func isRelative(path string) bool {
        
    if len(path) >= 2 && path[0:2] == "./" {
        return true
    } else
    
    if len(path) >= 3 && path[0:3] == "../" {
        return true
    }
    
    return false
}

func getFromMap(mapConfig *map[string]string, propName, pathBase string) (string, *errorStack.Error) {
    
    value, isValue := (*mapConfig)[propName]
    
    if isValue == false {
        return "", errorStack.Create("Brak zmiennej : " + propName)
    }
    
    if isRelative(value) {
        
        valueAbs, errApp := filepath.Abs(pathBase + "/" + value)

        if errApp != nil {
            return "", errorStack.From(errApp)
        }

        return valueAbs, nil
        
    } else {
        
        return value, nil
    }
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
    
    indexHash := strings.Index(line, "#")
    
    lineNoHash := line
    
    if indexHash >= 0 {
        lineNoHash = lineNoHash[0:indexHash]
    }
    
    
    lineTrim := strings.TrimSpace(lineNoHash)
    
    
    if lineTrim == "" {
        return "", "", false, nil
    }
    
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

