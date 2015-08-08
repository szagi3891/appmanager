package main


import (
	"os"
	"fmt"
	"./src/htmlBuilder"
)


func main() {

	html := htmlBuilder.Create("html")

	a := html.Create("a")
	a.Attr("href", "asdasdas").Text("jakiś link")
	
	
	html.Create("a").Attr("href", "ddasdas<<>>").Text("dasdasdas<<>>\"fdsdas")
	html.Create("a").Attr("href", "ddasdas<d>sadas").Html("dasdasdas<<>>\"fdsdas")
	
	value, err := html.ToString()

	if err != nil {
		
		fmt.Println(err)
		os.Exit(1)
	}
	
	fmt.Println("zawartość:")
	fmt.Println(value)
}