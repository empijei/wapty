package fuzz

import (
	"bytes"
	"fmt"
)

func nestedJSON(level int, fieldname string) (json string) {
	//TODO escape "
	var open = fmt.Sprintf("{\"%s\":", fieldname)
	outbuf := bytes.NewBuffer(nil)
	for i := 0; i < level; i++ {
		outbuf.WriteString(open)
	}
	outbuf.WriteString("{}")
	for i := 0; i < level; i++ {
		outbuf.WriteString("}")
	}
	return string(outbuf.Bytes())
}
