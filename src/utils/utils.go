package utils


import (
    "time"
    "bytes"
    "strconv"
    "os/exec"
    "strings"
    "../errorStack"
)

func GetCurrentTimeName() string {
    
    current           := time.Now()
    year, montch, day := current.Date()
    hour              := current.Hour()
    minute            := current.Minute()
    second            := current.Second()
    
    return frm(year, 4) + frm(int(montch), 2) + frm(day, 2) + frm(hour, 2) + frm(minute, 2) + frm(second, 2)
}

func frm(liczba int, digit int) string {
    
    out := strconv.FormatInt(int64(liczba), 10)
    
    for len(out) < digit {
        out = "0" + out
    }
    
    return out
}

func ExecCommand(confRun string) (string, *errorStack.Error) {
    
    confRunSlice := strings.Fields(confRun)
        
    cmd := exec.Command(confRunSlice[0], confRunSlice[1:]...)
    
    var bufOut bytes.Buffer
    var bufErr bytes.Buffer
    
	cmd.Stdout = &bufOut
	cmd.Stderr = &bufErr
    
	err := cmd.Run()
    
    if err != nil {
        return "", errorStack.From(err)
	}
    
    if bufErr.String() != "" {
        return "", errorStack.Create("Na stumieniu błędów znalazły się jakieś dane: " + bufErr.String())
    }
    
    return bufOut.String(), nil
}

func ExtractEnvVarible(cmd string) ([]string, string) {
    
    chunks := strings.Fields(cmd)
    
    env := []string{}
    
    for {
        
        if isMatch, first, body := getHead(chunks); isMatch {

            if strings.Index(first, "=") >= 0 {

                env    = append(env, first)
                chunks = body

            } else {

                break
            }
        
        } else {
            
            break
        }
    }
    
    return env, strings.Join(chunks, " ")
}

func getHead(chunks []string) (bool, string, []string) {
    
    if len(chunks) == 0 {
        return false, "", []string{}
    } else if len(chunks) == 1 {
        return true, chunks[0], []string{}
    } else {
        return true, chunks[0], chunks[1:]
    }
}

