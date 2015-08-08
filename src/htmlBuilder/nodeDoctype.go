package htmlBuilder



func createDoctype(value string) *nodeDoctype {
	
	node := createNode()
	
	newItem := &nodeDoctype{node, value}
	
	node.setParent(newItem)
	
	return newItem
}


type nodeDoctype struct {
	
	*node
	value string
}


func (self *nodeDoctype) toString(out *[]string) {
    
    *out = append(*out, "<!DOCTYPE " + self.value + ">")

    self.node.toStringChild(out)
}
