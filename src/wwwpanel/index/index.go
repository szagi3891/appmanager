package index

import (
    "net/http"
    "../layout"
    "../../htmlBuilder"
)

func GetResponse(req *http.Request, listBuild *[]string, lastCommitRepo string, isAvailableNewCommit bool) htmlBuilder.Node {
    
    root, body := layout.GetLayout()
    
    
    divMake := body.Tag("p").Text("Ostatni komit w repo : " + lastCommitRepo)
    body.Tag("p").Html("&nbsp;")
    
    if isAvailableNewCommit {
        
        divMake.Tag("span").Text(" - ")
        
        divMake.Tag("a").Attr("href", "/makebuild").Attr("onclick", "if (this.isclick != true) {this.isclick = true} else {return false}").Text("Twórz nowego builda")
    }
    
    
    
    body.Tag("p").Text("buildy :")
    
    for _, item := range *listBuild {
        body.Tag("p").Text(item)
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
