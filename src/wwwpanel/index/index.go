package index

import (
    "strconv"
    "net/http"
    "../layout"
    "../../htmlBuilder"
    "../../backend"
)

func int2string(value int) string {
    return strconv.FormatInt(int64(value), 10)
}

func GetResponse(req *http.Request, appInfo *[]*backend.AppInfo, mainPort, activePort int, listBuild *[]string, lastCommitRepo string, isAvailableNewCommit bool) htmlBuilder.Node {
    
    root, body := layout.GetLayout()
    
    body.Tag("p").Text("Aktywne przekierowanie ruchu: " + int2string(mainPort) + " -> " + int2string(activePort))
    body.Tag("p").Html("&nbsp;")
    
    
    body.Tag("p").Text("Działające aplikacje:")
    
    table := body.Tag("table")
    
    tableHeader := table.Tag("tr")
    
    tableHeader.Tag("td").Text("Nazwa builda")
    tableHeader.Tag("td").Text("Port")
    tableHeader.Tag("td").Text("Aktywnych połączeń")
    tableHeader.Tag("td").Text("-")
    tableHeader.Tag("td").Text("-")
    
    
    for _, appItem := range *appInfo {
        
        line := table.Tag("tr")
        
        firstCeill := line.Tag("td")
        firstCeill.Text(appItem.Name)
        
        line.Tag("td").Text(strconv.FormatInt(int64(appItem.Port), 10))
        line.Tag("td").Text(strconv.FormatInt(int64(appItem.Active), 10))
        
        if activePort == appItem.Port {
            firstCeill.Attr("style", "color:red")
        }
        
        tdSet := line.Tag("td")
        
        if appItem.Port == activePort {
            tdSet.Html("&nbsp;")
        } else {
            link := "/proxyset/" + appItem.Name + "/" + strconv.FormatInt(int64(appItem.Port), 10)
            tdSet.Tag("a").Attr("href", link).Text("Ustaw jako główną")
        }
        
        tdDown := line.Tag("td")
        
        if appItem.Port == activePort || appItem.Active > 0 {
            tdDown.Html("&nbsp;")
        } else {
            link := "/down/" + appItem.Name + "/" + strconv.FormatInt(int64(appItem.Port), 10)
            tdDown.Tag("a").Attr("href", link).Text("Wyłącz")
        }
    }
    
    body.Tag("p").Html("&nbsp;")
    
    
    
    divMake := body.Tag("p").Text("Ostatni komit w repo : " + lastCommitRepo)
    body.Tag("p").Html("&nbsp;")
    
    if isAvailableNewCommit {
        
        divMake.Tag("span").Text(" - ")
        
        divMake.Tag("a").Attr("href", "/makebuild").Attr("onclick", "if (this.isclick != true) {this.isclick = true} else {return false}").Text("Twórz nowego builda")
    }
    
    
    tableBuild := body.Tag("table")
    
    tableBuild.Tag("tr").Tag("td").Text("Buildy:")
    
    for _, item := range *listBuild {
        
        line := tableBuild.Tag("tr")
        line.Tag("td").Text(item)
        line.Tag("td").Tag("a").Attr("href", "/start/" + item).Text("app start")
    }
    
    return root
    
    /*
        pobierz listę buildów
            sha1 -> nazwa
        
        pobierz indormację na temat ostatniego sha1 komita w repo

        pobranie na jakim porcie działą proxy oraz na którego builda kieruje

        pobranie jaki buildy działają
    */
}
