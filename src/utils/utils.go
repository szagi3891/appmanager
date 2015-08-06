package utils


import (
    "time"
    "strconv"
)

func GetCurrentTimeName() string {
    
    current           := time.Now()
    year, montch, day := current.Date()
    hour              := current.Hour()
    minute            := current.Minute()
    second            := current.Second()
    
    return frm(year, 4) + frm(int(montch), 2) + frm(day, 2) + frm(hour, 2) + frm(minute, 2) + frm(second, 2)
}

//    //frm(year, 4) + frm(int(montch), 2) + frm(day, 2) + frm(hour, 2) + frm(minute, 2) + frm(second, 2)


func frm(liczba int, digit int) string {
    
    out := strconv.FormatInt(int64(liczba), 10)
    
    for len(out) < digit {
        out = "0" + out
    }
    
    return out
}

