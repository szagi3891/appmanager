package layout

import (
    "strconv"
    "../../htmlBuilder"
)

func GetRedirect(time int, url string) (htmlBuilder.Node, htmlBuilder.Node) {
    
    root := htmlBuilder.Text("")
	
	root.Doctype("html")
    
	head := root.Tag("head")
    
    contentValue := strconv.FormatInt(int64(time), 10) + ";URL=\"" + url + "\""
    head.Tag("meta").Attr("http-equiv", "refresh").Attr("content", contentValue)
    
    head.Tag("style").Attr("type", "text/css").Html(getStyle())
    
    body := root.Tag("body")
    
    return root, body
}

func GetLayout() (htmlBuilder.Node, htmlBuilder.Node) {
    
    
	root := htmlBuilder.Text("")
	
	root.Doctype("html")
	
	head := root.Tag("head")
	
    
    head.Tag("title").Html("panel zarządzający")
	
    
	head.Tag("meta").Attr("charset", "utf-8")
	
    head.Tag("style").Attr("type", "text/css").Html(getStyle())
    
    body := root.Tag("body")
    
    return root, body
}

func getStyle() string {
    return `

* {
    margin:0;
    padding:0;
    border:0;
    font-family: monospace;
}

html {
	width: 100%;
	height: 100%;
}

body {
	width: 100%;
	height: 100%;
}

a {
	text-decoration: none;
}

a:hover {
	text-decoration: underline;
}

p {
	margin-top: 2px;
    margin-bottom: 2px;
}


*:focus {
   outline: none;
}

/* zerowanie tabel */

table {
	margin: 0px;
	border-collapse:collapse;
	empty-cells: show;
}

td {
	border: 0px;
	padding:0px;
	vertical-align: top;
}

th {
	border: 0px;
	padding:0px;
	vertical-align: top;
}

/* wyzerowanie textarea po to aby nie byĹo scrolli w textarea w ie */

textarea {
    overflow: auto;
    border: none;
}

`
}