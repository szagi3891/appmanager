package wwwpanel


import (
    "net/http"
    "../errorStack"
    "../httpserver"
    urlmatchModule "../urlmatch"
    logrotorModule "../logrotor"
    "fmt"
    "sort"
    backendModule "../backend"
    layoutModule   "./layout"
    actionIndex    "./index"
    actionMessage  "./message"
)


func Start(port int64, logs *logrotorModule.Logs, manager *backendModule.Manager) *errorStack.Error {
    
    return httpserver.Start(port, func(out http.ResponseWriter, req *http.Request){
        
        //fmt.Fprint(out, resp.content)
        
        urlmatch, errUrlCreate := urlmatchModule.Create(req.RequestURI)
        
        if errUrlCreate != nil {
            
            fmt.Println(errUrlCreate)
            //sendResponse(out, req, log, nil, responseModule.CreatePage400(), timer.End(errUrlCreate))
            return
        }
        
        if urlmatch.IsIndex() {
            
            listApp := manager.GetAppList()
            
            //fmt.Println("listApp", listApp)
            //proxy.GetActiveBackend()
            
            //TODO - trzeba wyświetlić tą listę aplikacji
            mainPort   := manager.GetMainPort()
            activePort := manager.GetActiveBackend().Port()
            
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
            
            
            response := actionIndex.GetResponse(req, listApp, mainPort, activePort, listBuild, lastCommitRepo, isAvailableNewCommit)
            fmt.Fprint(out, response.String())
            
            //Akcja przełączająca beckend - trzeba zrobić mutex
            //przełączenie ma iść na konkretny numer portu oraz na konkretną nazwę builda
            return
        
        }
        
        if appName, nextUrl, isOk := urlmatch.MatchString("proxyset"); isOk {
        
            if port, _, isOk := nextUrl.MatchInt(); isOk {
                
                isOk := manager.SwitchByNameAndPort(appName, port)
                
                if isOk {
                    fmt.Fprint(out, layoutModule.GetRedirectMessage(2, "/", "Wyłączono poprawnie"))
                } else {
                    fmt.Fprint(out, layoutModule.GetRedirectMessage(2, "/", "Niewyłączono"))
                }
                
                return
            }
        }
        
        if appName, nextUrl, isOk := urlmatch.MatchString("down"); isOk {
        
            if port, _, isOk := nextUrl.MatchInt(); isOk {
                
                isOk := manager.DownByNameAndPort(appName, port)
                
                if isOk {
                    fmt.Fprint(out, layoutModule.GetRedirectMessage(2, "/", "Przełączono poprawnie"))
                } else {
                    fmt.Fprint(out, layoutModule.GetRedirectMessage(2, "/", "Nieprzełączono"))
                }
                
                return
            }
        }
        
        if appName, _, isOk := urlmatch.MatchString("start"); isOk {

            _, errStart := manager.New(appName)
            
            if errStart != nil {
                
                fmt.Println(errStart)
                panic("TODO")
            }
            
            layout, body := layoutModule.GetRedirect(2, "/")
            body.Tag("p").Text("Build wystartował")
                
            fmt.Fprint(out, layout.String())
            
            return
        }
        
        if _, isOk := urlmatch.Match("makebuild"); isOk {
            
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
                    
                    layout, body := layoutModule.GetRedirect(2, "/")
                    body.Tag("p").Text("Utworzono builda: " + newName)
                    
                    fmt.Fprint(out, layout.String())
                    
                } else {
                    
                    fmt.Fprint(out, actionMessage.GetResponse("Błąd budowania builda: " + errMake.String()).String())
                }
            
            } else {
                
                layout, body := layoutModule.GetRedirect(2, "/")
                body.Tag("p").Text("Build był juz zbudowany")
                
                fmt.Fprint(out, layout.String())
            }
            
            return
        
        }
        
        
        fmt.Fprint(out, "404")
    })
}


