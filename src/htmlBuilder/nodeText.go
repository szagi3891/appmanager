package htmlBuilder


func Text(value string) *nodeText {
    
	node := createNode()
	
	newItem := &nodeText{node, value}

	node.setParent(newItem)
	
	return newItem
}


type nodeText struct {
	
	*node
	value string
}


func (self *nodeText) toString(out *[]string) {
    
    *out = append(*out, Escape(self.value))
    
    self.node.toStringChild(out)
}
