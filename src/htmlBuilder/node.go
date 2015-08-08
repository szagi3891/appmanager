package htmlBuilder


import (
	"strings"
)


type Node interface {

	setParent(Node)
	getParent() Node
    toString(*[]string)
    Text(value string) *nodeText
    Tag(name string) *nodeTag
    Html(value string) *nodeHtml 
    Doctype(value string) *nodeDoctype
    String() string
}


func createNode() *node {
	
	return &node{
        parent : nil,
        child  : []Node{},
    }
}


type node struct {
    
	parent Node
	child  []Node
}


func (self *node) setParent(parent Node) {

	self.parent = parent
}


func (self *node) getParent() Node {

	return self.parent
}


func (self *node) Text(value string) *nodeText {
	
	newNode := Text(value)
	self.child = append(self.child, newNode)    
	return newNode
}


func (self *node) Tag(name string) *nodeTag {

	newNode := Tag(name)
	self.child = append(self.child, newNode)
	return newNode
}


func (self *node) Html(value string) *nodeHtml {
	
	newNode := createHtml(value)
	self.child = append(self.child, newNode)
	return newNode
}


func (self *node) Doctype(value string) *nodeDoctype {

	newNode := createDoctype(value)
	self.child = append(self.child, newNode)
	return newNode
}


func (self *node) toStringChild(out *[]string) {
	
	for _, child := range self.child {
		
		child.getParent().toString(out)
	}
}


func (self *node) String() (string) {
	
	out := []string{}
	
    self.getParent().toString(&out)
    
	return strings.Join(out, "")

}






