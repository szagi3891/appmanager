package htmlBuilder


func createHtml(value string) *nodeHtml {
	
	node := createNode()
	
	newItem := &nodeHtml{node, value}
	
	node.setParent(newItem)
	
	return newItem
}


type nodeHtml struct {
	
	*node
	value string
}


func (self *nodeHtml) toString(out *[]string) {
    
    *out = append(*out, self.value)

    self.node.toStringChild(out)
}
