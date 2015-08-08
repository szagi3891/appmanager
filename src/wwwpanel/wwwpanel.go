package wwwpanel


import (
    "net/http"
    "../errorStack"
    "../httpserver"
    urlmatchModule "../urlmatch"
    logrotorModule "../logrotor"
    "fmt"
    "sort"
    proxyModule   "../proxy"
    backendModule "../backend"
    layoutModule   "./layout"
    actionIndex    "./index"
    actionMessage  "./message"
)


func Start(port int64, appStderr *logrotorModule.LogWriter, manager *backendModule.Manager, proxy *proxyModule.Proxy) *errorStack.Error {
    
    return httpserver.Start(port, func(out http.ResponseWriter, req *http.Request){
        
        //fmt.Fprint(out, resp.content)
        
        urlmatch, errUrlCreate := urlmatchModule.Create(req.RequestURI)
        
        if errUrlCreate != nil {
            
            fmt.Println(errUrlCreate)
            //sendResponse(out, req, log, nil, responseModule.CreatePage400(), timer.End(errUrlCreate))
            return
        }
        
        if urlmatch.IsIndex() {
            
            _, errAppGet := manager.GetAppList()
            
            if errAppGet != nil {
                fmt.Println(errAppGet)
                panic("TODO")
            }
            
            //proxy.GetActiveBackend()
            
            //TODO - trzeba wyświetlić tą listę aplikacji
            
            listBuild, errList := manager.GetListBuild()
            
            if errList != nil {
                
                fmt.Println(errList)
                panic("TODO")
            }
            
            lastCommitRepo, errRepo := manager.GetSha1Repo()
            
            if errRepo != nil {
                fmt.Println(errRepo)
                panic("TODO")
            }
            
            
            isAvailableNewCommit := manager.IsAvailableNewCommit(listBuild, lastCommitRepo)
            
            
            //sortuj od największego stringa
            sort.Sort(sort.Reverse(sort.StringSlice(*listBuild)))
            
            
            response := actionIndex.GetResponse(req, listBuild, lastCommitRepo, isAvailableNewCommit)
            fmt.Fprint(out, response.String())
        
        
        //Akcja przełączająca beckend - trzeba zrobić mutex
        //przełączenie ma iść na konkretny numer portu oraz na konkretną nazwę builda
        
        } else if _, isOk := urlmatch.Match("makebuild"); isOk {
            
            listBuild, errList := manager.GetListBuild()
            
            if errList != nil {
                
                fmt.Println(errList)
                panic("TODO")
            }
            
            lastCommitRepo, errRepo := manager.GetSha1Repo()
            
            if errRepo != nil {
                fmt.Println(errRepo)
                panic("TODO")
            }
            
            isAvailableNewCommit := manager.IsAvailableNewCommit(listBuild, lastCommitRepo)
            
            if isAvailableNewCommit {
                
                newName, errMake := manager.MakeBuild()
                
                if errMake == nil {
                    
                    layout, body := layoutModule.GetRedirect(3, "/")
                    body.Tag("p").Text("Utworzono builda: " + newName)
                    
                    fmt.Fprint(out, layout.String())
                    
                } else {
                    
                    fmt.Fprint(out, actionMessage.GetResponse("Błąd budowania builda: " + errMake.String()).String())
                }
            
            } else {
                
                layout, body := layoutModule.GetRedirect(3, "/")
                body.Tag("p").Text("Build był juz zbudowany")
                
                fmt.Fprint(out, layout.String())
            }
            
        } else {
            
            fmt.Fprint(out, "404")
        }
    })
}


