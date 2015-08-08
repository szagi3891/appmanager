package message

import (
    "../layout"
    "../../htmlBuilder"
)

func GetResponse(message string) htmlBuilder.Node {
    
    root, body := layout.GetLayout()
    
    body.Tag("p").Tag("a").Attr("href", "/").Text("Powrót do strony głównej")
    body.Tag("pre").Text(message)
    
    return root
}
