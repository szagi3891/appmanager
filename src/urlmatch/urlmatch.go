package urlmatch


import (
    "fmt"
    "strings"
    "strconv"
    urlModule "net/url"
    "../errorStack"
)


func splitRequestURI(path string) (*[]string, *errorStack.Error) {
    
    
    					//wytnij początkowy "/"
	
	if len(path) > 0 && path[0] == "/"[0] {
        
		path = path[1:]
                                                //wycinaj wszystkie końcowe znaki
        
        for len(path) > 0 && path[len(path)-1] == "/"[0] {
            path = path[:len(path)-1]
        }
        
    } else {
        
        return nil, errorStack.Create("Spodziewano się że pierwszy znak będzie znakiem /")
    }
	
    
    if path == "" {
        return &[]string{}, nil
    }
    
    chunks := strings.Split(path, "/")
    
    return &chunks, nil
}


func Create(RequestURI string) (*Url, *errorStack.Error) {
    
    
	url, err := urlModule.Parse(RequestURI)
	
	if err != nil {
		return nil, errorStack.From(err)
	}
	
    chunks, errChunks := splitRequestURI(url.Path)
    
    if errChunks != nil {
        return nil, errChunks
    }
    
    query := url.Query()
    
    return &Url{
        path   : chunks,
        query  : &query,
    }, nil
}


type Url struct {
    path   *[]string
    query  *urlModule.Values
}

func (self *Url) Dump() {
    fmt.Println("Dump:")
    fmt.Println("len : ", len(*(self.path)))
    fmt.Println("Dump end")
}

func (self *Url) String() string {
    return "/" + strings.Join(*(self.path), "/")
}

func matchBegin(path, params *[]string) (bool, *[]string) {
    
    if len(*path) >= len(*params)  {

        for i:=0; i<len(*params); i++ {

            if (*params)[i] != (*path)[i] {
                return false, nil
            }
        }
        
        newPath := (*path)[len(*params):]
        
        return true, &newPath
        
    } else {
        
        return false, nil
    }
}

func (self *Url) Match(params... string) (*Url, bool) {

    
    if isMatch, nextPath := matchBegin(self.path, &params); isMatch {
    
        newUrl := &Url{path: nextPath, query: self.query}
        return newUrl, true
    
    } else {
        return nil, false
    }
}

func (self *Url) MatchString(params... string) (string, *Url, bool) {
    
    if isMatch, nextPath := matchBegin(self.path, &params); isMatch {
        
        if len(*nextPath) > 0 {
            
            paramOne := (*(nextPath))[0]
            newPath  := (*(nextPath))[1:]
            
            return paramOne, &Url{path: &newPath, query: self.query}, true
            
        } else {
            return "", nil, false
        }
    
    } else {
        return "", nil, false
    }
}

func (self *Url) IsIndex() bool {
    
    return len(*(self.path)) == 0
}

func (self *Url) MatchInt(params... string) (int, *Url, bool) {
    
    if isMatch, nextPath := matchBegin(self.path, &params); isMatch {
        
        if len(*nextPath) > 0 {
            
            paramOne := (*(nextPath))[0]
            newPath  := (*(nextPath))[1:]
            
            valueInt, errConv := strconv.ParseInt(paramOne, 10, 64)
            
            if errConv != nil {
                return 0, nil, false
            }
            
            return int(valueInt), &Url{path: &newPath, query: self.query}, true
            
        } else {
            return 0, nil, false
        }
    
    } else {
        return 0, nil, false
    }
}

func (self *Url) MatchParamInt(name string) (int, bool) {
    
    value := self.query.Get(name)
    
    valueInt, errConv := strconv.ParseInt(value, 10, 64)
    
    if errConv != nil {
        return 0, false
    }
    
    return int(valueInt), true
}





	/*
	u, err := url.Parse("http://wolnachata.edu.pl:8080/swf/zplayer.swf/?autoplay=0&c1=000000&down=0&mp3=http%3A%2F%2Fwww.archive.org%2Fdownload%2FRadioWolneMedia-Audycje25.09.2011-Cz.2%2Finfra-fakty_-_2011-09-25.mp3")
	
	if err != nil {
		log.Fatal(err)
	}
	
	
	u.Scheme = "https"
	u.Host   = "google.com"
	
	fmt.Println(u.Path, "path")
	
	q := u.Query()
	q.Set("q", "golang")
	
	u.RawQuery = q.Encode()
	
	fmt.Println(u)
	*/
