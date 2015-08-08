package htmlBuilder


import (
	"encoding/xml"
	"bytes"
)


/*
	escapowanie wartości
*/
func Escape(text string) string {

	escapeValue := &bytes.Buffer{}
	
	xml.Escape(escapeValue, []byte(text))
	
	return escapeValue.String()
}