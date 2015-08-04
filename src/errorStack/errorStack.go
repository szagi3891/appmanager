package errorStack


import (
	"runtime/debug"
	"strings"
	"reflect"
	"strconv"
)


type errorMessage struct {
	message string
	stack   string
}


func (self *errorMessage) String() string {
	
	return self.message + "\n\n" + self.stack
}





type Error struct {
	
	list *[]*errorMessage
}




func From(err error) *Error {
	
	mess := errorMessage{message : err.Error(), stack: string(debug.Stack())}

	errNew := Error{list : &[]*errorMessage{&mess}}
	
	return &errNew
}


func makeString(message string, params []interface{}) string {
	
		
	out := []string{message}
	
	for _, paramItem := range params {
		
		
		if value, isOk := paramItem.(string); isOk {
			
			out = append(out, value)
		
		} else if value, isOk := paramItem.(int); isOk {
			
			out = append(out, strconv.FormatInt(int64(value), 10))
		
		} else {
			
			panic("TODO - nieobsługiwany typ: " + reflect.TypeOf(paramItem).String())
		}
	}
	
	return strings.Join(out, ", ")
}


func CreateParams(message string, params... interface{}) *Error {
	
	mess := errorMessage{message : makeString(message, params), stack: string(debug.Stack())}

	err := Error{list : &[]*errorMessage{&mess}}
	
	return &err
}


func Create(message string) *Error {
	
	mess := errorMessage{message : makeString(message, []interface{}{}), stack: string(debug.Stack())}

	err := Error{list : &[]*errorMessage{&mess}}
	
	return &err
}


func (self *Error) Add(message string, params... interface{}) *Error {
	
	*(self.list) = append(*(self.list), &errorMessage{message : makeString(message, params), stack: string(debug.Stack())})
	
	return self
}


func (self *Error) String() string {
	
	
	out := []string{}
	
	
	for i := 0; i < len(*(self.list)); i++ {		
		out = append(out, (*(self.list))[i].String())
	}
	
	lineIn := getLine("-")
		
	return getLine("=") + strings.Join(out, lineIn) + getLine("=")
}


func getLine(char string) string {
	
	out := ""
	
	for i:=0; i<100; i++ {
		out = out + char
	}
	
	return "\n\n" + out + "\n\n"
}


