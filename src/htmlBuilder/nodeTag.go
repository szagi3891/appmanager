package htmlBuilder


import (
	"strings"
)


func Tag(name string) *nodeTag {

	node := createNode()
	
	newItem := &nodeTag{node, name, &map[string]string{}}

	node.setParent(newItem)
	
	return newItem
}



type nodeTag struct {
	
	*node
	name string
	attr *map[string]string
}


func (self *nodeTag) toString(out *[]string) {
    
    *out = append(*out, self.stringOpen())
    
    self.node.toStringChild(out)
    
    *out = append(*out, self.stringClose())
}


func (self *nodeTag) Attr(name, value string) *nodeTag {
	
	(*(self.attr))[name] = value

	return self
}




var singleTag map[string]bool = map[string]bool{

	"input" : true,
	"meta"  : true,
	"link"  : true,
	"br"    : true,
	"hr"    : true,
	"img"   : true,
	"param" : true,
}




func (self *nodeTag) stringOpen() string {
	
	selfClose := ""
	
	if _, isSingle := singleTag[self.name]; isSingle == true {
	
		selfClose = "/"
	}
	
	
	attrString := []string{}
	
	for name, value := range *(self.attr) {
		
		attrString = append(attrString, " " + name + "=\"" + Escape(value) + "\"")
	}
	
	return "<" + self.name + strings.Join(attrString, "") + selfClose + ">"
}



func (self *nodeTag) stringClose() string {
	
	_, isSingle := singleTag[self.name]
	
	if isSingle {
	
		return ""
	
	} else {
	
		return "</" + self.name + ">"
	}
}

